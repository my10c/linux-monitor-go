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

package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"

	myGlobal "github.com/my10c/linux-monitor-go/global"
)

var (
	// syslog need to this so configuration can use string instead of int
	SyslogPriority = map[string]int{
		"LOG_EMERG":   0,
		"LOG_ALERT":   1,
		"LOG_CRIT":    2,
		"LOG_ERR":     3,
		"LOG_WARNING": 4,
		"LOG_NOTICE":  5,
		"LOG_INFO":    6,
		"LOG_DEBUG":   7,
	}
	SyslogFacility = map[string]int{
		"LOG_MAIL":     0,
		"LOG_DAEMON":   1,
		"LOG_AUTH":     2,
		"LOG_SYSLOG":   3,
		"LOG_LPR":      4,
		"LOG_NEWS":     5,
		"LOG_UUCP":     6,
		"LOG_CRON":     7,
		"LOG_AUTHPRIV": 8,
		"LOG_FTP":      9,
		"LOG_LOCAL0":   10,
		"LOG_LOCAL1":   11,
		"LOG_LOCAL2":   12,
		"LOG_LOCAL3":   13,
		"LOG_LOCAL4":   14,
		"LOG_LOCAL5":   15,
		"LOG_LOCAL6":   16,
		"LOG_LOCAL7":   18,
	}
)

// Function to exit if an error occured
func ExitIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: "+fmt.Sprint(err))
		log.Printf("-< %s >-\n", fmt.Sprint(err))
		os.Exit(1)
	}
}

// Function to exit with the nagios standard exit code and word (WARNING, CRITICAL, UNKNOWN)
// if error was not nil
func ExitWithNagiosCode(exitValue int, err error) {
	if err != nil {
		var nagiosCode string
		switch exitValue {
		case 1:
			nagiosCode = "WARNING"
		case 2:
			nagiosCode = "CRITICAL"
		default:
			nagiosCode = "UNKNOWN"
		}
		fmt.Fprintln(os.Stderr, nagiosCode+" Error: "+fmt.Sprint(err))
		log.Printf("%s -< %s >-\n", nagiosCode, fmt.Sprint(err))
		os.Exit(exitValue)
	}
}

// Function to exit if pointer is nill
func ExitIfNill(ptr interface{}) {
	if ptr == nil {
		fmt.Fprintln(os.Stderr, "Error: got a nil pointer.")
		log.Printf("-< Error: got a nil pointer. >-\n")
		os.Exit(1)
	}
}

// Function to print the given message to stdout and log file
func StdOutAndLog(message string) {
	fmt.Printf("-< %s >-\n", message)
	log.Printf("-< %s >-\n", message)
	return
}

// Function to check if the user that runs the app is root
func IsRoot() {
	if os.Geteuid() != 0 {
		// since this checked the first time, there will be no log available yet
		fmt.Printf("-< %s must be run as root. >-\n", myGlobal.MyProgname)
		os.Exit(1)
	}
}

// Function to log of the nolog was not set to true
func LogMsg(message string) {
	if myGlobal.DefaultValues["nolog"] == "true" {
		return
	}
	log.Printf("%s\n", message)
}

// Function to log any reveived signal
func SignalHandler() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt)
	go func() {
		sigId := <-interrupt
		StdOutAndLog(fmt.Sprintf("received %v", sigId))
		os.Exit(0)
	}()
}

// Function to the MD5 value of the given file
func GetFileMD5(filePath string) (string, error) {
	var returnMD5String string
	// Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return returnMD5String, err
	}
	// Tell the program to call the following function when the current function returns
	defer file.Close()
	//Open a new hash interface to write to
	hash := md5.New()
	// Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("%s\n", err.Error())
		return returnMD5String, err
	}
	// Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	// Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

// Function to the MD5 value of the given file pointer
func GetFilePtrMD5(filePtr *os.File) (string, error) {
	var returnMD5String string
	//Open a new hash interface to write to
	hash := md5.New()
	// Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, filePtr); err != nil {
		log.Printf("%s\n", err.Error())
		return returnMD5String, err
	}
	// Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	// Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

