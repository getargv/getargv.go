//go:build darwin && cgo

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

// TestAsString calls getargv.AsString with pid, nuls, and skip fully exercised
func TestAsString(t *testing.T) {
	args := os.Args
	pid := uint(os.Getpid())
	for _, v := range [2]bool{true, false} {
		for i := range len(args) + 1 {
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

type asStringFailureCase struct {
	pid  pid_t
	skip uint
	nuls bool
	want error
}

// TestAsStringFailure calls getargv.AsString with pid, nuls, and skip set to cause expected errors
func TestAsStringFailure(t *testing.T) {
	this := uint(os.Getpid())
	num := uint(len(os.Args) + 1)

	for _, v := range []asStringFailureCase{
		{0, 0, true, syscall.EPERM},
		{99999 + 1, 0, true, syscall.ESRCH},
		{this, num, true, syscall.ERANGE},
	} {
		pid := v.pid
		skip := v.skip
		nuls := v.nuls
		want := v.want
		t.Run(fmt.Sprintf("pid=%d, skip=%d, nuls=%t", pid, skip, nuls), func(t *testing.T) {
			actual, err := AsString(pid, skip, nuls)
			if !errors.Is(err, want) {
				t.Fatalf(`%q, %v, should have been: %#q, nil`, actual, err, want)
			}
		})
	}
}

// TestAsBytes calls getargv.AsBytes with pid, nuls, and skip fully exercised
func TestAsBytes(t *testing.T) {
	args := os.Args
	pid := uint(os.Getpid())
	for _, nuls := range [2]bool{true, false} {
		for skip := range len(args)+1 {
			t.Run(fmt.Sprintf("skip=%d, nuls=%t", skip, nuls), func(t *testing.T) {
				var sep string
				if nuls {
					sep = " "
				} else {
					sep = "\x00"
				}
				want := []byte(strings.Join(args[skip:], sep))
				if skip < len(args) {
					want = append(want, 0)
				}
				actual, err := AsBytes(pid, uint(skip), nuls)
				if !bytes.Equal(want, actual) || err != nil {
					t.Fatalf(`%q, %v, should have been: %#q, nil`, actual, err, want)
				}
			})
		}
	}
}

type asBytesFailureCase struct {
	pid  pid_t
	skip uint
	nuls bool
	want error
}

// TestAsBytesFailure calls getargv.AsBytes with pid, nuls, and skip set to cause expected errors
func TestAsBytesFailure(t *testing.T) {
	this := uint(os.Getpid())
	num := uint(len(os.Args) + 1)

	for _, v := range []asBytesFailureCase{
		{0, 0, true, syscall.EPERM},
		{99999 + 1, 0, true, syscall.ESRCH},
		{this, num, true, syscall.ERANGE},
	} {
		pid := v.pid
		skip := v.skip
		nuls := v.nuls
		want := v.want
		t.Run(fmt.Sprintf("pid=%d, skip=%d, nuls=%t", pid, skip, nuls), func(t *testing.T) {
			actual, err := AsBytes(pid, skip, nuls)
			if !errors.Is(err, want) {
				t.Fatalf(`%q, %v, should have been: %#q, nil`, actual, err, want)
			}
		})
	}
}

// TestAsStrings calls getargv.AsStrings with test process' pid, checking for a valid return value.
func TestAsStrings(t *testing.T) {
	want := os.Args
	pid := uint(os.Getpid())
	actual, err := AsStrings(pid)
	if !slices.Equal(want, actual) || err != nil {
		t.Fatalf(`AsStrings(%d) = %q, %v, should have been: %#q, nil`, pid, actual, err, want)
	}
}

type asStringsFailureCase struct {
	pid uint
	err error
}

// TestAsStringsFailure calls getargv.AsStrings with various pids, checking for correct errors.
func TestAsStringsFailure(t *testing.T) {
	for _, v := range []asStringsFailureCase{{0, syscall.EPERM}, {99999 + 1, syscall.ESRCH}} {
		pid := v.pid
		want := v.err
		t.Run(fmt.Sprintf("pid = %d", pid), func(t *testing.T) {
			actual, err := AsStrings(pid)
			if !errors.Is(err, want) {
				t.Fatalf(`%q, %+v, should have been: [], %+v`, actual, err, want)
			}
		})
	}
}
