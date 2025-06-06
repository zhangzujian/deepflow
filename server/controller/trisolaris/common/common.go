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

package common

import (
	"github.com/deepflowio/deepflow/message/agent"
	api "github.com/deepflowio/deepflow/message/trident"
)

const (
	// VTAP
	VTAP_CONTROLLER_EXCEPTIONS_MASK = 0xFFFFFFFF00000000
	VTAP_TRIDENT_EXCEPTIONS_MASK    = 0x00000000FFFFFFFF
	VTAP_NO_REGISTER_EXCEPTION      = 0x20000000

	VTAP_TYPE_HYPER_V_NETWORK = 11

	SHUT_DOWN_STR  = "关闭"
	SHUT_DOWN_UINT = 0xffffffff

	ALL_CLUSTERS    = 0
	CLUSTER_OF_VTAP = 1

	ALL_SIMPLE_PLATFORM_DATA            = 1
	ALL_SIMPLE_PLATFORM_DATA_EXCEPT_POD = 2
	DOMAIN_TO_ALL_SIMPLE_PLATFORM_DATA  = 3
	DOMAIN_TO_PLATFORM_DATA_EXCEPT_POD  = 4
	DOMAIN_TO_PLATFORM_DATA_ONLY_POD    = 5
	ALL_SKIP_SIMPLE_PLATFORM_DATA       = 6
	PLATFORM_DATA_TYPE_1                = 6
	PLATFORM_DATA_TYPE_2                = 7
	PLATFORM_DATA_TYPE_3                = 8
	PLATFORM_DATA_BM_DEDICATED          = 9

	SKIP_PLATFORM_DATA_TYPE_1               = 10
	SKIP_PLATFORM_DATA_TYPE_2               = 11
	SKIP_PLATFORM_DATA_TYPE_3               = 12
	DOMAIN_TO_SKIP_ALL_SIMPLE_PLATFORM_DATA = 13
	DOMAIN_TO_SKIP_PLATFORM_DATA_EXCEPT_POD = 14
	DOMAIN_TO_SKIP_PLATFORM_DATA_ONLY_POD   = 15

	INGESTER_ALL_PLATFORM_DATA            = 16
	ALL_COMPLETE_PLATFORM_DATA_EXCEPT_POD = 17
	REGION_TO_PLATFORM_DATA_ONLY_POD      = 18
	AZ_TO_PLATFORM_DATA_ONLY_POD          = 19
	PLATFORM_DATA_FOR_INGESTER_1          = 20
	PLATFORM_DATA_FOR_INGESTER_2          = 21
	PLATFORM_DATA_FOR_INGESTER_MERGE      = 22

	NO_DOMAIN_TO_PLATFORM = 22

	DEFAULT_MAX_MEMORY         = 256
	DEFAULT_MAX_ESCAPE_SECONDS = 3600

	CONTROLLER_VTAP_MAX = 2000
	TSDB_VTAP_MAX       = 200

	// TRIDENT OS
	TRIDENT_LINUX   = 0
	TRIDENT_WINDOWS = 1

	TSDB_PROCESS_NAME = "ingester"

	CONN_DEFAULT_AZ     = "ALL"
	CONN_DEFAULT_REGION = "ffffffff-ffff-ffff-ffff-ffffffffffff"

	NPB_BUSINESS_ID  = 1
	PCAP_BUSINESS_ID = -3

	RESOURCE_GROUP_TYPE_ANONYMOUS_POD_GROUP                = 6
	RESOURCE_GROUP_TYPE_VM                                 = 1
	RESOURCE_GROUP_TYPE_ANONYMOUS_VM                       = 3
	RESOURCE_GROUP_TYPE_ANONYMOUS_POD                      = 5
	RESOURCE_GROUP_TYPE_ANONYMOUS_POD_SERVICE              = 8
	RESOURCE_GROUP_TYPE_ANONYMOUS_POD_GROUP_AS_POD_SERVICE = 81
	RESOURCE_GROUP_TYPE_ANONYMOUS_VL2                      = 14

	POD_SERVICE_TYPE_NODE_PORT = 2

	INTERNET_RESOURCE_GROUP_ID_UINT32 = -2 & 0xffffffff
	INTERNET_EPC_ID_UINT32            = -2 & 0xffffffff
	ANY_EPC_ID_UINT32                 = -1 & 0xffffffff
	RESOURCE_GROUP_TYPE_NONE          = 0
	RESOURCE_GROUP_TYPE_ANONYMOUS_IP  = 4

	APPLICATION_ANALYSIS          = 1
	APPLICATION_FLOW_BACKTRACKING = 2
	APPLICATION_NPB               = 6
	APPLICATION_PCAP              = 4

	ALL_VTAP_SHARE_POLICY_VERSION_OFFSET = 1000000000
	INGESTER_POLICY_VERSION_OFFSET       = 2000000000

	PROTOCOL_ALL = 256

	MAX_PAYLOAD_SLICE = 65535

	ALL_K8S_CLUSTER             = 0
	K8S_CLUSTER_IN_LOCAL_REGION = 1
	K8S_CLUSTER_IN_LOCAL_AZS    = 2

	DISABLED = 0
	ENABLED  = 1
)

const (
	CONFIG_KEY_NTP_ENABLED                 = "global.ntp.enabled"
	CONFIG_KEY_INGESTER_IP                 = "global.communication.ingester_ip"
	CONFIG_KEY_INGESTER_PORT               = "global.communication.ingester_port"
	CONFIG_KEY_PROXY_CONTROLLER_IP         = "global.communication.proxy_controller_ip"
	CONFIG_KEY_PROXY_CONTROLLER_PORT       = "global.communication.proxy_controller_port"
	CONFIG_KEY_CAPTURE_MODE                = "inputs.cbpf.common.capture_mode"
	CONFIG_KEY_DOMAIN_FILTER               = "inputs.resources.pull_resource_from_controller.domain_filter"
	CONFIG_KEY_HYPERVISOR_RESOURCE_ENABLED = "inputs.resources.private_cloud.hypervisor_resource_enabled"
)

var (
	KWP_NORMAL         = api.KubernetesWatchPolicy_KWP_NORMAL
	KWP_WATCH_ONLY     = api.KubernetesWatchPolicy_KWP_WATCH_ONLY
	KWP_WATCH_DISABLED = api.KubernetesWatchPolicy_KWP_WATCH_DISABLED
)

var (
	AGENT_KWP_NORMAL         = agent.KubernetesWatchPolicy_KWP_NORMAL
	AGENT_KWP_WATCH_ONLY     = agent.KubernetesWatchPolicy_KWP_WATCH_ONLY
	AGENT_KWP_WATCH_DISABLED = agent.KubernetesWatchPolicy_KWP_WATCH_DISABLED
)

type UpgradeData struct {
	Content  []byte
	TotalLen uint64
	PktCount uint32
	Md5Sum   string
	Step     uint64
	K8sImage string
}
