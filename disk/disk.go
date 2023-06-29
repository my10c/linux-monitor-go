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

package disk

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	myGlobal "github.com/my10c/linux-monitor-go/global"
	myThreshold "github.com/my10c/linux-monitor-go/threshold"
	myUtils "github.com/my10c/linux-monitor-go/utils"
)

const (
	PROCMOUNT = "/proc/mounts"
)

var (
	// valid device we support
	devRegex = `^(/dev/)(xvd|sd|disk|mapper)`
	symRegex = `^(/dev/)(disk|mapper)`
)

type parStruct struct {
	device     string `json:"device"`
	mountpoint string `json:"mount"`
	fsType     string `json:"fstype"`
	mountState string `json:"state"`
}

type diskType struct {
	totalSpace  uint64 `json:"total"`
	totalUse    uint64 `json:"used"`
	totalFree   uint64 `json:"free"`
	totalInodes uint64 `json:"inodes"`
	freeInodes  uint64 `json:"freeinodes"`
	mountPoint  string `json:"mount"`
	device      string `json:"device"`
	fsType      string `json:"fstype"`
	mountState  string `json:"state"`
}

func getPartitions() map[string]parStruct {
	// working variable
	var currPartition parStruct
	// get disks info from proc
	contents, err := ioutil.ReadFile(PROCMOUNT)
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// prep the regex, we ignore the errors
	expDev, _ := regexp.Compile(devRegex)
	expLogics, _ := regexp.Compile(symRegex)
	// create the return map
	detectedPartitions := make(map[string]parStruct)
	// get all lines and walk one at the time
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if line != "" {
			// we need the first 3 fields : device, mountpoint, type and first word of mount (rw or ro)
			currDevice := strings.Fields(line)[0]
			currMountPoint := strings.Fields(line)[1]
			currFSType := strings.Fields(line)[2]
			currState := strings.Split(strings.Fields(line)[3], ",")[0]
			// we only want those matching parRegex
			match := expDev.MatchString(currDevice)
			if match {
				// check is we have a possible symlink or fullpath
				match = expLogics.MatchString(currDevice)
				if match {
					currDevice, _ = filepath.EvalSymlinks(currDevice)
				}
				// get the disk/partion info
				currPartition.device = currDevice
				currPartition.mountpoint = currMountPoint
				currPartition.fsType = currFSType
				currPartition.mountState = currState
				detectedPartitions[currMountPoint] = currPartition
			}
		}
	}
	return detectedPartitions
}

// Function to get the given partition/mount point file system info
func getDiskinfo(path string) *diskType {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	}
	disk := &diskType{
		totalSpace:  fs.Blocks * uint64(fs.Bsize),
		totalFree:   fs.Bfree * uint64(fs.Bsize),
		totalUse:    (fs.Blocks * uint64(fs.Bsize)) - (fs.Bfree * uint64(fs.Bsize)),
		totalInodes: fs.Files,
		freeInodes:  fs.Ffree,
		mountPoint:  path,
	}
	return disk
}

// Function to get the available disks information
func New() map[string]*diskType {
	// create the disk/partition map
	detectedPart := make(map[string]*diskType)
	// will return empty map if no valid disk/partition was found
	for mntPoint, partInfo := range getPartitions() {
		currDisk := getDiskinfo(mntPoint)
		if currDisk == nil {
			return nil
		}
		currDisk.device = partInfo.device
		currDisk.fsType = partInfo.fsType
		currDisk.mountState = partInfo.mountState
		detectedPart[mntPoint] = currDisk
	}
	return detectedPart
}

// Functions to get disk/partitions element info
func (diskPtr *diskType) GetType() string {
	return diskPtr.fsType
}

func (diskPtr *diskType) GetSize(unit uint64) uint64 {
	return diskPtr.totalSpace / unit
}

func (diskPtr *diskType) GetUse(unit uint64) uint64 {
	return diskPtr.totalUse / unit
}

func (diskPtr *diskType) GetFree(unit uint64) uint64 {
	return diskPtr.totalFree / unit
}

func (diskPtr *diskType) GetInodes(unit uint64) uint64 {
	return diskPtr.totalInodes / unit
}

func (diskPtr *diskType) GetFreeInodes(unit uint64) uint64 {
	return diskPtr.freeInodes / unit
}

func (diskPtr *diskType) GetMountPoint() string {
	return diskPtr.mountPoint
}

func (diskPtr *diskType) GetDev() string {
	return diskPtr.device
}

func (diskPtr *diskType) GetState() string {
	return diskPtr.mountState
}

// Function to check available disk space
func (diskPtr *diskType) CheckFree(warn, crit string, unit uint64) int {
	warnThreshold, critThreshold, percent := myThreshold.SanityCheck(true, warn, crit)
	return myThreshold.CalculateUsage(true, percent, warnThreshold, critThreshold,
		diskPtr.GetFree(unit), diskPtr.GetSize(unit))
}

// Function to check available inodes if supported by the filesystem
func (diskPtr *diskType) CheckFreeInode(warn, crit string, unit uint64) int {
	// if both total inodes and free inodes are zero then the filesystem does use inodes
	if (diskPtr.GetInodes(unit) == 0) && (diskPtr.GetFreeInodes(unit) == 0) {
		return 0
	}
	warnThreshold, critThreshold, percent := myThreshold.SanityCheck(true, warn, crit)
	return myThreshold.CalculateUsage(true, percent, warnThreshold, critThreshold,
		diskPtr.GetInodes(unit), diskPtr.GetFreeInodes(unit))
}

// Function to check if the filesystem is mounted RO
func (diskPtr *diskType) CheckRO(mntPoint string) int {
	if diskPtr.mountState == "ro" {
		return 1
	}
	return 0
}

// Function wrapper to the checks
func (diskPtr *diskType) CheckIt(mode string, warn, crit string, unit uint64) int {
	var result int = 3
	switch mode {
	case "diskspace":
		result = diskPtr.CheckFree(warn, crit, unit)
	case "inode":
		result = diskPtr.CheckFreeInode(warn, crit, unit)
	case "ro":
		result = diskPtr.CheckRO(diskPtr.GetMountPoint())
	}
	return result
}

// Function to generate a disk info string
func (diskPtr *diskType) StatusMsg(mode string, unit uint64) string {
	var statusMsg string
	mntPoint := diskPtr.GetMountPoint()
	switch mode {
	case "diskspace":
		statusMsg = fmt.Sprintf("(%s:Free:%d)", mntPoint, diskPtr.GetFree(unit))
	case "inode":
		statusMsg = fmt.Sprintf("(%s:Free Inode:%d)", mntPoint, diskPtr.GetFreeInodes(unit))
	case "ro":
		statusMsg = fmt.Sprintf("(%s:mount state:%s)", mntPoint, diskPtr.GetState())
	}
	return statusMsg
}
