// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vppcalls

import (
	"net"
	"time"

	bfdapi "github.com/ligato/vpp-agent/plugins/vpp/binapi/bfd"
	"github.com/ligato/vpp-agent/plugins/vpp/model/bfd"
)

// BfdDetails is the wrapper structure for the BFD northbound API structure.
type BfdDetails struct {
	Bfd  *bfd.SingleHopBFD
	Meta *BfdMeta
}

// BfdMeta is combination of proto-modelled BFD data and VPP provided metadata
type BfdMeta struct {
	*BfdSessionMeta
	*BfdAuthKeyMeta
}

func (handler *bfdVppHandler) DumpBfdSingleHop() (*BfdDetails, error) {
	sessionDetails, err := handler.DumpBfdSessions()
	if err != nil {
		return nil, err
	}
	keyDetails, err := handler.DumpBfdAuthKeys()
	if err != nil {
		return nil, err
	}

	return &BfdDetails{
		Bfd: &bfd.SingleHopBFD{
			Sessions: sessionDetails.Session,
			Keys:     keyDetails.AuthKeys,
		},
		Meta: &BfdMeta{
			BfdSessionMeta: sessionDetails.Meta,
			BfdAuthKeyMeta: keyDetails.Meta,
		},
	}, nil
}

// BfdSessionDetails is the wrapper structure for the BFD session northbound API structure.
type BfdSessionDetails struct {
	Session []*bfd.SingleHopBFD_Session
	Meta    *BfdSessionMeta
}

// BfdSessionMeta is combination of proto-modelled BFD session data and session interface to index map
type BfdSessionMeta struct {
	SessionIfToIdx map[uint32]string
}

func (handler *bfdVppHandler) DumpBfdSessions() (*BfdSessionDetails, error) {
	defer func(t time.Time) {
		handler.stopwatch.TimeLog(bfdapi.BfdUDPSessionDump{}).LogTimeEntry(time.Since(t))
	}(time.Now())

	var sessions []*bfd.SingleHopBFD_Session
	meta := &BfdSessionMeta{
		SessionIfToIdx: make(map[uint32]string),
	}

	req := &bfdapi.BfdUDPSessionDump{}
	sessionsRequest := handler.callsChannel.SendMultiRequest(req)

	for {
		sessionDetails := &bfdapi.BfdUDPSessionDetails{}
		stop, err := sessionsRequest.ReceiveReply(sessionDetails)
		if stop {
			break
		}
		if err != nil {
			return nil, err
		}

		ifName, _, exists := handler.ifIndexes.LookupName(sessionDetails.SwIfIndex)
		if !exists {
			handler.log.Warnf("BFD session dump: interface name not found for index %d", sessionDetails.SwIfIndex)
		}
		var srcAddr, dstAddr net.IP = sessionDetails.LocalAddr, sessionDetails.PeerAddr

		// Put session info
		sessions = append(sessions, &bfd.SingleHopBFD_Session{
			Interface:             ifName,
			DestinationAddress:    dstAddr.String(),
			SourceAddress:         srcAddr.String(),
			DesiredMinTxInterval:  sessionDetails.DesiredMinTx,
			RequiredMinRxInterval: sessionDetails.RequiredMinRx,
			DetectMultiplier:      uint32(sessionDetails.DetectMult),
			Authentication: &bfd.SingleHopBFD_Session_Authentication{
				KeyId:           uint32(sessionDetails.BfdKeyID),
				AdvertisedKeyId: uint32(sessionDetails.ConfKeyID),
			},
		})
		// Put bfd interface info
		meta.SessionIfToIdx[sessionDetails.SwIfIndex] = ifName
	}

	return &BfdSessionDetails{
		Session: sessions,
		Meta:    meta,
	}, nil
}

func (handler *bfdVppHandler) DumpBfdUDPSessionsWithID(authKeyIndex uint32) (*BfdSessionDetails, error) {
	details, err := handler.DumpBfdSessions()
	if err != nil || len(details.Session) == 0 {
		return nil, err
	}

	var indexedSessions []*bfd.SingleHopBFD_Session
	for _, session := range details.Session {
		if session.Authentication != nil && session.Authentication.KeyId == authKeyIndex {
			indexedSessions = append(indexedSessions, session)
		}
	}

	return &BfdSessionDetails{
		Session: indexedSessions,
	}, nil
}

// BfdAuthKeyDetails is the wrapper structure for the BFD authentication key northbound API structure.
type BfdAuthKeyDetails struct {
	AuthKeys []*bfd.SingleHopBFD_Key
	Meta     *BfdAuthKeyMeta
}

// BfdAuthKeyMeta is combination of proto-modelled BFD session data and key-to-usage map
type BfdAuthKeyMeta struct {
	KeyIDToUseCount map[uint32]uint32
}

func (handler *bfdVppHandler) DumpBfdAuthKeys() (*BfdAuthKeyDetails, error) {
	defer func(t time.Time) {
		handler.stopwatch.TimeLog(bfdapi.BfdAuthKeysDump{}).LogTimeEntry(time.Since(t))
	}(time.Now())

	var authKeys []*bfd.SingleHopBFD_Key
	meta := &BfdAuthKeyMeta{
		KeyIDToUseCount: make(map[uint32]uint32),
	}

	req := &bfdapi.BfdAuthKeysDump{}
	keysRequest := handler.callsChannel.SendMultiRequest(req)

	for {
		keyDetails := &bfdapi.BfdAuthKeysDetails{}
		stop, err := keysRequest.ReceiveReply(keyDetails)
		if stop {
			break
		}
		if err != nil {
			return nil, err
		}

		// Put auth key info
		authKeys = append(authKeys, &bfd.SingleHopBFD_Key{
			AuthKeyIndex: keyDetails.ConfKeyID,
			Id:           keyDetails.ConfKeyID,
			AuthenticationType: func(authType uint8) bfd.SingleHopBFD_Key_AuthenticationType {
				if authType == 4 {
					return bfd.SingleHopBFD_Key_KEYED_SHA1
				}
				return bfd.SingleHopBFD_Key_METICULOUS_KEYED_SHA1
			}(keyDetails.AuthType),
		})
		// Put bfd key use count info
		meta.KeyIDToUseCount[keyDetails.ConfKeyID] = keyDetails.UseCount
	}

	return &BfdAuthKeyDetails{
		AuthKeys: authKeys,
		Meta:     meta,
	}, nil
}
