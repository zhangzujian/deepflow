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

package metadata

import (
	"fmt"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/golang/protobuf/proto"

	"github.com/deepflowio/deepflow/message/trident"
	. "github.com/deepflowio/deepflow/server/controller/common"
	models "github.com/deepflowio/deepflow/server/controller/db/metadb/model"
	. "github.com/deepflowio/deepflow/server/controller/trisolaris/utils"
)

type ServiceRawData struct {
	podServiceIDToPodServicePorts map[int][]*models.PodServicePort
	podGroupIDToPodGroupPorts     map[int][]*models.PodGroupPort
	podGroupIDToPodServiceID      map[int]int
	lbIDToVPCID                   map[int]int
	customServiceIDToIPs          map[int]mapset.Set[string]
	customServiceIPKeyToPorts     map[customServiceIPKey]mapset.Set[uint32]
}

type customServiceIPKey struct {
	vpcID int
	ip    string
}

func newCustomServiceIPKey(vpcID int, ip string) customServiceIPKey {
	return customServiceIPKey{
		vpcID: vpcID,
		ip:    ip,
	}
}

func newServiceRawData() *ServiceRawData {
	return &ServiceRawData{
		podServiceIDToPodServicePorts: make(map[int][]*models.PodServicePort),
		podGroupIDToPodGroupPorts:     make(map[int][]*models.PodGroupPort),
		podGroupIDToPodServiceID:      make(map[int]int),
		lbIDToVPCID:                   make(map[int]int),
	}
}

type ServiceDataOP struct {
	serviceRawData *ServiceRawData
	metaData       *MetaData
	services       []*trident.ServiceInfo
	ORGID
}

func newServiceDataOP(metaData *MetaData) *ServiceDataOP {
	return &ServiceDataOP{
		serviceRawData: newServiceRawData(),
		metaData:       metaData,
		services:       []*trident.ServiceInfo{},
		ORGID:          metaData.ORGID,
	}
}

