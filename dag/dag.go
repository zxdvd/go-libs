package dag

import (
	"errors"
)

var ErrCircle = errors.New("same node revisited, detect circle")

type Node interface {
	// return a list of all children/parents Nodes
	Nexts() []Node
}

type Dag struct {
	nodes  []Node
	sorted []Node
	// 	starts []*Node
	// 	ends   [][]*Node
}

func (d *Dag) Add(n ...Node) *Dag {
	d.nodes = append(d.nodes, n...)
	return d
}

func (d *Dag) Nodes() []Node {
	return d.nodes
}

func (d *Dag) CircleDetect() error {
	d.sorted = make([]Node, 0, len(d.nodes))
	finished := map[Node]bool{}

	for _, n := range d.nodes {
		// skip if the node is already visited
		if _, ok := finished[n]; ok {
			continue
		}
		// create a map to record current visiting nodes
		visiting := map[Node]bool{}
		err := d.visit(n, finished, visiting)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dag) visit(n Node, finished map[Node]bool, visiting map[Node]bool) error {
	if _, ok := finished[n]; ok {
		return nil
	}
	if _, ok := visiting[n]; ok {
		return ErrCircle
	}
	visiting[n] = true
	for _, child := range n.Nexts() {
		err := d.visit(child, finished, visiting)
		if err != nil {
			return err
		}
	}
	finished[n] = true
	d.sorted = append(d.sorted, n)
	return nil
}

func (d *Dag) Iterate(fn func(n Node) (bool, error)) error {
	for _, n := range d.sorted {
		if shouldContinue, err := fn(n); !shouldContinue {
			return err
		}
	}
	return nil
}
