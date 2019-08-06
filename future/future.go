package future

import (
	"errors"
	"reflect"
	"sync"
)

type state uint

const (
	PENDING state = 0 + iota
	RESOLVED
	REJECTED
)

type future struct {
	lock   sync.Mutex
	state  state
	value_ interface{}
	err_   error
	signal chan struct{}
}

func New() *future {
	return &future{
		state:  PENDING,
		signal: make(chan struct{}, 1),
	}
}

func NewN(n int) []*future {
	futures := make([]*future, n)
	for i, _ := range futures {
		futures[i] = New()
	}
	return futures
}

func (f *future) Get() (interface{}, error) {
	if f.state == RESOLVED {
		return f.value_, nil
	}
	if f.state == REJECTED {
		return nil, f.err_
	}
	<-f.signal
	return f.Get()
}

// resolve all futures, return err if get any error
func GetAll(futures []*future) ([]interface{}, error) {
	result := make([]interface{}, len(futures))
	chs := make([]reflect.SelectCase, len(futures))
	ids := make([]int, len(futures))
	for i, f := range futures {
		chs[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(f.signal),
		}
		ids[i] = i
	}
	for len(chs) > 0 {
		i, _, ok := reflect.Select(chs)
		if !ok {
			f := futures[ids[i]]
			if f.state == REJECTED {
				return nil, f.err_
			}
			result[ids[i]] = f.value_
			chs = append(chs[:i], chs[i+1:]...)
			ids = append(ids[:i], ids[i+1:]...)
		}
	}
	return result, nil
}

func (f *future) SetResult(val interface{}) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.state != PENDING {
		return errors.New("should not SetResult on resolved future")
	}
	f.state = RESOLVED
	f.value_ = val
	close(f.signal) // close to notify all waiters
	return nil
}

func (f *future) SetError(err error) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.state != PENDING {
		return errors.New("should not SetError on resolved future")
	}
	f.state = REJECTED
	f.err_ = err
	close(f.signal) // close to notify all waiters
	return nil
}
