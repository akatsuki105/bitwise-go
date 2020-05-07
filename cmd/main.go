package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
)

func main() {
	os.Exit(Run())
}

func Run() int {
	flag.Parse()

	argc := flag.NArg()

	switch argc {
	case 0:
		useTUI()
	case 1:
		target := flag.Arg(0)
		useCLI(target)
	default:
		fmt.Fprintf(os.Stderr, "Error: A number of argments must be 0 or 1\n")
		return ExitCodeError
	}

	return ExitCodeOK
}
