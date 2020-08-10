// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: ligato/vpp/vpp.proto

package vpp

import (
	proto "github.com/golang/protobuf/proto"
	abf "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/abf"
	acl "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/acl"
	interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	ipfix "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/ipfix"
	ipsec "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/ipsec"
	l2 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l2"
	l3 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l3"
	nat "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/nat"
	punt "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/punt"
	srv6 "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/srv6"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// ConfigData holds the entire VPP configuration.
type ConfigData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Interfaces             []*interfaces.Interface         `protobuf:"bytes,10,rep,name=interfaces,proto3" json:"interfaces,omitempty"`
	Spans                  []*interfaces.Span              `protobuf:"bytes,11,rep,name=spans,proto3" json:"spans,omitempty"`
	Acls                   []*acl.ACL                      `protobuf:"bytes,20,rep,name=acls,proto3" json:"acls,omitempty"`
	Abfs                   []*abf.ABF                      `protobuf:"bytes,21,rep,name=abfs,proto3" json:"abfs,omitempty"`
	BridgeDomains          []*l2.BridgeDomain              `protobuf:"bytes,30,rep,name=bridge_domains,json=bridgeDomains,proto3" json:"bridge_domains,omitempty"`
	Fibs                   []*l2.FIBEntry                  `protobuf:"bytes,31,rep,name=fibs,proto3" json:"fibs,omitempty"`
	XconnectPairs          []*l2.XConnectPair              `protobuf:"bytes,32,rep,name=xconnect_pairs,json=xconnectPairs,proto3" json:"xconnect_pairs,omitempty"`
	Routes                 []*l3.Route                     `protobuf:"bytes,40,rep,name=routes,proto3" json:"routes,omitempty"`
	Arps                   []*l3.ARPEntry                  `protobuf:"bytes,41,rep,name=arps,proto3" json:"arps,omitempty"`
	ProxyArp               *l3.ProxyARP                    `protobuf:"bytes,42,opt,name=proxy_arp,json=proxyArp,proto3" json:"proxy_arp,omitempty"`
	IpscanNeighbor         *l3.IPScanNeighbor              `protobuf:"bytes,43,opt,name=ipscan_neighbor,json=ipscanNeighbor,proto3" json:"ipscan_neighbor,omitempty"`
	Vrfs                   []*l3.VrfTable                  `protobuf:"bytes,44,rep,name=vrfs,proto3" json:"vrfs,omitempty"`
	L3Xconnects            []*l3.L3XConnect                `protobuf:"bytes,45,rep,name=l3xconnects,proto3" json:"l3xconnects,omitempty"`
	DhcpProxies            []*l3.DHCPProxy                 `protobuf:"bytes,46,rep,name=dhcp_proxies,json=dhcpProxies,proto3" json:"dhcp_proxies,omitempty"`
	TeibEntries            []*l3.TeibEntry                 `protobuf:"bytes,47,rep,name=teib_entries,json=teibEntries,proto3" json:"teib_entries,omitempty"`
	Nat44Global            *nat.Nat44Global                `protobuf:"bytes,50,opt,name=nat44_global,json=nat44Global,proto3" json:"nat44_global,omitempty"`
	Dnat44S                []*nat.DNat44                   `protobuf:"bytes,51,rep,name=dnat44s,proto3" json:"dnat44s,omitempty"`
	Nat44Interfaces        []*nat.Nat44Interface           `protobuf:"bytes,52,rep,name=nat44_interfaces,json=nat44Interfaces,proto3" json:"nat44_interfaces,omitempty"`
	Nat44Pools             []*nat.Nat44AddressPool         `protobuf:"bytes,53,rep,name=nat44_pools,json=nat44Pools,proto3" json:"nat44_pools,omitempty"`
	IpsecSpds              []*ipsec.SecurityPolicyDatabase `protobuf:"bytes,60,rep,name=ipsec_spds,json=ipsecSpds,proto3" json:"ipsec_spds,omitempty"`
	IpsecSas               []*ipsec.SecurityAssociation    `protobuf:"bytes,61,rep,name=ipsec_sas,json=ipsecSas,proto3" json:"ipsec_sas,omitempty"`
	IpsecTunnelProtections []*ipsec.TunnelProtection       `protobuf:"bytes,62,rep,name=ipsec_tunnel_protections,json=ipsecTunnelProtections,proto3" json:"ipsec_tunnel_protections,omitempty"`
	IpsecSps               []*ipsec.SecurityPolicy         `protobuf:"bytes,63,rep,name=ipsec_sps,json=ipsecSps,proto3" json:"ipsec_sps,omitempty"`
	PuntIpredirects        []*punt.IPRedirect              `protobuf:"bytes,70,rep,name=punt_ipredirects,json=puntIpredirects,proto3" json:"punt_ipredirects,omitempty"`
	PuntTohosts            []*punt.ToHost                  `protobuf:"bytes,71,rep,name=punt_tohosts,json=puntTohosts,proto3" json:"punt_tohosts,omitempty"`
	PuntExceptions         []*punt.Exception               `protobuf:"bytes,72,rep,name=punt_exceptions,json=puntExceptions,proto3" json:"punt_exceptions,omitempty"`
	Srv6Global             *srv6.SRv6Global                `protobuf:"bytes,83,opt,name=srv6_global,json=srv6Global,proto3" json:"srv6_global,omitempty"`
	Srv6Localsids          []*srv6.LocalSID                `protobuf:"bytes,80,rep,name=srv6_localsids,json=srv6Localsids,proto3" json:"srv6_localsids,omitempty"`
	Srv6Policies           []*srv6.Policy                  `protobuf:"bytes,81,rep,name=srv6_policies,json=srv6Policies,proto3" json:"srv6_policies,omitempty"`
	Srv6Steerings          []*srv6.Steering                `protobuf:"bytes,82,rep,name=srv6_steerings,json=srv6Steerings,proto3" json:"srv6_steerings,omitempty"`
	IpfixGlobal            *ipfix.IPFIX                    `protobuf:"bytes,90,opt,name=ipfix_global,json=ipfixGlobal,proto3" json:"ipfix_global,omitempty"`
	IpfixFlowprobeParams   *ipfix.FlowProbeParams          `protobuf:"bytes,91,opt,name=ipfix_flowprobe_params,json=ipfixFlowprobeParams,proto3" json:"ipfix_flowprobe_params,omitempty"`
	IpfixFlowprobes        []*ipfix.FlowProbeFeature       `protobuf:"bytes,92,rep,name=ipfix_flowprobes,json=ipfixFlowprobes,proto3" json:"ipfix_flowprobes,omitempty"`
}

