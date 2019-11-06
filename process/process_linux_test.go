// +build linux

package process

import (
	"os"
	"testing"
	"log"
)

func TestProcessLinux(t *testing.T) {
	selfpid := os.Getpid()
	p := NewProcess(selfpid)
	nspid := p.NSPid()
	log.Printf("nspid is %v\n", nspid)
}

func TestGetHostPid(t *testing.T) {
	selfpid := os.Getpid()
	hostpids, err := GetHostPids(selfpid)
	log.Printf("host pid is %v, err is %v\n", hostpids, err)
}
