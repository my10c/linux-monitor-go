// Copyright (c) 2017 - 2017 badassops
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//	* Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	* Redistributions in binary form must reproduce the above copyright
//	notice, this list of conditions and the following disclaimer in the
//	documentation and/or other materials provided with the distribution.
//	* Neither the name of the <organization> nor the
//	names of its contributors may be used to endorse or promote products
//	derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSEcw
// ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Version		:	0.1
//
// Date			:	July 3, 2017
//
// History	:
// 	Date:			Author:		Info:
//	June 4, 2017	LIS			First Go release
//
// TODO: add more functions

package zookeeper

import (
	"fmt"
	"strings"
	"time"

	myUtils		"github.com/my10c/linux-monitor-go/utils"
	myGlobal	"github.com/my10c/linux-monitor-go/global"
	zkClass		"github.com/samuel/go-zookeeper/zk"
)

const (
	StateUnknown			State = -1
	StateDisconnected		State = 0
	StateConnecting			State = 1
	StateAuthFailed			State = 4
	StateConnectedReadOnly	State = 5
	StateSaslAuthenticated	State = 6
	StateExpired			State = -112
	StateConnected			= State(100)
	StateHasSession			= State(101)
)

var (
	stateNames = map[State]string{
		StateUnknown:			"StateUnknown",
		StateDisconnected:		"StateDisconnected",
		StateConnectedReadOnly:	"StateConnectedReadOnly",
		StateSaslAuthenticated:	"StateSaslAuthenticated",
		StateExpired:			"StateExpired",
		StateAuthFailed:		"StateAuthFailed",
		StateConnecting:		"StateConnecting",
		StateConnected:			"StateConnected",
		StateHasSession:		"StateHasSession",
	}
)

type State int32

type ZkConn struct {
	zkConnecion *zkClass.Conn
}

// Function to create the zookeeper object then connect to the given zookeeper server
// we need this to be able to send any command to kafka
func New(zkServers []string) (*ZkConn, error) {
	if len(zkServers) == 0 {
		zkServers = []string{"127.0.0.1:2181"}
	}
	connectedZK, _, err := zkClass.Connect(zkServers, time.Second)
	if err != nil {
		myUtils.ExitWithNagiosCode(myGlobal.CRITICAL, err)
	}
	connection := &ZkConn{
		zkConnecion	: connectedZK,
	}
	return connection, nil
}

// Function to close a zookeeper connection
func (zkPtr *ZkConn) Close() {
	zkPtr.zkConnecion.Close()
}

// Function to see if a zookeeper is a leader
func (zkPtr *ZkConn) IsLeader() bool {
	if zkPtr.IsFollower() {
		return false
	}
	return true
}

// Function to see if a zookeeper is a leader
func (zkPtr *ZkConn) IsFollower() bool {
	// TODO
	return false
}

// Function to see if a zookeeper is get traffic
func (zkPtr *ZkConn) Status(wantStat string) (string, bool){
	if zkPtr.zkConnecion.State().String() == wantStat {
		return zkPtr.zkConnecion.State().String(), true
	}
	return zkPtr.zkConnecion.State().String(), false
}

// Function to check current kafka brokers
// this is not a zookeeper thing but for now in this class
func (zkPtr *ZkConn) CheckKafkaBroker(brokerList []string) (int, error) {
	kafkaBroker, _, err := zkPtr.zkConnecion.Children("/brokers/ids")
	if err != nil {
		return myGlobal.UNKNOWN, err
	}
	var hit bool
	var missingKey string = ""
	for _, givenKey := range brokerList {
		for _, detectedKey := range kafkaBroker {
			hit = false
			if detectedKey == givenKey {
				hit = true
				break
			}
		}
		// create the missing key string
		if hit == false {
			missingKey = strings.TrimSpace(fmt.Sprintf("%s %s", missingKey, givenKey))
		}
	}
	if len(missingKey) > 0 {
		err := fmt.Errorf("Missing brokers ids: %s", missingKey)
		return myGlobal.CRITICAL, err
	}
	return myGlobal.OK, nil
}

