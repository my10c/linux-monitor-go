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
// History	:
// 	Date:			Author:		Info:
//	July 1, 2017	LIS			First Go release
//
// TODO:
// NOTE We only get the fields that are importants

package procfs

const (

	// [KMG]Bytes units
	B  uint64 = 1
	KB uint64 = 1024 * B
	MB uint64 = 1024 * KB
	GB uint64 = 1024 * MB
	TB uint64 = 1024 * GB

	sysHZ  = 100
	minPid = 300
	// the prod directory
	PROC_DIR = "/proc"
	// system wide
	PROC_SYS_NETDEV  = "/proc/net/dev"
	PROC_SYS_LOADAVG = "/proc/loadavg"
	PROC_SYS_MEMINFO = "/proc/meminfo"
	PROC_SYS_MOUNTS  = "/proc/mounts"
	PROC_SYS_STAT    = "/proc/stat"
	PROC_SYS_UPTIME  = "/proc/uptime"
	// process base (uses pid)
	PROC_PID_CMDLINE = "cmdline"
	PROC_PID_COMM    = "comm"
	PROC_PID_LIMITS  = "limits"
	PROC_PID_SMAP    = "smaps"
	PROC_PID_STAT    = "stat"
)

// from the /proc/stat and then a single cpu line stats
type cpuStat struct {
	user      float64 `json:"user"`
	nice      float64 `json:"nice"`
	system    float64 `json:"system"`
	idle      float64 `json:"idle"`
	ioWait    float64 `json:"iowait"`
	irq       float64 `json:"irq"`
	softIRQ   float64 `json:"softirq"`
	steal     float64 `json:"steal"`
	guest     float64 `json:"guest"`
	guestNice float64 `json:"guestnice"`
}

type sysLoadavg struct {
	load1Avg  float64 `json:"load1navg"`
	load5Avg  float64 `json:"load5avg"`
	load10Avg float64 `json:"load10avg"`
	execProc  uint64  `json:"execproc"`
	execQueue uint64  `json:"execqueue"`
	lastPid   uint64  `json:"lastpid"`
}

type sysMemInfo struct {
	unit         string `json:"unit"`
	memTotal     uint64 `json:"memtotal"`
	memFree      uint64 `json:"memfree"`
	memAvailable uint64 `json:"memavailable"`
	buffers      uint64 `json:"buffers"`
	cached       uint64 `json:"cached"`
	swapCached   uint64 `json:"swapcached"`
	swapTotal    uint64 `json:"swapTotal"`
	swapFree     uint64 `json:"swapTotal"`
}

type sysMount struct {
	device      string `json:"device"`
	mountPoint  string `json:"mount"`
	fsType      string `json:"fstype"`
	mountState  string `json:"state"`
	totalSpace  uint64 `json:"total"`
	totalUse    uint64 `json:"used"`
	totalFree   uint64 `json:"free"`
	totalInodes uint64 `json:"inodes"`
	freeInodes  uint64 `json:"freeinodes"`
}

type sysStat struct {
	bootTime    uint64    `json:"boottime"`
	cpuTotal    cpuStat   `json:"cputotal"`
	cpu         []cpuStat `json:"cpu"`
	cntxtSwitch uint64    `json:"cntxtswitch"`
	procRunning uint64    `json:"procrunning"`
	procBlocked uint64    `json:"procblocked"`
}

type sysUptime struct {
	upTime   float64 `json:"uptime"`
	idleTime float64 `json:"idletime"`
}

type procCmdline struct {
	cmdArgs string `json:"cmdargs"`
}

type procComm struct {
	command string `json:"command"`
}

// Rss: resident memory usage, all memory the process uses,
//      including all memory this process shares with other processes. It does not include swap;
// Shared: memory that this process shares with other processes;
// Private: private memory used by this process, you can look for memory leaks here;
// Swap: swap memory used by the process;
// Pss: Proportional Set Size, a good overall memory indicator.
//      It is the Rss adjusted for sharing: if a process has 1MiB private and 20MiB shared
//      between other 10 processes, Pss is 1 + 20/10 = 3MiB
type procMem struct {
	comm    string `json:"com"`
	pid     uint64 `json:"pid"`
	rss     uint64 `json:"rss"`
	pss     uint64 `json:"pss"`
	shared  uint64 `json:"shared"`
	private uint64 `json:"private"`
	swap    uint64 `json:"swap"`
}

// -1 == unlimited
type procLimits struct {
	cpuTime          int64 `json:"cputime"`        // seconds
	fileSize         int64 `json:"filesize"`       // bytes
	dataSize         int64 `json:"datasize"`       // bytes
	stackSize        int64 `json:"stacKSize"`      // bytes
	coreFileSize     int64 `json:"corefilesize"`   // bytes
	residentSet      int64 `json:"residentset"`    // bytes
	processes        int64 `json:"processes"`      // processes
	openFiles        int64 `json:"openfiles"`      // files
	lockedMemory     int64 `json:"lockedmemory"`   // bytes
	addressSpace     int64 `json:"addressspace"`   // bytes
	fileLocks        int64 `json:"filelocks"`      // locks
	pendingSignals   int64 `json:"pendingsignals"` // signals
	msgqueueeSize    int64 `json:"msgqueueesize"`  // bytes
	nicePriority     int   `json:"nicepriority"`
	realtimePriority int   `json:"realtimepriority"`
	realtimeTimeout  int64 `json:"realtimetimeout"` // usecs
}