func (x *ConfigData) Reset() {
	*x = ConfigData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ligato_vpp_vpp_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigData) ProtoMessage() {}

func (x *ConfigData) ProtoReflect() protoreflect.Message {
	mi := &file_ligato_vpp_vpp_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigData.ProtoReflect.Descriptor instead.
func (*ConfigData) Descriptor() ([]byte, []int) {
	return file_ligato_vpp_vpp_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigData) GetInterfaces() []*interfaces.Interface {
	if x != nil {
		return x.Interfaces
	}
	return nil
}

func (x *ConfigData) GetSpans() []*interfaces.Span {
	if x != nil {
		return x.Spans
	}
	return nil
}

func (x *ConfigData) GetAcls() []*acl.ACL {
	if x != nil {
		return x.Acls
	}
	return nil
}

func (x *ConfigData) GetAbfs() []*abf.ABF {
	if x != nil {
		return x.Abfs
	}
	return nil
}

func (x *ConfigData) GetBridgeDomains() []*l2.BridgeDomain {
	if x != nil {
		return x.BridgeDomains
	}
	return nil
}

func (x *ConfigData) GetFibs() []*l2.FIBEntry {
	if x != nil {
		return x.Fibs
	}
	return nil
}

func (x *ConfigData) GetXconnectPairs() []*l2.XConnectPair {
	if x != nil {
		return x.XconnectPairs
	}
	return nil
}

func (x *ConfigData) GetRoutes() []*l3.Route {
	if x != nil {
		return x.Routes
	}
	return nil
}

func (x *ConfigData) GetArps() []*l3.ARPEntry {
	if x != nil {
		return x.Arps
	}
	return nil
}

func (x *ConfigData) GetProxyArp() *l3.ProxyARP {
	if x != nil {
		return x.ProxyArp
	}
	return nil
}

func (x *ConfigData) GetIpscanNeighbor() *l3.IPScanNeighbor {
	if x != nil {
		return x.IpscanNeighbor
	}
	return nil
}

func (x *ConfigData) GetVrfs() []*l3.VrfTable {
	if x != nil {
		return x.Vrfs
	}
	return nil
}

func (x *ConfigData) GetL3Xconnects() []*l3.L3XConnect {
	if x != nil {
		return x.L3Xconnects
	}
	return nil
}

func (x *ConfigData) GetDhcpProxies() []*l3.DHCPProxy {
	if x != nil {
		return x.DhcpProxies
	}
	return nil
}

func (x *ConfigData) GetTeibEntries() []*l3.TeibEntry {
	if x != nil {
		return x.TeibEntries
	}
	return nil
}

func (x *ConfigData) GetNat44Global() *nat.Nat44Global {
	if x != nil {
		return x.Nat44Global
	}
	return nil
}

func (x *ConfigData) GetDnat44S() []*nat.DNat44 {
	if x != nil {
		return x.Dnat44S
	}
	return nil
}

func (x *ConfigData) GetNat44Interfaces() []*nat.Nat44Interface {
	if x != nil {
		return x.Nat44Interfaces
	}
	return nil
}

func (x *ConfigData) GetNat44Pools() []*nat.Nat44AddressPool {
	if x != nil {
		return x.Nat44Pools
	}
	return nil
}

func (x *ConfigData) GetIpsecSpds() []*ipsec.SecurityPolicyDatabase {
	if x != nil {
		return x.IpsecSpds
	}
	return nil
}

func (x *ConfigData) GetIpsecSas() []*ipsec.SecurityAssociation {
	if x != nil {
		return x.IpsecSas
	}
	return nil
}

func (x *ConfigData) GetIpsecTunnelProtections() []*ipsec.TunnelProtection {
	if x != nil {
		return x.IpsecTunnelProtections
	}
	return nil
}

func (x *ConfigData) GetIpsecSps() []*ipsec.SecurityPolicy {
	if x != nil {
		return x.IpsecSps
	}
	return nil
}

func (x *ConfigData) GetPuntIpredirects() []*punt.IPRedirect {
	if x != nil {
		return x.PuntIpredirects
	}
	return nil
}

func (x *ConfigData) GetPuntTohosts() []*punt.ToHost {
	if x != nil {
		return x.PuntTohosts
	}
	return nil
}

func (x *ConfigData) GetPuntExceptions() []*punt.Exception {
	if x != nil {
		return x.PuntExceptions
	}
	return nil
}

func (x *ConfigData) GetSrv6Global() *srv6.SRv6Global {
	if x != nil {
		return x.Srv6Global
	}
	return nil
}

func (x *ConfigData) GetSrv6Localsids() []*srv6.LocalSID {
	if x != nil {
		return x.Srv6Localsids
	}
	return nil
}

func (x *ConfigData) GetSrv6Policies() []*srv6.Policy {
	if x != nil {
		return x.Srv6Policies
	}
	return nil
}

func (x *ConfigData) GetSrv6Steerings() []*srv6.Steering {
	if x != nil {
		return x.Srv6Steerings
	}
	return nil
}

func (x *ConfigData) GetIpfixGlobal() *ipfix.IPFIX {
	if x != nil {
		return x.IpfixGlobal
	}
	return nil
}

func (x *ConfigData) GetIpfixFlowprobeParams() *ipfix.FlowProbeParams {
	if x != nil {
		return x.IpfixFlowprobeParams
	}
	return nil
}

func (x *ConfigData) GetIpfixFlowprobes() []*ipfix.FlowProbeFeature {
	if x != nil {
		return x.IpfixFlowprobes
	}
	return nil
}

type Notification struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Interface *interfaces.InterfaceNotification `protobuf:"bytes,1,opt,name=interface,proto3" json:"interface,omitempty"`
}

