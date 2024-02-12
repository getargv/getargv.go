package getargv

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"syscall"
	"testing"
)

// TestAsString calls getargv.AsString with test process' pid, and nuls and skip fully exercised
func TestAsString(t *testing.T) {
	args := os.Args
	pid := uint(os.Getpid())
	for _, v := range [2]bool{true, false} {
		for i := 1; i <= len(args); i++ { // in go 1.22 use: for i := range len(args) {
			t.Run(fmt.Sprintf("skip=%d, nuls=%t", i, v), func(t *testing.T) {
				var sep string
				if v {
					sep = " "
				} else {
					sep = "\x00"
				}
				want := strings.Join(args[i:], sep)
				if i < len(args) {
					want += "\x00"
				}
				actual, err := AsString(pid, uint(i), v)
				if want != actual || err != nil {
					t.Fatalf(`%q, %v, should have been: %#q, nil`, actual, err, want)
				}
			})
		}
	}
}

// TestAsBytes calls getargv.AsBytes with test process' pid, and nuls and skip fully exercised
func TestAsBytes(t *testing.T) {
	args := os.Args
	pid := uint(os.Getpid())
	for _, v := range [2]bool{true, false} {
		for i := 1; i <= len(args); i++ { // in go 1.22 use: for i := range len(args) {
			t.Run(fmt.Sprintf("skip=%d, nuls=%t", i, v), func(t *testing.T) {
				var sep string
				if v {
					sep = " "
				} else {
					sep = "\x00"
				}
				want := []byte(strings.Join(args[i:], sep))
				if i < len(args) {
					want = append(want, 0)
				}
				actual, err := AsBytes(pid, uint(i), v)
				if !bytes.Equal(want, actual) || err != nil {
					t.Fatalf(`%q, %v, should have been: %#q, nil`, actual, err, want)
				}
			})
		}
	}
}

// TestAsStrings calls getargv.AsStrings with test process' pid,
// checking for a valid return value.
func TestAsStrings(t *testing.T) {
	want := os.Args
	pid := uint(os.Getpid())
	actual, err := AsStrings(pid)
	if !slices.Equal(want, actual) || err != nil {
		t.Fatalf(`AsStrings(%d) = %q, %v, should have been: %#q, nil`, pid, actual, err, want)
	}
}

type failureCase struct {
	pid uint
	err error
}

// TestAsStrings calls getargv.AsStrings with various pids,
// checking for correct errors.
func TestAsStringsFailure(t *testing.T) {
	for _, v := range []failureCase{{0,syscall.EPERM},{99999+1,syscall.ESRCH}} {
		pid := v.pid
		want := v.err
		t.Run(fmt.Sprintf("pid = %d, err = %v",pid,want), func(t *testing.T) {
		actual, err := AsStrings(pid)
		if !errors.Is(err, want) {
			t.Fatalf(`AsStrings(%d) = %q, %+v, should have been: [], %+v`, pid, actual, err, want)
		}
	})
	}
}
