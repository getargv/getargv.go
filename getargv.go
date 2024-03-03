// getargv.go
// Copyright (C) 2024 Camden Narzt
//
// Distributed under terms of the BSD-3 license.
//go:build darwin && cgo

// Package getargv fetches the arguments of other processes in multiple formats
//
// The getargv package can only be used on macOS, because other operating
// systems have other means of accessing these arguments.
package getargv

/*
#cgo pkg-config: getargv
#include "libgetargv.h"
*/
import "C"
import "unsafe"

type (
	argv        = C.struct_ArgvResult
	argvArgc    = C.struct_ArgvArgcResult
	argvOptions = C.struct_GetArgvOptions
	pid_t       = uint
)

func p2i(p *C.char) uintptr {
	return uintptr(unsafe.Pointer(p))
}

func getArgv(pid pid_t, skip uint, nuls bool) (*argv, error) {
	r := new(argv)
	o := argvOptions{
		skip: C.uint(skip),
		pid:  C.pid_t(pid),
		nuls: C.bool(nuls),
	}
	success, err := C.get_argv_of_pid(&o, r)
	if success {
		return r, nil
	} else {
		return nil, err
	}
}

// AsBytes gets the arguments of pid as a slice of bytes, with skip leading arguments skipped,
// and NUL bytes replaced with spaces if nuls is true
// it can return an error if:
//   - the caller does not have permission to view the arguments of the target pid
//   - the target pid does not exist
//   - the kernel returns the targeted pid's args in an invalid format
//   - the targeted pid's args are too long (somehow longer than ARG_MAX) and cannot be parsed safely.
//   - AsBytes was asked to skip more args than targeted pid has.
func AsBytes(pid pid_t, skip uint, nuls bool) ([]byte, error) {
	a, err := getArgv(pid, skip, nuls)
	if err != nil {
		return nil, err
	}
	defer C.free_ArgvResult(a)
	var null *C.char = nil
	if null == a.start_pointer || null == a.end_pointer {
		// length calc below doesn't work with no args
		return []byte{}, nil
	}
	end := p2i(a.end_pointer)
	start := p2i(a.start_pointer)
	return C.GoBytes(unsafe.Pointer(a.start_pointer), C.int(end-start+1)), nil
}

// AsString gets the arguments of pid as a string, with skip leading arguments skipped,
// and NUL bytes replaced with spaces if nuls is true
// it can return an error if:
//   - the caller does not have permission to view the arguments of the target pid
//   - the target pid does not exist
//   - the kernel returns the targeted pid's args in an invalid format
//   - the targeted pid's args are too long (somehow longer than ARG_MAX) and cannot be parsed safely.
//   - AsString was asked to skip more args than targeted pid has.
func AsString(pid pid_t, skip uint, nuls bool) (string, error) {
	a, err := getArgv(pid, skip, nuls)
	if err != nil {
		return "", err
	}
	defer C.free_ArgvResult(a)
	var null *C.char = nil
	if null == a.start_pointer || null == a.end_pointer {
		// length calc below doesn't work with no args
		return "", nil
	}
	end := p2i(a.end_pointer)
	start := p2i(a.start_pointer)
	return C.GoStringN(a.start_pointer, C.int(end-start+1)), nil
}

// AsStrings gets the arguments of pid as a slice of strings, it can return an error if:
//   - the caller does not have permission to view the arguments of the target pid
//   - the target pid does not exist
//   - the kernel returns the targeted pid's args in an invalid format
//   - the targeted pid's args are too long (somehow longer than ARG_MAX) and cannot be parsed safely.
func AsStrings(pid pid_t) ([]string, error) {
	a := new(argvArgc)
	success, err := C.get_argv_and_argc_of_pid(C.pid_t(pid), a)
	if !success {
		return nil, err
	}

	defer C.free_ArgvArgcResult(a)

	s := unsafe.Slice(a.argv, a.argc)
	slice := make([]string, a.argc)
	if a.argc > 0 {
		for i, v := range s {
			if i < int(a.argc-1) {
				this := p2i(v)
				next := p2i(s[i+1])
				len := C.int(next - this - 1) // exclude trailing NUL
				slice[i] = C.GoStringN(v, len)
			} else {
				slice[i] = C.GoString(v)
			}
		}
	}

	return slice, nil
}
