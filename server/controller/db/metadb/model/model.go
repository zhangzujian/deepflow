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

package model

import (
	"bytes"
	"compress/zlib"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"

	"github.com/deepflowio/deepflow/server/libs/logger"
)

var log = logger.MustGetLogger("db.metadb")

type Business struct {
	ID          int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name        string    `gorm:"column:name;type:char(64);default:null" json:"NAME"`
	Description string    `gorm:"column:description;type:varchar(256);default:null" json:"DESCRIPTION"`
	Type        int       `gorm:"column:type;type:int;default:1" json:"TYPE"`        // 1-data center; 2-ip; 3-vpc; 4-WAN; 5-NPB; 6-diagnose; 21-tmp ip; 31-tmp vpc
	VPCID       int       `gorm:"column:epc_id;type:int;default:null" json:"VPC_ID"` // for vpc type
	NetworkID   int       `gorm:"column:vl2_id;type:int;default:null" json:"VL2_ID"` // for ip type
	State       int       `gorm:"column:state;type:int;default:1" json:"STATE"`      // 0-disable; 1-enable
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
	Lcuuid      string    `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
}

type ResourceGroup struct {
	ID            int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	BusinessID    int       `gorm:"column:business_id;type:int;not null" json:"BUSINESS_ID"`
	Lcuuid        string    `gorm:"column:lcuuid;type:varchar(64);not null" json:"LCUUID"`
	Name          string    `gorm:"column:name;type:varchar(200);not null;default:''" json:"NAME"`
	Type          int       `gorm:"column:type;type:int;not null" json:"TYPE"`            // 1:vm, 2:ip, 3: anonymous vm, 4: anonymous ip, 5: reserved for pod_group, 6: anonymous pod_group, 7: reserved for pod_service, 8: anonymous pod_service, 81: anonymous pod_service as pod_group, 9：lb_bk_rule, 10：reserved for anonymous lb_bk_rule, 11: tmp vm, 21: tmp ip, 13: reserve for vl2, 14: anonymous vl2
	IPType        int       `gorm:"column:ip_type;type:int;default:null" json:"IP_TYPE"`  // 1: single ip, 2: ip range, 3: cidr, 4.mix [1, 2, 3]
	IPs           string    `gorm:"column:ips;type:text;default:null" json:"IPS"`         // ips separated by ,
	VMIDs         string    `gorm:"column:vm_ids;type:text;default:null" json:"VM_IDS"`   // vm ids separated by ,
	NetworkIDs    string    `gorm:"column:vl2_ids;type:text;default:null" json:"VL2_IDS"` // vl2 ids separated by ,
	VPCID         *int      `gorm:"column:epc_id;type:int;default:null" json:"VPC_ID"`
	PodClusterID  int       `gorm:"column:pod_cluster_id;type:int;default:null" json:"POD_CLUSTER_ID"`
	PodGroupIDs   string    `gorm:"column:pod_group_ids;type:text;default:null" json:"POD_GROUP_IDS"`     // pod group ids separated by ,
	PodServiceIDs string    `gorm:"column:pod_service_ids;type:text;default:null" json:"POD_SERVICE_IDS"` // pod service ids separated by ,
	LBID          int       `gorm:"column:lb_id;type:int;default:null" json:"LB_ID"`
	LBListenerID  int       `gorm:"column:lb_listener_id;type:int;default:null" json:"LB_LISTENER_ID"`
	ExtraInfoIDs  string    `gorm:"column:extra_info_ids;type:string;default:null" json:"EXTRA_INFO_IDS"`
	IconID        int       `gorm:"column:icon_id;type:int;default:-2" json:"ICON_ID"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
}

type ResourceGroupPort struct {
	ID              int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name            string    `gorm:"column:name;type:varchar(256);default:''" json:"NAME"`
	Ports           string    `gorm:"column:ports;type:text;default:null" json:"PORTS"` // Save server ports list when type is customize
	BusinessID      int       `gorm:"column:business_id;type:int;not null" json:"BUSINESS_ID"`
	ResourceGroupID int       `gorm:"column:rg_id;type:int;not null" json:"RG_ID"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
	Lcuuid          string    `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
}

