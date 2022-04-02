package flagwrap

import (
	"flag"
	"testing"
)

func TestParse(t *testing.T) {
	t.Setenv("C", "calendar")
	t.Setenv("D", "daylight")
	t.Setenv("e", "eternity")

	tcs := []struct {
		flagName, defaultValue, help string
		envName, envValue, expected  string
		val                          *string
	}{
		{
			"a", "aaa", "set from default",
			"", "", "aaa", nil,
		}, {
			"b", "bbb", "set from args",
			"", "", "bar", nil,
		}, {
			"c", "ccc", "set from env",
			"C", "calendar", "calendar", nil,
		}, {
			"d", "ddd", "flags over env",
			"D", "daylight", "doggo", nil,
		}, {
			"e", "eee", "case sensitivity",
			"e", "eternity", "eee", nil,
		}, {
			"f.g", "fff", "dots",
			"F_G", "ggg", "ggg", nil,
		},
	}

	fs := flag.NewFlagSet("test01", flag.ContinueOnError)
	for i := range tcs {
		tc := tcs[i]
		tc.val = fs.String(tc.flagName, tc.defaultValue, tc.help)
		tcs[i] = tc
		if tc.envName != "" {
			t.Setenv(tc.envName, tc.envValue)
		}
	}
	err := Parse(fs, []string{"-b=bar", "-d=doggo"})
	if err != nil {
		t.Fatal("parse error:", err)
		return
	}
	for i, tc := range tcs {
		if tc.expected != *tc.val {
			t.Logf("%d = %q, want %q", i, *tc.val, tc.expected)
			t.Fail()
		}
	}
}
