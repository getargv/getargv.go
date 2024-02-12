// getargv.go
// Copyright (C) 2024 Camden Narzt
//
// Distributed under terms of the BSD-3 license.
package getargv

/*
#cgo CFLAGS: -g
#cgo LDFLAGS: -lgetargv
#include "libgetargv.h"
*/
import "C"
import "unsafe"

type (
	Argv        = C.struct_ArgvResult
	ArgvArgc    = C.struct_ArgvArgcResult
	ArgvOptions = C.struct_GetArgvOptions
	pid_t       = uint
)

func p2i(p *C.char) uintptr {
	return uintptr(unsafe.Pointer(p))
}

func getArgv(pid pid_t, skip uint, nuls bool) (*Argv, error) {
	r := new(Argv)
	o := ArgvOptions{
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

func AsStrings(pid pid_t) ([]string, error) {
	a := new(ArgvArgc)
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
