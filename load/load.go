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
// Date			:	July 1, 2017
//
// History	:
// 	Date:			Author:		Info:
//	July 1, 2017	LIS			First Go release
//
// TODO:

package load

import (
	"fmt"
	myThreshold "github.com/my10c/linux-monitor-go/threshold"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	// system wide
	PROC_SYS_LOADAVG = "/proc/loadavg"
	PROC_SYS_STATS   = "/proc/stat"
)

var (
	// we need cpuXXX only
	cpuRegex = `^cpu\d.*`
)

type sysLoadavg struct {
	load1Avg  float64 `json:"load1navg"`
	load5Avg  float64 `json:"load5avg"`
	load10Avg float64 `json:"load10avg"`
	execProc  uint64  `json:"execproc"`
	execQueue uint64  `json:"execqueue"`
	lastPid   uint64  `json:"lastpid"`
	cpus      float64 `json:"cpus"`
}

func getCPU() float64 {
	var cnt float64 = 0
	contents, err := ioutil.ReadFile(PROC_SYS_STATS)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	lineMatch, _ := regexp.Compile(cpuRegex)
	// get all lines and walk one at the time
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		match := lineMatch.MatchString(line)
		if match {
			cnt++
		}
	}
	return cnt
}

// function to get current load info
func getLoadInfo() *sysLoadavg {
	contents, err := ioutil.ReadFile(PROC_SYS_LOADAVG)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	line := strings.Split(string(contents), " ")
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
		cpus:      getCPU(),
	}
	return currLoad
}

func New() *sysLoadavg {
	return getLoadInfo()
}

// function to update load stats
// NOTE: we do not call the indvidual function since we like the call to be as atomic as possible
func (loadPtr *sysLoadavg) Update() {
	getLoadInfo()
}

// the load1Avg value
func (loadPtr *sysLoadavg) Load1Avg() float64 {
	return loadPtr.load1Avg
}

// the load5Avg value
func (loadPtr *sysLoadavg) Load5Avg() float64 {
	return loadPtr.load5Avg
}

// the load10Avg value
func (loadPtr *sysLoadavg) Load10Avg() float64 {
	return loadPtr.load10Avg
}

// the execProc value
func (loadPtr *sysLoadavg) ExecProc() uint64 {
	return loadPtr.execProc
}

// the execQueue value
func (loadPtr *sysLoadavg) ExecQueue() uint64 {
	return loadPtr.execQueue
}

// the lastPid value
func (loadPtr *sysLoadavg) LastPid() uint64 {
	return loadPtr.lastPid
}

// get cpu count
func (loadPtr *sysLoadavg) CPUs() float64 {
	return loadPtr.cpus
}

// Function to check load
func (loadPtr *sysLoadavg) CheckLoad(warn, crit string) int {
	warnThreshold, critThreshold, percent := myThreshold.SanityCheck(false, warn, crit)
	return myThreshold.CalculateValue(false, percent, warnThreshold, critThreshold,
		loadPtr.Load1Avg(), loadPtr.CPUs())
}
