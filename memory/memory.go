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

package memory

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	myGlobal "github.com/my10c/linux-monitor-go/global"
	myThreshold "github.com/my10c/linux-monitor-go/threshold"
	myUtils "github.com/my10c/linux-monitor-go/utils"
)

const (
	PROCDIR    = "/proc"
	PROCMEM    = "/proc/meminfo"
	PROCESSCOM = "comm"
)

type sysMemStruct struct {
	memTotal     uint64 `json:"memTotal"`
	memFree      uint64 `json:"memFree"`
	memAvailable uint64 `json:"memAvailable"`
	buffers      uint64 `json:"buffers"`
	cached       uint64 `json:"cached"`
	swapCached   uint64 `json:"swapcached"`
	swapTotal    uint64 `json:"swapTotal"`
	swapFree     uint64 `json:"swapTotal"`
}

// Rss: resident memory usage, all memory the process uses,
//		including all memory this process shares with other processes. It does not include swap;
// Shared: memory that this process shares with other processes;
// Private: private memory used by this process, you can look for memory leaks here;
// Swap: swap memory used by the process;
// Pss: Proportional Set Size, a good overall memory indicator.
//		It is the Rss adjusted for sharing: if a process has 1MiB private and 20MiB shared
//		between other 10 processes, Pss is 1 + 20/10 = 3MiB

type processMemStruct struct {
	processName string `json:"procname"`
	rss         uint64 `json:"rss"`
	pss         uint64 `json:"pss"`
	shared      uint64 `json:"shared"`
	private     uint64 `json:"private"`
	swap        uint64 `json:"swap"`
}

var (
	memRegex     = `^(MemTotal|MemFree|MemAvailable|Buffers|Cached|SwapCached|SwapTotal|SwapFree)`
	procMemRegex = `^(Rss:|Pss:|Shared_Clean:|Shared_Dirty:|Private_Clean:|Private_Dirty:|Swap:)`
)

// function to get the system memory info
func getSysMemInfo(unit uint64) *sysMemStruct {
	// working variable
	var memTotal uint64 = 0
	var memFree uint64 = 0
	var memAvailable uint64 = 0
	var buffers uint64 = 0
	var cached uint64 = 0
	var swapCached uint64 = 0
	var swapTotal uint64 = 0
	var swapFree uint64 = 0
	// read the proc file
	contents, err := ioutil.ReadFile(PROCMEM)
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// prep the regex
	lineMatch, _ := regexp.Compile(memRegex)
	// get all lines and walk one at the time
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			// we only want those matching parRegex
			match := lineMatch.MatchString(line)
			if match {
				memName := myUtils.TrimLastChar(strings.Fields(line)[0], ":")
				memVal, _ := strconv.ParseUint(strings.Fields(line)[1], 10, 64)
				memUnit := strings.Fields(line)[2]
				// in case we reading in Mb or Gb, we need Kb
				switch memUnit {
				case "Mb":
					memVal = memVal * uint64(1024)
				case "Gb":
					memVal = memVal * uint64(1024*1024)
				}
				memVal = uint64(float64(memVal) / float64(unit))
				switch memName {
				case "MemTotal":
					memTotal = memVal
				case "MemFree":
					memFree = memVal
				case "MemAvailable":
					memAvailable = memVal
				case "Buffers":
					buffers = memVal
				case "Cached":
					cached = memVal
				case "SwapCached":
					swapCached = memVal
				case "SwapTotal":
					swapTotal = memVal
				case "SwapFree":
					swapFree = memVal
				}
			}
		}
	}
	sysMemValues := &sysMemStruct{
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

// Function to get a process memory usage info
func getPidMemInfo(processPid, name string, unit uint64) (*processMemStruct, error) {
	var err error
	var rss uint64 = 0
	var pss uint64 = 0
	var shared uint64 = 0
	var private uint64 = 0
	var swap uint64 = 0
	// read the smap file
	contents, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/smaps", PROCDIR, processPid))
	if err != nil {
		return nil, err
	}
	// prep the regex
	expKeys, _ := regexp.Compile(procMemRegex)
	// read the lines and add values
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		match := expKeys.MatchString(line)
		if match {
			keyName := myUtils.TrimLastChar(strings.Fields(line)[0], ":")
			keyValue, _ := strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			switch keyName {
			case "Rss":
				rss = rss + keyValue
			case "Pss":
				pss = pss + keyValue
			case "Shared_Clean":
				shared = shared + keyValue
			case "Shared_Dirty":
				shared = shared + keyValue
			case "Private_Clean":
				private = private + keyValue
			case "Private_Dirty":
				private = private + keyValue
			case "Swap":
				swap = swap + keyValue
			}
		}
	}
	procMem := &processMemStruct{
		processName: name,
		rss:         uint64(float64(rss) / float64(unit)),
		pss:         uint64(float64(pss) / float64(unit)),
		shared:      uint64(float64(shared) / float64(unit)),
		private:     uint64(float64(private) / float64(unit)),
		swap:        uint64(float64(swap) / float64(unit)),
	}
	return procMem, nil
}