type TapType struct {
	ID             int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name           string `gorm:"column:name;type:char(64);not null" json:"NAME"`
	Type           int    `gorm:"column:type;type:int;not null;default:1" json:"TYPE"` // 1:packet, 2:sFlow, 3:NetFlow V5 4:NetStream v5
	Region         string `gorm:"column:region;type:char(64);default:null" json:"REGION"`
	Value          int    `gorm:"column:value;type:int;not null" json:"VALUE"`
	VLAN           int    `gorm:"column:vlan;type:int;default:null" json:"VLAN"`
	SrcIP          string `gorm:"column:src_ip;type:char(64);default:null" json:"SRC_IP"`
	InterfaceIndex uint   `gorm:"column:interface_index;type:int unsigned;default:null" json:"INTERFACE_INDEX"` // 1 ~ 2^32-1
	InterfaceName  string `gorm:"column:interface_name;type:char(64);default:null" json:"INTERFACE_NAME"`
	SamplingRate   uint   `gorm:"column:sampling_rate;type:int unsigned;default:null" json:"SAMPLING_RATE"` // 1 ~ 2^32-1
	Description    string `gorm:"column:description;type:varchar(256);default:null" json:"DESCRIPTION"`
	Lcuuid         string `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
}

type Controller struct {
	ID                 int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	State              int       `gorm:"column:state;type:int;default:null" json:"STATE"` // 0.Temp 1.Creating 2.Complete 3.Modifying 4.Exception
	Name               string    `gorm:"column:name;type:char(64);default:null" json:"NAME"`
	Description        string    `gorm:"column:description;type:varchar(256);default:null" json:"DESCRIPTION"`
	IP                 string    `gorm:"column:ip;type:char(64);default:null" json:"IP"`
	NATIP              string    `gorm:"column:nat_ip;type:char(64);default:null" json:"NAT_IP"`
	CPUNum             int       `gorm:"column:cpu_num;type:int;default:0" json:"CPU_NUM"` // logical number of cpu
	MemorySize         int64     `gorm:"column:memory_size;type:bigint;default:0" json:"MEMORY_SIZE"`
	Arch               string    `gorm:"column:arch;type:varchar(256);default:null" json:"ARCH"`
	Os                 string    `gorm:"column:os;type:varchar(256);default:null" json:"OS"`
	KernelVersion      string    `gorm:"column:kernel_version;type:varchar(256);default:null" json:"KERNEL_VERSION"`
	VTapMax            int       `gorm:"column:vtap_max;type:int;default:2000" json:"VTAP_MAX"`
	SyncedAt           time.Time `gorm:"column:synced_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"SYNCED_AT"`
	NATIPEnabled       int       `gorm:"column:nat_ip_enabled;default:0" json:"NAT_IP_ENABLED"` // 0: disabled 1:enabled
	NodeType           int       `gorm:"column:node_type;type:int;default:2" json:"NODE_TYPE"`  // region node type 1.master 2.slave
	RegionDomainPrefix string    `gorm:"column:region_domain_prefix;type:varchar(256);default:''" json:"REGION_DOMAIN_PREFIX"`
	NodeName           string    `gorm:"column:node_name;type:char(64);default:null" json:"NODE_NAME"`
	PodIP              string    `gorm:"column:pod_ip;type:char(64);default:null" json:"POD_IP"`
	PodName            string    `gorm:"column:pod_name;type:char(64);default:null" json:"POD_NAME"`
	CAMD5              string    `gorm:"column:ca_md5;type:char(64);default:null" json:"CA_MD5"`
	Lcuuid             string    `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
}

type AZControllerConnection struct {
	ID           int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	AZ           string `gorm:"column:az;type:char(64);default:ALL" json:"AZ"`
	Region       string `gorm:"column:region;type:char(64);default:ffffffff-ffff-ffff-ffff-ffffffffffff" json:"REGION"`
	ControllerIP string `gorm:"column:controller_ip;type:char(64);default:null" json:"CONTROLLER_IP"`
	Lcuuid       string `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
}

func (AZControllerConnection) TableName() string {
	return "az_controller_connection"
}

