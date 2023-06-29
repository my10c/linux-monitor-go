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
// Date			:	June 4, 2017
//
// History	:
// 	Date:			Author:		Info:
//	June 4, 2017	LIS			First Go release
//
// TODO:

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	myGlobal "github.com/my10c/linux-monitor-go/global"
	myInit "github.com/my10c/linux-monitor-go/initialize"
	myLoad "github.com/my10c/linux-monitor-go/load"
	myStats "github.com/my10c/linux-monitor-go/stats"
	myUtils "github.com/my10c/linux-monitor-go/utils"
)

const (
	CheckVersion = "0.1"
)

var (
	cfgRequired = []string{"loadwarning", "loadcritical"}
	err         error
	exitVal     int = 0
)

func wrongMode(modeSelect string) {
	fmt.Printf("%s", myGlobal.MyInfo)
	if modeSelect == "help" {
		fmt.Printf("Supported modes\n")
	} else {
		fmt.Printf("Wrong mode, supported modes:\n")
	}
	fmt.Printf("\t load  : display current load.\n")
	fmt.Printf("\t queue : display curent queue.\n")
	os.Exit(3)
}

func checkMode(givenMode string) {
	switch givenMode {
	case "load", "queue":
	default:
		wrongMode(givenMode)
	}
}

func main() {
	// working variables
	var exitVal int = 0
	var exitMsg string
	// for stats
	var currStats string
	var stats *myStats.Stats
	// create emtpy error message
	err = fmt.Errorf("")
	// need to be root since the config file wil have passwords
	myUtils.IsRoot()
	// get and setup phase
	myUtils.IsLinuxSystem()
	myGlobal.MyVersion = CheckVersion
	cfgFile, givenMode := myInit.InitArgs(cfgRequired)
	cfgDict := myInit.InitConfig(cfgRequired, cfgFile)
	myInit.InitLog()
	myUtils.SignalHandler()
	checkMode(givenMode)
	if myGlobal.DefaultValues["stats"] == "true" {
		stats = myStats.New(cfgDict["statstid"], cfgDict["statstformat"])
	}
	thresHold := fmt.Sprintf(" (W:%s C:%s)", cfgDict["loadwarning"], cfgDict["loadcritical"])
	iter, _ := strconv.Atoi(cfgDict["iter"])
	iterWait, _ := time.ParseDuration(cfgDict["iterwait"])
	// Get the memory infos
	systemLoad := myLoad.New()
	if givenMode == "queue" {
		fmt.Printf("%s OK - Exec Process %v - Exec Queue %v\n",
			strings.ToUpper(myGlobal.MyProgname), systemLoad.ExecProc(), systemLoad.ExecQueue())
		os.Exit(exitVal)
	}
	for cnt := 0; cnt < iter; cnt++ {
		exitVal = systemLoad.CheckLoad(cfgDict["loadwarning"], cfgDict["loadcritical"])
		exitMsg = fmt.Sprintf("(Load 1min:%.2f 5min:%.2f 10min:%.2f)%s",
			systemLoad.Load1Avg(), systemLoad.Load5Avg(), systemLoad.Load10Avg(), thresHold)
		if myGlobal.DefaultValues["stats"] == "true" {
			currStats = fmt.Sprintf("\"cpu_load\": {\"cpu_count\": %.2f, \"load_curr\": %.2f, \"load_10_min\": %.2f, \"load_15_min\": %.2f}",
				systemLoad.CPUs(), systemLoad.Load1Avg(), systemLoad.Load5Avg(), systemLoad.Load10Avg())
			// write the stat record
			err := stats.Stats(currStats)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
		}
		// if we get an OK then exit no need to do all iterations
		if exitVal == myGlobal.OK {
			break
		}
		time.Sleep(iterWait * time.Second)
		// get new values for load
		systemLoad.Update()
	}

	// add the check name
	exitMsg = fmt.Sprintf("%s %s - %s\n",
		strings.ToUpper(myGlobal.MyProgname), myGlobal.Result[exitVal], exitMsg)
	fmt.Printf("%s", exitMsg)
	myUtils.LogMsg(fmt.Sprintf("%s", exitMsg))
	os.Exit(exitVal)
}
