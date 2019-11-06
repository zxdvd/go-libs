
// +build linux

package process

import (
	"io/ioutil"
	"os"
	"fmt"
	"bytes"
	"strings"
        "strconv"
)


const	trimStr = " \t\n\r"
type process struct {
	pid int
	statusKV map[string]string
}

func NewProcess(pid int) *process {
	return &process{
		pid: pid,
	}
}


func (p *process) readProcStatus() error {
	if p.statusKV != nil {
		return nil
	}
	p.statusKV = map[string]string{}
	filepath := fmt.Sprintf("/proc/%d/status", p.pid)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}
		key := string(bytes.Trim(parts[0], trimStr))
		value := string(bytes.Trim(parts[1], trimStr))
		p.statusKV[key] = value
	}
	return nil
}

func (p *process) NSPid() int {
	_ = p.readProcStatus()
	if pidStr, ok :=  p.statusKV["NSpid"]; ok {
		parts := strings.SplitN(pidStr, "\t", 2)
		if len(parts) != 2 {
			return -1
		}
		if nspid, err := strconv.Atoi(strings.Trim(parts[1], trimStr)); err == nil {
			return nspid
		}
	}
	return -1
}

func GetHostPids(nspid int) ([]int, error) {
	proc, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	names, err := proc.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	hostpids := make([]int, 0, 1)
	for _, n := range names {
		if !isUint(n) {
			continue
		}
		pid, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		p := NewProcess(pid)
		if nspid == p.NSPid() {
			hostpids = append(hostpids, pid)
		}
	}
	return hostpids, nil
}

func isUint(s string) bool {
	for i, ch := range []byte(s) {
		if ch < '0' || ch > '9' {
			return false
		}
		if ch == '0' && i == 0 {
			return false
		}
	}
	return true
}
