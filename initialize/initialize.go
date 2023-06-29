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
// Version		:	0.2
//
// Date			:	Jul 14 4, 2017
//
// History	:
// 	Date:			Author:		Info:
//	June 4, 2017	LIS			First Go release
//	July 14, 2017	LIS			fix for stats
//
// TODO:

package initialize

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	myGlobal "github.com/my10c/linux-monitor-go/global"
	myHelp "github.com/my10c/linux-monitor-go/help"
	myUtils "github.com/my10c/linux-monitor-go/utils"

	"github.com/my10c/simpleyaml"
	"gopkg.in/natefinch/lumberjack.v2"
)

// type used for flags in initArgs
type stringFlag struct {
	value string
	set   bool
}

// Function for the stringFlag struct, set the values
func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

// Function for the stringFlag struct, get the values
func (sf *stringFlag) String() string {
	return sf.value
}

// Function to return the yaml value, nil if error or nil if not found
func getYamlValue(yamFile *simpleyaml.Yaml, section string, key string) (string, error) {
	// Check if section exist and/or key, no point to go further if it doesn't exist
	keyExist := yamFile.GetPath(section, key)
	if keyExist.IsFound() == false {
		err := fmt.Errorf("Section %s and/or key %s not found\n", section, key)
		return "", err
	}
	// We need to ge the value and since we do not know what it is, we check
	// against the 3 supported type
	// check if value is a string
	if value, err := yamFile.Get(section).Get(key).String(); err == nil {
		return value, err
	}
	// check if value is a int
	if value, err := yamFile.Get(section).Get(key).Int(); err == nil {
		return strconv.Itoa(value), err
	}
	// check if value is a boolean
	if value, err := yamFile.Get(section).Get(key).Bool(); err == nil {
		return strconv.FormatBool(value), err
	}
	// we got here so this is neither a string, int or boolean
	err := fmt.Errorf("Unsupported value for section %s and key %s, suported are: string, int and bool\n", section, key)
	myUtils.LogMsg(fmt.Sprintf("%s\n", err.Error()))
	return "", err
}

// Function to get the configuration
func InitConfig(cfgList []string, argv ...string) map[string]string {
	// working variable
	var missingKeys []string
	dictCfg := make(map[string]string)
	// open given file and check that is a correct yaml file
	cfgFile, err := ioutil.ReadFile(argv[0])
	myUtils.ExitIfError(err)
	yamlFile, err := simpleyaml.NewYaml(cfgFile)
	myUtils.ExitIfError(err)
	// first check if the default common values need to be modify
	for defaultKey, _ := range myGlobal.DefaultValues {
		if newValue, err := getYamlValue(yamlFile, "common", defaultKey); err == nil {
			// replace the default value
			myGlobal.DefaultValues[defaultKey] = newValue
		}
	}
	// for log
	for defaultLog, _ := range myGlobal.DefaultLog {
		if newValue, err := getYamlValue(yamlFile, "log", defaultLog); err == nil {
			// replace the default value
			myGlobal.DefaultLog[defaultLog] = newValue
		}
	}
	// for tag
	for defaultTag, _ := range myGlobal.DefaultTag {
		if newValue, err := getYamlValue(yamlFile, "tag", defaultTag); err == nil {
			// replace the default value
			myGlobal.DefaultTag[defaultTag] = newValue
		}
	}
	// for email
	for defaultEmail, _ := range myGlobal.DefaultEmail {
		if newValue, err := getYamlValue(yamlFile, "email", defaultEmail); err == nil {
			// replace the default value
			myGlobal.DefaultEmail[defaultEmail] = newValue
		}
	}
	// for Syslog
	for defaultSyslog, _ := range myGlobal.DefaultSyslog {
		if newValue, err := getYamlValue(yamlFile, "syslog", defaultSyslog); err == nil {
			// replace the default value
			myGlobal.DefaultSyslog[defaultSyslog] = newValue
		}
	}
	// for Pagerduty
	for defaultPD, _ := range myGlobal.DefaultPD {
		if newValue, err := getYamlValue(yamlFile, "pagerduty", defaultPD); err == nil {
			// replace the default value
			myGlobal.DefaultPD[defaultPD] = newValue
		}
	}
	// for Slack
	for defaultSlack, _ := range myGlobal.DefaultSlack {
		if newValue, err := getYamlValue(yamlFile, "slack", defaultSlack); err == nil {
			// replace the default value
			myGlobal.DefaultSlack[defaultSlack] = newValue
		}
	}
	// set the config value
	// we set first the stats default values since these are optional
	for statsCnt := range myGlobal.StatsOptionalKeys {
		statsKey := myGlobal.StatsOptionalKeys[statsCnt]
		if statsValue, err := getYamlValue(yamlFile, myGlobal.MyProgname, statsKey); err == nil {
			dictCfg[statsKey] = statsValue
			// also update the global!
			myGlobal.DefaultStats[statsKey] = statsValue
		} else {
			dictCfg[statsKey] = myGlobal.DefaultStats[statsKey]
		}
	}
	// last we set the iter value. like stats these are optional
	for iterCnt := range myGlobal.IterOptionalKeys {
		iterKey := myGlobal.IterOptionalKeys[iterCnt]
		if iterValue, err := getYamlValue(yamlFile, myGlobal.MyProgname, iterKey); err == nil {
			dictCfg[iterKey] = iterValue
		} else {
			dictCfg[iterKey] = myGlobal.DefaultIter[iterKey]
		}
	}
	// get the check keys
	for cnt := range cfgList {
		keyName := cfgList[cnt]
		if cfgValue, err := getYamlValue(yamlFile, myGlobal.MyProgname, keyName); err == nil {
			// assign the value
			dictCfg[keyName] = cfgValue
		} else {
			// important part of the shared map is that the name of the keys can not be the
			// same across all checks unless the exact default value is valid for these checks
			// check shared keys
			if defaultValue, ok := myGlobal.SharedMap[keyName]; ok {
				dictCfg[keyName] = defaultValue
			} else {
				// key is really missing
				missingKeys = append(missingKeys, keyName)
			}
		}
	}
	// make sure we have all required configs
	if len(missingKeys) != 0 {
		fmt.Printf("Following keys are missing in the configration files: %s\n", missingKeys)
		myUtils.LogMsg(fmt.Sprintf("Following keys are missing in the configration files: %s", missingKeys))
		os.Exit(2)
	}
	return dictCfg
}

