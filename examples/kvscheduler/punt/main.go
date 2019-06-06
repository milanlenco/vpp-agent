//  Copyright (c) 2019 Cisco and/or its affiliates.
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
	"time"

	"github.com/ligato/cn-infra/agent"
	"github.com/ligato/vpp-agent/plugins/orchestrator"

	"github.com/ligato/vpp-agent/api/models/linux/interfaces"
	//"github.com/ligato/vpp-agent/api/models/linux/l3"
	"github.com/ligato/vpp-agent/api/models/vpp/interfaces"
	"github.com/ligato/vpp-agent/api/models/vpp/punt"
	"github.com/ligato/vpp-agent/api/models/vpp/l3"
	"github.com/ligato/vpp-agent/clientv2/linux/localclient"
	linux_ifplugin "github.com/ligato/vpp-agent/plugins/linux/ifplugin"
	linux_l3plugin "github.com/ligato/vpp-agent/plugins/linux/l3plugin"
	linux_nsplugin "github.com/ligato/vpp-agent/plugins/linux/nsplugin"
	vpp_ifplugin "github.com/ligato/vpp-agent/plugins/vpp/ifplugin"
	vpp_l3plugin "github.com/ligato/vpp-agent/plugins/vpp/l3plugin"
	vpp_puntplugin "github.com/ligato/vpp-agent/plugins/vpp/puntplugin"
	"github.com/ligato/vpp-agent/api/models/linux/namespace"
)

/*
	This example demonstrates punt plugin

	Deploy microservice for the gateway:

	term1$ docker run -it --rm  -e MICROSERVICE_LABEL=ms-gw lencomilan/ubuntu /bin/bash
	term1$ nc -l -p 8080 &
	term1$ nc -u -l -p 9090 &

	Deploy microservice for the client:

	term2$ docker run -it --rm  -e MICROSERVICE_LABEL=ms-client lencomilan/ubuntu /bin/bash
	term2$ nc 10.10.1.1 8080 # test TCP
	term2$ nc -u 10.10.1.1 9090 # test UDP
*/