func (r *ServiceRawData) ConvertDBData(md *MetaData) {
	dbDataCache := md.GetDBDataCache()
	for _, psp := range dbDataCache.GetPodServicePorts() {
		if _, ok := r.podServiceIDToPodServicePorts[psp.PodServiceID]; ok {
			r.podServiceIDToPodServicePorts[psp.PodServiceID] = append(
				r.podServiceIDToPodServicePorts[psp.PodServiceID], psp)
		} else {
			r.podServiceIDToPodServicePorts[psp.PodServiceID] = []*models.PodServicePort{psp}
		}
	}
	podGroupIDToPodServiceIDs := make(map[int][]int, len(dbDataCache.GetPodGroupPorts()))
	for _, pgp := range dbDataCache.GetPodGroupPorts() {
		if _, ok := r.podGroupIDToPodGroupPorts[pgp.PodGroupID]; ok {
			r.podGroupIDToPodGroupPorts[pgp.PodGroupID] = append(
				r.podGroupIDToPodGroupPorts[pgp.PodGroupID], pgp)
		} else {
			r.podGroupIDToPodGroupPorts[pgp.PodGroupID] = []*models.PodGroupPort{pgp}
		}
		if ids, ok := podGroupIDToPodServiceIDs[pgp.PodGroupID]; ok {
			podGroupIDToPodServiceIDs[pgp.PodGroupID] = append(
				podGroupIDToPodServiceIDs[pgp.PodGroupID], pgp.PodServiceID)
		} else {
			if Find[int](ids, pgp.PodServiceID) != true {
				podGroupIDToPodServiceIDs[pgp.PodGroupID] = []int{pgp.PodServiceID}
			}
		}
	}

	for _, lb := range dbDataCache.GetLBs() {
		r.lbIDToVPCID[lb.ID] = lb.VPCID
	}

	podServiceIDToName := make(map[int]string, len(dbDataCache.GetPodServices()))
	for _, podService := range dbDataCache.GetPodServices() {
		podServiceIDToName[podService.ID] = podService.Name
	}
	podGroupIDToPodServiceID := make(map[int]int, len(podGroupIDToPodServiceIDs))
	for groupId, podServiceIDs := range podGroupIDToPodServiceIDs {
		var (
			minName string
			minID   int
		)
		for _, serviceID := range podServiceIDs {
			name := podServiceIDToName[serviceID]
			if len(name) == 0 {
				continue
			}
			if minName == "" || minName > name {
				minName = name
				minID = serviceID
			}
		}
		podGroupIDToPodServiceID[groupId] = minID
	}
	r.podGroupIDToPodServiceID = podGroupIDToPodServiceID

	// VPC 范围内合并自定义服务：
	// 1. 去重 ip，ip 绑定第一个遍历到的服务
	// 2. 合并 ip 的端口，如果存在相同 ip 中的某个没有端口，那么端口全部丢弃
	customServiceIDToIPs := make(map[int]mapset.Set[string])
	customServiceIPKeyToPorts := make(map[customServiceIPKey]mapset.Set[uint32])
	for _, cs := range dbDataCache.GetCustomServices() {
		var ips []string
		var ports []uint32
		if cs.Type == CUSTOM_SERVICE_TYPE_IP {
			ips = strings.Split(cs.Resource, ",")
		} else {
			ipPorts := strings.Split(cs.Resource, ",")
			for _, ipPort := range ipPorts {
				ipPort = strings.TrimSpace(ipPort)
				separatorIndex := strings.LastIndex(ipPort, ":")
				if separatorIndex == -1 {
					log.Warningf("[ORG-%s] invalid ip port format: %s", md.ORGID, ipPort)
					continue
				}
				port, err := strconv.Atoi(ipPort[separatorIndex+1:])
				if err != nil {
					log.Warningf("[ORG-%s] invalid port format: %s", md.ORGID, ipPort[separatorIndex+1:])
					continue
				}
				ips = append(ips, ipPort[:separatorIndex])
				ports = append(ports, uint32(port))
			}
		}

		for i, ip := range ips {
			ipKey := newCustomServiceIPKey(cs.VPCID, ip)
			if existedPorts, ok := customServiceIPKeyToPorts[ipKey]; ok {
				// 该 ip 已经存在，那么 ip 不会跟本服务绑定
				// 已存在的 ip 没有端口，那么直接跳过
				if existedPorts.Cardinality() == 0 {
					continue
				}
				// 已存在的 ip 有端口，该 ip 没有端口，清空端口
				if len(ports) == 0 {
					customServiceIPKeyToPorts[ipKey] = mapset.NewSet[uint32]()
					continue
				}
				// 已存在的 ip 有端口，该 ip 有端口，合并端口
				customServiceIPKeyToPorts[ipKey].Add(ports[i])
			} else {
				// 该 ip 不存在，则 ip 会跟本服务绑定
				if _, ok := customServiceIDToIPs[cs.ID]; ok {
					customServiceIDToIPs[cs.ID].Add(ip)
				} else {
					customServiceIDToIPs[cs.ID] = mapset.NewSet[string](ip)
				}
				customServiceIPKeyToPorts[ipKey] = mapset.NewSet[uint32]()
				// 该 ip 没有端口，那么跳过
				if len(ports) == 0 {
					continue
				}
				// 该 ip 有端口
				customServiceIPKeyToPorts[ipKey].Add(ports[i])
			}
		}
	}
	r.customServiceIDToIPs = customServiceIDToIPs
	r.customServiceIPKeyToPorts = customServiceIPKeyToPorts
}

func (s *ServiceDataOP) updateServiceRawData(serviceRawData *ServiceRawData) {
	s.serviceRawData = serviceRawData
}

func ipToCidr(ip string) string {
	if strings.Contains(ip, "/") {
		return ip
	} else if strings.Contains(ip, ":") {
		return fmt.Sprintf("%s/128", ip)
	} else {
		return fmt.Sprintf("%s/32", ip)
	}

}

var protocol_string_to_int = map[string]trident.ServiceProtocol{
	"ALL":   trident.ServiceProtocol_ANY,
	"HTTP":  trident.ServiceProtocol_TCP_SERVICE,
	"HTTPS": trident.ServiceProtocol_TCP_SERVICE,
	"TCP":   trident.ServiceProtocol_TCP_SERVICE,
	"UDP":   trident.ServiceProtocol_UDP_SERVICE,
}

