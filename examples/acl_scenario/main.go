//  Copyright (c) 2018 Cisco and/or its affiliates.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at:
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package main

import (
	"fmt"
	"log"

	"github.com/ligato/cn-infra/agent"
	"github.com/ligato/vpp-agent/clientv2/linux/localclient"
	"github.com/ligato/vpp-agent/plugins/orchestrator"

	"github.com/ligato/vpp-agent/api/models/linux/interfaces"
	"github.com/ligato/vpp-agent/api/models/linux/l3"
	linux_ns "github.com/ligato/vpp-agent/api/models/linux/namespace"
	"github.com/ligato/vpp-agent/api/models/vpp/interfaces"
	"github.com/ligato/vpp-agent/api/models/vpp/l2"
	"github.com/ligato/vpp-agent/api/models/vpp/acl"
	linux_ifplugin "github.com/ligato/vpp-agent/plugins/linux/ifplugin"
	linux_l3plugin "github.com/ligato/vpp-agent/plugins/linux/l3plugin"
	linux_nsplugin "github.com/ligato/vpp-agent/plugins/linux/nsplugin"
	vpp_aclplugin "github.com/ligato/vpp-agent/plugins/vpp/aclplugin"
	vpp_ifplugin "github.com/ligato/vpp-agent/plugins/vpp/ifplugin"
	vpp_l2plugin "github.com/ligato/vpp-agent/plugins/vpp/l2plugin"
)

/*
	topology: https://wiki.fd.io/images/d/dc/Routing_and_Switching_Tutorial_Topology.jpg

	+ ACL that allows only ns0 to ping ns2

	Deploy microservices (namespaces) with:

	host-term1$ docker run -it --rm  -e MICROSERVICE_LABEL=ns0 lencomilan/ubuntu /bin/bash
	host-term2$ docker run -it --rm  -e MICROSERVICE_LABEL=ns1 lencomilan/ubuntu /bin/bash
	host-term3$ docker run -it --rm  -e MICROSERVICE_LABEL=ns2 lencomilan/ubuntu /bin/bash
*/

