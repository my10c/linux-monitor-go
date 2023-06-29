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

package main

import (
	"fmt"
	"os"
	//"strconv"
	"strings"
	//"time"

	myDisk "github.com/my10c/linux-monitor-go/disk"
	myGlobal "github.com/my10c/linux-monitor-go/global"
	myInit "github.com/my10c/linux-monitor-go/initialize"
	myUtils "github.com/my10c/linux-monitor-go/utils"
)

const (
	CheckVersion = "0.1"
)

var (
	err     error
	exitVal int = 0
)

func wrongMode(modeSelect string) {
	fmt.Printf("%s", myGlobal.MyInfo)
	if modeSelect == "help" {
		fmt.Printf("Supported modes\n")
	} else {
		fmt.Printf("Wrong mode, supported modes:\n")
	}
	fmt.Printf("\t diskspace : diskspace info\n")
	fmt.Printf("\t inodes    : inodes info.\n")
	fmt.Printf("\t both      : diskspace and inodes info.\n")
	fmt.Printf("\t status    : partition status, read-write or read-only.\n")
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
	case "both":
	case "status":
	default:
		wrongMode(givenMode)
	}
}

func main() {
	// working variables
	var exitVal int = 0
	// get and setup phase
	myUtils.IsLinuxSystem()
	// capture control-C
	myUtils.SignalHandler()
	// Set version
	myGlobal.MyVersion = CheckVersion
	// get and set configs
	_, givenMode := myInit.InitArgs(nil)
	checkMode(givenMode)
	myUnit := checkUnit(strings.ToUpper(myGlobal.DefaultValues["unit"]))
	// loop all found disk
	for mountPoint, diskPtr := range myDisk.New() {
		switch givenMode {
		case "diskspace":
			fmt.Printf("(%s) partition %s Free diskspace %d of %d total.\n",
				myGlobal.DefaultValues["unit"], mountPoint,
				diskPtr.GetFree(myUnit), diskPtr.GetSize(myUnit))
		case "inode":
			// inode unit always 1000
			fmt.Printf("(x%d) partition %s Free Inodes %d of %d total.\n",
				1000, mountPoint,
				diskPtr.GetFreeInodes(1000), diskPtr.GetInodes(1000))
		case "both":
			// inode unit always 1000
			fmt.Printf("(%s) partition %s Free diskspace %d of %d total, (x1000) Free Inodes %d of %d total.\n",
				myGlobal.DefaultValues["unit"], mountPoint,
				diskPtr.GetFree(myUnit), diskPtr.GetSize(myUnit),
				diskPtr.GetFreeInodes(1000), diskPtr.GetInodes(1000))
		case "status":
			fmt.Printf("partition %s mount state = %s.\n", mountPoint, diskPtr.GetState())
		}
	}
	os.Exit(exitVal)
}
