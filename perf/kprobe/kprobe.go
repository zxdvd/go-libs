package kprobe


import (
	"bufio"
	"bytes"
	"strings"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const traceDir = "/sys/kernel/debug/tracing"


var kprobeEnabled *bool

var ErrKprobeDisabled = errors.New("kprobe not enable, check kernel config or debugfs mount")
var ErrWrongProbe = errors.New("wrong probe string")


func absPath(path string) string {
	return filepath.Join(traceDir, path)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
}

func isKprobeEnabled() bool {
	if kprobeEnabled != nil {
		return *kprobeEnabled
	}
	if exists(absPath("kprobe_events")) {
		return true
	}
	return false
}


type probeType uint8
const (
	NORMAL_PROBE probeType = 0 + iota
	RETURN_PROBE
)

type Kprobe struct {
	probe string
	typ probeType
	name string
	event string
}

func NewKprobe(p string) (*Kprobe, error) {
	if !isKprobeEnabled() {
		return nil, ErrKprobeDisabled
	}
	var typ probeType
	if p[:2] == "p:" {
		typ = NORMAL_PROBE
	} else if p[:2] == "r:" {
		typ = RETURN_PROBE
	} else {
		return nil, ErrWrongProbe
	}
	left := p[2:]
	parts := strings.SplitN(left, " ", 2)
	if len(parts) < 2 {
		return nil, ErrWrongProbe
	}
	name, event := parts[0], parts[1]
	return &Kprobe{
		probe: p,
		typ: typ,
		name: name,
		event: event,
	}, nil
}

func (p *Kprobe) Add() error {
	return appendString(absPath("kprobe_events"), p.probe)
}

func (p *Kprobe) Remove() error {
	_ = p.Disable()
	removeCmd := "-:" + p.name + "\n"
	err := appendString(absPath("kprobe_events"), removeCmd)
	return err
}

func (p *Kprobe) Name() string {
	return p.name
}

func appendString(path string, line string) error {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}
	_, err = f.WriteString(line)
	f.Close()
	return err
}

func overwriteString(path string, line string) error {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.WriteString(line)
	f.Close()
	return err
}

func (p *Kprobe) Enable() error {
	path := fmt.Sprintf("events/kprobes/%s/enable", p.name)
	return overwriteString(absPath(path), "1")
}

func (p *Kprobe) Disable() error {
	path := fmt.Sprintf("events/kprobes/%s/enable", p.name)
	return overwriteString(absPath(path), "0")
}

func Events() ([][]byte, error) {
	content, err := ioutil.ReadFile(absPath("kprobe_events"))
	return bytes.Split(content, []byte("\n")), err
}

func RemoveAll() error {
	// this file won't exists if no kprobes
	if !exists(absPath("events/kprobes")) {
		return nil
	}
	// this will disable all kprobes
	if err := overwriteString(absPath("events/kprobes/enable"), "0"); err != nil {
		fmt.Printf("failed to disable all kprobes\n")
	}
	return overwriteString(absPath("kprobe_events"), "")
}


func ReadEvents() (chan string, error) {
	f, err := os.OpenFile(absPath("trace_pipe"), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	ch := make(chan string, 30)
	buf := bufio.NewReader(f)
	go func() {
		for {
			s, err := buf.ReadString('\n')
			if err != nil {
				panic(err)
			}
			ch <- s
		}
	}()
	return ch, nil
}

type handleEvent func(p *Kprobe, event string)

func KprobeStream(probe string, fn handleEvent, stop chan struct{}) error {
	p, err := NewKprobe(probe)
	if err != nil {
		return err
	}
	if err := p.Add(); err != nil {
		return err
	}
	defer p.Remove()

	evtCh, err := ReadEvents()
	if err != nil {
		return err
	}
	if err := p.Enable(); err != nil {
		return err
	}
	for {
		select {
			case <- stop:
				return nil
			case evt := <- evtCh:
				fn(p, evt)
		}
	}
	return nil
}