func main() {
	// Set inter-dependency between VPP & Linux plugins
	vpp_ifplugin.DefaultPlugin.LinuxIfPlugin = &linux_ifplugin.DefaultPlugin
	vpp_ifplugin.DefaultPlugin.NsPlugin = &linux_nsplugin.DefaultPlugin
	linux_ifplugin.DefaultPlugin.VppIfPlugin = &vpp_ifplugin.DefaultPlugin

	ep := &App{
		Orchestrator:  &orchestrator.DefaultPlugin,
		LinuxIfPlugin: &linux_ifplugin.DefaultPlugin,
		LinuxL3Plugin: &linux_l3plugin.DefaultPlugin,
		VPPIfPlugin:   &vpp_ifplugin.DefaultPlugin,
		VPPL2Plugin:   &vpp_l2plugin.DefaultPlugin,
		VPPAclPlugin:  &vpp_aclplugin.DefaultPlugin,
	}

	a := agent.NewAgent(
		agent.AllPlugins(ep),
	)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

type App struct {
	LinuxIfPlugin *linux_ifplugin.IfPlugin
	LinuxL3Plugin *linux_l3plugin.L3Plugin
	VPPIfPlugin   *vpp_ifplugin.IfPlugin
	VPPL2Plugin   *vpp_l2plugin.L2Plugin
	VPPAclPlugin  *vpp_aclplugin.ACLPlugin
	Orchestrator  *orchestrator.Plugin
}

func (app *App) String() string {
	return "ACL-scenario"
}

func (app *App) Init() error {
	return nil
}

func (app *App) AfterInit() error {
	go deployTopology()
	return nil
}

func (app *App) Close() error {
	return nil
}

func deployTopology() {
	txn := localclient.DataResyncRequest("acl-scenario")
	err := txn.
		LinuxInterface(ns2LinuxTap).
		LinuxInterface(ns1Veth1).
		LinuxInterface(ns1Veth2).
		LinuxInterface(ns0Veth1).
		LinuxInterface(ns0Veth2).
		LinuxRoute(ns2LinuxRoute).
		LinuxRoute(ns1LinuxRoute).
		LinuxRoute(ns0LinuxRoute).
		VppInterface(ns2VppTap).
		VppInterface(ns1AfPacket).
		VppInterface(ns0AfPacket).
		VppInterface(bdLoopback).
		BD(bd).
		ACL(ns2Acl).
		Send().ReceiveReply()
	if err != nil {
		fmt.Println(err)
		return
	}
}

var (
	/* ns2 <-> VPP */
	ns2LinuxTap = &linux_interfaces.Interface{
		Name:    "ns2-linux-tap",
		Type:    linux_interfaces.Interface_TAP_TO_VPP,
		Enabled: true,
		IpAddresses: []string{
			"10.0.1.1/24",
		},
		HostIfName: "tap0",
		Link: &linux_interfaces.Interface_Tap{
			Tap: &linux_interfaces.TapLink{VppTapIfName: "ns2-vpp-tap"},
		},
		Namespace: &linux_ns.NetNamespace{
			Type:      linux_ns.NetNamespace_MICROSERVICE,
			Reference: "ns2",
		},
	}

	ns2LinuxRoute = &linux_l3.Route{
		OutgoingInterface: ns2LinuxTap.Name,
		Scope:             linux_l3.Route_GLOBAL,
		DstNetwork:        "10.0.0.0/16",
		GwAddr:            "10.0.1.10",
	}

	ns2VppTap = &vpp_interfaces.Interface{
		Name:    "ns2-vpp-tap",
		Type:    vpp_interfaces.Interface_TAP,
		Enabled: true,
		IpAddresses: []string{
			"10.0.1.10/24",
		},
		Link: &vpp_interfaces.Interface_Tap{
			Tap: &vpp_interfaces.TapLink{
				Version:        2,
				ToMicroservice: "ns2",
			},
		},
	}

	/* ns1 <-> VPP */
	ns1Veth1 = &linux_interfaces.Interface{
		Name:    "ns1-veth1",
		Type:    linux_interfaces.Interface_VETH,
		Enabled: true,
		IpAddresses: []string{
			"10.0.0.2/24",
		},
		HostIfName: "vpp1",
		Link: &linux_interfaces.Interface_Veth{
			Veth: &linux_interfaces.VethLink{PeerIfName: "ns1-veth2"},
		},
		Namespace: &linux_ns.NetNamespace{
			Type:      linux_ns.NetNamespace_MICROSERVICE,
			Reference: "ns1",
		},
	}

	ns1Veth2 = &linux_interfaces.Interface{
		Name:       "ns1-veth2",
		Type:       linux_interfaces.Interface_VETH,
		Enabled:    true,
		HostIfName: "vethns1",
		Link: &linux_interfaces.Interface_Veth{
			Veth: &linux_interfaces.VethLink{PeerIfName: "ns1-veth1"},
		},
	}

	ns1LinuxRoute = &linux_l3.Route{
		OutgoingInterface: ns1Veth1.Name,
		Scope:             linux_l3.Route_GLOBAL,
		DstNetwork:        "10.0.0.0/16",
		GwAddr:            "10.0.0.10",
	}

	ns1AfPacket = &vpp_interfaces.Interface{
		Name:    "ns1-afpacket",
		Type:    vpp_interfaces.Interface_AF_PACKET,
		Enabled: true,
		Link: &vpp_interfaces.Interface_Afpacket{
			Afpacket: &vpp_interfaces.AfpacketLink{
				HostIfName: ns1Veth2.HostIfName,
			},
		},
	}

	/* ns0 <-> VPP */
	ns0Veth1 = &linux_interfaces.Interface{
		Name:    "ns0-veth1",
		Type:    linux_interfaces.Interface_VETH,
		Enabled: true,
		IpAddresses: []string{
			"10.0.0.1/24",
		},
		HostIfName: "vpp0",
		Link: &linux_interfaces.Interface_Veth{
			Veth: &linux_interfaces.VethLink{PeerIfName: "ns0-veth2"},
		},
		Namespace: &linux_ns.NetNamespace{
			Type:      linux_ns.NetNamespace_MICROSERVICE,
			Reference: "ns0",
		},
	}

	ns0Veth2 = &linux_interfaces.Interface{
		Name:       "ns0-veth2",
		Type:       linux_interfaces.Interface_VETH,
		Enabled:    true,
		HostIfName: "vethns0",
		Link: &linux_interfaces.Interface_Veth{
			Veth: &linux_interfaces.VethLink{PeerIfName: "ns0-veth1"},
		},
	}

	ns0LinuxRoute = &linux_l3.Route{
		OutgoingInterface: ns0Veth1.Name,
		Scope:             linux_l3.Route_GLOBAL,
		DstNetwork:        "10.0.0.0/16",
		GwAddr:            "10.0.0.10",
	}

	ns0AfPacket = &vpp_interfaces.Interface{
		Name:    "ns0-afpacket",
		Type:    vpp_interfaces.Interface_AF_PACKET,
		Enabled: true,
		Link: &vpp_interfaces.Interface_Afpacket{
			Afpacket: &vpp_interfaces.AfpacketLink{
				HostIfName: ns0Veth2.HostIfName,
			},
		},
	}

	/* bridge domain */

	bdLoopback = &vpp_interfaces.Interface{
		Name:    "bd-loopback",
		Type:    vpp_interfaces.Interface_SOFTWARE_LOOPBACK,
		Enabled: true,
		IpAddresses: []string{
			"10.0.0.10/24",
		},
	}

	bd = &vpp_l2.BridgeDomain{
		Name:    "bd",
		Flood:   true,
		Forward: true,
		Learn:   true,
		Interfaces: []*vpp_l2.BridgeDomain_Interface{
			{
				Name:                    bdLoopback.Name,
				BridgedVirtualInterface: true,
			},
			{
				Name: ns0AfPacket.Name,
			},
			{
				Name: ns1AfPacket.Name,
			},
		},
	}

	/* policy */

	ns2Acl = &vpp_acl.ACL{
		Name: "ns2-acl",
		Interfaces: &vpp_acl.ACL_Interfaces{
			Egress: []string{ns2VppTap.Name},
		},
		Rules: []*vpp_acl.ACL_Rule{
			{ // allow any access from ns0
				Action: vpp_acl.ACL_Rule_PERMIT,
				IpRule: &vpp_acl.ACL_Rule_IpRule{
					Ip: &vpp_acl.ACL_Rule_IpRule_Ip{
						SourceNetwork: "10.0.0.1/32",
					},
				},
			},
			{ // deny everything/ICMP within BD, except for ns0, from accessing ns2
				Action: vpp_acl.ACL_Rule_DENY,
				IpRule: &vpp_acl.ACL_Rule_IpRule{
					Ip: &vpp_acl.ACL_Rule_IpRule_Ip{
						SourceNetwork: "10.0.0.0/24",
					},
					/* XXX Uncomment if only ICMP protocol from ns1 should be blocked
					Icmp: &vpp_acl.ACL_Rule_IpRule_Icmp{
						// ANY
						IcmpCodeRange: &vpp_acl.ACL_Rule_IpRule_Icmp_Range{
							First: 0,
							Last:  255,
						},
						IcmpTypeRange: &vpp_acl.ACL_Rule_IpRule_Icmp_Range{
							First: 0,
							Last:  255,
						},
					},
					*/
				},
			},
			{ // allows ns2 communicating with the rest of the world
				Action: vpp_acl.ACL_Rule_PERMIT,
				IpRule: &vpp_acl.ACL_Rule_IpRule{
					Ip: &vpp_acl.ACL_Rule_IpRule_Ip{
						// ANY
					},
				},
			},
		},
	}
)