// Function to get memory usage of the given processes
func getProcessesMemInfo(unit uint64) map[string]*processMemStruct {
	// create the map
	allProcessMemInfo := make(map[string]*processMemStruct)
	procIDs, _ := ioutil.ReadDir(PROCDIR)
	for _, f := range procIDs {
		// make sure we its a directory
		if f.IsDir() {
			baseName := f.Name()
			// the directory has to be an int and greater then 300
			// the info about the "RESERVED_PIDS" with default value of 300 can be found in kernel/pid.c
			if pidVal, err := strconv.Atoi(baseName); err == nil {
				if pidVal > 300 {
					// get the process name,
					if pidCommFile, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/comm", PROCDIR, baseName)); err == nil {
						commLines := strings.Split(string(pidCommFile), "\n")
						for _, commName := range commLines {
							if len(commName) > 0 {
								// any process name with '/' we skip since these are memory save
								if !strings.ContainsAny(commName, "/") {
									// get the process meminfo
									processMemInfo, err := getPidMemInfo(baseName, commName, unit)
									if err == nil {
										allProcessMemInfo[commName] = processMemInfo
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return allProcessMemInfo
}

// Function to get the system and processes memory info
func New(unit uint64) *sysMemStruct {
	return getSysMemInfo(unit)
}

// Function to get the process commmand name
func (memPtr *processMemStruct) Name() string {
	return memPtr.processName
}

// Function to get the RSS usage
func (memPtr *processMemStruct) Rss() uint64 {
	return memPtr.rss
}

// Function to get the PSS usage
func (memPtr *processMemStruct) Pss() uint64 {
	return memPtr.pss
}

// Function to get the SHARED usage
func (memPtr *processMemStruct) Shared() uint64 {
	return memPtr.shared
}

// Function to get the PRIVATE usage
func (memPtr *processMemStruct) Private() uint64 {
	return memPtr.private
}

// Function to get the SWAP usage
func (memPtr *processMemStruct) Swap() uint64 {
	return memPtr.swap
}

// Function to get a process's memory type usage
func (memPtr *processMemStruct) GetVal(memType string) (uint64, error) {
	switch memType {
	case "rss", "memory":
		return memPtr.rss, nil
	case "pss":
		return memPtr.pss, nil
	case "shared":
		return memPtr.shared, nil
	case "private":
		return memPtr.private, nil
	case "swap":
		return memPtr.swap, nil
	}
	err := fmt.Errorf("memType not supported: %s", memType)
	return 666, err
}

// function to get the top memory usage by type limit by thge given count
func GetTop(count int, memType string, unit uint64) string {
	var workList []*processMemStruct
	var topList string
	// generate the process list
	allProcs := getProcessesMemInfo(unit)
	for _, val := range allProcs {
		workList = append(workList, val)
	}

	sort.Slice(workList, func(i, j int) bool {
		switch memType {
		case "rss", "memory":
			return workList[i].Rss() > workList[j].Rss()
		case "pss":
			return workList[i].Pss() > workList[j].Pss()
		case "shared":
			return workList[i].Shared() > workList[j].Shared()
		case "private":
			return workList[i].Private() > workList[j].Private()
		case "swap":
			return workList[i].Swap() > workList[j].Swap()
		}
		return false
	})

	for cnt := 0; cnt < count; cnt++ {
		val, _ := workList[cnt].GetVal(memType)
		topList = fmt.Sprintf("%s %s:%v ", topList, workList[cnt].Name(), val)
	}
	return strings.TrimSpace(topList)
}

// Functions to get the system memory type current value
func (sysMemPtr *sysMemStruct) Total() uint64 {
	return sysMemPtr.memTotal
}

func (sysMemPtr *sysMemStruct) Free() uint64 {
	return sysMemPtr.memFree
}

func (sysMemPtr *sysMemStruct) Available() uint64 {
	return sysMemPtr.memAvailable
}

func (sysMemPtr *sysMemStruct) Buffers() uint64 {
	return sysMemPtr.buffers
}

func (sysMemPtr *sysMemStruct) Cached() uint64 {
	return sysMemPtr.cached
}

func (sysMemPtr *sysMemStruct) CachedSwaped() uint64 {
	return sysMemPtr.swapCached
}

func (sysMemPtr *sysMemStruct) Swap() uint64 {
	return sysMemPtr.swapTotal
}

func (sysMemPtr *sysMemStruct) FreeSwap() uint64 {
	return sysMemPtr.swapFree
}

func (sysMemPtr *sysMemStruct) SwapUsage() uint64 {
	return sysMemPtr.swapTotal - sysMemPtr.swapFree
}

// function to show current system memory
func (sysMemPtr *sysMemStruct) Show() string {
	var sysMemInfo string
	sysMemInfo = fmt.Sprintf("\nTotal        %s\n", strconv.FormatUint(sysMemPtr.Total(), 10))
	sysMemInfo = fmt.Sprintf("%sFree         %s\n", sysMemInfo, strconv.FormatUint(sysMemPtr.Free(), 10))
	sysMemInfo = fmt.Sprintf("%sAvailable    %s\n", sysMemInfo, strconv.FormatUint(sysMemPtr.Available(), 10))
	sysMemInfo = fmt.Sprintf("%sBuffers      %s\n", sysMemInfo, strconv.FormatUint(sysMemPtr.Buffers(), 10))
	sysMemInfo = fmt.Sprintf("%sCached       %s\n", sysMemInfo, strconv.FormatUint(sysMemPtr.Cached(), 10))
	sysMemInfo = fmt.Sprintf("%sSwapCached   %s\n", sysMemInfo, strconv.FormatUint(sysMemPtr.CachedSwaped(), 10))
	sysMemInfo = fmt.Sprintf("%sTotalSwap    %s\n", sysMemInfo, strconv.FormatUint(sysMemPtr.Swap(), 10))
	sysMemInfo = fmt.Sprintf("%sFreeSwap     %s\n", sysMemInfo, strconv.FormatUint(sysMemPtr.FreeSwap(), 10))
	return sysMemInfo
}

// Calculate the real fee == Free + Cached
func (sysMemPtr *sysMemStruct) RealFree() uint64 {
	return sysMemPtr.Free() + sysMemPtr.Cached()
}

// Calculate the real usage == Total - RealFree
func (sysMemPtr *sysMemStruct) RealUsage() uint64 {
	return sysMemPtr.Total() - sysMemPtr.RealFree()
}

// Function to check available memory
func (sysMemPtr *sysMemStruct) CheckFreeMem(warn, crit string) int {
	warnThreshold, critThreshold, percent := myThreshold.SanityCheck(false, warn, crit)
	return myThreshold.CalculateUsage(false, percent, warnThreshold, critThreshold,
		sysMemPtr.RealUsage(), sysMemPtr.Total())
}

// Function to check available swap
func (sysMemPtr *sysMemStruct) CheckFreeSwap(warn, crit string) int {
	warnThreshold, critThreshold, percent := myThreshold.SanityCheck(false, warn, crit)
	return myThreshold.CalculateUsage(false, percent, warnThreshold, critThreshold,
		sysMemPtr.SwapUsage(), sysMemPtr.Swap())
}

// Function to get the usage in percent
func (sysMemPtr *sysMemStruct) UsagePercent() int {
	return int((float64(sysMemPtr.RealUsage()) / float64(sysMemPtr.Total())) * 100)
}

// Function to get the swapusage in percent
func (sysMemPtr *sysMemStruct) SwapUsagePercent() int {
	// capture devided by nul id there is no swap setup
	if sysMemPtr.Swap() == 0 {
		return 0
	}
	return int((float64(sysMemPtr.SwapUsage()) / float64(sysMemPtr.Swap())) * 100)
}
