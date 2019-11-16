package uprobe

import (
	"testing"
	"fmt"
	"time"
)

func TestUprobe(t *testing.T) {
	probe := "p:test_bash /bin/bash:readline"
	done := make(chan struct{}, 1)
	go func() {
		time.Sleep(3 * time.Second)
		done <- struct{}{}
	}()
	fn := func(p *Uprobe, event string) {
		fmt.Println("event: ", event)
	}
	if err := UprobeStream(probe, fn, done); err != nil {
		panic(err)
	}
}


func TestRemoveAll(t *testing.T) {
	if err := RemoveAll(); err != nil {
		t.Fatal(err)
	}
}

