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
// Date			:	July 9, 2017
//
// History		:
//	Date:			Author:		Info:
//	July 9, 2017	LIS			First release
//
// TODO:

package procfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// function to get current network devices info
func getNetInfo() map[string]*netDevice {
	contents, err := ioutil.ReadFile(PROC_SYS_NETDEV)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	// create the map
	netDevices := make(map[string]*netDevice)
	// prep the regex
	lineMatch, _ := regexp.Compile(netRegex)
	// read per line
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			// we only want those that matching parRegex
			match := lineMatch.MatchString(line)
			if match {
				deviceName := strings.Fields(line)[0]
				rxBytes, _ := strconv.ParseUint(strings.Fields(line)[1], 10, 64)
				rxPackets, _ := strconv.ParseUint(strings.Fields(line)[2], 10, 64)
				rxErrors, _ := strconv.ParseUint(strings.Fields(line)[3], 10, 64)
				rxDropped, _ := strconv.ParseUint(strings.Fields(line)[4], 10, 64)
				txBytes, _ := strconv.ParseUint(strings.Fields(line)[9], 10, 64)
				txPackets, _ := strconv.ParseUint(strings.Fields(line)[10], 10, 64)
				txErrors, _ := strconv.ParseUint(strings.Fields(line)[11], 10, 64)
				txDropped, _ := strconv.ParseUint(strings.Fields(line)[12], 10, 64)
				collisions, _ := strconv.ParseUint(strings.Fields(line)[14], 10, 64)
				carrier, _ := strconv.ParseUint(strings.Fields(line)[15], 10, 64)
				netDevices[deviceName] = &netDevice{
					ifName:     deviceName,
					rxBytes:    rxBytes,
					rxPackets:  rxPackets,
					rxErrors:   rxErrors,
					rxDropped:  rxDropped,
					txBytes:    txBytes,
					txPackets:  txPackets,
					txErrors:   txErrors,
					txDropped:  txDropped,
					collisions: collisions,
					carrier:    carrier,
				}
			}
		}
	}
	return netDevices
}

func NewNet() map[string]*netDevice {
	return getNetInfo()
}

// function to update network device stats
// NOTE: we do not call the indvidual () function since we like the call to bne as atomic as possible
func (netPtr *netDevice) Update() {
	contents, err := ioutil.ReadFile(PROC_SYS_NETDEV)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	// prep the regex
	lineMatch, _ := regexp.Compile(netRegex)
	// read per line
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			// we only want those that matching parRegex
			match := lineMatch.MatchString(line)
			if match {
				deviceName := strings.Fields(line)[0]
				if deviceName == netPtr.ifName {
					netPtr.rxBytes, _ = strconv.ParseUint(strings.Fields(line)[1], 10, 64)
					netPtr.rxPackets, _ = strconv.ParseUint(strings.Fields(line)[2], 10, 64)
					netPtr.rxErrors, _ = strconv.ParseUint(strings.Fields(line)[3], 10, 64)
					netPtr.rxDropped, _ = strconv.ParseUint(strings.Fields(line)[4], 10, 64)
					netPtr.txBytes, _ = strconv.ParseUint(strings.Fields(line)[9], 10, 64)
					netPtr.txPackets, _ = strconv.ParseUint(strings.Fields(line)[10], 10, 64)
					netPtr.txErrors, _ = strconv.ParseUint(strings.Fields(line)[11], 10, 64)
					netPtr.txDropped, _ = strconv.ParseUint(strings.Fields(line)[12], 10, 64)
					netPtr.collisions, _ = strconv.ParseUint(strings.Fields(line)[14], 10, 64)
					netPtr.carrier, _ = strconv.ParseUint(strings.Fields(line)[15], 10, 64)
				}
			}
		}
	}
}

// update the rxBytes value
func (netPtr *netDevice) RxBytes() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/rx_bytes", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.rxBytes, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.rxBytes
}

// update the rxPackets value
func (netPtr *netDevice) RxPackets() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/rx_packets", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.rxPackets, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.rxPackets
}

// update the rxErrors value
func (netPtr *netDevice) RxErrors() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/rx_errors", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.rxErrors, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.rxErrors
}

// update the rxDropped value
func (netPtr *netDevice) RxDropped() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/rx_dropped", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.rxDropped, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.rxDropped
}

// update the txBytes value
func (netPtr *netDevice) TxBytes() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/tx_bytes", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.txBytes, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.txBytes
}

// update the txPackets value
func (netPtr *netDevice) TxPackets() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/tx_packets", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.txPackets, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.txPackets
}

// update the txErrors value
func (netPtr *netDevice) TxErrors() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/tx_errors", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.txErrors, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.txErrors
}

// update the txDropped value
func (netPtr *netDevice) TxdRopped() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/tx_dropped", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.txDropped, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.txDropped
}

// update the collisions value
func (netPtr *netDevice) Collisions() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/collisions", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.collisions, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.collisions
}

// update the carrier value
func (netPtr *netDevice) Carrier() uint64 {
	procFile := fmt.Sprintf("/sys/class/net/%s/statistics/tx_carrier_errors", netPtr.ifName)
	contents, err := ioutil.ReadFile(procFile)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	netPtr.carrier, _ = strconv.ParseUint(string(line[0]), 10, 64)
	return netPtr.carrier
}
