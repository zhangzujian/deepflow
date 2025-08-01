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

package dbcache

import (
	"gorm.io/gorm"

	"github.com/deepflowio/deepflow/server/agent_config"
	. "github.com/deepflowio/deepflow/server/controller/common"
	models "github.com/deepflowio/deepflow/server/controller/db/metadb/model" // FIXME: To avoid ambiguity, name the package either mysql_model or db_model.
	"github.com/deepflowio/deepflow/server/controller/trisolaris/config"
	dbmgr "github.com/deepflowio/deepflow/server/controller/trisolaris/dbmgr"
	. "github.com/deepflowio/deepflow/server/controller/trisolaris/utils"
	"github.com/deepflowio/deepflow/server/libs/logger"
)

var log = logger.MustGetLogger("trisolaris.dbcache")

type DBDataCache struct {
	networks                []*models.Network
	vms                     []*models.VM
	vRouters                []*models.VRouter
	vpcs                    []*models.VPC
	subnets                 []*models.Subnet
	dhcpPorts               []*models.DHCPPort
	pods                    []*models.Pod
	vInterfaces             []*models.VInterface
	skipVTaps               []*models.VTap
	azs                     []*models.AZ
	hostDevices             []*models.Host
	podNodes                []*models.PodNode
	domains                 []*models.Domain
	subDomains              []*models.SubDomain
	chVTapPorts             []*models.ChVTapPort
	sysConfigurations       []*models.SysConfiguration
	wanIPs                  []*models.WANIP
	lanIPs                  []*models.LANIP
	floatingIPs             []*models.FloatingIP
	regions                 []*models.Region
	peerConnections         []*models.PeerConnection
	podServices             []*models.PodService
	podServicePorts         []*models.PodServicePort
	redisInstances          []*models.RedisInstance
	rdsInstances            []*models.RDSInstance
	podGroupPorts           []*models.PodGroupPort
	podGroups               []*models.PodGroup
	podClusters             []*models.PodCluster
	lbs                     []*models.LB
	lbTargetServers         []*models.LBTargetServer
	lbListeners             []*models.LBListener
	nats                    []*models.NATGateway
	vmPodNodeConns          []*models.VMPodNodeConnection
	vipDomains              []*models.Domain
	tapTypes                []*models.TapType
	resourceGroups          []*models.ResourceGroup
	resourceGroupExtraInfos []*models.ResourceGroupExtraInfo
	acls                    []*models.ACL
	npbTunnels              []*models.NpbTunnel
	npbPolicies             []*models.NpbPolicy
	pcapPolicies            []*models.PcapPolicy
	cens                    []*models.CEN
	processes               []*models.Process
	vips                    []*models.VIP
	customServices          []*models.CustomService

	podNSs     []*models.PodNamespace
	vtaps      []*models.VTap
	vtapGroups []*models.VTapGroup
	chDevices  []*models.ChDevice

	config *config.Config

	ORGID
}

func NewDBDataCache(orgID ORGID, cfg *config.Config) *DBDataCache {
	return &DBDataCache{
		config: cfg,
		ORGID:  orgID,
	}
}

func (d *DBDataCache) GetVms() []*models.VM {
	return d.vms
}

func (d *DBDataCache) GetVRouters() []*models.VRouter {
	return d.vRouters
}

func (d *DBDataCache) GetNetworks() []*models.Network {
	return d.networks
}

func (d *DBDataCache) GetVPCs() []*models.VPC {
	return d.vpcs
}

func (d *DBDataCache) GetSubnets() []*models.Subnet {
	return d.subnets
}

func (d *DBDataCache) GetDhcpPorts() []*models.DHCPPort {
	return d.dhcpPorts
}

func (d *DBDataCache) GetPods() []*models.Pod {
	return d.pods
}

func (d *DBDataCache) GetVInterfaces() []*models.VInterface {
	return d.vInterfaces
}

func (d *DBDataCache) GetAZs() []*models.AZ {
	return d.azs
}

func (d *DBDataCache) GetHostDevices() []*models.Host {
	return d.hostDevices
}

func (d *DBDataCache) GetPodNodes() []*models.PodNode {
	return d.podNodes
}