type Analyzer struct {
	ID                int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	State             int       `gorm:"column:state;type:int;default:null" json:"STATE"`    // 0.Temp 1.Creating 2.Complete 3.Modifying 4.Exception
	HaState           int       `gorm:"column:ha_state;type:int;default:1" json:"HA_STATE"` // 1.master 2.backup
	Name              string    `gorm:"column:name;type:char(64);default:null" json:"NAME"`
	Description       string    `gorm:"column:description;type:varchar(256);default:null" json:"DESCRIPTION"`
	IP                string    `gorm:"column:ip;type:char(64);default:null" json:"IP"`
	NATIP             string    `gorm:"column:nat_ip;type:char(64);default:null" json:"NAT_IP"`
	Agg               int       `gorm:"column:agg;type:int;default:1" json:"AGG"`
	CPUNum            int       `gorm:"column:cpu_num;type:int;default:0" json:"CPU_NUM"` // logical number of cpu
	MemorySize        int64     `gorm:"column:memory_size;type:bigint;default:0" json:"MEMORY_SIZE"`
	Arch              string    `gorm:"column:arch;type:varchar(256);default:null" json:"ARCH"`
	Os                string    `gorm:"column:os;type:varchar(256);default:null" json:"OS"`
	KernelVersion     string    `gorm:"column:kernel_version;type:varchar(256);default:null" json:"KERNEL_VERSION"`
	PcapDataMountPath string    `gorm:"column:pcap_data_mount_path;type:varchar(256);default:null" json:"PCAP_DATA_MOUNT_PATH"`
	VTapMax           int       `gorm:"column:vtap_max;type:int;default:200" json:"VTAP_MAX"`
	SyncedAt          time.Time `gorm:"column:synced_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"SYNCED_AT"`
	NATIPEnabled      int       `gorm:"column:nat_ip_enabled;default:0" json:"NAT_IP_ENABLED"` // 0: disabled 1:enabled
	PodIP             string    `gorm:"column:pod_ip;type:char(64);default:null" json:"POD_IP"`
	PodName           string    `gorm:"column:pod_name;type:char(64);default:null" json:"pod_name"`
	CAMD5             string    `gorm:"column:ca_md5;type:char(64);default:null" json:"CA_MD5"`
	Lcuuid            string    `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
}

type AZAnalyzerConnection struct {
	ID         int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	AZ         string `gorm:"column:az;type:char(64);default:ALL" json:"AZ"`
	Region     string `gorm:"column:region;type:char(64);default:ffffffff-ffff-ffff-ffff-ffffffffffff" json:"REGION"`
	AnalyzerIP string `gorm:"column:analyzer_ip;type:char(64);default:null" json:"ANALYZER_IP"`
	Lcuuid     string `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
}

func (AZAnalyzerConnection) TableName() string {
	return "az_analyzer_connection"
}

