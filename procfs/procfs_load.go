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
	"strconv"
	"strings"
)

// function to get current load info
func getLoadInfo() *sysLoadavg {
	contents, err := ioutil.ReadFile(PROC_SYS_LOADAVG)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), "\n")
	execVals := strings.Split(string(line[3]), "/")
	load1Avg, _ := strconv.ParseFloat(string(line[0]), 64)
	load5Avg, _ := strconv.ParseFloat(string(line[1]), 64)
	load10Avg, _ := strconv.ParseFloat(string(line[2]), 64)
	execProc, _ := strconv.ParseUint(execVals[0], 10, 64)
	execQueue, _ := strconv.ParseUint(execVals[1], 10, 64)
	lastPid, _ := strconv.ParseUint(string(line[4]), 10, 64)
	currLoad := &sysLoadavg{
		load1Avg:  load1Avg,
		load5Avg:  load5Avg,
		load10Avg: load10Avg,
		execProc:  execProc,
		execQueue: execQueue,
		lastPid:   lastPid,
	}
	return currLoad
}

func NewLoad() *sysLoadavg {
	return getLoadInfo()
}

// function to update network device stats
// NOTE: we do not call the indvidual () function since we like the call to bne as atomic as possible
func (loadPtr *sysLoadavg) Update() {
	getLoadInfo()
}

// update the load1Avg value
func (loadPtr *sysLoadavg) Load1Avg() float64 {
	loadPtr.Update()
	return loadPtr.load1Avg
}

// update the load5Avg value
func (loadPtr *sysLoadavg) Load5Avg() float64 {
	loadPtr.Update()
	return loadPtr.load5Avg
}

// update the load10Avg value
func (loadPtr *sysLoadavg) Load10Avg() float64 {
	loadPtr.Update()
	return loadPtr.load10Avg
}

// update the execProc value
func (loadPtr *sysLoadavg) ExecProc() uint64 {
	loadPtr.Update()
	return loadPtr.execProc
}

// update the execQueue value
func (loadPtr *sysLoadavg) ExecQueue() uint64 {
	loadPtr.Update()
	return loadPtr.execQueue
}

// update the lastPid value
func (loadPtr *sysLoadavg) LastPid() uint64 {
	loadPtr.Update()
	return loadPtr.lastPid
}
