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
// Date			:	Jul 14, 2017
//
// History	:
// 	Date:			Author:		Info:
//	June 4, 2017	LIS			First Go release
//	Jul 14, 2017	LIS			Add more options for stats
//
// TODO:

package global

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	HR       = "__________________________________________________"
	OK       = 0
	WARNING  = 1
	CRITICAL = 2
	UNKNOWN  = 3

	// [KMG]Bytes units
	B  uint64 = 1
	KB uint64 = 1024 * B
	MB uint64 = 1024 * KB
	GB uint64 = 1024 * MB
	TB uint64 = 1024 * GB
)

var (
	MyVersion   = "0.1"
	now         = time.Now()
	MyProgname  = path.Base(os.Args[0])
	myAuthor    = "Luc Suryo"
	myCopyright = "Copyright 2014 - " + strconv.Itoa(now.Year()) + " ©badassops"
	myLicense   = "License BSD, http://www.freebsd.org/copyright/freebsd-license.html ♥"
	myEmail     = "<luc@badassops.com>"

	// set the check info, we do not show the version as teh global version is the framework
	// version and not the check's
	MyInfo = fmt.Sprintf("%s\n%s\n%s\nWritten by %s %s\n",
		MyProgname, myCopyright, myLicense, myAuthor, myEmail)

	// Global variables
	Logfile  string
	ConfFile string
	MyUnit   string = "KB"
	MyTop    string = "10"

	// special for addtional setup
	ExtraInfo string

	// defaults
	DefaultValues     map[string]string
	DefaultConfDir    = "/etc/nagios-plugins-go"
	DefaultConfigFile = fmt.Sprintf("%s/nagios-plugins-go.yaml", DefaultConfDir)
	// alert, stats logging and debuging mode
	DefaultNoAlert     = "false"
	DefaultCreateStats = "false"
	DefaultNoLog       = "false"
	DefaultDebug       = "false"

	// for logging
	DefaultLog           map[string]string
	DefaultLogsDir       = "/var/log/nagios-plugins-go"
	DefaultLogFile       = fmt.Sprintf("%s.log", MyProgname)
	DefaultLogMaxSize    = 128 // megabytes
	DefaultLogMaxBackups = 14  // 14 files
	DefaultLogMaxAge     = 14  // 14 days

	// email
	DefaultEmail               map[string]string
	DefaultEmailFrom           = ""
	DefaultEmailFromName       = ""
	DefaultEmailTo             = ""
	DefaultEmailToName         = ""
	DefaultEmailUser           = ""
	DefaultEmailpassword       = ""
	DefaultEmailhost           = "localhost"
	DefaultEmailHostPort       = 25
	DefaultEmailHostSubjectTag = "[MONITOR]"

	// tag
	DefaultTag     map[string]string
	DefaultTagfile = ""
	DefaultTagKey  = ""

	// syslog
	DefaultSyslog         map[string]string
	DefaultSyslogTag      = fmt.Sprintf("[%s]", MyProgname)
	DefaultSyslogPriority = "LOG_INFO"
	DefaultSyslogFacility = "LOG_SYSLOG"

	// pagerdutry
	DefaultPD            map[string]string
	DefaultPDServiceKey  = ""
	DefaultPDServiceName = ""
	DefaultPDValidUnit   = "hour"
	DefaultPDEvent       = "MONITOR ALERT"

	// slack
	DefaultSlack           map[string]string
	DefaultSlackServiceKey = ""
	DefaultSlackChannel    = ""
	DefaultSlackUser       = "MONITOR"
	DefaultSlackIconEmoji  = ":bangbang:"

	// result wording
	Result = []string{"OK", "WARNING", "CRITICAL", "UNKNOWN"}

	// stats is always optional but has a default value. so this is hardcoded!
	DefaultStats        map[string]string
	StatsOptionalKeys   = []string{"statsdir", "statsfile", "statstid", "statstformat"}
	DefaultStatsDir     = DefaultLogsDir
	DefaultStatsFile    = fmt.Sprintf("%s", MyProgname)
	DefaultStatsTId     = "_t"
	DefaultStatsTFormat = "2006-01-02T15:04:05Z"

	DefaultIter      map[string]string
	IterOptionalKeys = []string{"iter", "iterwait"}
	DefaultIterCnt   = 3
	DefaultIterWait  = 10 // seconds

	// Shared map between checks
	SharedMap map[string]string
)

