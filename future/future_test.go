package future

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestFuture(t *testing.T) {
	f := New()
	go func() {
		time.Sleep(1 * time.Second)
		err := f.SetResult(10)
		if err != nil {
			panic(err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(4)
	for i := 0; i < 4; i++ {
		log.Println("loop i", i)
		go func() {
			defer wg.Done()
			val, err := f.Get()
			if err != nil {
				panic(err)
			}
			fmt.Println("val", val)
		}()
	}
	wg.Wait()
}
