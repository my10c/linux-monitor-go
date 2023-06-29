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

package process

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	// myUtils "github.com/my10c/linux-monitor-go/utils"
)

const (
	PROCDIR	= "/proc"
	PROCESSCOM = "comm"
	PROCESSCMDLINE = "cmdline"
)

type procInfo struct {
	comm string `json:"command"`
	cmdline string `json:"commline"`
	count int `jason:"count"`
}

// Function to get a process memory usage info
func getPidProcessInfo(processPid, cmd string) (*procInfo, error) {
	var cmdline string
	contents, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/cmdline", PROCDIR, processPid))
	if err != nil {
		return nil, err
	}
	// read the lines and add values, there should be only one line but we will need to be sure
	lines := strings.Split(string(contents), "\n")
	for _, cmdline = range lines {;}
	procInfo := &procInfo{
		comm: cmd,
		cmdline: cmdline,
		count: 1,
	}
	return procInfo, nil
}

// Function to get all process command and their command-line
func getProcessesInfo() map[string]*procInfo {
	// create the map
	allProcessInfo := make(map[string]*procInfo)
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
								// any process name with '/' we skip
								if !strings.ContainsAny(commName, "/") {
									// get the process meminfo
									processInfo, err := getPidProcessInfo(baseName, commName)
									if err == nil {
										// need to count if the same process already was found
										// could be a child, parent or different cmdline!
										if _, ok := allProcessInfo[commName]; ok {
												// overwrite the count TODO: what if cmdline is different!
												allProcessInfo[commName].count = allProcessInfo[commName].count + 1
										} else {
											allProcessInfo[commName] = processInfo
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return allProcessInfo
}
