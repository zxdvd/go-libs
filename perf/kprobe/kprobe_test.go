package kprobe

import (
	"testing"
	"fmt"

)

func TestAdd(t *testing.T) {
	p, err := NewKprobe("p:test_probe SyS_clone")
	fmt.Printf("%v %v\n", p, err)
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Add(); err != nil {
		t.Fatal(err)
	}
	events, err := Events()
	if err != nil {
		t.Fatal(err)
	}
	for _, evt := range events {
		fmt.Printf("existed event: %s\n", evt)
	}
	if err := p.Enable(); err != nil {
		fmt.Println("failed to enable event")
		t.Fatal(err)
	}

	readerCh, err := ReadEvents()
	if err != nil {
		fmt.Println("failed to read events")
		t.Fatal(err)
	}
	for s := range readerCh {
		fmt.Println("event: ", s)
	}

	if err := p.Remove(); err != nil {
		fmt.Println("failed to remove event")
		t.Fatal(err)
	}
	events, _ = Events()
	for _, evt := range events {
		fmt.Printf("existed event: %s\n", evt)
	}

}


func TestRemoveAll(t *testing.T) {
	if err := RemoveAll(); err != nil {
		t.Fatal(err)
	}
}
