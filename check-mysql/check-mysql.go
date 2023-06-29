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

	myInit		"github.com/my10c/linux-monitor-go/initialize"
	myUtils		"github.com/my10c/linux-monitor-go/utils"
	myMySQL		"github.com/my10c/linux-monitor-go/mysql"
	myGlobal	"github.com/my10c/linux-monitor-go/global"
	myThreshold	"github.com/my10c/linux-monitor-go/threshold"
	myAlert		"github.com/my10c/linux-monitor-go/alert"
	//myStats		"github.com/my10c/nagios-plugins-go/stats"
)

const (
	table = "MONITOR"
	field = "timestamp"
	extraInfo = "Requires the table to have a field named `timestamp` and format `varchar(128)`"
	CheckVersion = "0.2"
)

var (
	cfgRequired = []string{"username", "password", "database", "hostname", "port"}
	err error
	exitVal int = 0
)

func wrongMode(modeSelect string) {
	fmt.Printf("%s", myGlobal.MyInfo)
	if modeSelect == "help" {
		fmt.Printf("Supported modes\n")
	} else {
		fmt.Printf("Wrong mode, supported modes:\n")
	}
	fmt.Printf("\t basic       : checks select/insert/delete.\n")
	fmt.Printf("\t readonly    : checks select.\n")
	fmt.Printf("\t slavestatus : checks if slave is running.\n")
	fmt.Printf("\t slavelag    : checks slave lag, requires the configs: `lagwarning` and `lagcritical`.\n")
	fmt.Printf("\t process     : checks the processes count, requires the configs: `processwarning` and `processcritical`.\n")
	fmt.Printf("\t dropcreate  : checks drop and create tables, requires the config: `tablename`.\n")
	fmt.Printf("\t showconfig  : show the current configuration and then exit.\n")
	os.Exit(3)
}

func main() {
	// need to be root since the config file wil have passwords
	myUtils.IsRoot()
	var thresHold string = ""
	var exitMsg string
	// add the extra setup info
	myGlobal.ExtraInfo = extraInfo
	myGlobal.MyVersion = CheckVersion
	cfgFile, checkMode := myInit.InitArgs(cfgRequired)
	switch checkMode {
		case "slavelag":
			cfgRequired = append(cfgRequired, "lagwarning")
			cfgRequired = append(cfgRequired, "lagcritical")
		case "process":
			cfgRequired = append(cfgRequired, "processwarning" )
			cfgRequired = append(cfgRequired, "processcritical" )
		case "dropcreate":
			cfgRequired = append(cfgRequired, "tablename" )
	}
	cfgDict := myInit.InitConfig(cfgRequired, cfgFile)
	myInit.InitLog()
	myUtils.SignalHandler()
	dbCheck := myMySQL.New(cfgDict)
	//--> stats := myStats.New()
	data := time.Now().Format(time.RFC3339)
	iter, _ := strconv.Atoi(cfgDict["iter"])
	iterWait, _ := time.ParseDuration(cfgDict["iterwait"])
	for cnt :=0 ; cnt < iter ; cnt++ {
		switch checkMode {
			case "basic":
				exitVal, err = dbCheck.BasisCheck(table, field, data)
			case "readonly":
				exitVal, err = dbCheck.ReadCheck(table, field)
			case "slavestatus":
				exitVal, err = dbCheck.SlaveStatusCheck()
			case "slavelag":
				warning, critical, _ := myThreshold.SanityCheck(false, cfgDict["lagwarning"], cfgDict["lagcritical"])
				exitVal, err = dbCheck.SlaveLagCheck(warning, critical)
				thresHold = fmt.Sprintf(" (W:%d C:%d )", warning, critical)
			case "process":
				warning, critical, _ := myThreshold.SanityCheck(false, cfgDict["processwarning"], cfgDict["processcritical"])
				exitVal, err = dbCheck.ProcessStatusCheck(warning, critical)
				thresHold = fmt.Sprintf(" (W:%d C:%d )", warning, critical)
			case "dropcreate":
				exitVal, err = dbCheck.DropCreateCheck(cfgDict["tablename"])
			case "showconfig":
				myUtils.ShowMap(cfgDict)
				myUtils.ShowMap(nil)
				os.Exit(0)
			default:
				wrongMode(checkMode)
		}
		// if we get an OK then exit no need to do all iterations
		if exitVal == myGlobal.OK {
			break
		}
		time.Sleep(iterWait * time.Second)
	}
	//
	// TODO write stats here
	// We only need 1 entry in the stats instead of all iteration like other check need
	//
	if exitVal != myGlobal.OK {
		if myGlobal.DefaultValues["noalert"] == "false" {
			myAlert.SendAlert(exitVal, checkMode, err.Error())
		}
		exitMsg = fmt.Sprintf("%s %s - Check running mode: %s - Error: %s %s\n",
			strings.ToUpper(myGlobal.MyProgname), myGlobal.Result[exitVal], checkMode, err.Error(), thresHold)
	} else {
		exitMsg = fmt.Sprintf("%s %s - Check running mode: %s - %s %s \n",
		strings.ToUpper(myGlobal.MyProgname), myGlobal.Result[exitVal], checkMode, err, thresHold)
	}
	fmt.Printf("%s", exitMsg)
	myUtils.LogMsg(fmt.Sprintf("%s", exitMsg))
	os.Exit(exitVal)
}
