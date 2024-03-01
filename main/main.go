package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strconv"
	"github.com/getargv/getargv.go"
)

func main() {
	var skip = flag.Uint("s", 0, "number of leading args to skip")
	var nuls = flag.Bool("0", false, "convert nuls to spaces")
	flag.Parse()
	if (flag.NArg() != 1) {
		fmt.Println("a single pid must be provided as an argument")
		return
	}
	var err error
	pid, err := strconv.ParseUint(flag.Arg(0), 10, 32)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return
	}
	byts, err := getargv.AsBytes(uint(pid), *skip, *nuls)
	err = binary.Write(os.Stdout, binary.NativeEndian, byts)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return
	}
}
