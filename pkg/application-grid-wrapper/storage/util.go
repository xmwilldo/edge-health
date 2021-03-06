/*
Copyright 2020 The SuperEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"encoding/json"
	"net"
	"strconv"

	"github.com/hashicorp/serf/serf"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

const (
	TopologyAnnotationsKey = "topologyKeys"

	EdgeLocalEndpoint  = "superedge.io/local-endpoint"
	EdgeLocalPort      = "superedge.io/local-port"
	MasterEndpointName = "kubernetes"
)

func getTopologyKeys(objectMeta *metav1.ObjectMeta) []string {
	if !hasTopologyKey(objectMeta) {
		return nil
	}

	var keys []string
	keyData := objectMeta.Annotations[TopologyAnnotationsKey]
	if err := json.Unmarshal([]byte(keyData), &keys); err != nil {
		klog.Errorf("can't parse topology keys %s, %v", keyData, err)
		return nil
	}

	return keys
}

func hasTopologyKey(objectMeta *metav1.ObjectMeta) bool {
	if objectMeta.Annotations == nil {
		return false
	}

	_, ok := objectMeta.Annotations[TopologyAnnotationsKey]
	return ok
}

func genLocalEndpoints(eps *v1.Endpoints) *v1.Endpoints {
	if eps.Namespace != metav1.NamespaceDefault || eps.Name != MasterEndpointName {
		return eps
	}

	klog.V(4).Infof("begin to gen local ep %v", eps)
	ipAddress, e := eps.Annotations[EdgeLocalEndpoint]
	if !e {
		return eps
	}

	portStr, e := eps.Annotations[EdgeLocalPort]
	if !e {
		return eps
	}

	klog.V(4).Infof("get local endpoint %s:%s", ipAddress, portStr)
	port, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		klog.Errorf("parse int %s err %v", portStr, err)
		return eps
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		klog.Warningf("parse ip %s nil", ipAddress)
		return eps
	}

	nep := eps.DeepCopy()
	nep.Subsets = []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: ipAddress,
				},
			},
			Ports: []v1.EndpointPort{
				{
					Protocol: v1.ProtocolTCP,
					Port:     int32(port),
					Name:     "https",
				},
			},
		},
	}

	klog.V(4).Infof("gen new endpoint complete %v", nep)
	return nep
}

func pruneEndpoints(services map[types.NamespacedName]*serviceContainer,
	eps *v1.Endpoints, localAppInfo map[types.NamespacedName][]serf.Member,
	wrapperInCluster, serviceAutonomyEnhancementEnabled bool) *v1.Endpoints {

	epsKey := types.NamespacedName{Namespace: eps.Namespace, Name: eps.Name}

	if wrapperInCluster {
		eps = genLocalEndpoints(eps)
	}

	// dangling endpoints
	svc, ok := services[epsKey]
	if !ok {
		klog.V(4).Infof("Dangling endpoints %s, %+#v", eps.Name, eps.Subsets)
		return eps
	}

	// normal service
	if len(svc.keys) == 0 {
		klog.V(4).Infof("Normal endpoints %s, %+#v", eps.Name, eps.Subsets)
		if eps.Namespace == metav1.NamespaceDefault && eps.Name == MasterEndpointName {
			return eps
		}
		if serviceAutonomyEnhancementEnabled {
			newEps := eps.DeepCopy()
			for si := range newEps.Subsets {
				subnet := &newEps.Subsets[si]
				subnet.Addresses = filterLocalAppInfoConcernedAddresses(subnet.Addresses, localAppInfo[epsKey])
				subnet.NotReadyAddresses = filterLocalAppInfoConcernedAddresses(subnet.NotReadyAddresses, localAppInfo[epsKey])
			}
			klog.V(4).Infof("Normal endpoints after LocalNodeInfo filter %s: subnets from %+#v to %+#v", eps.Name, eps.Subsets, newEps.Subsets)
			return newEps
		}
		return eps
	}

	return eps
}

func filterLocalAppInfoConcernedAddresses(addresses []v1.EndpointAddress, members []serf.Member) []v1.EndpointAddress {
	localAppInfo := make(map[string]serf.Member)
	for _, member := range members {
		localAppInfo[member.Addr.String()] = member
	}

	filteredEndpointAddresses := make([]v1.EndpointAddress, 0)
	for i := range addresses {
		addr := addresses[i]
		_, found := localAppInfo[addr.IP]
		if found && localAppInfo[addr.IP].Status == serf.StatusAlive {
			filteredEndpointAddresses = append(filteredEndpointAddresses, addr)
			delete(localAppInfo, addr.IP)
		}
	}

	// TODO ?????????????????????alive??????????????????????????????ip
	if len(localAppInfo) != 0 {
		for _, member := range localAppInfo {
			if member.Status == serf.StatusAlive {
				epa := v1.EndpointAddress{
					IP: member.Addr.String(),
					TargetRef: &v1.ObjectReference{
						Kind:      "Pod",
						Namespace: member.Tags["namespace"],
						Name:      member.Name,
					},
				}
				filteredEndpointAddresses = append(filteredEndpointAddresses, epa)
			}
		}
	}
	return filteredEndpointAddresses
}