func main() {
	// Set inter-dependency between VPP & Linux plugins
	vpp_ifplugin.DefaultPlugin.LinuxIfPlugin = &linux_ifplugin.DefaultPlugin
	vpp_ifplugin.DefaultPlugin.NsPlugin = &linux_nsplugin.DefaultPlugin
	linux_ifplugin.DefaultPlugin.VppIfPlugin = &vpp_ifplugin.DefaultPlugin

	ep := &ExamplePlugin{
		Orchestrator:  &orchestrator.DefaultPlugin,
		LinuxIfPlugin: &linux_ifplugin.DefaultPlugin,
		LinuxL3Plugin: &linux_l3plugin.DefaultPlugin,
		VPPIfPlugin:   &vpp_ifplugin.DefaultPlugin,
		VPPL3Plugin:   &vpp_l3plugin.DefaultPlugin,
		VPPPuntPlugin: &vpp_puntplugin.DefaultPlugin,
	}

	a := agent.NewAgent(
		agent.AllPlugins(ep),
	)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

// ExamplePlugin is the main plugin which
// handles resync and changes in this example.
type ExamplePlugin struct {
	LinuxIfPlugin *linux_ifplugin.IfPlugin
	LinuxL3Plugin *linux_l3plugin.L3Plugin
	VPPIfPlugin   *vpp_ifplugin.IfPlugin
	VPPL3Plugin   *vpp_l3plugin.L3Plugin
	VPPPuntPlugin *vpp_puntplugin.PuntPlugin
	Orchestrator  *orchestrator.Plugin
}

// String returns plugin name
func (p *ExamplePlugin) String() string {
	return "vpp-punt-example"
}

// Init handles initialization phase.
func (p *ExamplePlugin) Init() error {
	return nil
}

// AfterInit handles phase after initialization.
func (p *ExamplePlugin) AfterInit() error {
	go testLocalClientWithScheduler()
	return nil
}

// Close cleans up the resources.
func (p *ExamplePlugin) Close() error {
	return nil
}

func testLocalClientWithScheduler() {
	// initial resync
	time.Sleep(time.Second * 2)
	fmt.Println("=== RESYNC ===")

	txn := localclient.DataResyncRequest("example")
	err := txn.
		LinuxInterface(clientLinuxTap).
		LinuxInterface(gwLinuxTap).
		VppInterface(clientVPPTap).
		VppInterface(gwVPPTap).
		PuntIPRedirect(puntRule).
		ProxyArp(proxyArp).
		Arp(arpForGw).
		//LinuxArpEntry(arpForClient).
		Send().ReceiveReply()
	if err != nil {
		fmt.Println(err)
		return
	}
}

var (
	/* gw <-> VPP */

	gwLinuxTap = &linux_interfaces.Interface{
		Name:    "linux-tap-gw",
		Type:    linux_interfaces.Interface_TAP_TO_VPP,
		Enabled: true,
		IpAddresses: []string{
			"10.10.1.1/24",
		},
		PhysAddress: "aa:aa:aa:aa:aa:aa",
		HostIfName: "tap_to_vpp",
		Link: &linux_interfaces.Interface_Tap{
			Tap: &linux_interfaces.TapLink{
				VppTapIfName: "vpp-tap-gw",
			},
		},
		Namespace: &linux_namespace.NetNamespace{
			Type:      linux_namespace.NetNamespace_MICROSERVICE,
			Reference: "ms-gw",
		},
	}

	gwVPPTap = &vpp_interfaces.Interface{
		Name:    "vpp-tap-gw",
		Type:    vpp_interfaces.Interface_TAP,
		Enabled: true,
		Unnumbered: &vpp_interfaces.Interface_Unnumbered{
			InterfaceWithIp: "vpp-tap-client",
		},
		PhysAddress: "bb:bb:bb:bb:bb:bb",
		Link: &vpp_interfaces.Interface_Tap{
			Tap: &vpp_interfaces.TapLink{
				Version: 2,
				ToMicroservice: "ms-gw",
			},
		},
	}

	/* client <-> VPP */

	clientLinuxTap = &linux_interfaces.Interface{
		Name:    "linux-tap-client",
		Type:    linux_interfaces.Interface_TAP_TO_VPP,
		Enabled: true,
		IpAddresses: []string{
			"10.10.1.20/24",
		},
		PhysAddress: "cc:cc:cc:cc:cc:cc",
		HostIfName: "tap_to_vpp",
		Link: &linux_interfaces.Interface_Tap{
			Tap: &linux_interfaces.TapLink{
				VppTapIfName: "vpp-tap-client",
			},
		},
		Namespace: &linux_namespace.NetNamespace{
			Type:      linux_namespace.NetNamespace_MICROSERVICE,
			Reference: "ms-client",
		},
	}

	clientVPPTap = &vpp_interfaces.Interface{
		Name:    "vpp-tap-client",
		Type:    vpp_interfaces.Interface_TAP,
		Enabled: true,
		PhysAddress: "dd:dd:dd:dd:dd:dd",
		IpAddresses: []string{
			"10.10.1.1/24",
		},
		Link: &vpp_interfaces.Interface_Tap{
			Tap: &vpp_interfaces.TapLink{
				Version: 2,
				ToMicroservice: "ms-client",
			},
		},
	}


	/* punting */

	puntRule = &vpp_punt.IPRedirect{
		L3Protocol:  vpp_punt.L3Protocol_ALL,
		RxInterface: "vpp-tap-client",
		TxInterface: "vpp-tap-gw",
		NextHop:     "10.10.1.1",
	}

	/* static ARPs */

	proxyArp = &vpp_l3.ProxyARP{
		Interfaces: []*vpp_l3.ProxyARP_Interface{
			{
				Name:"vpp-tap-gw",
			},
			/*
			{
				Name:"vpp-tap-client",
			},
			*/
		},
		Ranges: []*vpp_l3.ProxyARP_Range{
			{
				FirstIpAddr: "10.10.1.1",
				LastIpAddr:  "10.10.1.255",
			},
		},
	}

	arpForGw = &vpp_l3.ARPEntry{
		Interface:   "vpp-tap-gw",
		IpAddress:   "10.10.1.1",
		PhysAddress: "aa:aa:aa:aa:aa:aa",
		Static:      true,
	}

	/*
	arpForClient = &linux_l3.ARPEntry{
		Interface: "linux-tap-client",
		IpAddress: "10.10.1.1",
		HwAddress: "dd:dd:dd:dd:dd:dd",
	}
	*/
)
