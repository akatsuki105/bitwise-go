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

const (
	mode32Bit = iota
	modeBinary
)

const (
	BinaryScale = "    " + "[31 - 24](fg:cyan)" + "           " + "[23 - 16](fg:cyan)" + "           " + "[15 -  8](fg:cyan)" + "           " + "[ 7 -  0](fg:cyan)"
)

type TUIState struct {
	mode          int
	cursor        [2]int
	dec           int
	bin, oct, hex string
}

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
	state := TUIState{
		mode:   mode32Bit,
		cursor: [2]int{0, 31},
		dec:    0,
		bin:    dec2BitArray(0),
		oct:    "0",
		hex:    "0",
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p0 := widgets.NewParagraph()
	p0.Title = "Number"
	p0.PaddingLeft = 1
	header := "[Decimal:](fg:green)            [Hexdecimal:](fg:green)            [Octal:](fg:green)            "
	p0.Text = header + "\n\n"
	p0.SetRect(0, 0, 73, 5)
	p0.BorderStyle.Fg = ui.ColorBlue

	p1 := widgets.NewParagraph()
	p1.Title = "Binary"
	p1.PaddingLeft = 1
	{
		header := state.bin
		footer := fmt.Sprintf("bit %d", state.cursor[1])
		p1.Text = header + "\n" + BinaryScale + "\n" + footer
	}
	p1.SetRect(0, 5, 73, 10)

	ui.Render(p0, p1)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {

		case "q", "<C-c>":
			return

		case "<Tab>":
			switch state.mode {
			case mode32Bit:
				state.mode = modeBinary
				p0.BorderStyle.Fg = ui.ColorWhite
				p1.BorderStyle.Fg = ui.ColorBlue
			case modeBinary:
				state.mode = mode32Bit
				p0.BorderStyle.Fg = ui.ColorBlue
				p1.BorderStyle.Fg = ui.ColorWhite
			}
			ui.Render(p0, p1)

		case "<Left>":
			switch state.mode {
			case modeBinary:
				state.cursor[1]++
				if state.cursor[1] > 31 {
					state.cursor[1] = 0
				}
				header := state.bin
				footer := fmt.Sprintf("bit %d", state.cursor[1])
				p1.Text = header + "\n" + BinaryScale + "\n" + footer
			}
			ui.Render(p0, p1)

		case "<Right>":
			switch state.mode {
			case modeBinary:
				state.cursor[1]--
				if state.cursor[1] < 0 {
					state.cursor[1] = 31
				}
				header := state.bin
				footer := fmt.Sprintf("bit %d", state.cursor[1])
				p1.Text = header + "\n" + BinaryScale + "\n" + footer
			}
			ui.Render(p0, p1)
		}
	}
}

func useCLI(target string) {
	var decimal int64

	switch {
	case strings.HasPrefix(target, "0x"):
		tmp, err := strconv.ParseInt(target[2:], 16, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid hex value.")
			return
		}
		decimal = tmp
	case strings.HasPrefix(target, "0b"):
		target = strings.Replace(target, "_", "", -1)
		tmp, err := strconv.ParseInt(target[2:], 2, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid binary value.")
			return
		}
		decimal = tmp
	case strings.HasPrefix(target, "0"):
		tmp, err := strconv.ParseInt(target[1:], 8, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid octal value.")
			return
		}
		decimal = tmp
	default:
		tmp, err := strconv.ParseInt(target, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid decimal value.")
			return
		}
		decimal = tmp
	}

	bin, oct, dec, hex := toString(decimal)
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("bin[2]  :"), color.CyanString(bin))
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("oct[8]  :"), color.CyanString(oct))
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("dec[10] :"), color.CyanString(dec))
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("hex[16] :"), color.CyanString(hex))
}

func toString(decimal int64) (bin, oct, dec, hex string) {
	binTmp := strconv.FormatInt(decimal, 2)
	var padding int
	if len(binTmp)%4 == 0 {
		padding = 0
	} else {
		padding = 4 - (len(binTmp) % 4)
	}

	bin = ""
	bin += strings.Repeat("0", padding)
	for index, char := range binTmp {
		bin += string(char)
		if (padding+index)%4 == 3 && index+1 != len(binTmp) {
			bin += " "
		}
	}

	oct = strconv.FormatInt(decimal, 8)
	dec = strconv.FormatInt(decimal, 10)
	hex = strconv.FormatInt(decimal, 16)
	return bin, oct, dec, hex
}

func dec2BitArray(decimal int64) string {
	bin := strconv.FormatInt(decimal, 2)
	for i := len(bin); i <= 32; i++ {
		bin = "0" + bin
	}

	result := ""
	for i := 0; i < 32; i++ {
		result += string(bin[i])

		if i == 31 {
			break
		}
		result += " "

		if i > 0 && i%8 == 7 {
			result += "| "
		}
	}
	return result
}
