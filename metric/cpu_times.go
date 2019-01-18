package metric

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type cpuTimesStat struct {
	*procPidTimesStat
	*procTimesStat
}

type procPidTimesStat struct {
	Utime uint64 `json:"utime"`
	Stime uint64 `json:"stime"`
}

type procTimesStat struct {
	User   uint64 `json:"user"`
	System uint64 `json:"system"`
	Idle   uint64 `json:"idle"`
	Nice   uint64 `json:"nice"`
	Iowait uint64 `json:"iowait"`
}

func (t *cpuTimesStat) total() uint64 {
	return t.User + t.Nice + t.System + t.Idle + t.Iowait
}

func (t *cpuTimesStat) sysUsed() uint64 {
	return t.User + t.Nice + t.System
}

func (t *cpuTimesStat) procUsed() uint64 {
	return t.Utime + t.Stime
}

func sampleCPUtimesStat() *cpuTimesStat {
	pps := getProcPidStat()
	ps := getProcStat()

	if pps == nil || ps == nil {
		return nil
	}
	
	return &cpuTimesStat{
		procPidTimesStat: pps,
		procTimesStat:    ps,
	}
}

// Reads utime and stime from /proc/[pid]/stat file
func getProcPidStat() *procPidTimesStat {
	contents, err := ioutil.ReadFile("/proc/" + pid + "/stat")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fields := strings.Fields(string(contents))
	utime, err := strconv.ParseUint(fields[13], 10, 64)
	if err != nil {
		fmt.Println("procStat[13]: ", err.Error())
	}
	stime, err := strconv.ParseUint(fields[14], 10, 64)
	if err != nil {
		fmt.Println("procStat[13]: ", err.Error())
	}
	return &procPidTimesStat{
		Utime: utime,
		Stime: stime,
	}
}

// Reads stats from /proc/stat file
func getProcStat() *procTimesStat {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	fields := strings.Fields(string(contents))
	user, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		fmt.Println("procStat[0] ", err.Error())
	}
	nice, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		fmt.Println("procStat[1] ", err.Error())
	}
	system, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		fmt.Println("procStat[2] ", err.Error())
	}
	idle, err := strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		fmt.Println("procStat[3] ", err.Error())
	}
	iowait, err := strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		fmt.Println("procStat[4] ", err.Error())
	}
	return &procTimesStat{
		User:   user,
		System: system,
		Idle:   idle,
		Nice:   nice,
		Iowait: iowait,
	}
}
