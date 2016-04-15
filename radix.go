package radix

type Trie struct {
	root *node
	dump bool
}

type noneType int

const vNONE noneType = 0

type node struct {
	child    []*node
	childidx []byte // first byte of each child
	s        string
	v        interface{}
}

// NewTrie creates a new trie.
func NewTrie(dump bool) *Trie {
	return &Trie{
		dump: dump,
	}
}

func dump(s string) string {
	ss := make([]byte, len(s))
	copy(ss, s)
	return string(ss)
}

func (t *Trie) newSubTree(s string, v interface{}) *node {
	n := &node{
		v: v,
	}
	t.setS(n, s)
	return n
}

func (t *Trie) setS(n *node, s string) {
	if t.dump {
		n.s = dump(s)
	} else {
		n.s = s
	}
}

func (n *node) setV(v interface{}) (ov interface{}, ok bool) {
	ov = n.v
	if _, ok = ov.(noneType); ok {
		ov, ok = nil, false
	}
	n.v = v
	return
}

// Add inserts a key-value pair into trie. If there is an old value for the
// key, old value will be returned and 'has' will be true.
func (t *Trie) Add(s string, v interface{}) (ov interface{}, has bool) {
	if s == "" {
		return
	}
	if t.root == nil {
		t.root = t.newSubTree(s, v)
		return
	}
	n := t.root
INSERT:
	for {
		var l int
		min := len(s)
		if min > len(n.s) {
			min = len(n.s)
		}
		for ; l < min; l++ {
			if s[l] != n.s[l] {
				break
			}
		}
		switch {
		case l == len(n.s): // totally match this node
			s = s[l:]
			if len(s) == 0 { // end
				return n.setV(v)
			}
			first := 0
			if len(s[first:]) > 0 {
				for i := 0; i < len(n.childidx); i++ {
					if n.childidx[i] == s[first] {
						n = n.child[i]
						continue INSERT
					}
				}
			}
		default: // split
			prefix, suffix := n.s[:l], n.s[l:]
			child := &node{
				child:    n.child,
				childidx: n.childidx,
				v:        n.v,
			}
			t.setS(child, suffix)
			*n = node{}
			t.setS(n, prefix)
			n.child = []*node{child}
			n.childidx = []byte{child.s[0]}
			n.v = vNONE
			s = s[l:]
			if len(s) == 0 { // end
				return n.setV(v)
			}
		}
		// construct a new subtree using rest of pattern and
		// append it to the child list of this node
		child := t.newSubTree(s, v)
		n.child = append(n.child, child)
		n.childidx = append(n.childidx, child.s[0])
		return
	}
}

// Lookup searchs the trie to find an exact match and returns the
// associated value.  If not found, ok will be false.
func (t *Trie) Lookup(s string) (v interface{}, ok bool) {
	n := t.root.lookup(s)
	if n == nil {
		v, ok = nil, false
	} else {
		v, ok = n.v, true
		if _, ok = v.(noneType); ok {
			v, ok = nil, false
		}
	}
	return
}

func (n *node) lookup(s string) *node {
	if n == nil {
		return nil
	}

	minLen := len(s)
	if minLen > len(n.s) {
		minLen = len(n.s)
	}
	var l int // length of longest common prefix
	for l = 0; l < minLen && s[l] == n.s[l]; l++ {
	} // at the end of loop: pattern[:l] == n.s[:l]
	switch l {
	case len(n.s): // totally match this node
		s = s[l:]
		if len(s) == 0 { // end
			return n
		}
		// go down
		var k int
		for k = 0; k < len(n.childidx); k++ {
			if n.childidx[k] == s[0] {
				if end := n.child[k].lookup(s); end != nil {
					return end
				}
				break
			}
		}
		fallthrough
	default:
		return nil
	}
}
