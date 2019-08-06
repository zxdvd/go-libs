package dag

import (
	"log"
	"testing"
)

func detectAndPanic(d *Dag) {
	if err := d.CircleDetect(); err != nil {
		panic(err)
	}
}

func TestDag(t *testing.T) {
	var a, b, c, d, e Node
	dag := &Dag{}
	dag.Adds(&a, &b, &c, &d, &e)
	// add self to make a circle
	a.AddParent(&a)
	err := dag.CircleDetect()
	if err != ErrCircle {
		panic("should have circle here")
	}
	a.removeLast()

	a.AddParent(&b)
	b.AddParent(&c)
	a.AddParent(&c)
	c.AddParent(&e)
	detectAndPanic(dag)
	e.AddParent(&a)
	err = dag.CircleDetect()
	if err != ErrCircle {
		panic("should have circle here")
	}
	e.removeLast()
	log.Printf("dag %v\n", dag)
}