func getProtocol(protocol string) trident.ServiceProtocol {
	if protocol, ok := protocol_string_to_int[protocol]; ok {
		return protocol
	}
	return trident.ServiceProtocol_ANY
}

func serviceToProto(
	vpcID int, ips []string, protocol trident.ServiceProtocol, serverPorts []uint32,
	serviceType trident.ServiceType, serviceID int) *trident.ServiceInfo {

	return &trident.ServiceInfo{
		EpcId:       proto.Uint32(uint32(vpcID)),
		Ips:         ips,
		Protocol:    &protocol,
		ServerPorts: serverPorts,
		Type:        &serviceType,
		Id:          proto.Uint32(uint32(serviceID)),
	}
}

func podGroupToProto(
	podGroupId int, protocol trident.ServiceProtocol, serverPorts []uint32,
	serviceType trident.ServiceType, serviceID int) *trident.ServiceInfo {

	return &trident.ServiceInfo{
		PodGroupId:  proto.Uint32(uint32(podGroupId)),
		Protocol:    &protocol,
		ServerPorts: serverPorts,
		Type:        &serviceType,
		Id:          proto.Uint32(uint32(serviceID)),
	}
}

func NodeToProto(
	podClusterId int, protocol trident.ServiceProtocol, serverPorts []uint32,
	serviceType trident.ServiceType, serviceID int) *trident.ServiceInfo {

	return &trident.ServiceInfo{
		PodClusterId: proto.Uint32(uint32(podClusterId)),
		Protocol:     &protocol,
		ServerPorts:  serverPorts,
		Type:         &serviceType,
		Id:           proto.Uint32(uint32(serviceID)),
	}
}

type GroupKey struct {
	protocol     trident.ServiceProtocol
	podServiceID int
}

type NodeKey struct {
	podClusterID int
	protocol     trident.ServiceProtocol
	podServiceID int
}