// Function to initialize logging
func InitLog() {
	// is nolog was requested then we return
	if myGlobal.DefaultValues["nolog"] == "true" {
		return
	}
	if len(myGlobal.DefaultLog["logfile"]) > 0 {
		// create directory
		err := os.MkdirAll(myGlobal.DefaultLog["logdir"], 0755)
		if err != nil {
			fmt.Printf("Unable to create the Log directory, logs are send to console!\n")
			myUtils.LogMsg(fmt.Sprintf("%s\n", err.Error()))
			return
		}
		logFileFullPath := fmt.Sprintf("%s/%s", myGlobal.DefaultLog["logdir"], myGlobal.DefaultLog["logfile"])
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		MaxSize, _ := strconv.Atoi(myGlobal.DefaultLog["logmaxsize"])
		MaxBackups, _ := strconv.Atoi(myGlobal.DefaultLog["logmaxbackups"])
		MaxAge, _ := strconv.Atoi(myGlobal.DefaultLog["logmaxage"])
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFileFullPath,
			MaxSize:    MaxSize,
			MaxBackups: MaxBackups,
			MaxAge:     MaxAge,
		})
	}
}

// Function to process the given args
func InitArgs(cfg []string) (string, string) {
	var myConfigFile stringFlag
	var myMode stringFlag
	var myUnit stringFlag
	var myTop stringFlag
	var setup *bool
	var noalert *bool
	var stats *bool
	var nolog *bool
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", myGlobal.MyProgname)
		flag.PrintDefaults()
	}
	// common flags
	version := flag.Bool("version", false, "Prints current version and exit.")
	flag.Var(&myMode, "mode", "check mode, use `-mode help` to see available modes.")
	// any check script has these options
	if strings.HasPrefix(myGlobal.MyProgname, "check-") == true {
		flag.Var(&myConfigFile, "config", "Configuration file to be used.")
		setup = flag.Bool("setup", false, "Show the setup information and exit.")
		noalert = flag.Bool("noalert", false, "Send no alert.")
		stats = flag.Bool("stats", false, "Create stats if set.")
		nolog = flag.Bool("nolog", false, "Do not log result.")
	}
	// any get script has these options
	if strings.HasPrefix(myGlobal.MyProgname, "get-") == true {
		flag.Var(&myUnit, "unit", "Unit KB(default), MB, GB or TB.")
		flag.Var(&myTop, "top", "Show top usage.")
	}
	flag.Parse()
	// commons
	if *version {
		fmt.Printf("%s\n", myGlobal.MyVersion)
		os.Exit(0)
	}
	if !myMode.set {
		myHelp.Help(1)
	}
	// its a check script
	if strings.HasPrefix(myGlobal.MyProgname, "check-") == true {
		if *setup {
			myHelp.SetupHelp(cfg)
		}
		// if not set we use the default
		if !myConfigFile.set {
			myConfigFile.Set(myGlobal.DefaultConfigFile)
		}
		// set the noalert and nolog
		myGlobal.DefaultValues["noalert"] = strconv.FormatBool(*noalert)
		myGlobal.DefaultValues["nolog"] = strconv.FormatBool(*nolog)
		myGlobal.DefaultValues["stats"] = strconv.FormatBool(*stats)
	}
	// its a get script
	if strings.HasPrefix(myGlobal.MyProgname, "get-") == true {
		// if unit was given as a arg, override default
		if myUnit.set {
			myGlobal.DefaultValues["unit"] = myUnit.value
		}
		// if top was given as a arg, override default
		if myTop.set {
			myGlobal.DefaultValues["top"] = myTop.value
		}
	}
	return myConfigFile.value, myMode.value
}
