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

package procfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// function to get the system memory info
func getSysMemInfo() *sysMemInfo {
	// read the proc file
	contents, err := ioutil.ReadFile(PROC_SYS_MEMINFO)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	// working variables
	var memTotal, memFree, memAvailable, buffers, cached, swapCached, swapTotal, swapFree uint64
	// prep the regex
	lineMatch, _ := regexp.Compile(sysMeminfoRegex)
	// get all lines and walk one at the time
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			// we only want those matching parRegex
			match := lineMatch.MatchString(line)
			if match {
				memName := strings.Fields(line)[0]
				memVal, _ := strconv.ParseUint(strings.Fields(line)[1], 10, 64)
				memUnit := strings.Fields(line)[2]
				// in case we reading in Mb or Gb, we need Kb
				switch memUnit {
				case "Mb":
					memVal = memVal * uint64(1024)
				case "Gb":
					memVal = memVal * uint64(1024*1024)
				}
				switch memName {
				case "MemTotal:":
					memTotal = memVal
				case "MemFree:":
					memFree = memVal
				case "MemAvailable":
					memAvailable = memVal
				case "Buffers:":
					buffers = memVal
				case "Cached":
					cached = memVal
				case "SwapCached:":
					swapCached = memVal
				case "SwapTotal:":
					swapTotal = memVal
				case "SwapFree:":
					swapFree = memVal
				}
			}
		}
	}
	sysMemValues := &sysMemInfo{
		memTotal:     memTotal,
		memFree:      memFree,
		memAvailable: memAvailable,
		buffers:      buffers,
		cached:       cached,
		swapCached:   swapCached,
		swapTotal:    swapTotal,
		swapFree:     swapFree,
	}
	return sysMemValues
}

// function to get the system memory info
func NewMem() *sysMemInfo {
	return getSysMemInfo()
}

// function to update the system memory info
func (sysMemPtr *sysMemInfo) Update() {
	// read the proc file
	contents, err := ioutil.ReadFile(PROC_SYS_MEMINFO)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	// prep the regex
	lineMatch, _ := regexp.Compile(sysMeminfoRegex)
	// get all lines and walk one at the time
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			// we only want those matching parRegex
			match := lineMatch.MatchString(line)
			if match {
				memName := strings.Fields(line)[0]
				memVal, _ := strconv.ParseUint(strings.Fields(line)[1], 10, 64)
				memUnit := strings.Fields(line)[2]
				// in case we reading in Mb or Gb, we need Kb
				switch memUnit {
				case "Mb":
					memVal = memVal * uint64(1024)
				case "Gb":
					memVal = memVal * uint64(1024*1024)
				}
				switch memName {
				case "MemTotal:":
					sysMemPtr.memTotal = memVal
				case "MemFree:":
					sysMemPtr.memFree = memVal
				case "MemAvailable":
					sysMemPtr.memAvailable = memVal
				case "Buffers:":
					sysMemPtr.buffers = memVal
				case "Cached":
					sysMemPtr.cached = memVal
				case "SwapCached:":
					sysMemPtr.swapCached = memVal
				case "SwapTotal:":
					sysMemPtr.swapTotal = memVal
				case "SwapFree:":
					sysMemPtr.swapFree = memVal
				}
			}
		}
	}
}

// functions to get the system memory type current value
func (sysMemPtr *sysMemInfo) Total() uint64 {
	return sysMemPtr.memTotal
}

func (sysMemPtr *sysMemInfo) Free() uint64 {
	return sysMemPtr.memFree
}

func (sysMemPtr *sysMemInfo) Available() uint64 {
	return sysMemPtr.memAvailable
}

func (sysMemPtr *sysMemInfo) Buffers() uint64 {
	return sysMemPtr.buffers
}

func (sysMemPtr *sysMemInfo) Cached() uint64 {
	return sysMemPtr.cached
}

func (sysMemPtr *sysMemInfo) CachedSwaped() uint64 {
	return sysMemPtr.swapCached
}

func (sysMemPtr *sysMemInfo) Swap() uint64 {
	return sysMemPtr.swapTotal
}

func (sysMemPtr *sysMemInfo) FreeSwap() uint64 {
	return sysMemPtr.swapFree
}

// calculate the real fee == Free + Cached
func (sysMemPtr *sysMemInfo) RealFree() uint64 {
	return sysMemPtr.Free() + sysMemPtr.Cached()
}

// calculate the real usage == Total - RealFree
func (sysMemPtr *sysMemInfo) RealUsage() uint64 {
	return sysMemPtr.Total() - sysMemPtr.RealFree()
}

// calculate the usage in percent
func (sysMemPtr *sysMemInfo) UsagePercent() int {
	return int((float64(sysMemPtr.RealUsage()) / float64(sysMemPtr.Total())) * 100)
}

// calculate the swap usage
func (sysMemPtr *sysMemInfo) SwapUsage() uint64 {
	return sysMemPtr.swapTotal - sysMemPtr.swapFree
}

// calculate the swapusage in percent
func (sysMemPtr *sysMemInfo) SwapUsagePercent() int {
	// capture devided by null if there is no swap setup
	if sysMemPtr.Swap() == 0 {
		return 0
	}
	return int((float64(sysMemPtr.SwapUsage()) / float64(sysMemPtr.Swap())) * 100)
}

// function to show current system memory
func (sysMemPtr *sysMemInfo) Show() {
	fmt.Printf("\nTotal        %s\n", strconv.FormatUint(sysMemPtr.Total(), 10))
	fmt.Printf("%sFree         %s\n", strconv.FormatUint(sysMemPtr.Free(), 10))
	fmt.Printf("%sAvailable    %s\n", strconv.FormatUint(sysMemPtr.Available(), 10))
	fmt.Printf("%sBuffers      %s\n", strconv.FormatUint(sysMemPtr.Buffers(), 10))
	fmt.Printf("%sCached       %s\n", strconv.FormatUint(sysMemPtr.Cached(), 10))
	fmt.Printf("%sSwapCached   %s\n", strconv.FormatUint(sysMemPtr.CachedSwaped(), 10))
	fmt.Printf("%sTotalSwap    %s\n", strconv.FormatUint(sysMemPtr.Swap(), 10))
	fmt.Printf("%sFreeSwap     %s\n", strconv.FormatUint(sysMemPtr.FreeSwap(), 10))
}
