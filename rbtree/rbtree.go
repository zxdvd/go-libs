package rbtree

const (
	red = 1
	black = 2
)

type NodeData interface {
	Less(other NodeData) bool
	Equal(other NodeData) bool
}

type rbtNode struct {
	parent *rbtNode
	left *rbtNode
	right *rbtNode
	data NodeData
	color uint8
}

var NIL = &rbtNode{}

func NewNode(val NodeData) *rbtNode {
	return &rbtNode{
		parent: NIL,
		left: NIL,
		right: NIL,
		data: val,
		color: red,
	}
}

func (n *rbtNode) isLeft() bool {
	return n == n.parent.left
}

func (n *rbtNode) isRight() bool {
	return n == n.parent.right
}

func (n *rbtNode) Uncle() *rbtNode {
	if n.parent.isLeft() {
		return n.parent.right
	}
	return n.parent.left
}

type rbtTree struct {
	root *rbtNode
	size uint
}

func NewTree() *rbtTree {
	return &rbtTree{
		root: NIL,
	}
}

func (t *rbtTree) Find(val NodeData) *rbtNode {
	for cur := t.root; cur != nil && cur != NIL; {
		if cur.data.Equal(val) {
			return cur
		}
		if cur.data.Less(val) {
			cur = cur.right
		} else {
			cur = cur.left
		}
	}
	return nil
}

func (t *rbtTree) Insert(val NodeData) (newNode *rbtNode, inserted bool) {
	cur := t.root
	var parent *rbtNode
	for cur != NIL {
		if cur.data.Equal(val) {
			return cur, false
		}
		parent = cur
		if cur.data.Less(val) {
			cur = cur.right
		} else {
			cur = cur.left
		}
	}
	// begin to insert new node
	inserted = true
	newNode = NewNode(val)
	if parent != nil {
		newNode.parent = parent
		if parent.data.Less(val) {
			parent.right = newNode
		} else {
			parent.left = newNode
		}
	} else {
	// first inserted, set as root, root is black
		newNode.color = black
		t.root = newNode	
	}
	t.insertFixup(newNode)
	return newNode, inserted
}

func (t *rbtTree) insertFixup(cur *rbtNode) {
	// rbtree allow to add a red under a black
	if cur.parent.color == black {
		return
	}
	if cur == t.root {
		cur.color = black
		return
	}
	// now both current and parent are red,grandparent is black
	uncle := cur.Uncle()
	// if parent and uncle are both red, simply swap their color with grandparent
	// and continue fixup 
	if uncle.color == red {
		cur.parent.color = black
		uncle.color = black
		cur.parent.parent.color = red
		t.insertFixup(cur.parent.parent)
		return
	}
	// now parent is red and uncle is black or not exists
	// TODO need to deal with uncle not exists
	if cur.parent.isLeft() {
		if cur.isRight() {
			t.rotateLeft(cur)
			cur = cur.left
		}
		cur.parent.color = black
		cur.parent.parent.color = red
		t.rotateRight(cur.parent)
	} else {
		if cur.isLeft() {
			t.rotateRight(cur)
			cur = cur.right
		}
		cur.parent.color = black
		cur.parent.parent.color = red
		t.rotateLeft(cur.parent)
	}
}

func (t *rbtTree) rotateLeft(n *rbtNode) {
	return
}
func (t *rbtTree) rotateRight(n *rbtNode) {
	return
}