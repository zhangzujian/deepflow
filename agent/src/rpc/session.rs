/*
 * Copyright (c) 2024 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

use std::sync::{
    atomic::{AtomicU64, Ordering},
    Arc, Weak,
};
use std::time::{Duration, Instant};

use anyhow::{anyhow, Result};
use log::{debug, error, info};
use md5::{Digest, Md5};
use parking_lot::RwLock;
use tonic::transport::Channel;

use crate::{
    common::{DEFAULT_CONTROLLER_PORT, DEFAULT_CONTROLLER_TLS_PORT},
    exception::ExceptionHandler,
    trident::AgentId,
    utils::stats::{self, AtomicTimeStats},
};
use grpc::dial as grpc_dial;

use public::{
    counter::{Countable, Counter, CounterType, CounterValue, RefCountable},
    proto::agent::{self, Exception, PluginType, Status},
};

pub const DEFAULT_TIMEOUT: Duration = Duration::from_secs(5);
pub const SESSION_TIMEOUT: Duration = Duration::from_secs(30);

const GRPC_CALL_ENDPOINTS: [&str; 10] = [
    "push",
    "ntp",
    "upgrade",
    "sync",
    "genesis_sync",
    "kubernetes_api_sync",
    "get_kubernetes_cluster_id",
    "gpid_sync",
    "plugin",
    "prometheus_api_sync",
];

const PUSH_ENDPOINT: usize = 0;
const NTP_ENDPOINT: usize = 1;
const UPGRADE_ENDPOINT: usize = 2;
const SYNC_ENDPOINT: usize = 3;
const GENESIS_SYNC_ENDPOINT: usize = 4;
const KUBERNETES_API_SYNC_ENDPOINT: usize = 5;
const GET_KUBERNETES_CLUSTER_ID_ENDPOINT: usize = 6;
const GPID_SYNC_ENDPOINT: usize = 7;
const PLUGIN_ENDPOINT: usize = 8;

macro_rules! response_size {
    (push, $($_:ident),*) => {
        "is stream"
    };
    (upgrade, $($_:ident),*) => {
        "is stream"
    };
    (plugin, $($_:ident),*) => {
        "is stream"
    };
    ($_:ident,  $response:ident) => {
        format!(
            "{}B",
            $response
                .as_ref()
                .map(|r| r.get_ref().encoded_len())
                .unwrap_or_default(),
        )
    };
}

macro_rules! sync_grpc_call {
    ($self:ident, $func:ident, $request:ident, $enpoint:ident) => {{
        use prost::Message;

        let prefix = std::concat!("grpc ", stringify!($func));

        log::trace!("{} prepare client", prefix);
        $self.update_current_server().await;
        let (channel, rx_size) = match $self.get_client() {
            Some(c) => c,
            None => {
                $self.set_request_failed(true);
                return Err(tonic::Status::cancelled("grpc client not connected"));
            }
        };
        let mut client = agent::synchronizer_client::SynchronizerClient::new(channel)
            .max_decoding_message_size(rx_size);

        let request_len = $request.encoded_len();
        let now = Instant::now();
        log::trace!("{} send request", prefix);
        let response = client.$func($request).await;
        log::trace!("{} receive response", prefix);
        let now_elapsed = now.elapsed();
        $self.counters[$enpoint].delay.update(now_elapsed);
        if log::log_enabled!(log::Level::Debug) {
            debug!(
                "{} latency {:?}ms request {}B response {}",
                prefix,
                now_elapsed.as_millis(),
                request_len,
                response_size!($func, response),
            );
        }
        response
    }};
}

struct Client {
    channel: Option<Channel>,
    // max receiving message size
    rx_size: usize,
}

pub struct Session {
    controller_cert_file_prefix: String,

    server_dispatcher: RwLock<ServerDispatcher>,

    version: AtomicU64,
    client: RwLock<Client>,
    exception_handler: ExceptionHandler,
    counters: Vec<Arc<GrpcCallCounter>>,
}

impl Session {
    pub fn new(
        port: u16,
        tls_port: u16,
        timeout: Duration,
        controller_cert_file_prefix: String,
        controller_ips: Vec<String>,
        exception_handler: ExceptionHandler,
        stats_collector: &stats::Collector,
    ) -> Session {
        let counters = (0..GRPC_CALL_ENDPOINTS.len())
            .into_iter()
            .map(|_| Arc::new(GrpcCallCounter::default()))
            .collect::<Vec<_>>();

        for (endpoint, counter) in counters.iter().enumerate() {
            stats_collector.register_countable(
                &stats::SingleTagModule("grpc_call", "endpoint", GRPC_CALL_ENDPOINTS[endpoint]),
                Countable::Ref(Arc::downgrade(&counter) as Weak<dyn RefCountable>),
            );
        }

        let config = Config {
            ips: controller_ips,
            port,
            tls_port,
            timeout,
            enable_tls: controller_cert_file_prefix.len() > 0,
            ..Default::default()
        };

        Session {
            server_dispatcher: RwLock::new(ServerDispatcher::new(config)),
            version: AtomicU64::new(0),
            client: RwLock::new(Client {
                channel: None,
                rx_size: crate::config::config::Communication::default().grpc_buffer_size,
            }),
            exception_handler,
            counters,
            controller_cert_file_prefix,
        }
    }

    pub fn reset_server_ip(&self, controller_ips: Vec<String>) {
        self.server_dispatcher
            .write()
            .update_controller_ips(controller_ips);
    }

    pub fn reset(&self) {
        self.close();
        self.server_dispatcher.write().reset();
    }

    async fn dial(&self, remote: &str, remote_port: u16, controller_cert_file_prefix: &str) {
        match grpc_dial(remote, remote_port, controller_cert_file_prefix).await {
            Ok(channel) => {
                self.client.write().channel.replace(channel);
            }
            Err(e) => {
                self.exception_handler.set(Exception::ControllerSocketError);
                self.set_request_failed(true);
                error!("{}", e);
            }
        }
    }

    pub fn get_client(&self) -> Option<(Channel, usize)> {
        let c = self.client.read();
        match &c.channel {
            Some(channel) => Some((channel.clone(), c.rx_size)),
            _ => None,
        }
    }

    pub fn set_rx_size(&self, size: usize) {
        let mut c = self.client.write();
        if c.rx_size != size {
            c.rx_size = size;
            self.version.fetch_add(1, Ordering::SeqCst);
        }
    }

    pub fn get_current_server(&self) -> (String, u16) {
        self.server_dispatcher.read().get_current_ip()
    }

    pub async fn update_current_server(&self) -> bool {
        let changed = self.server_dispatcher.write().update_current_ip();
        if changed || self.get_client().is_none() {
            let (ip, port) = self.server_dispatcher.read().get_current_ip();
            self.dial(&ip, port, &self.controller_cert_file_prefix)
                .await;
            self.version.fetch_add(1, Ordering::SeqCst);
        }
        changed
    }

    pub fn get_version(&self) -> u64 {
        self.version.load(Ordering::SeqCst)
    }

    pub fn close(&self) {
        self.client.write().channel.take();
    }

    pub fn get_request_failed(&self) -> bool {
        self.server_dispatcher.read().get_request_failed()
    }

    pub fn set_request_failed(&self, failed: bool) {
        self.server_dispatcher.write().set_request_failed(failed);
    }

    pub fn get_proxy_server(&self) -> (Option<String>, u16) {
        let d = self.server_dispatcher.read();
        (d.get_proxy_ip(), d.get_proxy_port())
    }

    pub fn set_proxy_server(&self, ip: Option<String>, port: u16) {
        let mut d = self.server_dispatcher.write();
        d.set_proxy_ip(ip);
        d.set_proxy_port(port);
    }

    pub async fn grpc_push_with_statsd(
        &self,
        request: agent::SyncRequest,
    ) -> Result<tonic::Response<tonic::codec::Streaming<agent::SyncResponse>>, tonic::Status> {
        sync_grpc_call!(self, push, request, PUSH_ENDPOINT)
    }

    async fn grpc_sync_inner(
        &self,
        request: agent::SyncRequest,
        with_statsd: bool,
    ) -> Result<tonic::Response<agent::SyncResponse>, tonic::Status> {
        log::trace!("grpc sync prepare client");
        self.update_current_server().await;
        let (channel, rx_size) = match self.get_client() {
            Some(c) => c,
            None => {
                self.set_request_failed(true);
                return Err(tonic::Status::cancelled("grpc client not connected"));
            }
        };
        let mut client = agent::synchronizer_client::SynchronizerClient::new(channel)
            .max_decoding_message_size(rx_size);

        if !with_statsd {
            log::trace!("grpc sync send request");
            let response = client.sync(request).await;
            log::trace!("grpc sync receive response");
            response
        } else {
            let now = Instant::now();
            log::trace!("grpc sync send request");
            let response = client.sync(request).await;
            log::trace!("grpc sync receive response");
            let now_elapsed = now.elapsed();
            self.counters[SYNC_ENDPOINT].delay.update(now_elapsed);
            debug!("grpc sync latency {:?}ms", now_elapsed.as_millis());
            response
        }
    }

    pub async fn grpc_sync(
        &self,
        request: agent::SyncRequest,
    ) -> Result<tonic::Response<agent::SyncResponse>, tonic::Status> {
        self.grpc_sync_inner(request, false).await
    }

    pub async fn grpc_sync_with_statsd(
        &self,
        request: agent::SyncRequest,
    ) -> Result<tonic::Response<agent::SyncResponse>, tonic::Status> {
        self.grpc_sync_inner(request, true).await
    }

    pub async fn grpc_upgrade_with_statsd(
        &self,
        request: agent::UpgradeRequest,
    ) -> Result<tonic::Response<tonic::codec::Streaming<agent::UpgradeResponse>>, tonic::Status>
    {
        sync_grpc_call!(self, upgrade, request, UPGRADE_ENDPOINT)
    }

    pub async fn grpc_ntp_with_statsd(
        &self,
        request: agent::NtpRequest,
    ) -> Result<tonic::Response<agent::NtpResponse>, tonic::Status> {
        // Ntp rpc name is `query`
        sync_grpc_call!(self, query, request, NTP_ENDPOINT)
    }

    pub async fn grpc_genesis_sync_with_statsd(
        &self,
        request: agent::GenesisSyncRequest,
    ) -> Result<tonic::Response<agent::GenesisSyncResponse>, tonic::Status> {
        sync_grpc_call!(self, genesis_sync, request, GENESIS_SYNC_ENDPOINT)
    }

    pub async fn grpc_kubernetes_api_sync_with_statsd(
        &self,
        request: agent::KubernetesApiSyncRequest,
    ) -> Result<tonic::Response<agent::KubernetesApiSyncResponse>, tonic::Status> {
        sync_grpc_call!(
            self,
            kubernetes_api_sync,
            request,
            KUBERNETES_API_SYNC_ENDPOINT
        )
    }

    pub async fn grpc_get_kubernetes_cluster_id_with_statsd(
        &self,
        request: agent::KubernetesClusterIdRequest,
    ) -> Result<tonic::Response<agent::KubernetesClusterIdResponse>, tonic::Status> {
        sync_grpc_call!(
            self,
            get_kubernetes_cluster_id,
            request,
            GET_KUBERNETES_CLUSTER_ID_ENDPOINT
        )
    }

    pub async fn grpc_gpid_sync(
        &self,
        request: agent::GpidSyncRequest,
    ) -> Result<tonic::Response<agent::GpidSyncResponse>, tonic::Status> {
        sync_grpc_call!(self, gpid_sync, request, GPID_SYNC_ENDPOINT)
    }

    pub async fn grpc_plugin(
        &self,
        request: agent::PluginRequest,
    ) -> Result<tonic::Response<tonic::codec::Streaming<agent::PluginResponse>>, tonic::Status>
    {
        sync_grpc_call!(self, plugin, request, PLUGIN_ENDPOINT)
    }

    pub async fn grpc_get_plugin(
        &self,
        name: &str,
        plugin_type: PluginType,
        agent_id: &AgentId,
    ) -> Result<Vec<u8>> {
        let s = self
            .grpc_plugin(agent::PluginRequest {
                ctrl_ip: Some(agent_id.ipmac.ip.to_string()),
                ctrl_mac: Some(agent_id.ipmac.mac.to_string()),
                plugin_type: Some(plugin_type as i32),
                plugin_name: Some(name.into()),
                team_id: Some(agent_id.team_id.clone()),
            })
            .await?;

        let mut data = vec![];
        let mut msg = s.into_inner();
        let mut total_len = 0u64;
        let mut msg_md5 = String::new();
        while let Some(message) = msg.message().await? {
            if message.status.unwrap_or_default() != Status::Success as i32 {
                return Err(anyhow!("fetch wasm prog fail, server return non success"));
            }
            if let Some(d) = message.content {
                data.extend(d);
            }
            total_len = message.total_len.unwrap_or_default();
            if msg_md5.is_empty() {
                msg_md5 = message.md5.unwrap_or_default();
            }
        }
        if data.is_empty() || data.len() != total_len as usize {
            return Err(anyhow!("fetch wasm prog fail, length incorrect"));
        }
        let md5_digest = Md5::new().chain_update(&data[..]).finalize();
        match hex::decode(msg_md5.as_bytes()) {
            Ok(bs) if &bs[..] != md5_digest.as_slice() => {
                return Err(anyhow!("fetch wasm prog fail, md5 checksum incorrect"))
            }
            Err(_) => {
                return Err(anyhow!(
                    "fetch wasm prog fail, invalid md5 checksum in message"
                ))
            }
            _ => (),
        }
        debug!(
            "pulled {:?} plugin {} with len {} checksum {}",
            plugin_type, name, total_len, msg_md5
        );
        Ok(data)
    }
}

struct Config {
    ips: Vec<String>,
    port: u16,
    tls_port: u16,
    proxy_ip: Option<String>,
    proxy_port: u16,
    timeout: Duration,
    enable_tls: bool,
}

impl Default for Config {
    fn default() -> Self {
        Config {
            ips: vec![],
            proxy_ip: None,
            port: DEFAULT_CONTROLLER_PORT,
            tls_port: DEFAULT_CONTROLLER_TLS_PORT,
            proxy_port: DEFAULT_CONTROLLER_PORT,
            timeout: DEFAULT_TIMEOUT,
            enable_tls: false,
        }
    }
}

struct ServerDispatcher {
    config: Config,

    current_ip: String,
    current_port: u16,
    current_ip_index: usize,

    proxied: bool,
    request_failed: bool,
}

impl ServerDispatcher {
    fn new(config: Config) -> ServerDispatcher {
        ServerDispatcher {
            config,

            current_ip_index: 0,
            current_ip: String::new(),
            current_port: 0,

            proxied: false,
            request_failed: false,
        }
    }

    fn reset(&mut self) {
        self.current_ip_index = 0;
        self.current_ip = String::new();
        self.current_port = 0;
        self.proxied = false;
        self.request_failed = false;
    }

    fn update_controller_ips(&mut self, controller_ips: Vec<String>) {
        self.reset();
        self.config.ips = controller_ips;
    }

    fn get_current_ip(&self) -> (String, u16) {
        (self.current_ip.clone(), self.current_port)
    }

    fn set_current_ip(&mut self, ip: String) {
        self.current_ip = ip;
    }

    fn get_proxy_ip(&self) -> Option<String> {
        self.config.proxy_ip.clone()
    }

    fn set_proxy_ip(&mut self, ip: Option<String>) {
        self.config.proxy_ip = ip;
    }

    fn is_proxy_ip(&self) -> bool {
        return self.proxied;
    }

    fn get_port(&self, is_proxy: bool) -> u16 {
        if is_proxy {
            return self.config.proxy_port;
        }
        if self.config.enable_tls {
            return self.config.tls_port;
        }
        return self.config.port;
    }

    pub fn set_proxy_port(&mut self, port: u16) {
        self.config.proxy_port = port;
    }

    fn get_proxy_port(&self) -> u16 {
        return self.config.proxy_port;
    }

    fn get_request_failed(&self) -> bool {
        self.request_failed
    }

    fn set_request_failed(&mut self, failed: bool) {
        self.request_failed = failed;
    }

    fn get_current_controller_ip(&self) -> String {
        // controller_ips一定不为空
        self.config.ips[self.current_ip_index].clone()
    }

    fn next_controller_ip(&mut self) {
        self.current_ip_index += 1;
        if self.current_ip_index >= self.config.ips.len() {
            self.current_ip_index = 0;
        }
    }

    fn update_current_ip(&mut self) -> bool {
        if self.current_ip.len() == 0 {
            self.current_ip = self.get_current_controller_ip();
            self.current_port = self.get_port(false);
            // 第一次访问，直接返回
            return true;
        }

        match (self.proxied, self.request_failed) {
            // 访问代理控制器失败，重连控制器
            (true, true) => {
                let new_ip = self.get_current_controller_ip();
                info!(
                    "rpc IP changed to controller {} from unavailable proxy {}",
                    new_ip, self.current_ip
                );
                self.current_ip = new_ip;
                self.current_port = self.get_port(false);
                self.proxied = false;
                true
            }
            // 成功访问代理控制器
            (true, false) => {
                let proxy_port = self.get_proxy_port();
                let proxy_ip = self.get_proxy_ip().unwrap();
                if proxy_port != self.current_port || self.current_ip != proxy_ip {
                    info!(
                        "rpc Proxy changed to proxy {} {} from proxy {} {}",
                        proxy_ip, proxy_port, self.current_ip, self.current_port
                    );
                    // 配置变更需要更新
                    self.current_port = proxy_port;
                    self.current_ip = proxy_ip;
                    true
                } else {
                    false
                }
            }
            // 访问控制器失败，更新控制器IP地址
            (false, true) => {
                self.next_controller_ip();
                let port = self.get_port(false);
                let ip = self.get_current_controller_ip();
                info!(
                    "rpc IP changed to controller {} {} from unavailable controller {} {}",
                    ip, port, self.current_ip, self.current_port
                );
                self.current_port = port;
                self.current_ip = ip;

                true
            }
            // 访问控制器成功，切换为代理控制器
            (false, false) => {
                let Some(proxy_ip) = self.get_proxy_ip() else {
                    return false;
                };
                let proxy_port = self.get_proxy_port();
                self.proxied = true;

                if self.current_port != proxy_port || self.current_ip != proxy_ip {
                    info!(
                        "rpc IP changed to proxy {} {} from controller {} {}",
                        proxy_ip, proxy_port, self.current_ip, self.current_port
                    );
                    self.current_port = proxy_port;
                    self.current_ip = proxy_ip;
                    true
                } else {
                    info!(
                        "rpc proxy {} {} and controller {} {} are the same, not updated",
                        proxy_ip, proxy_port, proxy_ip, proxy_port
                    );
                    false
                }
            }
        }
    }
}

#[derive(Default)]
pub struct GrpcCallCounter {
    pub delay: AtomicTimeStats,
}

impl RefCountable for GrpcCallCounter {
    fn get_counters(&self) -> Vec<Counter> {
        let max_delay = self.delay.max_ns.swap(0, Ordering::Relaxed) / 1000;
        let delay_count = self.delay.count.swap(0, Ordering::Relaxed) as u64;
        let sum = self.delay.sum_ns.swap(0, Ordering::Relaxed) / 1000;
        let avg_delay = if delay_count == 0 {
            0
        } else {
            sum / delay_count
        };
        vec![
            (
                "max_delay",
                CounterType::Gauged,
                CounterValue::Unsigned(max_delay),
            ),
            (
                "avg_delay",
                CounterType::Gauged,
                CounterValue::Unsigned(avg_delay),
            ),
            (
                "delay_count",
                CounterType::Gauged,
                CounterValue::Unsigned(delay_count),
            ),
        ]
    }
}