func (d *DBDataCache) GetPodGroupPorts() []*models.PodGroupPort {
	return d.podGroupPorts
}

func (d *DBDataCache) GetPodGroups() []*models.PodGroup {
	return d.podGroups
}

func (d *DBDataCache) GetPodClusters() []*models.PodCluster {
	return d.podClusters
}

func (d *DBDataCache) GetLBs() []*models.LB {
	return d.lbs
}

func (d *DBDataCache) GetLBTargetServers() []*models.LBTargetServer {
	return d.lbTargetServers
}

func (d *DBDataCache) GetLBListeners() []*models.LBListener {
	return d.lbListeners
}

func (d *DBDataCache) GetNats() []*models.NATGateway {
	return d.nats
}

func (d *DBDataCache) GetVmPodNodeConns() []*models.VMPodNodeConnection {
	return d.vmPodNodeConns
}

func (d *DBDataCache) GetVipDomains() []*models.Domain {
	return d.vipDomains
}

func (d *DBDataCache) GetAgentGroupUserConfigsFromDB(db *gorm.DB) []*agent_config.MySQLAgentGroupConfiguration {
	agentGroupConfigs, err := dbmgr.DBMgr[agent_config.MySQLAgentGroupConfiguration](db).Gets()
	if err != nil {
		log.Error(d.Log(err.Error()))
	}
	return agentGroupConfigs
}

func (d *DBDataCache) GetDomains() []*models.Domain {
	return d.domains
}

func (d *DBDataCache) GetSubDomains() []*models.SubDomain {
	return d.subDomains
}

func (d *DBDataCache) GetChVTapPorts() []*models.ChVTapPort {
	return d.chVTapPorts
}

func (d *DBDataCache) GetSysConfigurations() []*models.SysConfiguration {
	return d.sysConfigurations
}

func (d *DBDataCache) GetWANIPs() []*models.WANIP {
	return d.wanIPs
}

func (d *DBDataCache) GetLANIPs() []*models.LANIP {
	return d.lanIPs
}

func (d *DBDataCache) GetRegions() []*models.Region {
	return d.regions
}

func (d *DBDataCache) GetFloatingIPs() []*models.FloatingIP {
	return d.floatingIPs
}

func (d *DBDataCache) GetPodServices() []*models.PodService {
	return d.podServices
}

func (d *DBDataCache) GetPodServicePorts() []*models.PodServicePort {
	return d.podServicePorts
}

func (d *DBDataCache) GetRedisInstances() []*models.RedisInstance {
	return d.redisInstances
}

func (d *DBDataCache) GetRdsInstances() []*models.RDSInstance {
	return d.rdsInstances
}

func (d *DBDataCache) GetPeerConnections() []*models.PeerConnection {
	return d.peerConnections
}

func (d *DBDataCache) GetSkipVTaps() []*models.VTap {

	return d.skipVTaps
}

func (d *DBDataCache) GetResourceGroups() []*models.ResourceGroup {
	return d.resourceGroups
}

func (d *DBDataCache) GetResourceGroupExtraInfos() []*models.ResourceGroupExtraInfo {
	return d.resourceGroupExtraInfos
}

func (d *DBDataCache) GetACLs() []*models.ACL {
	return d.acls
}

func (d *DBDataCache) GetNpbTunnels() []*models.NpbTunnel {
	return d.npbTunnels
}

func (d *DBDataCache) GetNpbPolicies() []*models.NpbPolicy {
	return d.npbPolicies
}

func (d *DBDataCache) GetPcapPolicies() []*models.PcapPolicy {
	return d.pcapPolicies
}

func (d *DBDataCache) GetCENs() []*models.CEN {
	return d.cens
}

func (d *DBDataCache) GetProcesses() []*models.Process {
	return d.processes
}

func (d *DBDataCache) GetPodNSsIDAndName() []*models.PodNamespace {
	return d.podNSs
}

func (d *DBDataCache) GetVTapsIDAndName() []*models.VTap {
	return d.vtaps
}

func (d *DBDataCache) GetVTapGroupsIDAndLcuuid() []*models.VTapGroup {
	return d.vtapGroups
}

func (d *DBDataCache) GetChDevicesIDTypeAndName() []*models.ChDevice {
	return d.chDevices
}

