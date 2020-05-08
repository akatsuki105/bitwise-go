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

		switch target {
		case "v", "version":
			useVersion()
		case "h", "help", "Help":
			useHelp()
		default:
			if err := useCLI(target); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return ExitCodeError
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "command error: a number of argments must be 0 or 1\n")
		return ExitCodeError
	}

	return ExitCodeOK
}

func useHelp() {
	fmt.Println(`bitwise Help: `)
	fmt.Println(`
CLI: 
$ bitwise 0b0100 // binary
$ bitwise 0777	 // octal
$ bitwise 1234   // decimal
$ bitwise 0xff   // hexdecimal`)
	fmt.Print(`
TUI: 
$ bitwise`)
}

func useVersion() {
	fmt.Print(`bitwise
Version: v0.0.0`)
}
