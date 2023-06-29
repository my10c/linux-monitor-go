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
// Date			:	Jul 14, 2017
//
// History	:
// 	Date:			Author:		Info:
//	June 30, 2017	LIS			First Go release
//	Jul 14, 2017	LIS			Added stats
//

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	myAlert		"github.com/my10c/linux-monitor-go/alert"
	myDisk		"github.com/my10c/linux-monitor-go/disk"
	myGlobal	"github.com/my10c/linux-monitor-go/global"
	myInit		"github.com/my10c/linux-monitor-go/initialize"
	myUtils		"github.com/my10c/linux-monitor-go/utils"
	myStats		"github.com/my10c/linux-monitor-go/stats"
)

const (
	field        = "timestamp"
	extraInfo    = "Requires the warning and critical thresholds\n\t\tEmpty unit defaults to MB and empty disk defaults to all disks"
	CheckVersion = "0.1"
)

var (
	cfgRequired = []string{"critical", "warning", "disk", "unit"}
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
	fmt.Printf("\t diskspace : checks diskspace.\n")
	fmt.Printf("\t inodes    : checks inodes.\n")
	fmt.Printf("\t ro        : checks partition is not read-only mode.\n")
	os.Exit(3)
}

func wrongUnit(confUnit string) {
	fmt.Printf("%s", myGlobal.MyInfo)
	fmt.Printf("Wrong unit %s, supported unit:\n", confUnit)
	fmt.Printf("\t KB	: KiloBytes, most accurate.\n")
	fmt.Printf("\t MB	: MegaBytes, good accuracy.\n")
	fmt.Printf("\t GB	: GigaBytes, less accurate.\n")
	fmt.Printf("\t TB	: TerraBytes, worst accuracy.\n")
	os.Exit(3)
}

func checkUnit(unit string) uint64 {
	var unitBytes uint64
	switch unit {
	case "":
		unitBytes = myGlobal.MB
	case "KB":
		unitBytes = myGlobal.KB
	case "MB":
		unitBytes = myGlobal.MB
	case "GB":
		unitBytes = myGlobal.GB
	case "TB":
		unitBytes = myGlobal.TB
	default:
		wrongUnit(unit)
	}
	return unitBytes
}

func checkMode(givenMode string) {
	switch givenMode {
	case "diskspace":
	case "inode":
	case "ro":
	default:
		wrongMode(givenMode)
	}
}