type VTap struct {
	ID                  int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name                string    `gorm:"column:name;type:varchar(256);not null" json:"NAME"`
	RawHostname         string    `gorm:"column:raw_hostname;type:varchar(256);" json:"RAW_HOSTNAME"`
	Owner               string    `gorm:"column:owner;type:varchar(64);default:''" json:"OWNER"`
	State               int       `gorm:"column:state;type:int;default:1" json:"STATE"`   // 0.not-connected 1.normal
	Enable              int       `gorm:"column:enable;type:int;default:1" json:"ENABLE"` // 0: stop 1: running
	Type                int       `gorm:"column:type;type:int;default:0" json:"TYPE"`     // 1: process 2: vm 3: public cloud 4: analyzer 5: physical machine 6: dedicated physical machine 7: host pod 8: vm pod
	CtrlIP              string    `gorm:"column:ctrl_ip;type:char(64);not null" json:"CTRL_IP"`
	CtrlMac             string    `gorm:"column:ctrl_mac;type:char(64);default:null" json:"CTRL_MAC"`
	TapMac              string    `gorm:"column:tap_mac;type:char(64);default:null" json:"TAP_MAC"`
	AnalyzerIP          string    `gorm:"column:analyzer_ip;type:char(64);not null" json:"ANALYZER_IP"`
	CurAnalyzerIP       string    `gorm:"column:cur_analyzer_ip;type:char(64);not null" json:"CUR_ANALYZER_IP"`
	ControllerIP        string    `gorm:"column:controller_ip;type:char(64);not null" json:"CONTROLLER_IP"`
	CurControllerIP     string    `gorm:"column:cur_controller_ip;type:char(64);not null" json:"CUR_CONTROLLER_IP"`
	LaunchServer        string    `gorm:"column:launch_server;type:char(64);not null" json:"LAUNCH_SERVER"`
	LaunchServerID      int       `gorm:"column:launch_server_id;type:int;default:null" json:"LAUNCH_SERVER_ID"`
	AZ                  string    `gorm:"column:az;type:char(64);default:''" json:"AZ"`
	Region              string    `gorm:"column:region;type:char(64);default:''" json:"REGION"`
	Revision            string    `gorm:"column:revision;type:varchar(256);default:null" json:"REVISION"`
	SyncedControllerAt  time.Time `gorm:"column:synced_controller_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"SYNCED_CONTROLLER_AT"`
	SyncedAnalyzerAt    time.Time `gorm:"column:synced_analyzer_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"SYNCED_ANALYZER_AT"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	BootTime            int       `gorm:"column:boot_time;type:int;default:0" json:"BOOT_TIME"`
	Exceptions          int64     `gorm:"column:exceptions;type:bigint unsigned;default:0" json:"EXCEPTIONS"`
	VTapLcuuid          string    `gorm:"column:vtap_lcuuid;type:char(64);default:null" json:"VTAP_LCUUID"`
	VtapGroupLcuuid     string    `gorm:"column:vtap_group_lcuuid;type:char(64);default:null" json:"VTAP_GROUP_LCUUID"`
	CPUNum              int       `gorm:"column:cpu_num;type:int;default:0" json:"CPU_NUM"` // logical number of cpu
	MemorySize          int64     `gorm:"column:memory_size;type:bigint;default:0" json:"MEMORY_SIZE"`
	Arch                string    `gorm:"column:arch;type:varchar(256);default:null" json:"ARCH"`
	Os                  string    `gorm:"column:os;type:varchar(256);default:null" json:"OS"`
	KernelVersion       string    `gorm:"column:kernel_version;type:varchar(256);default:null" json:"KERNEL_VERSION"`
	ProcessName         string    `gorm:"column:process_name;type:varchar(256);default:null" json:"PROCESS_NAME"`
	CurrentK8sImage     string    `gorm:"column:current_k8s_image;type:varchar(512);default:null" json:"CURRENT_K8S_IMAGE"`
	LicenseType         int       `gorm:"column:license_type;type:int;default:null" json:"LICENSE_TYPE"`           // 1: A类 2: B类 3: C类
	LicenseFunctions    string    `gorm:"column:license_functions;type:char(64)" json:"LICENSE_FUNCTIONS"`         // separated by ,; 1: 流量分发 2: 网络监控 3: 应用监控
	EnableFeatures      string    `gorm:"column:enable_features;type:char(64)" json:"ENABLE_FEATURES"`             // separated by ,
	DisableFeatures     string    `gorm:"column:disable_features;type:char(64)" json:"DISABLE_FEATURES"`           // separated by ,
	FollowGroupFeatures string    `gorm:"column:follow_group_features;type:char(64)" json:"FOLLOW_GROUP_FEATURES"` // separated by ,
	TapMode             int       `gorm:"column:tap_mode;type:int;default:null" json:"TAP_MODE"`
	ExpectedRevision    string    `gorm:"column:expected_revision;type:text;default null" json:"EXPECTED_REVISION"`
	UpgradePackage      string    `gorm:"column:upgrade_package;type:text;default null" json:"UPGRADE_PACKAGE"`
	TeamID              int       `gorm:"column:team_id;type:int;default:0" json:"TEAM_ID"`
	Lcuuid              string    `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
}

func (VTap) TableName() string {
	return "vtap"
}

func (v VTap) GetID() int {
	return v.ID
}

func (v VTap) GetLcuuid() string {
	return v.Lcuuid
}

type VTapGroup struct {
	ID               int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name             string    `gorm:"column:name;type:varchar(64);not null" json:"NAME"`
	CreatedAt        time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
	Lcuuid           string    `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
	ShortUUID        string    `gorm:"column:short_uuid;type:char(32);default:null" json:"SHORT_UUID"`
	TeamID           int       `gorm:"column:team_id;type:int;default:0" json:"TEAM_ID"`
	UserID           int       `gorm:"column:user_id;type:int;default:null" json:"USER_ID"`
	LicenseFunctions string    `gorm:"column:license_functions;type:char(64)" json:"LICENSE_FUNCTIONS"` // separated by ,
}

func (VTapGroup) TableName() string {
	return "vtap_group"
}

type LicenseFuncLog struct {
	ID                  int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	TeamID              int       `gorm:"column:team_id;type:int;default:1" json:"TEAM_ID"`
	AgentID             int       `gorm:"column:agent_id;type:int" json:"AGENT_ID"`
	AgentName           string    `gorm:"column:agent_name;type:varchar(256)" json:"AGENT_NAME"`
	UserID              int       `gorm:"column:user_id;type:int" json:"USER_ID"`
	LicenseFunction     int       `gorm:"column:license_function;type:int;" json:"ENABLED_FEATURE"`
	Enabled             int       `gorm:"column:enabled;type:int" json:"ENABLED"`
	AgentGroupName      string    `gorm:"column:agent_group_name;type:varchar(64)" json:"AGENT_GROUP_NAME"`
	AgentGroupOperation int       `gorm:"column:agent_group_operation;type:int" json:"AGENT_GROUP_OPERATION"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
}

