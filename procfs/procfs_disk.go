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
// Date			:	July 1, 2017
//
// History		:
//	Date:			Author:		Info:
//	July 1, 2017	LIS			First Go release
//
// TODO:

package procfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

// function to get all disk info based on those mounted
func getMountInfo() map[string]*sysMount {
	contents, err := ioutil.ReadFile(PROC_SYS_MOUNTS)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	// create the map
	devMounted := make(map[string]*sysMount)
	// prep the regex
	lineMatch, _ := regexp.Compile(sysMountsRegex)
	devSymlink, _ := regexp.Compile(symRegex)
	// read per line
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		// skip empty line
		if len(line) > 0 {
			// we only want those that matching parRegex
			match := lineMatch.MatchString(line)
			if match {
				// we need the first 3 fields : device, mountpoint, type and first word of mount (rw or ro)
				currDevice := strings.Fields(line)[0]
				currMountPoint := strings.Fields(line)[1]
				currFSType := strings.Fields(line)[2]
				currState := strings.Split(strings.Fields(line)[3], ",")[0]
				// check is we have a possible symlink or fullpath
				match = devSymlink.MatchString(currDevice)
				if match {
					currDevice, _ = filepath.EvalSymlinks(currDevice)
				}
				// create the sysMount info for this mountpoint
				devSysMount := &sysMount{
					device:     currDevice,
					mountPoint: currMountPoint,
					fsType:     currFSType,
					mountState: currState,
				}
				// get the stats
				devSysMount.Update()
				// add to list
				devMounted[currMountPoint] = devSysMount
			}
		}
	}
	return devMounted
}

func NewDisk() map[string]*sysMount {
	return getMountInfo()
}

// function to get the given disk stats
func (mountPtr *sysMount) Update() {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(mountPtr.mountPoint, &fs)
	if err != nil {
		fmt.Printf("Errored: %s\n", err.Error())
		os.Exit(1)
	}
	mountPtr.totalSpace = fs.Blocks * uint64(fs.Bsize)
	mountPtr.totalUse = fs.Bfree * uint64(fs.Bsize)
	mountPtr.totalFree = (fs.Blocks * uint64(fs.Bsize)) - (fs.Bfree * uint64(fs.Bsize))
	mountPtr.totalInodes = fs.Files
	mountPtr.freeInodes = fs.Ffree
}

// functions to get disk/partitions element info
func (mountPtr *sysMount) GetType() string {
	return mountPtr.fsType
}

func (mountPtr *sysMount) GetSize(unit uint64) uint64 {
	return mountPtr.totalSpace / unit
}

func (mountPtr *sysMount) GetUse(unit uint64) uint64 {
	return mountPtr.totalUse / unit
}

func (mountPtr *sysMount) GetFree(unit uint64) uint64 {
	return mountPtr.totalFree / unit
}

func (mountPtr *sysMount) GetInodes(unit uint64) uint64 {
	return mountPtr.totalInodes / unit
}

func (mountPtr *sysMount) GetFreeInodes(unit uint64) uint64 {
	return mountPtr.freeInodes / unit
}

func (mountPtr *sysMount) GetMountPoint() string {
	return mountPtr.mountPoint
}

func (mountPtr *sysMount) GetDev() string {
	return mountPtr.device
}

func (mountPtr *sysMount) GetState() string {
	return mountPtr.mountState
}
