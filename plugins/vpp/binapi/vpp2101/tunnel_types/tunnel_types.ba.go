// Code generated by GoVPP's binapi-generator. DO NOT EDIT.

// Package tunnel_types contains generated bindings for API file tunnel_types.api.
//
// Contents:
//   2 enums
//
package tunnel_types

import (
	"strconv"

	api "git.fd.io/govpp.git/api"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the GoVPP api package it is being compiled against.
// A compilation error at this line likely means your copy of the
// GoVPP api package needs to be updated.
const _ = api.GoVppAPIPackageIsVersion2

// TunnelEncapDecapFlags defines enum 'tunnel_encap_decap_flags'.
type TunnelEncapDecapFlags uint8

const (
	TUNNEL_API_ENCAP_DECAP_FLAG_NONE             TunnelEncapDecapFlags = 0
	TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_DF    TunnelEncapDecapFlags = 1
	TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_SET_DF     TunnelEncapDecapFlags = 2
	TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_DSCP  TunnelEncapDecapFlags = 4
	TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_ECN   TunnelEncapDecapFlags = 8
	TUNNEL_API_ENCAP_DECAP_FLAG_DECAP_COPY_ECN   TunnelEncapDecapFlags = 16
	TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_INNER_HASH TunnelEncapDecapFlags = 32
)

var (
	TunnelEncapDecapFlags_name = map[uint8]string{
		0:  "TUNNEL_API_ENCAP_DECAP_FLAG_NONE",
		1:  "TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_DF",
		2:  "TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_SET_DF",
		4:  "TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_DSCP",
		8:  "TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_ECN",
		16: "TUNNEL_API_ENCAP_DECAP_FLAG_DECAP_COPY_ECN",
		32: "TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_INNER_HASH",
	}
	TunnelEncapDecapFlags_value = map[string]uint8{
		"TUNNEL_API_ENCAP_DECAP_FLAG_NONE":             0,
		"TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_DF":    1,
		"TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_SET_DF":     2,
		"TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_DSCP":  4,
		"TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_COPY_ECN":   8,
		"TUNNEL_API_ENCAP_DECAP_FLAG_DECAP_COPY_ECN":   16,
		"TUNNEL_API_ENCAP_DECAP_FLAG_ENCAP_INNER_HASH": 32,
	}
)

func (x TunnelEncapDecapFlags) String() string {
	s, ok := TunnelEncapDecapFlags_name[uint8(x)]
	if ok {
		return s
	}
	str := func(n uint8) string {
		s, ok := TunnelEncapDecapFlags_name[uint8(n)]
		if ok {
			return s
		}
		return "TunnelEncapDecapFlags(" + strconv.Itoa(int(n)) + ")"
	}
	for i := uint8(0); i <= 8; i++ {
		val := uint8(x)
		if val&(1<<i) != 0 {
			if s != "" {
				s += "|"
			}
			s += str(1 << i)
		}
	}
	if s == "" {
		return str(uint8(x))
	}
	return s
}

// TunnelMode defines enum 'tunnel_mode'.
type TunnelMode uint8

const (
	TUNNEL_API_MODE_P2P TunnelMode = 0
	TUNNEL_API_MODE_MP  TunnelMode = 1
)

var (
	TunnelMode_name = map[uint8]string{
		0: "TUNNEL_API_MODE_P2P",
		1: "TUNNEL_API_MODE_MP",
	}
	TunnelMode_value = map[string]uint8{
		"TUNNEL_API_MODE_P2P": 0,
		"TUNNEL_API_MODE_MP":  1,
	}
)

func (x TunnelMode) String() string {
	s, ok := TunnelMode_name[uint8(x)]
	if ok {
		return s
	}
	return "TunnelMode(" + strconv.Itoa(int(x)) + ")"
}