// Function to write a log if debug was enabled
func WriteDebug(debug string, messsage string) {
	debugMode, err := strconv.ParseBool(debug)
	if err != nil {
		return
	}
	if debugMode == true {
		log.Printf("Debug -< %s >-\n", messsage)
	}
	return
}

// Function to check if the system is the given OS
func IsOS(osName string) (string, bool) {
	if runtime.GOOS == osName {
		return runtime.GOOS, true
	}
	return runtime.GOOS, false
}

// Function to check if we are on a Linux system
func IsLinuxSystem() {
	if osName, ok := IsOS("linux"); !ok {
		fmt.Printf("%s", myGlobal.MyInfo)
		fmt.Printf("OS (%s) not supported, this check can only be run on a Linux system.\n", osName)
		os.Exit(1)
	}
	return
}

// Function print helper to show config value
func showSectionValue(section string, sectionDict map[string]string) {
	fmt.Printf("%s:\n", section)
	for mapKey, mapValue := range sectionDict {
		fmt.Printf("\t%s: %s\n", mapKey, mapValue)
	}
	return
}

// Function to show a map entry, key and value and displayed in yaml
func ShowMap(cfgDict map[string]string) {
	if cfgDict == nil {
		// display the common values
		showSectionValue("common", myGlobal.DefaultValues)
		// display the log values
		showSectionValue("log", myGlobal.DefaultLog)
		// display the email values
		showSectionValue("email", myGlobal.DefaultEmail)
		// display the tag values
		showSectionValue("tag", myGlobal.DefaultTag)
		// display the syslog values
		showSectionValue("syslog", myGlobal.DefaultSyslog)
		// display the pagerduty values
		showSectionValue("pagerduty", myGlobal.DefaultPD)
		// display the slack values
		showSectionValue("slack", myGlobal.DefaultSlack)
	}
	if cfgDict != nil {
		fmt.Printf("%s:\n", myGlobal.MyProgname)
		for mapKey, mapValue := range cfgDict {
			fmt.Printf("\t%s: %s\n", mapKey, mapValue)
		}
	}
	return
}

// Function to map string to int for syslog
func GetSyslog(priority string, facility string) (int, int, error) {
	var priorityValue int
	var facilityValue int
	var err error = nil
	if mapVal, ok := SyslogPriority[priority]; ok {
		priorityValue = mapVal
	} else {
		err = fmt.Errorf("Given Syslog Priority is incorrect: %s\n", priority)
	}
	if mapVal, ok := SyslogFacility[facility]; ok {
		facilityValue = mapVal
	} else {
		if err == nil {
			err = fmt.Errorf("Given Syslog Facility is incorrect: %s\n", facility)
		} else {
			err = fmt.Errorf("%sGiven Syslog Facility is incorrect: %s\n", err.Error(), facility)
		}
	}
	return priorityValue, facilityValue, err
}

// Function to convert a strings to a slice
func StringToSlice(args ...string) []string {
	return strings.Fields(strings.Join(args, " "))
}

// Function to convert slice of string to a single string
func SliceToString(array []*string) string {
	var buffer bytes.Buffer
	for idx := range array {
		buffer.WriteString(*array[idx])
	}
	return buffer.String()
}

// Function to remove the last char (given) of a given string
func TrimLastChar(givenString, suffix string) string {
	if strings.HasSuffix(givenString, suffix) {
		givenString = givenString[:len(givenString)-len(suffix)]
	}
	return givenString
}

// function to convert network speed from bytes/sec to X/sec
func ConvertSpeed(bytes uint64, unit string) float64 {
	// speed unit is 1000 and not 1024 for network speed!
	var converted float64 = 0
	switch unit {
	case "KB":
		converted = float64(bytes) / float64(1000)
	case "MB":
		converted = float64(bytes) / float64(1000*1000)
	case "GB":
		converted = float64(bytes) / float64(1000*1000*1000)
	case "TB":
		converted = float64(bytes) / float64(1000*1000*1000*1000)
	}
	return converted
}
