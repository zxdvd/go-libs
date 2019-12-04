package dag

import (
	"testing"
)

type node struct {
	children []Node
}

func (n *node) Nexts() []Node {
	return n.children
}

func detectAndPanic(d *Dag) {
	if err := d.CircleDetect(); err != nil {
		panic(err)
	}
}

func TestDag(t *testing.T) {
	a := &node{}
	b := &node{}
	dag := &Dag{}
	dag.Add(a, b)
	detectAndPanic(dag)
	// add self to make a circle
	a.children = []Node{a}
	err := dag.CircleDetect()
	if err != ErrCircle {
		panic("should have circle here")
	}
}
