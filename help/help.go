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

package help

import (
	"fmt"
	"os"
	"strings"

	myGlobal "github.com/my10c/linux-monitor-go/global"
	myUtils "github.com/my10c/linux-monitor-go/utils"

	"github.com/fatih/color"
)

// Function to print a list of configurable keys and their default values if any
func printCfgValues(sectioName string, disableKey string, cfgDict map[string]string) {
	if len(disableKey) > 0 {
		if strings.Contains(disableKey, ":") {
			color.HiYellow("\t# to disable set `%s`, if shown empty, then its disable by default.\n", disableKey)
		} else {
			color.HiYellow("\t# to disable set an empty `%s`, if shown empty, then its disable by default.\n", disableKey)
		}
	}
	fmt.Printf("\t%s:\n", sectioName)
	for defaultKey, defaultValue := range cfgDict {
		if checkDefaultValue, ok := myGlobal.SharedMap[defaultKey]; ok {
			defaultValue = checkDefaultValue
		}
		fmt.Printf("\t  %s: %s\n", defaultKey, defaultValue)
	}
}

// Function to show how to setup the check
func SetupHelp(cfg []string) {
	fmt.Printf("%s", myGlobal.MyInfo)
	fmt.Printf("Setup the configuration file:\n")
	fmt.Printf("# Create a configuration file, any name would do, as long its in yaml fornmat.\n")
	fmt.Printf("# Default to %s\n", myGlobal.DefaultConfigFile)
	fmt.Printf("# Add the following key/pair values, these are required:")
	color.Yellow(" (used the `-mode help` to see other required key)\n")
	fmt.Printf("\t%s:\n", myGlobal.MyProgname)
	for cnt := range cfg {
		if defaultValue, ok := myGlobal.SharedMap[cfg[cnt]]; ok {
			fmt.Printf("\t  %s: %s\n", cfg[cnt], defaultValue)
		} else {
			fmt.Printf("\t  %s:\n", cfg[cnt])
		}
	}
	// stats are hardcoded
	color.HiGreen("# When enable stats (-stats) then the keys below are required, below are the default values.\n")
	fmt.Printf("\t  statsdir: %s.stat\n", myGlobal.DefaultLog["logdir"])
	fmt.Printf("\t  statsfile: %s.stat\n", myGlobal.MyProgname)
	// iter are also hardcoded
	color.HiGreen("# How many iteration to perform and how much to wait between (seconds) them before its an issue.\n")
	// print the sections configs
	fmt.Printf("\t  iter: %s\n", myGlobal.DefaultIter["iter"])
	fmt.Printf("\t  iterwait: %s\n", myGlobal.DefaultIter["iterwait"])
	fmt.Printf("# Values shown are the default values. If a section is ommitted, then it will use the default values.\n")
	printCfgValues("common", "", myGlobal.DefaultValues)
	printCfgValues("log", "logfile", myGlobal.DefaultLog)
	printCfgValues("email", "emailto", myGlobal.DefaultEmail)
	printCfgValues("tag", "tagfile", myGlobal.DefaultTag)
	printCfgValues("syslog", "syslogtag: off", myGlobal.DefaultSyslog)
	printCfgValues("pagerduty", "pdservicekey", myGlobal.DefaultPD)
	printCfgValues("slack", "slackservicekey", myGlobal.DefaultSlack)
	fmt.Printf("\nNOTE\n")
	if len(myGlobal.ExtraInfo) > 0 {
		color.HiRed("\t* %s\n", myGlobal.ExtraInfo)
	}
	color.HiRed("\t* The key must be all lowercase!\n")
	fmt.Printf("\t* Any key value that contains any of these charaters: ':#[]()*' must be double quoted!\n")
	fmt.Printf("\t* The key `logmaxsize` value unit is megabytes.\n")
	fmt.Printf("\t* The `tagfile` and `tagkeyname` keys are use to get a tag; useful in AWS, info by looking for\n")
	fmt.Printf("\t  the keyword `tagkeyname` in the configured file `tagfile`, line format: 'keyname value', nothing fancy!\n")
	fmt.Printf("\t* The pagerduty `pdvalidunit` is the unit used to create an event-id so no duplicate is created.\n")
	fmt.Printf("\t  Valid choices are hour or minute. If an event was create at hour X (or minute X) then pagerduty\n")
	fmt.Printf("\t  will not create a new event until the next hour, it sees it as an update to an existing event,.\n")
	fmt.Printf("\t  because it has the same event-id, but do realize there always the possiblity that it could\n")
	fmt.Printf("\t  overlap, certainly if it set to minute, you could get alert every minute!.\n")
	fmt.Printf("\t  If the `pdvalidunit` value is invalid then it defaults to hour, valid options are `hour` or `minute`.\n")
	fmt.Printf("\t* The key `emailsubjecttag` is use for email filtering.\n")
	fmt.Printf("\t* Syslog Valid `syslogpriority`: ")
	for keyPriority, _ := range myUtils.SyslogPriority {
		fmt.Printf("%s ", keyPriority)
	}
	cnt := 0
	fmt.Printf("\n\t* Syslog Valid `syslogfacility`: ")
	for keyFacility, _ := range myUtils.SyslogFacility {
		if cnt > 5 {
			fmt.Printf("\n\t\t%s ", keyFacility)
			cnt = 0
		} else {
			fmt.Printf("%s ", keyFacility)
			cnt += 1
		}
	}
	fmt.Printf("\n")
	os.Exit(0)
}

// Function to show the help information
func Help(exitVal int) {
	fmt.Printf("%s", myGlobal.MyInfo)
	if strings.HasPrefix(myGlobal.MyProgname, "check-") == true {
		optionList := "<-config config file> <-noalert> <-nostats> <-check mode> <-setup> <-version> <-help>"
		fmt.Printf("\nUsage : %s\n\tflags: %s\n", myGlobal.MyProgname, optionList)
		fmt.Printf("\t*config: the configuration file to use, should be full path, use --setup for more information.\n")
		fmt.Printf("\t*check: mode, this is defined per check, use 'mode help' to see valid modes.\n")
		fmt.Printf("\tnoalert: do not send alert.\n")
		fmt.Printf("\tstats: create the stats.\n")
		fmt.Printf("\tsetup: show the setup guide.\n")
		fmt.Printf("\tversion: print %s version.\n", myGlobal.MyProgname)
		fmt.Printf("\thelp: short version of this help page.\n")
		fmt.Printf("\n\t* == required flag.\n")
	}
	if strings.HasPrefix(myGlobal.MyProgname, "get-") == true {
		optionList := "<-mode> <-unit> <-top> <-version> <-help>"
		fmt.Printf("\nUsage : %s\n\tflags: %s\n", myGlobal.MyProgname, optionList)
		fmt.Printf("\tmode: this is defined per script, use 'mode help' to see valid modes.\n")
		fmt.Printf("\tunit: the unit to be use, KB, MB, GB or TB, default to %s.\n", myGlobal.MyUnit)
		fmt.Printf("\t*top: show the given tops resource usage, default show %s\n", myGlobal.MyTop)
		fmt.Printf("\tversion: print %s version.\n", myGlobal.MyProgname)
		fmt.Printf("\thelp: short version of this help page.\n")
		fmt.Printf("\n\t* == only available on certain script.\n")
	}
	os.Exit(exitVal)
}
