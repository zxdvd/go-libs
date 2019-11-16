package uprobe


import (
	"io"
	"bufio"
	"bytes"
	"strings"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/zxdvd/go-libs/std-helper/fs"
	"github.com/zxdvd/go-libs/std-helper/debug"
)

const traceDir = "/sys/kernel/debug/tracing"


var uprobeEnabled *bool

var ErrUprobeDisabled = errors.New("uprobe not enable, check kernel config or debugfs mount")
var ErrWrongProbe = errors.New("wrong probe string")


func absPath(path string) string {
	return filepath.Join(traceDir, path)
}


func isUprobeEnabled() bool {
	if uprobeEnabled != nil {
		return *uprobeEnabled
	}
	if fs.Exists(absPath("uprobe_events")) {
		return true
	}
	return false
}


type probeType uint8
const (
	NORMAL_PROBE probeType = 0 + iota
	RETURN_PROBE
)

type Uprobe struct {
	name string
	cmdPath string
	symbol string
	probe string
	typ probeType
}

func NewUprobe(p string) (*Uprobe, error) {
	if !isUprobeEnabled() {
		return nil, ErrUprobeDisabled
	}

	buf := bytes.NewBufferString(p)
	tok, _ := buf.ReadString(':')

	var typ probeType
	if tok == "p:" {
		typ = NORMAL_PROBE
	} else if tok == "r:" {
		typ = RETURN_PROBE
	} else {
		return nil, ErrWrongProbe
	}
	probeName, err := buf.ReadString(' ')
	if err != nil {
		return nil, err
	}
	probeName = strings.TrimSpace(probeName)

	cmdPath, err := buf.ReadString(':')
	if err != nil {
		return nil, err
	}
	cmdPath = strings.Trim(cmdPath, " :")

	symbol, err := buf.ReadString(' ')
	if err != nil && err != io.EOF {
		return nil, err
	}
	symbol = strings.Trim(symbol, " ")
	// if the symbol is a not a offset, then it is a symbol name
	if !strings.HasPrefix(symbol, "0x") {
		symOffset, err := debug.GetSymbolOffset(cmdPath, symbol)
		if err != nil {
			return nil, err
		}
		symbol = fmt.Sprintf("0x%x", symOffset)
	}

	probe := fmt.Sprintf("%s%s %s:%s %s", tok, probeName, cmdPath, symbol, buf.String())
	fmt.Println("final probe string is: ", probe)

	return &Uprobe{
		name: probeName,
		cmdPath: cmdPath,
		symbol: symbol,
		probe: probe,
		typ: typ,
	}, nil
}

func (p *Uprobe) Add() error {
	return fs.AppendFile(absPath("uprobe_events"), p.probe)
}

func (p *Uprobe) Remove() error {
	_ = p.Disable()
	removeCmd := "-:" + p.name + "\n"
	err := fs.AppendFile(absPath("uprobe_events"), removeCmd)
	return err
}

func (p *Uprobe) Name() string {
	return p.name
}

func (p *Uprobe) Enable() error {
	path := fmt.Sprintf("events/uprobes/%s/enable", p.name)
	return fs.TruncFileWithString(absPath(path), "1")
}

func (p *Uprobe) Disable() error {
	path := fmt.Sprintf("events/uprobes/%s/enable", p.name)
	return fs.TruncFileWithString(absPath(path), "0")
}

func Events() ([][]byte, error) {
	content, err := ioutil.ReadFile(absPath("uprobe_events"))
	return bytes.Split(content, []byte("\n")), err
}

func RemoveAll() error {
	// this file won't exists if no uprobes
	if !fs.Exists(absPath("events/uprobes")) {
		return nil
	}
	// this will disable all uprobes
	if err := fs.TruncFileWithString(absPath("events/uprobes/enable"), "0"); err != nil {
		fmt.Printf("failed to disable all uprobes\n")
	}
	return fs.TruncFileWithString(absPath("uprobe_events"), "")
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

type handleEvent func(p *Uprobe, event string)

func UprobeStream(probe string, fn handleEvent, stop chan struct{}) error {
	p, err := NewUprobe(probe)
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

