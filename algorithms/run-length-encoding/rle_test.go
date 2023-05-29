package run_length_encoding

import "testing"

func TestEncode(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"AAAABBBCCD", "4A3B2C1D"},
		{"A", "1A"},
		{"AA", "2A"},
		{"AB", "1A1B"},
		{"ABC", "1A1B1C"},
		{"", ""},
	}
	for _, c := range cases {
		got := Encode(c.in)
		if got != c.want {
			t.Errorf("encode(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDecode(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"10A2B", "AAAAAAAAAABB"},
		{"4A3B2C1D", "AAAABBBCCD"},
		{"1A", "A"},
		{"2A", "AA"},
		{"1A1B", "AB"},
		{"1A1B1C", "ABC"},
		{"", ""},
	}
	for _, c := range cases {
		got := Decode(c.in)
		if got != c.want {
			t.Errorf("decode(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