type procSmaps struct {
	rss          uint64 `json:"rss"`
	pss          uint64 `json:"pss"`
	shared       uint64 `json:"shared"`
	sharedClean  uint64 `json:"sharedclean"`
	sharedDirty  uint64 `json:"shareddirty"`
	private      uint64 `json:"private"`
	privateClean uint64 `json:"privateclean"`
	privateDirty uint64 `json:"privatedirty"`
	swap         uint64 `json:"swap"`
}

// has 52 fields we only want these
type procStat struct {
	pid   uint   `json:"pid"`   // 1
	comm  string `json:"comm"`  // 2
	state string `json:"state"` // 3
	// R Running
	// S Sleeping in an interruptible wait
	// D Waiting in uninterruptible disk sleep
	// Z Zombie
	// T Stopped (on a signal) or (before Linux 2.6.33)
	// t Tracing stop (Linux 2.6.33 onward)
	// W Paging (only before Linux 2.6.0)
	// X Dead (from Linux 2.6.0 onward)
	// x Dead (Linux 2.6.33 to 3.13 only)
	// K Wakekill (Linux 2.6.33 to 3.13 only)
	// W Waking (Linux 2.6.33 to 3.13 only)
	// P Parked (Linux 3.9 to 3.13 only)
	ppid        uint   `json:"ppid"`       // 4
	tty_nr      uint   `json:"ttynr"`      // 7
	minflt      uint64 `json:"minflt"`     // 10
	cminflt     uint64 `json:"cminflt"`    // 11
	majflt      uint64 `json:"majflt"`     // 12
	cmajflt     uint64 `json:"cmajflt"`    // 13
	priority    uint64 `json:"priority"`   // 18
	nice        uint64 `json:"noce"`       // 19
	num_threads uint64 `json:"numthreads"` // 20
	starttime   uint64 `json:"cstarttime"` // 22
	vsize       uint64 `json:"vsize"`      // 23
	rss         uint64 `json:"rss"`        // 24
	rsslim      uint64 `json:"rsslim"`     // 25
}

// System
type systemProc struct {
	stat    *sysStat
	meminfo *sysMemInfo
	uptime  *sysUptime
	loadavg *sysLoadavg
	mounts  *map[string]sysMount
	process *map[string]*processProc
}

// Single process
type processProc struct {
	comm    *procComm
	cmdline *procCmdline
	smaps   *procSmaps
	stat    *procStat
	limit   *procLimits
}

// All process
type allProcMem struct {
	procs map[string]*procMem
}

// network devices
// either read /proc/net/dev or /sys/class/net/*/statistics/*
type netDevice struct {
	ifName     string `json:"iFName"`     // 1
	rxBytes    uint64 `json:"rxbytes"`    // 2	== statistics file: rx_bytes
	rxPackets  uint64 `json:"rxpackets"`  // 3	== statistics file: rx_packets
	rxErrors   uint64 `json:"rxerrors"`   // 4	== statistics file: rx_errors
	rxDropped  uint64 `json:"rxdropped"`  // 5	== statistics file: rx_dropped
	txBytes    uint64 `json:"txbytes"`    // 10	== statistics file: tx_bytes
	txPackets  uint64 `json:"txpackets"`  // 11	== statistics file: tx_packets
	txErrors   uint64 `json:"txerrors"`   // 12	== statistics file: tx_errors
	txDropped  uint64 `json:"txdropped"`  // 13	== statistics file: tx_dropped
	collisions uint64 `json:"collisions"` // 15	== statistics file: collisions
	carrier    uint64 `json:"carrier"`    // 16	== statistics file: tx_carrier_errors
}

var (
	// for system
	sysStatRegex    = `^(btime|cpu|ctxt|procs_running|procs_blocked)`
	sysMeminfoRegex = `^(MemTotal|MemFree|MemAvailable|Buffers|Cached|SwapCached|SwapTotal|SwapFree)`
	sysMountsRegex  = `^(/dev/)(xvd|sd|disk|mapper)`
	// for process
	smapsRegex  = `^(Rss:|Pss:|Shared_Clean:|Shared_Dirty:|Private_Clean:|Private_Dirty:|Swap:)`
	limitsRegex = `^Max(cpu time|file size|data size|stack size|core file size|resident set|processes|open files|locked memory|address space|file locks|pending signals|msgqueue size|nice priority|realtime priority|realtime timeout)`
	// for disks
	symRegex = `^(/dev/)(disk|mapper)`
	// for network. so beside chars and digits it might contain '-', local interface is still lo
	// the regex is require for /proc/net/dev
	netRegex = `^([A-Za-z].*[0-9A-Za-z\-]*|\s.*lo):`
)
