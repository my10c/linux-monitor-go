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
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// function to get a process memory usage info
func getProcMemInfo(pidStr, pidComm string, pidInt uint64) (*procMem, error) {
	// working variables
	var err error
	var rss, pss, shared, private, swap uint64
	// read the smap file
	contents, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/smaps", PROC_DIR, pidStr))
	if err != nil {
		return nil, err
	}
	// prep the regex
	expKeys, _ := regexp.Compile(smapsRegex)
	// read the lines and add values
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		match := expKeys.MatchString(line)
		if match {
			keyName := strings.Fields(line)[0]
			keyValue, _ := strconv.ParseUint(strings.Fields(line)[1], 10, 64)
			memUnit := strings.Fields(line)[2]
			// stricky part, is should be Kb but just in case we reading in Mb or Gb, we need Kb
			switch memUnit {
			case "Mb":
				keyValue = keyValue * uint64(1024)
			case "Gb":
				keyValue = keyValue * uint64(1024*1024)
			}
			switch keyName {
			case "Rss:":
				rss = rss + keyValue
			case "Pss:":
				pss = pss + keyValue
			case "Shared_Clean:":
				shared = shared + keyValue
			case "Shared_Dirty:":
				shared = shared + keyValue
			case "Private_Clean:":
				private = private + keyValue
			case "Private_Dirty:":
				private = private + keyValue
			case "Sawp:":
				swap = swap + keyValue
			}
		}
	}
	procMem := &procMem{
		comm:    pidStr,
		pid:     pidInt,
		rss:     rss,
		pss:     pss,
		shared:  shared,
		private: private,
		swap:    swap,
	}
	return procMem, nil
}

// function to get memory usage of all processes
func getAllProcMemInfo() map[string]*procMem {
	// create the map
	allProcessMemInfo := make(map[string]*procMem)
	procPids, _ := ioutil.ReadDir(PROC_DIR)
	for _, f := range procPids {
		// make sure we its a directory
		if f.IsDir() {
			pidStr := f.Name()
			// the directory has to be an int and greater then 300
			// the info about the "RESERVED_PIDS" with default value of 300 can be found in kernel/pid.c
			if pidInt, err := strconv.ParseUint(pidStr, 10, 64); err == nil {
				if pidInt > 300 {
					// get the process name, if fail we ignore since the process could have finished (race)
					if pidCommFile, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/comm", PROC_DIR, pidStr)); err == nil {
						lines := strings.Split(string(pidCommFile), "\n")
						for _, pidComm := range lines {
							if len(pidComm) > 0 {
								// any process name with '/' we skip
								if !strings.ContainsAny(pidComm, "/") {
									// get the process meminfo
									processMemInfo, err := getProcMemInfo(pidStr, pidComm, pidInt)
									if err == nil {
										allProcessMemInfo[pidComm] = processMemInfo
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

// function to get all the processes memory info
func NewMemProc() *allProcMem {
	allInfo := &allProcMem{
		procs: getAllProcMemInfo(),
	}
	return allInfo
}

// function to get the process commmand name
func (memPtr *procMem) Comm() string {
	return memPtr.comm
}

// function to get the RSS usage
func (memPtr *procMem) Rss() uint64 {
	return memPtr.rss
}

// function to get the PSS usage
func (memPtr *procMem) Pss() uint64 {
	return memPtr.pss
}

// function to get the SHARED usage
func (memPtr *procMem) Shared() uint64 {
	return memPtr.shared
}

// function to get the PRIVATE usage
func (memPtr *procMem) Private() uint64 {
	return memPtr.private
}

// function to get the SWAP usage
func (memPtr *procMem) Swap() uint64 {
	return memPtr.swap
}

// function to get a process's memory type usage
func (memPtr *procMem) GetVal(memType string) (uint64, error) {
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
	return 0, err
}

// function to get the top memory usage by type limit by the given count
func (procPtr *allProcMem) GetTop(count int, memType string) string {
	var workList []*procMem
	var topList string
	for _, val := range procPtr.procs {
		workList = append(workList, val)
	}
	// requires Go 1.8+
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
		topList = fmt.Sprintf("%s %s:%v ", topList, workList[cnt].Comm(), val)
	}
	return strings.TrimSpace(topList)
}