func (LicenseFuncLog) TableName() string {
	return "license_func_log"
}

type DataSource struct {
	ID                        int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	DisplayName               string    `gorm:"column:display_name;type:char(64);default:''" json:"DISPLAY_NAME"`
	DataTableCollection       string    `gorm:"column:data_table_collection;type:char(64);default:''" json:"DATA_TABLE_COLLECTION"`
	State                     int       `gorm:"column:state;type:int;default:1" json:"STATE"`
	BaseDataSourceID          int       `gorm:"column:base_data_source_id;type:int" json:"BASE_DATA_SOURCE_ID"`
	IntervalTime              int       `gorm:"column:interval_time;type:int" json:"INTERVAL"`
	RetentionTime             int       `gorm:"column:retention_time;type:int" json:"RETENTION_TIME"` // unit: hour
	QueryTime                 int       `gorm:"column:query_time;type:int" json:"QUERY_TIME"`         // unit: minute
	SummableMetricsOperator   string    `gorm:"column:summable_metrics_operator;type:char(64)" json:"SUMMABLE_METRICS_OPERATOR"`
	UnSummableMetricsOperator string    `gorm:"column:unsummable_metrics_operator;type:char(64)" json:"UNSUMMABLE_METRICS_OPERATOR"`
	UpdatedAt                 time.Time `gorm:"column:updated_at" json:"UPDATED_AT"`
	Lcuuid                    string    `gorm:"column:lcuuid;type:char(64)" json:"LCUUID"`
}

// SysConfiguration [...]
type SysConfiguration struct {
	ID        int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	ParamName string `gorm:"column:param_name;type:char(64);not null" json:"PARAM_NAME"`
	Value     string `gorm:"column:value;type:varchar(256);default:null" json:"VALUE"`
	Comments  string `gorm:"column:comments;type:text;default:null" json:"COMMENTS"`
	Lcuuid    string `gorm:"column:lcuuid;type:char(64);default:null" json:"LCUUID"`
}

type KubernetesCluster struct {
	ID          int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	ClusterID   string    `gorm:"column:cluster_id;type:varchar(256);" json:"CLUSTER_ID"`
	Value       string    `gorm:"column:value;type:varchar(256);" json:"VALUE"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	SyncedAt    time.Time `gorm:"column:synced_at;type:datetime" json:"SYNCED_AT"`
	UpdatedTime time.Time `gorm:"column:updated_time;type:datetime;default:null" json:"UPDATED_TIME"`
}

type ACL struct {
	ID           int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	BusinessID   int       `gorm:"column:business_id;type:int;not null" json:"BUSINESS_ID"`
	Name         string    `gorm:"column:name;type:char(64);default:null" json:"NAME"`
	Type         int       `gorm:"column:type;type:int;default:2" json:"TYPE"`                     // 1-epc; 2-custom
	TapType      int       `gorm:"column:tap_type;type:int;default:3" json:"TAP_TYPE"`             // 1-WAN; 3-LAN
	State        int       `gorm:"column:state;type:int;default:null;default:0" json:"STATE"`      // 0-disable; 1-enable
	Applications string    `gorm:"column:applications;type:char(64);not null" json:"APPLICATIONS"` // separated by , (1-performance analysis; 2-backpacking; 6-npb)
	EpcID        int       `gorm:"column:epc_id;type:int;default:null" json:"EPC_ID"`
	SrcGroupIDs  string    `gorm:"column:src_group_ids;type:text;default:null" json:"SRC_GROUP_IDS"` // separated by ,
	DstGroupIDs  string    `gorm:"column:dst_group_ids;type:text;default:null" json:"DST_GROUP_IDS"` // separated by ,
	Protocol     *int      `gorm:"column:protocol;type:int;default:null" json:"PROTOCOL"`
	SrcPorts     string    `gorm:"column:src_ports;type:text;default:null" json:"SRC_PORTS"` // separated by ,
	DstPorts     string    `gorm:"column:dst_ports;type:text;default:null" json:"DST_PORTS"` // separated by ,
	Vlan         int       `gorm:"column:vlan;type:int;default:null" json:"VLAN"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
	Lcuuid       string    `gorm:"column:lcuuid;type:char(64);default:null" json:"LCUUID"`
}