func (x *Notification) Reset() {
	*x = Notification{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ligato_vpp_vpp_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Notification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Notification) ProtoMessage() {}

func (x *Notification) ProtoReflect() protoreflect.Message {
	mi := &file_ligato_vpp_vpp_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Notification.ProtoReflect.Descriptor instead.
func (*Notification) Descriptor() ([]byte, []int) {
	return file_ligato_vpp_vpp_proto_rawDescGZIP(), []int{1}
}

func (x *Notification) GetInterface() *interfaces.InterfaceNotification {
	if x != nil {
		return x.Interface
	}
	return nil
}

type Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Interface *interfaces.InterfaceStats `protobuf:"bytes,1,opt,name=interface,proto3" json:"interface,omitempty"`
}

func (x *Stats) Reset() {
	*x = Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ligato_vpp_vpp_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stats) ProtoMessage() {}

func (x *Stats) ProtoReflect() protoreflect.Message {
	mi := &file_ligato_vpp_vpp_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stats.ProtoReflect.Descriptor instead.
func (*Stats) Descriptor() ([]byte, []int) {
	return file_ligato_vpp_vpp_proto_rawDescGZIP(), []int{2}
}

func (x *Stats) GetInterface() *interfaces.InterfaceStats {
	if x != nil {
		return x.Interface
	}
	return nil
}

var File_ligato_vpp_vpp_proto protoreflect.FileDescriptor

var file_ligato_vpp_vpp_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x76, 0x70, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76,
	0x70, 0x70, 0x1a, 0x18, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x61,
	0x62, 0x66, 0x2f, 0x61, 0x62, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x6c, 0x69,
	0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x61, 0x63, 0x6c, 0x2f, 0x61, 0x63, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76,
	0x70, 0x70, 0x2f, 0x69, 0x70, 0x66, 0x69, 0x78, 0x2f, 0x69, 0x70, 0x66, 0x69, 0x78, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70,
	0x2f, 0x69, 0x70, 0x66, 0x69, 0x78, 0x2f, 0x66, 0x6c, 0x6f, 0x77, 0x70, 0x72, 0x6f, 0x62, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76,
	0x70, 0x70, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x6c,
	0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x73, 0x2f, 0x73, 0x70, 0x61, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x21, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1c, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x69,
	0x70, 0x73, 0x65, 0x63, 0x2f, 0x69, 0x70, 0x73, 0x65, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x21, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x6c, 0x32, 0x2f,
	0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f,
	0x6c, 0x32, 0x2f, 0x66, 0x69, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x6c, 0x69,
	0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x6c, 0x32, 0x2f, 0x78, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x6c, 0x69, 0x67, 0x61,
	0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x6c, 0x33, 0x2f, 0x61, 0x72, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x16, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f,
	0x6c, 0x33, 0x2f, 0x6c, 0x33, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x6c, 0x69, 0x67,
	0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x6c, 0x33, 0x2f, 0x6c, 0x33, 0x78, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70,
	0x70, 0x2f, 0x6c, 0x33, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x18, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x6c, 0x33, 0x2f,
	0x74, 0x65, 0x69, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x6c, 0x69, 0x67, 0x61,
	0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x6c, 0x33, 0x2f, 0x76, 0x72, 0x66, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f,
	0x6e, 0x61, 0x74, 0x2f, 0x6e, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x6c,
	0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x70, 0x75, 0x6e, 0x74, 0x2f, 0x70,
	0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x6c, 0x69, 0x67, 0x61, 0x74,
	0x6f, 0x2f, 0x76, 0x70, 0x70, 0x2f, 0x73, 0x72, 0x76, 0x36, 0x2f, 0x73, 0x72, 0x76, 0x36, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x10, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x40, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63,
	0x65, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74,
	0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73,
	0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x12, 0x31, 0x0a, 0x05, 0x73, 0x70, 0x61, 0x6e, 0x73, 0x18,
	0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76,
	0x70, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x2e, 0x53, 0x70,
	0x61, 0x6e, 0x52, 0x05, 0x73, 0x70, 0x61, 0x6e, 0x73, 0x12, 0x27, 0x0a, 0x04, 0x61, 0x63, 0x6c,
	0x73, 0x18, 0x14, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f,
	0x2e, 0x76, 0x70, 0x70, 0x2e, 0x61, 0x63, 0x6c, 0x2e, 0x41, 0x43, 0x4c, 0x52, 0x04, 0x61, 0x63,
	0x6c, 0x73, 0x12, 0x27, 0x0a, 0x04, 0x61, 0x62, 0x66, 0x73, 0x18, 0x15, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x61, 0x62,
	0x66, 0x2e, 0x41, 0x42, 0x46, 0x52, 0x04, 0x61, 0x62, 0x66, 0x73, 0x12, 0x42, 0x0a, 0x0e, 0x62,
	0x72, 0x69, 0x64, 0x67, 0x65, 0x5f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x1e, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70,
	0x2e, 0x6c, 0x32, 0x2e, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e,
	0x52, 0x0d, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x12,
	0x2b, 0x0a, 0x04, 0x66, 0x69, 0x62, 0x73, 0x18, 0x1f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6c, 0x32, 0x2e, 0x46, 0x49,
	0x42, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x66, 0x69, 0x62, 0x73, 0x12, 0x42, 0x0a, 0x0e,
	0x78, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x5f, 0x70, 0x61, 0x69, 0x72, 0x73, 0x18, 0x20,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70,
	0x70, 0x2e, 0x6c, 0x32, 0x2e, 0x58, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x69,
	0x72, 0x52, 0x0d, 0x78, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x50, 0x61, 0x69, 0x72, 0x73,
	0x12, 0x2c, 0x0a, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18, 0x28, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6c, 0x33,
	0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x06, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x12, 0x2b,
	0x0a, 0x04, 0x61, 0x72, 0x70, 0x73, 0x18, 0x29, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6c,
	0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6c, 0x33, 0x2e, 0x41, 0x52, 0x50,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x61, 0x72, 0x70, 0x73, 0x12, 0x34, 0x0a, 0x09, 0x70,
	0x72, 0x6f, 0x78, 0x79, 0x5f, 0x61, 0x72, 0x70, 0x18, 0x2a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6c, 0x33, 0x2e, 0x50,
	0x72, 0x6f, 0x78, 0x79, 0x41, 0x52, 0x50, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x41, 0x72,
	0x70, 0x12, 0x46, 0x0a, 0x0f, 0x69, 0x70, 0x73, 0x63, 0x61, 0x6e, 0x5f, 0x6e, 0x65, 0x69, 0x67,
	0x68, 0x62, 0x6f, 0x72, 0x18, 0x2b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x6c, 0x69, 0x67,
	0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6c, 0x33, 0x2e, 0x49, 0x50, 0x53, 0x63, 0x61,
	0x6e, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x52, 0x0e, 0x69, 0x70, 0x73, 0x63, 0x61,
	0x6e, 0x4e, 0x65, 0x69, 0x67, 0x68, 0x62, 0x6f, 0x72, 0x12, 0x2b, 0x0a, 0x04, 0x76, 0x72, 0x66,
	0x73, 0x18, 0x2c, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f,
	0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6c, 0x33, 0x2e, 0x56, 0x72, 0x66, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x04, 0x76, 0x72, 0x66, 0x73, 0x12, 0x3b, 0x0a, 0x0b, 0x6c, 0x33, 0x78, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x73, 0x18, 0x2d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6c, 0x69,
	0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6c, 0x33, 0x2e, 0x4c, 0x33, 0x58, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x0b, 0x6c, 0x33, 0x78, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x73, 0x12, 0x3b, 0x0a, 0x0c, 0x64, 0x68, 0x63, 0x70, 0x5f, 0x70, 0x72, 0x6f, 0x78,
	0x69, 0x65, 0x73, 0x18, 0x2e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6c, 0x69, 0x67, 0x61,
	0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6c, 0x33, 0x2e, 0x44, 0x48, 0x43, 0x50, 0x50, 0x72,
	0x6f, 0x78, 0x79, 0x52, 0x0b, 0x64, 0x68, 0x63, 0x70, 0x50, 0x72, 0x6f, 0x78, 0x69, 0x65, 0x73,
	0x12, 0x3b, 0x0a, 0x0c, 0x74, 0x65, 0x69, 0x62, 0x5f, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x18, 0x2f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e,
	0x76, 0x70, 0x70, 0x2e, 0x6c, 0x33, 0x2e, 0x54, 0x65, 0x69, 0x62, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0b, 0x74, 0x65, 0x69, 0x62, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x3e, 0x0a,
	0x0c, 0x6e, 0x61, 0x74, 0x34, 0x34, 0x5f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x18, 0x32, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70,
	0x2e, 0x6e, 0x61, 0x74, 0x2e, 0x4e, 0x61, 0x74, 0x34, 0x34, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c,
	0x52, 0x0b, 0x6e, 0x61, 0x74, 0x34, 0x34, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x12, 0x30, 0x0a,
	0x07, 0x64, 0x6e, 0x61, 0x74, 0x34, 0x34, 0x73, 0x18, 0x33, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6e, 0x61, 0x74, 0x2e,
	0x44, 0x4e, 0x61, 0x74, 0x34, 0x34, 0x52, 0x07, 0x64, 0x6e, 0x61, 0x74, 0x34, 0x34, 0x73, 0x12,
	0x49, 0x0a, 0x10, 0x6e, 0x61, 0x74, 0x34, 0x34, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61,
	0x63, 0x65, 0x73, 0x18, 0x34, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6c, 0x69, 0x67, 0x61,
	0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6e, 0x61, 0x74, 0x2e, 0x4e, 0x61, 0x74, 0x34, 0x34,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x52, 0x0f, 0x6e, 0x61, 0x74, 0x34, 0x34,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x12, 0x41, 0x0a, 0x0b, 0x6e, 0x61,
	0x74, 0x34, 0x34, 0x5f, 0x70, 0x6f, 0x6f, 0x6c, 0x73, 0x18, 0x35, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x6e, 0x61, 0x74,
	0x2e, 0x4e, 0x61, 0x74, 0x34, 0x34, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x50, 0x6f, 0x6f,
	0x6c, 0x52, 0x0a, 0x6e, 0x61, 0x74, 0x34, 0x34, 0x50, 0x6f, 0x6f, 0x6c, 0x73, 0x12, 0x47, 0x0a,
	0x0a, 0x69, 0x70, 0x73, 0x65, 0x63, 0x5f, 0x73, 0x70, 0x64, 0x73, 0x18, 0x3c, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x28, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69,
	0x70, 0x73, 0x65, 0x63, 0x2e, 0x53, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x50, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x52, 0x09, 0x69, 0x70, 0x73,
	0x65, 0x63, 0x53, 0x70, 0x64, 0x73, 0x12, 0x42, 0x0a, 0x09, 0x69, 0x70, 0x73, 0x65, 0x63, 0x5f,
	0x73, 0x61, 0x73, 0x18, 0x3d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6c, 0x69, 0x67, 0x61,
	0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69, 0x70, 0x73, 0x65, 0x63, 0x2e, 0x53, 0x65, 0x63,
	0x75, 0x72, 0x69, 0x74, 0x79, 0x41, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x08, 0x69, 0x70, 0x73, 0x65, 0x63, 0x53, 0x61, 0x73, 0x12, 0x5c, 0x0a, 0x18, 0x69, 0x70,
	0x73, 0x65, 0x63, 0x5f, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x3e, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6c,
	0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69, 0x70, 0x73, 0x65, 0x63, 0x2e,
	0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x50, 0x72, 0x6f, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x16, 0x69, 0x70, 0x73, 0x65, 0x63, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x50, 0x72, 0x6f,
	0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3d, 0x0a, 0x09, 0x69, 0x70, 0x73, 0x65,
	0x63, 0x5f, 0x73, 0x70, 0x73, 0x18, 0x3f, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x6c, 0x69,
	0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69, 0x70, 0x73, 0x65, 0x63, 0x2e, 0x53,
	0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x08, 0x69,
	0x70, 0x73, 0x65, 0x63, 0x53, 0x70, 0x73, 0x12, 0x46, 0x0a, 0x10, 0x70, 0x75, 0x6e, 0x74, 0x5f,
	0x69, 0x70, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x73, 0x18, 0x46, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x70,
	0x75, 0x6e, 0x74, 0x2e, 0x49, 0x50, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x0f,
	0x70, 0x75, 0x6e, 0x74, 0x49, 0x70, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x73, 0x12,
	0x3a, 0x0a, 0x0c, 0x70, 0x75, 0x6e, 0x74, 0x5f, 0x74, 0x6f, 0x68, 0x6f, 0x73, 0x74, 0x73, 0x18,
	0x47, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76,
	0x70, 0x70, 0x2e, 0x70, 0x75, 0x6e, 0x74, 0x2e, 0x54, 0x6f, 0x48, 0x6f, 0x73, 0x74, 0x52, 0x0b,
	0x70, 0x75, 0x6e, 0x74, 0x54, 0x6f, 0x68, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x43, 0x0a, 0x0f, 0x70,
	0x75, 0x6e, 0x74, 0x5f, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x48,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70,
	0x70, 0x2e, 0x70, 0x75, 0x6e, 0x74, 0x2e, 0x45, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x0e, 0x70, 0x75, 0x6e, 0x74, 0x45, 0x78, 0x63, 0x65, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x3c, 0x0a, 0x0b, 0x73, 0x72, 0x76, 0x36, 0x5f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x18,
	0x53, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76,
	0x70, 0x70, 0x2e, 0x73, 0x72, 0x76, 0x36, 0x2e, 0x53, 0x52, 0x76, 0x36, 0x47, 0x6c, 0x6f, 0x62,
	0x61, 0x6c, 0x52, 0x0a, 0x73, 0x72, 0x76, 0x36, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x12, 0x40,
	0x0a, 0x0e, 0x73, 0x72, 0x76, 0x36, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x73, 0x69, 0x64, 0x73,
	0x18, 0x50, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e,
	0x76, 0x70, 0x70, 0x2e, 0x73, 0x72, 0x76, 0x36, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x53, 0x49,
	0x44, 0x52, 0x0d, 0x73, 0x72, 0x76, 0x36, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x73, 0x69, 0x64, 0x73,
	0x12, 0x3c, 0x0a, 0x0d, 0x73, 0x72, 0x76, 0x36, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65,
	0x73, 0x18, 0x51, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f,
	0x2e, 0x76, 0x70, 0x70, 0x2e, 0x73, 0x72, 0x76, 0x36, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x52, 0x0c, 0x73, 0x72, 0x76, 0x36, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x12, 0x40,
	0x0a, 0x0e, 0x73, 0x72, 0x76, 0x36, 0x5f, 0x73, 0x74, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x73,
	0x18, 0x52, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e,
	0x76, 0x70, 0x70, 0x2e, 0x73, 0x72, 0x76, 0x36, 0x2e, 0x53, 0x74, 0x65, 0x65, 0x72, 0x69, 0x6e,
	0x67, 0x52, 0x0d, 0x73, 0x72, 0x76, 0x36, 0x53, 0x74, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x73,
	0x12, 0x3a, 0x0a, 0x0c, 0x69, 0x70, 0x66, 0x69, 0x78, 0x5f, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c,
	0x18, 0x5a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e,
	0x76, 0x70, 0x70, 0x2e, 0x69, 0x70, 0x66, 0x69, 0x78, 0x2e, 0x49, 0x50, 0x46, 0x49, 0x58, 0x52,
	0x0b, 0x69, 0x70, 0x66, 0x69, 0x78, 0x47, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x12, 0x57, 0x0a, 0x16,
	0x69, 0x70, 0x66, 0x69, 0x78, 0x5f, 0x66, 0x6c, 0x6f, 0x77, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x5f,
	0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x5b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x6c,
	0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69, 0x70, 0x66, 0x69, 0x78, 0x2e,
	0x46, 0x6c, 0x6f, 0x77, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x52,
	0x14, 0x69, 0x70, 0x66, 0x69, 0x78, 0x46, 0x6c, 0x6f, 0x77, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x4d, 0x0a, 0x10, 0x69, 0x70, 0x66, 0x69, 0x78, 0x5f, 0x66,
	0x6c, 0x6f, 0x77, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x73, 0x18, 0x5c, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69, 0x70, 0x66,
	0x69, 0x78, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x46, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x52, 0x0f, 0x69, 0x70, 0x66, 0x69, 0x78, 0x46, 0x6c, 0x6f, 0x77, 0x70, 0x72,
	0x6f, 0x62, 0x65, 0x73, 0x22, 0x5a, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4a, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f,
	0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x73, 0x2e,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
	0x22, 0x4c, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x43, 0x0a, 0x09, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6c,
	0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x76, 0x70, 0x70, 0x2e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66,
	0x61, 0x63, 0x65, 0x73, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x73, 0x52, 0x09, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65, 0x42, 0x2c,
	0x5a, 0x2a, 0x67, 0x6f, 0x2e, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2e, 0x69, 0x6f, 0x2f, 0x76,
	0x70, 0x70, 0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x33, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x6c, 0x69, 0x67, 0x61, 0x74, 0x6f, 0x2f, 0x76, 0x70, 0x70, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ligato_vpp_vpp_proto_rawDescOnce sync.Once
	file_ligato_vpp_vpp_proto_rawDescData = file_ligato_vpp_vpp_proto_rawDesc
)

func file_ligato_vpp_vpp_proto_rawDescGZIP() []byte {
	file_ligato_vpp_vpp_proto_rawDescOnce.Do(func() {
		file_ligato_vpp_vpp_proto_rawDescData = protoimpl.X.CompressGZIP(file_ligato_vpp_vpp_proto_rawDescData)
	})
	return file_ligato_vpp_vpp_proto_rawDescData
}

var file_ligato_vpp_vpp_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_ligato_vpp_vpp_proto_goTypes = []interface{}{
	(*ConfigData)(nil),                       // 0: ligato.vpp.ConfigData
	(*Notification)(nil),                     // 1: ligato.vpp.Notification
	(*Stats)(nil),                            // 2: ligato.vpp.Stats
	(*interfaces.Interface)(nil),             // 3: ligato.vpp.interfaces.Interface
	(*interfaces.Span)(nil),                  // 4: ligato.vpp.interfaces.Span
	(*acl.ACL)(nil),                          // 5: ligato.vpp.acl.ACL
	(*abf.ABF)(nil),                          // 6: ligato.vpp.abf.ABF
	(*l2.BridgeDomain)(nil),                  // 7: ligato.vpp.l2.BridgeDomain
	(*l2.FIBEntry)(nil),                      // 8: ligato.vpp.l2.FIBEntry
	(*l2.XConnectPair)(nil),                  // 9: ligato.vpp.l2.XConnectPair
	(*l3.Route)(nil),                         // 10: ligato.vpp.l3.Route
	(*l3.ARPEntry)(nil),                      // 11: ligato.vpp.l3.ARPEntry
	(*l3.ProxyARP)(nil),                      // 12: ligato.vpp.l3.ProxyARP
	(*l3.IPScanNeighbor)(nil),                // 13: ligato.vpp.l3.IPScanNeighbor
	(*l3.VrfTable)(nil),                      // 14: ligato.vpp.l3.VrfTable
	(*l3.L3XConnect)(nil),                    // 15: ligato.vpp.l3.L3XConnect
	(*l3.DHCPProxy)(nil),                     // 16: ligato.vpp.l3.DHCPProxy
	(*l3.TeibEntry)(nil),                     // 17: ligato.vpp.l3.TeibEntry
	(*nat.Nat44Global)(nil),                  // 18: ligato.vpp.nat.Nat44Global
	(*nat.DNat44)(nil),                       // 19: ligato.vpp.nat.DNat44
	(*nat.Nat44Interface)(nil),               // 20: ligato.vpp.nat.Nat44Interface
	(*nat.Nat44AddressPool)(nil),             // 21: ligato.vpp.nat.Nat44AddressPool
	(*ipsec.SecurityPolicyDatabase)(nil),     // 22: ligato.vpp.ipsec.SecurityPolicyDatabase
	(*ipsec.SecurityAssociation)(nil),        // 23: ligato.vpp.ipsec.SecurityAssociation
	(*ipsec.TunnelProtection)(nil),           // 24: ligato.vpp.ipsec.TunnelProtection
	(*ipsec.SecurityPolicy)(nil),             // 25: ligato.vpp.ipsec.SecurityPolicy
	(*punt.IPRedirect)(nil),                  // 26: ligato.vpp.punt.IPRedirect
	(*punt.ToHost)(nil),                      // 27: ligato.vpp.punt.ToHost
	(*punt.Exception)(nil),                   // 28: ligato.vpp.punt.Exception
	(*srv6.SRv6Global)(nil),                  // 29: ligato.vpp.srv6.SRv6Global
	(*srv6.LocalSID)(nil),                    // 30: ligato.vpp.srv6.LocalSID
	(*srv6.Policy)(nil),                      // 31: ligato.vpp.srv6.Policy
	(*srv6.Steering)(nil),                    // 32: ligato.vpp.srv6.Steering
	(*ipfix.IPFIX)(nil),                      // 33: ligato.vpp.ipfix.IPFIX
	(*ipfix.FlowProbeParams)(nil),            // 34: ligato.vpp.ipfix.FlowProbeParams
	(*ipfix.FlowProbeFeature)(nil),           // 35: ligato.vpp.ipfix.FlowProbeFeature
	(*interfaces.InterfaceNotification)(nil), // 36: ligato.vpp.interfaces.InterfaceNotification
	(*interfaces.InterfaceStats)(nil),        // 37: ligato.vpp.interfaces.InterfaceStats
}
var file_ligato_vpp_vpp_proto_depIdxs = []int32{
	3,  // 0: ligato.vpp.ConfigData.interfaces:type_name -> ligato.vpp.interfaces.Interface
	4,  // 1: ligato.vpp.ConfigData.spans:type_name -> ligato.vpp.interfaces.Span
	5,  // 2: ligato.vpp.ConfigData.acls:type_name -> ligato.vpp.acl.ACL
	6,  // 3: ligato.vpp.ConfigData.abfs:type_name -> ligato.vpp.abf.ABF
	7,  // 4: ligato.vpp.ConfigData.bridge_domains:type_name -> ligato.vpp.l2.BridgeDomain
	8,  // 5: ligato.vpp.ConfigData.fibs:type_name -> ligato.vpp.l2.FIBEntry
	9,  // 6: ligato.vpp.ConfigData.xconnect_pairs:type_name -> ligato.vpp.l2.XConnectPair
	10, // 7: ligato.vpp.ConfigData.routes:type_name -> ligato.vpp.l3.Route
	11, // 8: ligato.vpp.ConfigData.arps:type_name -> ligato.vpp.l3.ARPEntry
	12, // 9: ligato.vpp.ConfigData.proxy_arp:type_name -> ligato.vpp.l3.ProxyARP
	13, // 10: ligato.vpp.ConfigData.ipscan_neighbor:type_name -> ligato.vpp.l3.IPScanNeighbor
	14, // 11: ligato.vpp.ConfigData.vrfs:type_name -> ligato.vpp.l3.VrfTable
	15, // 12: ligato.vpp.ConfigData.l3xconnects:type_name -> ligato.vpp.l3.L3XConnect
	16, // 13: ligato.vpp.ConfigData.dhcp_proxies:type_name -> ligato.vpp.l3.DHCPProxy
	17, // 14: ligato.vpp.ConfigData.teib_entries:type_name -> ligato.vpp.l3.TeibEntry
	18, // 15: ligato.vpp.ConfigData.nat44_global:type_name -> ligato.vpp.nat.Nat44Global
	19, // 16: ligato.vpp.ConfigData.dnat44s:type_name -> ligato.vpp.nat.DNat44
	20, // 17: ligato.vpp.ConfigData.nat44_interfaces:type_name -> ligato.vpp.nat.Nat44Interface
	21, // 18: ligato.vpp.ConfigData.nat44_pools:type_name -> ligato.vpp.nat.Nat44AddressPool
	22, // 19: ligato.vpp.ConfigData.ipsec_spds:type_name -> ligato.vpp.ipsec.SecurityPolicyDatabase
	23, // 20: ligato.vpp.ConfigData.ipsec_sas:type_name -> ligato.vpp.ipsec.SecurityAssociation
	24, // 21: ligato.vpp.ConfigData.ipsec_tunnel_protections:type_name -> ligato.vpp.ipsec.TunnelProtection
	25, // 22: ligato.vpp.ConfigData.ipsec_sps:type_name -> ligato.vpp.ipsec.SecurityPolicy
	26, // 23: ligato.vpp.ConfigData.punt_ipredirects:type_name -> ligato.vpp.punt.IPRedirect
	27, // 24: ligato.vpp.ConfigData.punt_tohosts:type_name -> ligato.vpp.punt.ToHost
	28, // 25: ligato.vpp.ConfigData.punt_exceptions:type_name -> ligato.vpp.punt.Exception
	29, // 26: ligato.vpp.ConfigData.srv6_global:type_name -> ligato.vpp.srv6.SRv6Global
	30, // 27: ligato.vpp.ConfigData.srv6_localsids:type_name -> ligato.vpp.srv6.LocalSID
	31, // 28: ligato.vpp.ConfigData.srv6_policies:type_name -> ligato.vpp.srv6.Policy
	32, // 29: ligato.vpp.ConfigData.srv6_steerings:type_name -> ligato.vpp.srv6.Steering
	33, // 30: ligato.vpp.ConfigData.ipfix_global:type_name -> ligato.vpp.ipfix.IPFIX
	34, // 31: ligato.vpp.ConfigData.ipfix_flowprobe_params:type_name -> ligato.vpp.ipfix.FlowProbeParams
	35, // 32: ligato.vpp.ConfigData.ipfix_flowprobes:type_name -> ligato.vpp.ipfix.FlowProbeFeature
	36, // 33: ligato.vpp.Notification.interface:type_name -> ligato.vpp.interfaces.InterfaceNotification
	37, // 34: ligato.vpp.Stats.interface:type_name -> ligato.vpp.interfaces.InterfaceStats
	35, // [35:35] is the sub-list for method output_type
	35, // [35:35] is the sub-list for method input_type
	35, // [35:35] is the sub-list for extension type_name
	35, // [35:35] is the sub-list for extension extendee
	0,  // [0:35] is the sub-list for field type_name
}

func init() { file_ligato_vpp_vpp_proto_init() }
func file_ligato_vpp_vpp_proto_init() {
	if File_ligato_vpp_vpp_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ligato_vpp_vpp_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ligato_vpp_vpp_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Notification); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ligato_vpp_vpp_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ligato_vpp_vpp_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ligato_vpp_vpp_proto_goTypes,
		DependencyIndexes: file_ligato_vpp_vpp_proto_depIdxs,
		MessageInfos:      file_ligato_vpp_vpp_proto_msgTypes,
	}.Build()
	File_ligato_vpp_vpp_proto = out.File
	file_ligato_vpp_vpp_proto_rawDesc = nil
	file_ligato_vpp_vpp_proto_goTypes = nil
	file_ligato_vpp_vpp_proto_depIdxs = nil
}
