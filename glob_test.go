package glob

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTree(t *testing.T) {
	patterns := []struct {
		s string
		i interface{}
	}{
		{"*abcd*ef*", 1},
		{"*.google.com", 2},
		{"http://example.com/books/*", 3},
		{"*://example.com/movies", 4},
		{`http://example.com/\*`, 5},
		{`http://example.com/*`, 6},
		{"你好*世界*", 7},
		{`foo\`, 8},
		{`b\ar`, 9},
	}
	data := []struct {
		s string
		v interface{}
	}{
		{"abcdef", 1},
		{"abcdefef", 1},
		{"abcabcdefgef", 1},
		{"google.com", nil},
		{"www.google.com", 2},
		{"http://example.com/books/", 3},
		{"http://example.com/", 6},
		{"http://example.com/*", 5},
		{"你好世界", 7},
		{"你你好世界", nil},
		{"你好世界世界界界", 7},
		{"你好,世界", 7},
		{"你好,世界。", 7},
		{`foo\`, nil},
		{`foo`, 8},
		{`b\ar`, nil},
		{`bar`, 9},
	}

	tr := &Trie{}
	for _, p := range patterns {
		tr.Add(p.s, p.i)
	}

	for _, data := range data {
		v, ok := tr.Lookup(data.s)
		if data.v == nil {
			assert.False(t, ok)
			assert.Nil(t, v)
		} else {
			assert.True(t, ok)
			assert.Equal(t, data.v, v)
		}
	}

}

func printSibling(node *node) {
	fmt.Printf("%s: ", node.s)
	for _, n := range node.child {
		fmt.Printf("%s ", n.s)
	}
	if node.wcard != nil {
		fmt.Printf("*-")
	}
	fmt.Println("")
	for _, n := range node.child {
		printSibling(n)
	}
	if node.wcard != nil {
		printSibling(node.wcard)
	}
}