func (ACL) TableName() string {
	return "acl"
}

type GroupACL struct {
	ID      int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	GroupID int    `gorm:"column:group_id;type:int;not null" json:"GROUP_ID"`
	ACLID   int    `gorm:"column:acl_id;type:int;not null" json:"ACL_ID"`
	Lcuuid  string `gorm:"column:lcuuid;type:char(64);default:null" json:"LCUUID"`
}

func (GroupACL) TableName() string {
	return "group_acl"
}

type PolicyACLGroup struct {
	ID     int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	ACLIDs string `gorm:"column:acl_ids;type:text;not null" json:"ACL_IDS"` // separated by ,
	COUNT  int    `gorm:"column:count;type:int;not null" json:"COUNT"`
}

func (PolicyACLGroup) TableName() string {
	return "policy_acl_group"
}

// ResourceGroupExtraInfo [...]
type ResourceGroupExtraInfo struct {
	ID           int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	ResourceType int    `gorm:"column:resource_type;type:int;not null" json:"RESOURCE_TYPE"` // 1: epc, 2: vm, 3: pod_group, 4: pod_service
	ResourceID   int    `gorm:"column:resource_id;type:int;not null" json:"RESOURCE_ID"`
	ResourceName string `gorm:"column:resource_name;type:char(64);not null" json:"RESOURCE_NAME"`
}

func (ResourceGroupExtraInfo) TableName() string {
	return "resource_group_extra_info"
}

// NpbPolicy [...]
type NpbPolicy struct {
	ID               int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name             string    `gorm:"column:name;type:char(64);default:null" json:"NAME"`
	State            int       `gorm:"column:state;type:int;default:0" json:"STATE"` // 0-disable; 1-enable
	BusinessID       int       `gorm:"column:business_id;type:int;not null" json:"BUSINESS_ID"`
	Direction        int       `gorm:"column:direction;type:int;default:1" json:"DIRECTION"` // 1-two way; 2-server to client
	Vni              *int      `gorm:"column:vni;type:int;default:null" json:"VNI"`
	NpbTunnelID      int       `gorm:"column:npb_tunnel_id;type:int;default:null" json:"NPB_TUNNEL_ID"`
	Distribute       int       `gorm:"column:distribute;type:int;default:0" json:"distribute"` // 0-drop, 1-distribute
	PayloadSlice     *int      `gorm:"column:payload_slice;type:int;default:null" json:"PAYLOAD_SLICE"`
	ACLID            int       `gorm:"column:acl_id;type:int;default:null" json:"ACL_ID"`
	PolicyACLGroupID int       `gorm:"column:policy_acl_group_id;type:int;default:null" json:"POLICY_ACL_GROUP_ID"`
	VtapType         int       `gorm:"column:vtap_type;type:type:tinyint(1);default:null" json:"VTAP_TYPE"` // 1-vtap; 2-vtap_group
	VtapIDs          string    `gorm:"column:vtap_ids;type:text" json:"VTAP_IDS"`                           // separated by ,
	VtapGroupIDs     string    `gorm:"column:vtap_group_ids;type:text" json:"VTAP_GROUP_IDS"`               // separated by ,
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
	Lcuuid           string    `gorm:"column:lcuuid;type:char(64);default:null" json:"LCUUID"`
	TeamID           int       `gorm:"column:team_id;type:int;default:1" json:"TEAM_ID"`
}

func (NpbPolicy) TableName() string {
	return "npb_policy"
}

// NpbTunnel [...]
type NpbTunnel struct {
	ID           int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name         string    `gorm:"column:name;type:char(64);not null" json:"NAME"`
	IP           string    `gorm:"column:ip;type:char(64);default:null" json:"IP"`
	Type         int       `gorm:"column:type;type:int;default:0" json:"TYPE"`                            // (0-VXLAN；1-ERSPAN)
	VNIInputType int       `gorm:"column:vni_input_type;type:tinyint(1);default:1" json:"VNI_INPUT_TYPE"` // 1: entire one, 2: two parts
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
	Lcuuid       string    `gorm:"column:lcuuid;type:char(64);default:null" json:"LCUUID"`
	TeamID       int       `gorm:"column:team_id;type:int;default:1" json:"TEAM_ID"`
}

