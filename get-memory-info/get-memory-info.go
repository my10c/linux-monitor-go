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
// Date			:	Nov 5, 2017
//
// History	:
// 	Date:			Author:		Info:
//	Nov 5, 2017		LIS			First release
//
// TODO:

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	//	"time"

	myGlobal "github.com/my10c/linux-monitor-go/global"
	myInit "github.com/my10c/linux-monitor-go/initialize"
	myMemory "github.com/my10c/linux-monitor-go/memory"
	myUtils "github.com/my10c/linux-monitor-go/utils"
)

const (
	CheckVersion = "0.1"
)

var (
	exitVal  int = 0
	dblSpace     = `[\s\p{Zs}]{2,}`
)

func wrongMode(modeSelect string) {
	fmt.Printf("%s", myGlobal.MyInfo)
	if modeSelect == "help" {
		fmt.Printf("Supported modes\n")
	} else {
		fmt.Printf("Wrong mode, supported modes:\n")
	}
	fmt.Printf("\t system      : show the current system memory status.\n")
	fmt.Printf("\t top-memory  : show top %d process memory usage.\n", myGlobal.MyTop)
	fmt.Printf("\t top-rss     : show top %d process memory usage.\n", myGlobal.MyTop)
	fmt.Printf("\t top-private : show top %d process private memory usage.\n", myGlobal.MyTop)
	fmt.Printf("\t top-swap    : show top %d process swap memory usage.\n", myGlobal.MyTop)
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
	// since everything is in KB
	var unitBytes uint64 = 1
	switch unit {
	case "", "KB":
		unitBytes = 1
	case "MB":
		unitBytes = myGlobal.KB
	case "GB":
		unitBytes = myGlobal.MB
	case "TB":
		unitBytes = myGlobal.GB
	default:
		wrongUnit(unit)
	}
	return unitBytes
}

func checkMode(givenMode string) {
	switch givenMode {
	case "system", "top-memory", "top-rss", "top-private", "top-swap":
	default:
		wrongMode(givenMode)
	}
}

func main() {
	// working variables
	var exitVal int = 0
	var exitMsg string
	// works only on Linux system
	myUtils.IsLinuxSystem()
	// capture control-C
	myUtils.SignalHandler()
	// Set version
	myGlobal.MyVersion = CheckVersion
	// get and set configs
	_, givenMode := myInit.InitArgs(nil)
	topCount, _ := strconv.Atoi(myGlobal.DefaultValues["top"])
	checkMode(givenMode)
	givenUnit := checkUnit(strings.ToUpper(myGlobal.DefaultValues["unit"]))
	regexRemove := regexp.MustCompile(dblSpace)
	// Get the memory infos
	systemMemInfo := myMemory.New(givenUnit)
	switch givenMode {
	case "system":
		exitMsg = fmt.Sprintf("\t(Unit %s)%s", myGlobal.DefaultValues["unit"], systemMemInfo.Show())
	case "top-memory", "top-rss", "top-private", "top-swap":
		// remove the top- string as that is an invalid option for the memory class
		cleanedMode := strings.Replace(givenMode, "top-", "", -1)
		// output one line per process info, todo so remove the double space then replace single space with carriage-return
		cleanedInfo := strings.Replace(regexRemove.ReplaceAllString(myMemory.GetTop(topCount, cleanedMode, givenUnit), " "), " ", "\n", -1)
		exitMsg = fmt.Sprintf("(Unit %s)\n%s", myGlobal.DefaultValues["unit"], cleanedInfo)
	default:
		wrongMode(givenMode)
	}
	exitMsg = fmt.Sprintf("%d %s - %s\n", topCount, givenMode, exitMsg)
	fmt.Printf("%s", exitMsg)
	os.Exit(exitVal)
}
