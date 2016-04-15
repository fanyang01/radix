package radix

type Pattern struct {
	trie PatternTrie
}

// Compile compiles several alternative patterns into one.
func Compile(patterns ...string) *Pattern {
	p := &Pattern{PatternTrie{}}
	for _, pattern := range patterns {
		p.trie.Add(pattern, struct{}{})
	}
	return p
}

// Match tests whether s matches any patterns in p.
func (p *Pattern) Match(s string) bool {
	_, ok := p.trie.Lookup(s)
	return ok
}

// Match tests whether s matches pattern.
func Match(pattern, s string) bool {
	return Compile(pattern).Match(s)
}
