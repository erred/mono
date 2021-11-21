package main

import "testing"

func TestCanonicalPath(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"foo/bar.md", "/foo/bar/"},
		{"foo/bar.html", "/foo/bar/"},
		{"foo/index.md", "/foo/"},
		{"index.html", "/"},
		{"other.html", "/other/"},
	}
	for _, c := range cases {
		got := canonicalPath(c.in)
		if got != c.out {
			t.Errorf("canonicalPath(%s)=%s expected=%s", c.in, got, c.out)
		}
	}
}