// All traversals are guaranteed to be in order
func (s *ServiceDataOP) generateService() {
	dbDataCache := s.metaData.GetDBDataCache()
	services := []*trident.ServiceInfo{}
	rData := s.serviceRawData
	groupKeys := []GroupKey{}
	for _, podGroup := range dbDataCache.GetPodGroups() {
		ports, ok := rData.podGroupIDToPodGroupPorts[podGroup.ID]
		if ok == false {
			continue
		}
		groupKeys = groupKeys[:0]
		podServiceID := rData.podGroupIDToPodServiceID[podGroup.ID]
		keyToPorts := make(map[GroupKey][]uint32)
		for _, port := range ports {
			key := GroupKey{
				protocol:     getProtocol(port.Protocol),
				podServiceID: podServiceID,
			}
			if ports, ok := keyToPorts[key]; ok {
				if Find[uint32](ports, uint32(port.Port)) != true {
					keyToPorts[key] = append(keyToPorts[key], uint32(port.Port))
				}
			} else {
				groupKeys = append(groupKeys, key)
				keyToPorts[key] = []uint32{uint32(port.Port)}
			}
		}
		for index := range groupKeys {
			if valuse, ok := keyToPorts[groupKeys[index]]; ok {
				service := podGroupToProto(
					podGroup.ID,
					groupKeys[index].protocol,
					valuse,
					trident.ServiceType_POD_SERVICE_POD_GROUP,
					groupKeys[index].podServiceID,
				)
				services = append(services, service)
			}
		}
	}

	nodeKeys := []NodeKey{}
	nodeServiceInfo := make(map[NodeKey][]uint32)
	for _, podService := range dbDataCache.GetPodServices() {
		podServiceports, ok := rData.podServiceIDToPodServicePorts[podService.ID]
		if ok == false {
			continue
		}
		if podService.ServiceClusterIP == "" {
			log.Debugf(s.Logf("pod service(id=%d) has no service_cluster_ip", podService.ID))
			continue
		}
		protocols := []trident.ServiceProtocol{}
		protocolToPorts := make(map[trident.ServiceProtocol][]uint32)
		for _, podServiceport := range podServiceports {
			protocol := getProtocol(podServiceport.Protocol)
			if _, ok := protocolToPorts[protocol]; ok {
				protocolToPorts[protocol] = append(protocolToPorts[protocol], uint32(podServiceport.Port))
			} else {
				protocols = append(protocols, protocol)
				protocolToPorts[protocol] = []uint32{uint32(podServiceport.Port)}
			}
			if podService.Type == POD_SERVICE_TYPE_NODEPORT || podService.Type == POD_SERVICE_TYPE_LOADBALANCER {
				key := NodeKey{
					podClusterID: podService.PodClusterID,
					protocol:     protocol,
					podServiceID: podServiceport.PodServiceID,
				}
				if _, ok := nodeServiceInfo[key]; ok {
					nodeServiceInfo[key] = append(nodeServiceInfo[key], uint32(podServiceport.NodePort))
				} else {
					nodeKeys = append(nodeKeys, key)
					nodeServiceInfo[key] = []uint32{uint32(podServiceport.NodePort)}
				}
			}
		}

		ips := []string{podService.ServiceClusterIP}
		for index := range protocols {
			if ports, ok := protocolToPorts[protocols[index]]; ok {
				service := serviceToProto(
					podService.VPCID,
					ips,
					protocols[index],
					ports,
					trident.ServiceType_POD_SERVICE_IP,
					podService.ID,
				)
				services = append(services, service)
			}
		}
	}
	for index := range nodeKeys {
		if values, ok := nodeServiceInfo[nodeKeys[index]]; ok {
			service := NodeToProto(
				nodeKeys[index].podClusterID,
				nodeKeys[index].protocol,
				values,
				trident.ServiceType_POD_SERVICE_NODE,
				nodeKeys[index].podServiceID,
			)
			services = append(services, service)
		}
	}

	for _, lbListener := range dbDataCache.GetLBListeners() {
		vpcID := rData.lbIDToVPCID[lbListener.LBID]
		var ips []string
		if lbListener.IPs != "" {
			ips = strings.Split(lbListener.IPs, ",")
		}
		service := serviceToProto(
			vpcID,
			ips,
			getProtocol(lbListener.Protocol),
			[]uint32{uint32(lbListener.Port)},
			trident.ServiceType_LB_SERVICE,
			lbListener.ID,
		)
		services = append(services, service)
	}

	for _, lbts := range dbDataCache.GetLBTargetServers() {
		vpcID := rData.lbIDToVPCID[lbts.LBID]
		var ips []string
		if lbts.IP == "" {
			log.Debugf(s.Logf("lb_target_server(id=%d) has no ips", lbts.ID))
			continue
		} else {
			ips = []string{lbts.IP}
		}

		service := serviceToProto(
			vpcID,
			ips,
			getProtocol(lbts.Protocol),
			[]uint32{uint32(lbts.Port)},
			trident.ServiceType_LB_SERVICE,
			lbts.LBListenerID,
		)
		services = append(services, service)
	}

	serviceTypeCustomService := trident.ServiceType_CUSTOM_SERVICE
	for _, customService := range dbDataCache.GetCustomServices() {
		ips, ok := s.serviceRawData.customServiceIDToIPs[customService.ID]
		if !ok {
			continue
		}
		for ip := range ips.Iter() {
			service := &trident.ServiceInfo{
				Type:        &serviceTypeCustomService,
				Id:          proto.Uint32(uint32(customService.ID)),
				EpcId:       proto.Uint32(uint32(customService.VPCID)),
				Ips:         []string{ip},
				ServerPorts: s.serviceRawData.customServiceIPKeyToPorts[newCustomServiceIPKey(customService.VPCID, ip)].ToSlice(),
			}
			services = append(services, service)
		}
	}

	s.services = services
	log.Debugf(s.Logf("service have %d", len(s.services)))
}

func (s *ServiceDataOP) GetServiceData() []*trident.ServiceInfo {
	return s.services
}

func (s *ServiceDataOP) generateRawData() {
	serviceRawData := newServiceRawData()
	serviceRawData.ConvertDBData(s.metaData)
	s.updateServiceRawData(serviceRawData)
}

func (s *ServiceDataOP) GenerateServiceData() {
	s.generateRawData()
	s.generateService()
}