func (NpbTunnel) TableName() string {
	return "npb_tunnel"
}

// PcapPolicy [...]
type PcapPolicy struct {
	ID               int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name             string    `gorm:"column:name;type:char(64);default:null" json:"NAME"`
	State            int       `gorm:"column:state;type:int;default:0" json:"STATE"` // 0-disable; 1-enable
	BusinessID       int       `gorm:"column:business_id;type:int;not null" json:"BUSINESS_ID"`
	ACLID            int       `gorm:"column:acl_id;type:int;default:null" json:"ACL_ID"`
	VtapType         int       `gorm:"column:vtap_type;type:type:tinyint(1);default:null" json:"VTAP_TYPE"` // 1-vtap; 2-vtap_group
	VtapIDs          string    `gorm:"column:vtap_ids;type:text;default:null" json:"VTAP_IDS"`              // separated by ,
	VtapGroupIDs     string    `gorm:"column:vtap_group_ids;type:text;default:null" json:"VTAP_GROUP_IDS"`  // separated by ,
	PayloadSlice     *int      `gorm:"column:payload_slice;type:int;default:null" json:"PAYLOAD_SLICE"`
	PolicyACLGroupID int       `gorm:"column:policy_acl_group_id;type:int;default:null" json:"POLICY_ACL_GROUP_ID"`
	UserID           int       `gorm:"column:user_id;type:int;default:null" json:"USER_ID"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
	Lcuuid           string    `gorm:"column:lcuuid;type:char(64);default:null" json:"LCUUID"`
	TeamID           int       `gorm:"column:team_id;type:int;default:1" json:"TEAM_ID"`
}

func (PcapPolicy) TableName() string {
	return "pcap_policy"
}

type DialTestTask struct {
	ID            int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name          string    `gorm:"column:name;type:varchar(256);not null" json:"NAME"`
	Protocol      int       `gorm:"column:protocol;type:int;not null" json:"PROTOCOL"` // 1.ICMP
	Host          string    `gorm:"column:host;type:varchar(256);not null" json:"HOST"`
	OvertimeTime  int       `gorm:"column:overtime_time;type:int;default:2000" json:"OVERTIME_TIME"`
	Payload       int       `gorm:"column:payload;type:int;default:64" json:"PAYLOAD"`
	TTL           int       `gorm:"column:ttl;type:smallint;default:64" json:"TTL"`
	DialLocation  string    `gorm:"column:dial_location;type:varchar(256);not null" json:"DIAL_LOCATION"`
	DialFrequency int       `gorm:"column:dial_frequency;type:int;default:1000" json:"DIAL_FREQUENCY"`
	PCAP          []byte    `gorm:"column:pcap;type:mediumblob" json:"PCAP"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
}

func (DialTestTask) TableName() string {
	return "dial_test_task"
}

type VTapRepo struct {
	ID        int             `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name      string          `gorm:"column:name;type:char(64);not null" json:"NAME"`
	Arch      string          `gorm:"column:arch;type:varchar(256);default:''" json:"ARCH"`
	OS        string          `gorm:"column:os;type:varchar(256);default:''" json:"OS"`
	Branch    string          `gorm:"column:branch;type:varchar(256);default:''" json:"BRANCH"`
	RevCount  string          `gorm:"column:rev_count;type:varchar(256);default:''" json:"REV_COUNT"`
	CommitID  string          `gorm:"column:commit_id;type:varchar(256);default:''" json:"COMMIT_ID"`
	Image     compressedBytes `gorm:"column:image;type:logblob" json:"IMAGE"`
	K8sImage  string          `gorm:"column:k8s_image;type:varchar(512);default:''" json:"K8S_IMAGE"`
	CreatedAt time.Time       `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt time.Time       `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
}

type compressedBytes []byte

// Scan scan decompress value into compressedBytes, implements sql.Scanner interface
func (c *compressedBytes) Scan(value interface{}) error {
	// decompress
	compressedData, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("failed to decompress compressedImage value:", value))
	}

	var b bytes.Buffer
	b.Write(compressedData)
	r, err := zlib.NewReader(&b)
	if err != nil {
		return err
	}
	defer r.Close()

	originData, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	*c = originData
	return nil
}

// Value return compress value, implement driver.Valuer interface
func (c compressedBytes) Value() (driver.Value, error) {
	// compress
	t1 := time.Now()
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(c)
	if err != nil {
		return nil, fmt.Errorf("failed to write compressed data: %v", err)
	}
	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zlib writer: %v", err)
	}
	log.Infof("compress time comsumed: %v", time.Since(t1))
	return b.String(), nil
}

func (VTapRepo) TableName() string {
	return "vtap_repo"
}

type Plugin struct {
	ID        int             `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name      string          `gorm:"column:name;type:varchar(256);not null" json:"NAME"`
	Type      int             `gorm:"column:type;type:int" json:"TYPE"`                // 1: wasm 2: so 3: lua
	UserName  int             `gorm:"column:user_name;type:int;default:1" json:"USER"` // 1: agent 2: server
	Image     compressedBytes `gorm:"column:image;type:logblob;not null" json:"IMAGE"`
	CreatedAt time.Time       `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"CREATED_AT"`
	UpdatedAt time.Time       `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"UPDATED_AT"`
}