func (d *DBDataCache) GetVIPs() []*models.VIP {
	return d.vips
}

func (d *DBDataCache) GetCustomServices() []*models.CustomService {
	return d.customServices
}

// SetCustomServices sets the custom services for testing purposes
func (d *DBDataCache) SetCustomServices(services []*models.CustomService) {
	d.customServices = services
}

func GetTapTypesFromDB(db *gorm.DB) []*models.TapType {
	tapTypes, err := dbmgr.DBMgr[models.TapType](db).Gets()
	if err != nil {
		log.Error(err)
	}

	return tapTypes
}

func (d *DBDataCache) GetDataCacheFromDB(db *gorm.DB) {
	if db == nil {
		log.Error("no db connect")
		return
	}
	networks, err := dbmgr.DBMgr[models.Network](db).Gets()
	if err == nil {
		d.networks = networks
	} else {
		log.Error(d.Log(err.Error()))
	}
	vms, err := dbmgr.DBMgr[models.VM](db).Gets()
	if err == nil {
		d.vms = vms
	} else {
		log.Error(d.Log(err.Error()))
	}

	vRouters, err := dbmgr.DBMgr[models.VRouter](db).Gets()
	if err == nil {
		d.vRouters = vRouters
	} else {
		log.Error(d.Log(err.Error()))
	}

	vpcs, err := dbmgr.DBMgr[models.VPC](db).Gets()
	if err == nil {
		d.vpcs = vpcs
	} else {
		log.Error(d.Log(err.Error()))
	}

	subnets, err := dbmgr.DBMgr[models.Subnet](db).Gets()
	if err == nil {
		d.subnets = subnets
	} else {
		log.Error(d.Log(err.Error()))
	}

	dhcpPorts, err := dbmgr.DBMgr[models.DHCPPort](db).Gets()
	if err == nil {
		d.dhcpPorts = dhcpPorts
	} else {
		log.Error(d.Log(err.Error()))
	}

	podFields := []string{}
	if d.config.ExportersEnabled {
		podFields = []string{"id", "name", "epc_id", "container_ids", "pod_cluster_id", "pod_node_id", "pod_namespace_id", "pod_group_id", "az", "domain", "label"}
	} else {
		podFields = []string{"id", "name", "epc_id", "container_ids", "pod_cluster_id", "pod_node_id", "pod_namespace_id", "pod_group_id", "az", "domain"}
	}
	pods, err := dbmgr.DBMgr[models.Pod](db).GetFields(podFields)
	if err == nil {
		d.pods = pods
	} else {
		log.Error(d.Log(err.Error()))
	}

	vInterfaces, err := dbmgr.DBMgr[models.VInterface](db).GetFields([]string{
		"id", "devicetype", "deviceid", "iftype", "mac", "vmac", "subnetid", "name", "region", "lcuuid", "domain", "sub_domain",
	})
	if err == nil {
		d.vInterfaces = vInterfaces
	} else {
		log.Error(d.Log(err.Error()))
	}

	azs, err := dbmgr.DBMgr[models.AZ](db).Gets()
	if err == nil {
		d.azs = azs
	} else {
		log.Error(d.Log(err.Error()))
	}

	hostDevices, err := dbmgr.DBMgr[models.Host](db).Gets()
	if err == nil {
		d.hostDevices = hostDevices
	} else {
		log.Error(d.Log(err.Error()))
	}

	podNodes, err := dbmgr.DBMgr[models.PodNode](db).Gets()
	if err == nil {
		d.podNodes = podNodes
	} else {
		log.Error(d.Log(err.Error()))
	}

	domains, err := dbmgr.DBMgr[models.Domain](db).Gets()
	if err == nil {
		d.domains = domains
	} else {
		log.Error(d.Log(err.Error()))
	}

	subDomains, err := dbmgr.DBMgr[models.SubDomain](db).Gets()
	if err == nil {
		d.subDomains = subDomains
	} else {
		log.Error(d.Log(err.Error()))
	}

	chVTapPorts, err := dbmgr.DBMgr[models.ChVTapPort](db).GetFields([]string{"vtap_id", "tap_port"})
	if err == nil {
		d.chVTapPorts = chVTapPorts
	} else {
		log.Error(d.Log(err.Error()))
	}

	sysConfigurations, err := dbmgr.DBMgr[models.SysConfiguration](db).Gets()
	if err == nil {
		d.sysConfigurations = sysConfigurations
	} else {
		log.Error(d.Log(err.Error()))
	}

	wanIPs, err := dbmgr.DBMgr[models.WANIP](db).GetFields([]string{"ip", "vifid", "netmask", "domain"})
	if err == nil {
		d.wanIPs = wanIPs
	} else {
		log.Error(d.Log(err.Error()))
	}

	lanIPs, err := dbmgr.DBMgr[models.LANIP](db).GetFields([]string{"ip", "vifid", "vl2id", "domain"})
	if err == nil {
		d.lanIPs = lanIPs
	} else {
		log.Error(d.Log(err.Error()))
	}

	floatingIPs, err := dbmgr.DBMgr[models.FloatingIP](db).Gets()
	if err == nil {
		d.floatingIPs = floatingIPs
	} else {
		log.Error(d.Log(err.Error()))
	}

	regions, err := dbmgr.DBMgr[models.Region](db).Gets()
	if err == nil {
		d.regions = regions
	} else {
		log.Error(d.Log(err.Error()))
	}

	peerConnections, err := dbmgr.DBMgr[models.PeerConnection](db).Gets()
	if err == nil {
		d.peerConnections = peerConnections
	} else {
		log.Error(d.Log(err.Error()))
	}

	podServices, err := dbmgr.DBMgr[models.PodService](db).GetFields([]string{
		"id", "name", "type", "service_cluster_ip", "pod_namespace_id", "pod_cluster_id", "epc_id", "az",
	})
	if err == nil {
		d.podServices = podServices
	} else {
		log.Error(d.Log(err.Error()))
	}

	podServicePorts, err := dbmgr.DBMgr[models.PodServicePort](db).Gets()
	if err == nil {
		d.podServicePorts = podServicePorts
	} else {
		log.Error(d.Log(err.Error()))
	}

	redisInstances, err := dbmgr.DBMgr[models.RedisInstance](db).Gets()
	if err == nil {
		d.redisInstances = redisInstances
	} else {
		log.Error(d.Log(err.Error()))
	}

	rdsInstances, err := dbmgr.DBMgr[models.RDSInstance](db).Gets()
	if err == nil {
		d.rdsInstances = rdsInstances
	} else {
		log.Error(d.Log(err.Error()))
	}

	podGroupPorts, err := dbmgr.DBMgr[models.PodGroupPort](db).Gets()
	if err == nil {
		d.podGroupPorts = podGroupPorts
	} else {
		log.Error(d.Log(err.Error()))
	}
	podGroups, err := dbmgr.DBMgr[models.PodGroup](db).GetFields([]string{"id", "name", "type"})
	if err == nil {
		d.podGroups = podGroups
	} else {
		log.Error(d.Log(err.Error()))
	}
	podClusters, err := dbmgr.DBMgr[models.PodCluster](db).Gets()
	if err == nil {
		d.podClusters = podClusters
	} else {
		log.Error(d.Log(err.Error()))
	}

	lbs, err := dbmgr.DBMgr[models.LB](db).Gets()
	if err == nil {
		d.lbs = lbs
	} else {
		log.Error(d.Log(err.Error()))
	}

	lbTargetServers, err := dbmgr.DBMgr[models.LBTargetServer](db).Gets()
	if err == nil {
		d.lbTargetServers = lbTargetServers
	} else {
		log.Error(d.Log(err.Error()))
	}

	lbListeners, err := dbmgr.DBMgr[models.LBListener](db).Gets()
	if err == nil {
		d.lbListeners = lbListeners
	} else {
		log.Error(d.Log(err.Error()))
	}

	nats, err := dbmgr.DBMgr[models.NATGateway](db).Gets()
	if err == nil {
		d.nats = nats
	} else {
		log.Error(d.Log(err.Error()))
	}

	vmPodNodeConns, err := dbmgr.DBMgr[models.VMPodNodeConnection](db).Gets()
	if err == nil {
		d.vmPodNodeConns = vmPodNodeConns
	} else {
		log.Error(d.Log(err.Error()))
	}

	DomainMgr := dbmgr.DBMgr[models.Domain](db)
	vipDomains, err := DomainMgr.GetBatchFromTypes([]int{QINGCLOUD_PRIVATE, CMB_CMDB})
	if err == nil {
		d.vipDomains = vipDomains
	} else {
		log.Error(d.Log(err.Error()))
	}

	skipVTaps, err := dbmgr.DBMgr[models.VTap](db).GetBatchFromTypes([]int{
		VTAP_TYPE_KVM, VTAP_TYPE_WORKLOAD_V, VTAP_TYPE_POD_VM})
	if err == nil {
		d.skipVTaps = skipVTaps
	} else {
		log.Error(d.Log(err.Error()))
	}

	resourceGroups, err := dbmgr.DBMgr[models.ResourceGroup](db).Gets()
	if err == nil {
		d.resourceGroups = resourceGroups
	} else {
		log.Error(d.Log(err.Error()))
	}

	resourceGroupExtraInfos, err := dbmgr.DBMgr[models.ResourceGroupExtraInfo](db).Gets()
	if err == nil {
		d.resourceGroupExtraInfos = resourceGroupExtraInfos
	} else {
		log.Error(d.Log(err.Error()))
	}

	acls, err := dbmgr.DBMgr[models.ACL](db).GetBatchFromState(ACL_STATE_ENABLE)
	if err == nil {
		d.acls = acls
	} else {
		log.Error(d.Log(err.Error()))
	}

	npbTunnels, err := dbmgr.DBMgr[models.NpbTunnel](db).Gets()
	if err == nil {
		d.npbTunnels = npbTunnels
	} else {
		log.Error(d.Log(err.Error()))
	}

	npbPolicies, err := dbmgr.DBMgr[models.NpbPolicy](db).GetBatchFromState(ACL_STATE_ENABLE)
	if err == nil {
		d.npbPolicies = npbPolicies
	} else {
		log.Error(d.Log(err.Error()))
	}

	pcapPolicies, err := dbmgr.DBMgr[models.PcapPolicy](db).GetBatchFromState(ACL_STATE_ENABLE)
	if err == nil {
		d.pcapPolicies = pcapPolicies
	} else {
		log.Error(d.Log(err.Error()))
	}

	cens, err := dbmgr.DBMgr[models.CEN](db).Gets()
	if err == nil {
		d.cens = cens
	} else {
		log.Error(d.Log(err.Error()))
	}

	processes, err := dbmgr.DBMgr[models.Process](db).Gets()
	if err == nil {
		d.processes = processes
	} else {
		log.Error(d.Log(err.Error()))
	}

	podNSs, err := dbmgr.DBMgr[models.PodNamespace](db).GetFields([]string{"id", "name"})
	if err == nil {
		d.podNSs = podNSs
	} else {
		log.Error(d.Log(err.Error()))
	}

	vtaps, err := dbmgr.DBMgr[models.VTap](db).GetFields([]string{"id", "name", "launch_server_id", "type", "vtap_group_lcuuid"})
	if err == nil {
		d.vtaps = vtaps
	} else {
		log.Error(d.Log(err.Error()))
	}

	vtapGroups, err := dbmgr.DBMgr[models.VTapGroup](db).GetFields([]string{"id", "lcuuid"})
	if err == nil {
		d.vtapGroups = vtapGroups
	} else {
		log.Error(d.Log(err.Error()))
	}

	chDevices, err := dbmgr.DBMgr[models.ChDevice](db).GetFields([]string{"devicetype", "deviceid", "name"})
	if err == nil {
		d.chDevices = chDevices
	} else {
		log.Error(d.Log(err.Error()))
	}

	vips, err := dbmgr.DBMgr[models.VIP](db).Gets()
	if err == nil {
		d.vips = vips
	} else {
		log.Error(d.Log(err.Error()))
	}

	var customServices []*models.CustomService
	err = db.Order("type asc, id asc").Find(&customServices).Error
	if err == nil {
		d.customServices = customServices
	} else {
		log.Error(d.Log(err.Error()))
	}
}
