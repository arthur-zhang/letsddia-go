# Rolling Hash (Rabin-Karp Algorithm)

implementation based on 

- [wikipedia Rabin–Karp_algorithm](https://en.wikipedia.org/wiki/Rabin%E2%80%93Karp_algorithm)
- [Rolling Hash (Rabin-Karp Algorithm)](http://courses.csail.mit.edu/6.006/spring11/rec/rec06.pdf)

# Usage

```go
func TestRabinKarp(t *testing.T) {
	tests := []struct {
		s       string
		pattern string
		wantIdx int
		wantOk  bool
	}{
		{"hello world", "world", 6, true},
		{"hello world", "world!", 0, false},
		{"hello world", "o", 4, true},
		{"hello world", "hello world", 0, true},
		{"hello world", "hello world!", 0, false},
		{"", "", 0, true},
		{"", "a", 0, false},
		{"a", "", 0, true},
		{"abc", "abc", 0, true},
		{"abcabc", "abc", 0, true},
		{"abcabc", "acb", 0, false},
	}

	for _, test := range tests {
		gotIdx, gotOk := RabinKarpSearch(test.s, test.pattern)
		assert.Equal(t, test.wantIdx, gotIdx)
		assert.Equal(t, test.wantOk, gotOk)
	}
}
```
