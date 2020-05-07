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
}

func useTUI() {
	state := TUIState{
		mode:   mode32Bit,
		cursor: [2]int{0, 31},
		bits:   *newBitArray(),
		dec:    0,
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
				p1.Text = getP1Text(&state)
			}
			ui.Render(p0, p1)

		case "<Right>":
			switch state.mode {
			case modeBinary:
				state.cursor[1]--
				if state.cursor[1] < 0 {
					state.cursor[1] = 31
				}
				p1.Text = getP1Text(&state)
			}
			ui.Render(p0, p1)

		case "<Up>":
			switch state.mode {
			case modeBinary:
				bit := state.bits.buf[state.cursor[1]]
				if bit == 0 {
					bit = 1

					state.bits.buf[state.cursor[1]] = bit
					dec := state.bits.decimal()

					state.oct = strconv.FormatInt(dec, 8)
					state.dec = int(dec)
					state.hex = strconv.FormatInt(dec, 16)
					state.bits.update(dec)
					p0.Text = getP0Text(&state)
					p1.Text = getP1Text(&state)
				}
			}
			ui.Render(p0, p1)

		case "<Down>":
			switch state.mode {
			case modeBinary:
				bit := state.bits.buf[state.cursor[1]]
				if bit == 1 {
					bit = 0

					state.bits.buf[state.cursor[1]] = bit
					dec := state.bits.decimal()

					state.oct = strconv.FormatInt(dec, 8)
					state.dec = int(dec)
					state.hex = strconv.FormatInt(dec, 16)
					state.bits.update(dec)
					p0.Text = getP0Text(&state)
					p1.Text = getP1Text(&state)
				}
			}
			ui.Render(p0, p1)
		}
	}
}

func getP0Text(state *TUIState) string {
	header := "[Decimal:](fg:green)              [Hexdecimal:](fg:green)            [Octal:](fg:green)            "
	footer := fmt.Sprintf("%d", state.dec)
	length := len(footer)
	for i := 0; i < 22-length; i++ {
		footer += " "
	}
	footer += state.hex
	for i := 0; i < 23-len(state.hex); i++ {
		footer += " "
	}
	footer += state.oct
	text := header + "\n\n" + footer
	return text
}

func getP1Text(state *TUIState) string {
	header := state.bits.toString()
	footer := fmt.Sprintf("bit %d", state.cursor[1])
	text := header + "\n" + BinaryScale + "\n" + footer
	return text
}
