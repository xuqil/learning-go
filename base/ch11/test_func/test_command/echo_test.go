package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEcho(t *testing.T) {
	testCases := []struct {
		newline bool
		sep     string
		args    []string

		want string
	}{
		{true, "", []string{}, "\n"},
		{false, "", []string{}, ""},
		{true, "\t", []string{"one", "two", "three"}, "one\ttwo\tthree\n"},
		{true, ",", []string{"a", "b", "c"}, "a,b,c\n"},
		{false, ":", []string{"1", "2", "3"}, "1:2:3"},
	}
	for _, tc := range testCases {
		descr := fmt.Sprintf("echo(%v, %q, %q)",
			tc.newline, tc.sep, tc.args)

		// 增加了一个全局名为 out 的变量来替代直接使用 os.Stdout
		out = new(bytes.Buffer)
		if err := echo(tc.newline, tc.sep, tc.args); err != nil {
			t.Errorf("%s failed: %v", descr, err)
			continue
		}
		got := out.(*bytes.Buffer).String()
		if got != tc.want {
			t.Errorf("%s = %q, want %q", descr, got, tc.want)
		}
	}
}