func main() {
	// working variables
	var resultVal int
	var exitVal int = 0
	var exitMsg string
	var currMount string
	var currStats string
	var stats *myStats.Stats
	// use to stop stats creation in case not check all disks
	var doStats bool = true
	// create emtpy error message
	err = fmt.Errorf("")
	// need to be root since the config file wil have passwords
	myUtils.IsRoot()
	// get and setup phase
	myUtils.IsLinuxSystem()
	myGlobal.ExtraInfo = extraInfo
	myGlobal.MyVersion = CheckVersion
	cfgFile, givenMode := myInit.InitArgs(cfgRequired)
	cfgDict := myInit.InitConfig(cfgRequired, cfgFile)
	myInit.InitLog()
	myUtils.SignalHandler()
	givenUnit := checkUnit(cfgDict["unit"])
	checkMode(givenMode)
	// we do onlt stats on diskspace!
	if  myGlobal.DefaultValues["stats"] == "true" && givenMode == "diskspace" {
		stats = myStats.New(cfgDict["statstid"], cfgDict["statstformat"])
		currStats = fmt.Sprintf("\"partitions\": {")
	}
	thresHold := fmt.Sprintf(" (W:%s C:%s Unit:%s)", cfgDict["warning"], cfgDict["critical"], cfgDict["unit"])
	iter, _ := strconv.Atoi(cfgDict["iter"])
	iterWait, _ := time.ParseDuration(cfgDict["iterwait"])
	// loop all found disk
	for mountPoint, diskPtr := range myDisk.New() {
		// loop times required iterations if errored
		for cnt := 0; cnt < iter; cnt++ {
			if len(cfgDict["disk"]) == 0 {
				// need to do all partitions
				resultVal = diskPtr.CheckIt(givenMode, cfgDict["warning"], cfgDict["critical"], givenUnit)
			} else {
				doStats = false
				// if disk is set we stop as soon we have a hit, can be mountpoint or device name
				if (diskPtr.GetDev() == cfgDict["disk"]) || (mountPoint == cfgDict["disk"]) {
					resultVal = diskPtr.CheckIt(givenMode, cfgDict["warning"], cfgDict["critical"], givenUnit)
					doStats = true
				}
			}
			// create the stats record
			if myGlobal.DefaultValues["stats"] == "true" && givenMode == "diskspace" && doStats {
				currMount = diskPtr.GetMountPoint()
				if  currMount == "/" {
					currMount = "root"
				}
				currStats = fmt.Sprintf("%s\"%v\": {\"total\": %v, \"free\": %v, \"use\": %v}", currStats,
					currMount, diskPtr.GetSize(givenUnit), diskPtr.GetFree(givenUnit), diskPtr.GetUse(givenUnit))
				currStats = fmt.Sprintf("%s, ", currStats)
			}
			// got OK, break and go to next partition
			if resultVal == myGlobal.OK {
				break
			} else {
				// we set the value of exitVal to to the highest just once
				// so that we get critical if any result is critical
				// BUG: since test it done `iter` time we should reset if we get an ok
				if exitVal < resultVal {
					exitVal = resultVal
				}
			}
			time.Sleep(iterWait * time.Second)
		}
		// break if not checking all disks
		if (diskPtr.GetDev() == cfgDict["disk"]) || (mountPoint == cfgDict["disk"]) {
			// create the disk message only for the disk
			exitMsg = fmt.Sprintf("%s%s ",
				myGlobal.Result[resultVal], diskPtr.StatusMsg(givenMode, givenUnit))
			if resultVal != myGlobal.OK {
				err = fmt.Errorf("%s%s%s ",
					err.Error(), myGlobal.Result[resultVal], diskPtr.StatusMsg(givenMode, givenUnit))
			}
			break
		}
		// create the disk message appended
		exitMsg = fmt.Sprintf("%s%s%s ",
			exitMsg, myGlobal.Result[resultVal], diskPtr.StatusMsg(givenMode, givenUnit))
		if resultVal != myGlobal.OK {
			err = fmt.Errorf("%s%s%s ",
				err.Error(), myGlobal.Result[resultVal], diskPtr.StatusMsg(givenMode, givenUnit))
		}
	}
	// create final message
	if givenMode != "ro" {
		// add the check name and threshold info to the message
		exitMsg = fmt.Sprintf("%s %s - %s%s\n",
			strings.ToUpper(myGlobal.MyProgname), myGlobal.Result[exitVal], exitMsg, thresHold)
		err = fmt.Errorf("%s%s", err.Error(), thresHold)
	} else {
		// add only the check name
		exitMsg = fmt.Sprintf("%s %s - %s\n",
			strings.ToUpper(myGlobal.MyProgname), myGlobal.Result[exitVal], exitMsg)
		err = fmt.Errorf("%s", err.Error())
	}
	// alert on error
	if exitVal != myGlobal.OK {
		if myGlobal.DefaultValues["noalert"] == "false" {
			// add threshold to error message
			myAlert.SendAlert(exitVal, givenMode, err.Error())
		}
	}
	fmt.Printf("%s", exitMsg)
	myUtils.LogMsg(fmt.Sprintf("%s", exitMsg))
	// write the stat after final record cleanup
	if myGlobal.DefaultValues["stats"] == "true" && givenMode == "diskspace" {
		currStats = fmt.Sprintf("%s}", myUtils.TrimLastChar(currStats, ", "))
		err := stats.Stats(currStats)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
	}
	os.Exit(exitVal)
}
