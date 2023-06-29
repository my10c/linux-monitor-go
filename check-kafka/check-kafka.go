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
	//myKafka		"github.com/my10c/linux-monitor-go/kafka"
	myZoo		"github.com/my10c/linux-monitor-go/zookeeper"
	myGlobal	"github.com/my10c/linux-monitor-go/global"
	// myThreshold	"github.com/my10c/linux-monitor-go/threshold"
	myAlert		"github.com/my10c/linux-monitor-go/alert"
	//myStats		"github.com/my10c/linux-monitor-go/stats"
)

const (
	extraInfo = "the `brokers` values should be kafka-ids (numbers) separated by spaces."
	CheckVersion = "0.1"
)

var (
	cfgRequired = []string{"zookeeperhost", "kafkahost"}
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
	fmt.Printf("\t broker		: check if all kafka brokers are up.\n")
	fmt.Printf("\t topic		: check create and delete topic (TODO).\n")
	fmt.Printf("\t pubsub		: check publish and consume a messaage. (TODO).\n")
	fmt.Printf("\t showconfig	: show the current configuration and then exit.\n")
	os.Exit(3)
}

func main() {
	// need to be root since the config file wil have passwords
	myUtils.IsRoot()
	// var thresHold string = ""
	var exitMsg string
	var extraMsg string
	var zKConn *myZoo.ZkConn
	// add the extra setup info
	myGlobal.ExtraInfo = extraInfo
	myGlobal.MyVersion = CheckVersion
	cfgFile, checkMode := myInit.InitArgs(cfgRequired)
	switch checkMode {
		case "broker":
			cfgRequired = append(cfgRequired, "brokers")
		case "topic":
			cfgRequired = append(cfgRequired, "broker", "topic")
		case "pubsub":
			cfgRequired = append(cfgRequired, "broker", "topic")
	}
	cfgDict := myInit.InitConfig(cfgRequired, cfgFile)
	myInit.InitLog()
	myUtils.SignalHandler()
	// we need to switch again as broker is a pure zookeeeper call
	switch checkMode {
		case "broker":
			zKConn, _ = myZoo.New(myUtils.StringToSlice(cfgDict["zookeeperhost"]))
			defer zKConn.Close()
		case "topic":
		case "pubsub":
	}
	//--> stats := myStats.New()
	// data := time.Now().Format(time.RFC3339)
	iter, _ := strconv.Atoi(cfgDict["iter"])
	iterWait, _ := time.ParseDuration(cfgDict["iterwait"])
	for cnt :=0 ; cnt < iter ; cnt++ {
		switch checkMode {
			case "broker":
				exitVal, err = zKConn.CheckKafkaBroker(myUtils.StringToSlice(cfgDict["brokers"]))
				extraMsg = fmt.Sprintf("brokers: %s", cfgDict["brokers"])
			case "topic":
				fmt.Printf("mode not available yet....\n")
				os.Exit(0)
			case "pubsub":
				fmt.Printf("mode not available yet....\n")
				os.Exit(0)
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
		exitMsg = fmt.Sprintf("%s %s - Check running mode: %s - Error: %s\n",
			strings.ToUpper(myGlobal.MyProgname), myGlobal.Result[exitVal], checkMode, err.Error())
	} else {
		exitMsg = fmt.Sprintf("%s %s - Check running mode: %s - %s \n",
		strings.ToUpper(myGlobal.MyProgname), myGlobal.Result[exitVal], checkMode, extraMsg)
	}

	fmt.Printf("%s", exitMsg)
	myUtils.LogMsg(fmt.Sprintf("%s", exitMsg))
	os.Exit(exitVal)
}
