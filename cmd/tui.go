package main

import (
	"fmt"
	"log"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	BinaryScale = "    " + "[31 - 24](fg:cyan)" + "           " + "[23 - 16](fg:cyan)" + "           " + "[15 -  8](fg:cyan)" + "           " + "[ 7 -  0](fg:cyan)"
)

const (
	mode32Bit = iota
	modeBinary
)

type TUIState struct {
	mode     int
	cursor   [2]int
	bits     bitArray
	dec      int
	oct, hex string
	count    int
}

func useTUI() {
	state := TUIState{
		mode:   mode32Bit,
		cursor: [2]int{2, 31},
		bits:   *newBitArray(-1),
		dec:    0,
		oct:    "0",
		hex:    "0",
		count:  0,
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p0 := widgets.NewParagraph()
	p0.Title = "Number"
	p0.PaddingLeft = 1
	p0.Text = getP0Text(&state)
	p0.SetRect(0, 0, 64, 5)
	p0.BorderStyle.Fg = ui.ColorBlue

	p1 := widgets.NewParagraph()
	p1.Title = "Binary"
	p1.PaddingLeft = 1
	p1.Text = getP1Text(&state)
	p1.SetRect(0, 5, 73, 10)

	ui.Render(p0, p1)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>", "<Escape>":
			return

		case "<Up>", "k", "<Down>", "j":
			switch state.mode {
			case mode32Bit:
				state.mode = modeBinary
				p0.BorderStyle.Fg = ui.ColorWhite
				p1.BorderStyle.Fg = ui.ColorBlue
				state.bits.update(int64(state.dec), state.cursor[1])
				p0.Text = getP0Text(&state)
				p1.Text = getP1Text(&state)
			case modeBinary:
				state.mode = mode32Bit
				p0.BorderStyle.Fg = ui.ColorBlue
				p1.BorderStyle.Fg = ui.ColorWhite
				state.bits.update(int64(state.dec), -1)
				p0.Text = getP0Text(&state)
				p1.Text = getP1Text(&state)
			}
			ui.Render(p0, p1)

		case "<Left>", "h":
			switch state.mode {
			case mode32Bit:
				state.cursor[0]++
				if state.cursor[0] > 2 {
					state.cursor[0] = 0
				}
				p0.Text = getP0Text(&state)
			case modeBinary:
				state.cursor[1]++
				if state.cursor[1] > 31 {
					state.cursor[1] = 0
				}
				state.bits.update(int64(state.dec), state.cursor[1])
				p1.Text = getP1Text(&state)
			}
			ui.Render(p0, p1)

		case "<Right>", "l":
			switch state.mode {
			case mode32Bit:
				state.cursor[0]--
				if state.cursor[0] < 0 {
					state.cursor[0] = 2
				}
				p0.Text = getP0Text(&state)
			case modeBinary:
				state.cursor[1]--
				if state.cursor[1] < 0 {
					state.cursor[1] = 31
				}
				state.bits.update(int64(state.dec), state.cursor[1])
				p1.Text = getP1Text(&state)
			}
			ui.Render(p0, p1)

		case "<Space>":
			switch state.mode {
			case modeBinary:
				bit := state.bits.buf[state.cursor[1]]
				bit = (bit + 1) % 2
				state.bits.buf[state.cursor[1]] = bit
				dec := state.bits.decimal()

				state.oct = strconv.FormatInt(dec, 8)
				state.dec = int(dec)
				state.hex = strconv.FormatInt(dec, 16)
				state.bits.update(dec, state.cursor[1])
				p0.Text = getP0Text(&state)
				p1.Text = getP1Text(&state)
			}
			ui.Render(p0, p1)
		}

		if state.mode == mode32Bit {
			state.count++
			switch state.cursor[0] {
			case 2: // dec
				dec := int64(state.dec)
				switch e.ID {
				case "0":
					if state.dec > 0 {
						dec *= 10
					}
				case "1", "2", "3", "4", "5", "6", "7", "8", "9":
					num, _ := strconv.Atoi(e.ID)
					dec = dec*10 + int64(num)
				case "<C-<Backspace>>":
					dec /= 10
				}
				state.oct = strconv.FormatInt(dec, 8)
				state.dec = int(dec)
				state.hex = strconv.FormatInt(dec, 16)
				state.bits.update(dec, -1)
				p0.Text = getP0Text(&state)
				p1.Text = getP1Text(&state)
				ui.Render(p0, p1)
			case 1: // hex
				dec := int64(state.dec)
				switch e.ID {
				case "0":
					if state.dec > 0 {
						dec *= 16
					}
				case "1", "2", "3", "4", "5", "6", "7", "8", "9":
					num, _ := strconv.Atoi(e.ID)
					dec = dec*16 + int64(num)
				case "a", "b", "c", "d", "e", "f":
					r := []rune(e.ID)
					ascii := int(r[0])
					num := ascii - 87
					dec = dec*16 + int64(num)
				case "<C-<Backspace>>":
					dec /= 16
				}
				state.oct = strconv.FormatInt(dec, 8)
				state.dec = int(dec)
				state.hex = strconv.FormatInt(dec, 16)
				state.bits.update(dec, -1)
				p0.Text = getP0Text(&state)
				p1.Text = getP1Text(&state)
				ui.Render(p0, p1)
			case 0: // oct
				dec := int64(state.dec)
				switch e.ID {
				case "0":
					if state.dec > 0 {
						dec *= 8
					}
				case "1", "2", "3", "4", "5", "6", "7":
					num, _ := strconv.Atoi(e.ID)
					dec = dec*8 + int64(num)
				case "<C-<Backspace>>":
					dec /= 8
				}
				state.oct = strconv.FormatInt(dec, 8)
				state.dec = int(dec)
				state.hex = strconv.FormatInt(dec, 16)
				state.bits.update(dec, -1)
				p0.Text = getP0Text(&state)
				p1.Text = getP1Text(&state)
				ui.Render(p0, p1)
			}
		}
	}
}

func getP0Text(state *TUIState) string {
	header := "[Decimal:              Hexdecimal:           Octal:           ](fg:green)"
	footer := fmt.Sprintf("%d", state.dec)
	length := len(footer)
	if state.mode == mode32Bit && state.cursor[0] == 2 {
		footer = fmt.Sprintf("[%d](fg:blue)", state.dec)
	}
	for i := 0; i < 22-length; i++ {
		footer += " "
	}

	if state.mode == mode32Bit && state.cursor[0] == 1 {
		footer += fmt.Sprintf("[%s](fg:blue)", state.hex)
	} else {
		footer += state.hex
	}
	for i := 0; i < 22-len(state.hex); i++ {
		footer += " "
	}

	if state.mode == mode32Bit && state.cursor[0] == 0 {
		footer += fmt.Sprintf("[%s](fg:blue)", state.oct)
	} else {
		footer += state.oct
	}
	text := header + "\n \n" + footer
	return text
}

func getP1Text(state *TUIState) string {
	header := state.bits.toString()
	footer := fmt.Sprintf("bit %d", state.cursor[1])
	text := header + "\n" + BinaryScale + "\n" + footer
	return text
}