func init() {
	// setup the default value, these are hardcoded.
	ExtraInfo = ""
	// the common section
	DefaultValues = make(map[string]string)
	// use by check scripts
	DefaultValues["noalert"] = DefaultNoAlert
	DefaultValues["stats"] = DefaultCreateStats
	DefaultValues["nolog"] = DefaultNoLog
	DefaultValues["debug"] = DefaultDebug
	// use ny get scripts
	DefaultValues["unit"] = MyUnit
	DefaultValues["top"] = MyTop
	// for Log
	DefaultLog = make(map[string]string)
	DefaultLog["logdir"] = DefaultLogsDir
	DefaultLog["logfile"] = DefaultLogFile
	DefaultLog["logmaxsize"] = strconv.Itoa(DefaultLogMaxSize)
	DefaultLog["logmaxbackups"] = strconv.Itoa(DefaultLogMaxBackups)
	DefaultLog["logmaxage"] = strconv.Itoa(DefaultLogMaxAge)
	// for email
	DefaultEmail = make(map[string]string)
	DefaultEmail["emailfrom"] = DefaultEmailFrom
	DefaultEmail["emailfromname"] = DefaultEmailFromName
	DefaultEmail["emailto"] = DefaultEmailTo
	DefaultEmail["emailtoname"] = DefaultEmailToName
	DefaultEmail["emailsubjecttag"] = DefaultEmailHostSubjectTag
	DefaultEmail["emailuser"] = DefaultEmailUser
	DefaultEmail["emailpass"] = DefaultEmailpassword
	DefaultEmail["emailhost"] = DefaultEmailhost
	DefaultEmail["emailhostport"] = strconv.Itoa(DefaultEmailHostPort)
	// these are for getting a instance/system tag
	DefaultTag = make(map[string]string)
	DefaultTag["tagfile"] = DefaultTagfile
	DefaultTag["tagkeyname"] = DefaultTagKey
	// for syslog
	DefaultSyslog = make(map[string]string)
	DefaultSyslog["syslogtag"] = DefaultSyslogTag
	DefaultSyslog["syslogpriority"] = DefaultSyslogPriority
	DefaultSyslog["syslogfacility"] = DefaultSyslogFacility
	// for pagerduty
	DefaultPD = make(map[string]string)
	DefaultPD["pdservicekey"] = DefaultPDServiceKey
	DefaultPD["pdservicename"] = DefaultPDServiceName
	DefaultPD["pdvalidunit"] = DefaultPDValidUnit
	DefaultPD["pdevent"] = DefaultPDEvent
	// for slack
	DefaultSlack = make(map[string]string)
	DefaultSlack["slackservicekey"] = DefaultSlackServiceKey
	DefaultSlack["slackchannel"] = DefaultSlackChannel
	DefaultSlack["slackuser"] = DefaultSlackUser
	DefaultSlack["iconemoji"] = DefaultSlackIconEmoji
	// for stat
	DefaultStats = make(map[string]string)
	DefaultStats["statsdir"] = DefaultStatsDir
	DefaultStats["statsfile"] = DefaultStatsFile
	DefaultStats["statstid"] = DefaultStatsTId
	DefaultStats["statstformat"] = DefaultStatsTFormat
	// for iter
	DefaultIter = make(map[string]string)
	DefaultIter["iter"] = strconv.Itoa(DefaultIterCnt)
	DefaultIter["iterwait"] = strconv.Itoa(DefaultIterWait)
	// the shared map
	SharedMap = make(map[string]string)
	SharedMap["unit"] = "KB"
	SharedMap["topcount"] = "10"
	SharedMap["zookeeperhost"] = "localhost:2181"
	SharedMap["kafkahost"] = "localhost:9092"
}