func (Plugin) TableName() string {
	return "plugin"
}

type MailServer struct {
	ID           int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Status       int    `gorm:"column:status;type:int;not null" json:"STATUS"`
	Host         string `gorm:"column:host;type:text;not null" json:"HOST"`
	Port         int    `gorm:"column:port;type:int;not null" json:"PORT"`
	UserName     string `gorm:"column:user_name;type:text;not null" json:"USER"`
	Password     string `gorm:"column:password;type:text;not null" json:"PASSWORD"`
	Security     string `gorm:"column:security;type:text;not null" json:"SECURITY"`
	NtlmEnabled  int    `gorm:"column:ntlm_enabled;type:int" json:"NTLM_ENABLED"`
	NtlmName     string `gorm:"column:ntlm_name;type:text" json:"NTLM_NAME"`
	NtlmPassword string `gorm:"column:ntlm_password;type:text" json:"NTLM_PASSWORD"`
	Lcuuid       string `gorm:"unique;column:lcuuid;type:char(64)" json:"LCUUID"`
}

func (MailServer) TableName() string {
	return "mail_server"
}

type AlarmPolicy struct {
	ID     int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name   string `gorm:"column:name;type:char(128)" json:"NAME"`
	UserID int    `gorm:"column:user_id;type:int" json:"USER_ID"`
	TeamID int    `gorm:"column:team_id;type:int;default:1" json:"TEAM_ID"`
}

func (AlarmPolicy) TableName() string {
	return "alarm_policy"
}

type ORG struct {
	ID          int            `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name        string         `gorm:"column:name;type:char(128);default:''" json:"NAME"`
	ORGID       int            `gorm:"column:org_id;type:int;default:0" json:"ORG_ID"`
	Lcuuid      string         `gorm:"column:lcuuid;type:char(64);not null" json:"LCUUID"`
	OwnerUserID int            `gorm:"column:owner_user_id;type:int;default:0" json:"OWNER_USER_ID"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;default:null" json:"DELETED_AT" mapstructure:"DELETED_AT"`
}

func (o ORG) GetID() int {
	return o.ORGID
}

func (o ORG) GetLcuuid() string {
	return o.Lcuuid
}

type Team struct {
	ID          int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name        string `gorm:"column:name;type:char(128);default:''" json:"NAME"`
	TeamID      int    `gorm:"column:team_id;type:int;default:0" json:"TEAM_ID"`
	ShortLcuuid string `gorm:"column:short_lcuuid;type:char(64);default:''" json:"SHORT_LCUUID"`
	ORGID       int    `gorm:"column:org_id;type:int;default:0" json:"ORG_ID"`
}

type User struct {
	ID       int    `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	UserName string `gorm:"column:username;type:char(128)" json:"USERNAME"`
}

type ResourceVersion struct {
	ID        int       `gorm:"primaryKey;column:id;type:int;not null" json:"ID"`
	Name      string    `gorm:"column:name;type:varchar(255)" json:"RESOURCE"`
	Version   uint32    `gorm:"column:version;type:int unsigned" json:"VERSION"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;type:datetime" json:"CREATED_AT"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;type:datetime" json:"UPDATED_AT"`
}
