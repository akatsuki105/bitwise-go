package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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

func useTUI() {
	const (
		mode32Bit = iota
		modeBinary
	)
	mode := mode32Bit

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p0 := widgets.NewParagraph()
	p0.Title = "32bit"
	p0.Text = "xxxx xxxx"
	p0.SetRect(0, 0, 60, 5)
	p0.BorderStyle.Fg = ui.ColorBlue

	p1 := widgets.NewParagraph()
	p1.Title = "Binary"
	p1.Text = "Simple colored text\nwith label. It [can be](fg:red) multilined with \\n or [break automatically](fg:red,fg:bold)"
	p1.SetRect(0, 5, 60, 10)

	ui.Render(p0, p1)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Tab>":
			switch mode {
			case mode32Bit:
				mode = modeBinary
				p0.BorderStyle.Fg = ui.ColorWhite
				p1.BorderStyle.Fg = ui.ColorBlue
			case modeBinary:
				mode = mode32Bit
				p0.BorderStyle.Fg = ui.ColorBlue
				p1.BorderStyle.Fg = ui.ColorWhite
			}
			ui.Render(p0, p1)
		}
	}
}

func useCLI(target string) {
	var i int64

	switch {
	case strings.HasPrefix(target, "0x"):
		tmp, err := strconv.ParseInt(target[2:], 16, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid hex value.")
			return
		}
		i = tmp
	case strings.HasPrefix(target, "0b"):
		target = strings.Replace(target, "_", "", -1)
		tmp, err := strconv.ParseInt(target[2:], 2, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid binary value.")
			return
		}
		i = tmp
	case strings.HasPrefix(target, "0"):
		tmp, err := strconv.ParseInt(target[1:], 8, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid octal value.")
			return
		}
		i = tmp
	default:
		tmp, err := strconv.ParseInt(target, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid decimal value.")
			return
		}
		i = tmp
	}

	binTmp := strconv.FormatInt(i, 2)
	var padding int
	if len(binTmp)%4 == 0 {
		padding = 0
	} else {
		padding = 4 - (len(binTmp) % 4)
	}

	bin := ""
	bin += strings.Repeat("0", padding)
	for index, char := range binTmp {
		bin += string(char)
		if (padding+index)%4 == 3 && index+1 != len(binTmp) {
			bin += " "
		}
	}

	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("bin[2]  :"), color.CyanString(bin))
	oct := strconv.FormatInt(i, 8)
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("oct[8]  :"), color.CyanString(oct))
	dec := strconv.FormatInt(i, 10)
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("dec[10] :"), color.CyanString(dec))
	hex := strconv.FormatInt(i, 16)
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("hex[16] :"), color.CyanString(hex))
}
