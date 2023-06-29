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

package alerts

import (
	"fmt"
	"os"
	"strings"

	myGlobal	"github.com/my10c/linux-monitor-go/global"
	myTag		"github.com/my10c/linux-monitor-go/tag"
	myUtils		"github.com/my10c/linux-monitor-go/utils"
)

// Function to sent alerts
func SendAlert(exitVal int, checkMode string, checkErr string) error {
	var hostName string
	var message string
	var err error = nil
	var result error = nil
	// create the full message and subject
	errWord := myGlobal.Result[exitVal]
	hostName, hostOK := os.Hostname()
	if hostOK != nil {
		hostName = "Unable to get hostname"
	}
	hostName = strings.TrimSpace(hostName)
	tagInfo, tagOK := myTag.GetTagInfo()
	if tagOK != nil {
		message = fmt.Sprintf("TAG: no tag found\nHost: %s\n%s: %s\nCheck running mode: %s\nError: %s\n",
				hostName, myGlobal.MyProgname, errWord, checkMode, checkErr)
	} else {
		message = fmt.Sprintf("TAG: %s\nHost: %s\n%s: %s\nCheck running mode: %s\nError: %s\n",
			tagInfo, hostName, myGlobal.MyProgname, errWord, checkMode, checkErr)
	}
	// is any of these fails we capture that, hence err could be a set of errors!
	// Syslog : only if syslog tag was not set to of
	if myGlobal.DefaultSyslog["syslogtag"] != "off" {
		result = alertSyslog(message)
		if result != nil {
			err = fmt.Errorf("Syslog %s", result.Error())
		}
	}
	// Email : only if emailto is not empty
	if len(myGlobal.DefaultEmail["emailto"]) > 0 {
		errSubject := fmt.Sprintf("%s %s: %s : %s ",
			myGlobal.DefaultEmail["emailsubjecttag"], errWord, hostName, myGlobal.MyProgname)
		result = alertEmail(message, errSubject)
		if result != nil {
			if err != nil {
				// append error
				err = fmt.Errorf("%s\nEmail %s", err.Error(), result.Error())
			} else {
				err = fmt.Errorf("%s", result.Error())
			}
		}
	}
	// Pagerduty : only if key is empty
	 if len(myGlobal.DefaultPD["pdservicekey"]) > 0 {
		result = alertPD(message, tagInfo, hostName)
		if result != nil {
			if err != nil {
				// append error
				err = fmt.Errorf("%s\nPagerduty %s", err.Error(), result.Error())
			} else {
				err = fmt.Errorf("%s", result.Error())
			}
		}
	 }
	// Slack : only if key is empty {
	if len(myGlobal.DefaultSlack["slackservicekey"]) > 0 {
		result = alertSlack(message)
		if result != nil {
			if err != nil {
				// append error
				err = fmt.Errorf("%s\nSlack %s", err.Error(), result.Error())
			} else {
				err = fmt.Errorf("%s", result.Error())
			}
		}
	}
	if err != nil {
		myUtils.LogMsg(fmt.Sprintf("%s\n", err.Error()))
	}
	return err
}
