// +build darwin

package process

import "errors"

var ErrNoPidNamespace = errors.New("darwin do not support pid namespace")

type process struct {
	pid int
}

func NewProcess(pid int) *process {
	return &process{
		pid: pid,
	}
}

func GetHostPids(nspid int) ([]int, error) {
	return nil, ErrNoPidNamespace
}
