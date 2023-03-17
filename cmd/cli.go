package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func useCLI(target string) error {
	var decimal int64

	switch {
	case strings.HasPrefix(target, "0x"):
		target = strings.Replace(target, "_", "", -1)
		tmp, err := strconv.ParseInt(target[2:], 16, 64)
		if err != nil {
			return fmt.Errorf("parse error: invalid hexdecimal value")
		}
		decimal = tmp
	case strings.HasPrefix(target, "0b") || strings.HasPrefix(target, "b") || strings.HasPrefix(target, "%"):
		target = strings.Replace(target, "_", "", -1)
		tmp, err := strconv.ParseInt(target[2:], 2, 64)
		if err != nil {
			return fmt.Errorf("parse error: invalid binary value")
		}
		decimal = tmp
	case strings.HasPrefix(target, "0"):
		target = strings.Replace(target, "_", "", -1)
		tmp, err := strconv.ParseInt(target[1:], 8, 64)
		if err != nil {
			return fmt.Errorf("parse error: invalid octal value")
		}
		decimal = tmp
	default:
		target = strings.Replace(target, "_", "", -1)
		tmp, err := strconv.ParseInt(target, 10, 64)
		if err != nil {
			return fmt.Errorf("parse error: invalid decimal value")
		}
		decimal = tmp
	}

	bin, oct, dec, hex := toString(decimal)
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("bin[2]  :"), color.CyanString(bin))
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("oct[8]  :"), color.CyanString(oct))
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("dec[10] :"), color.CyanString(dec))
	fmt.Fprintf(color.Output, "%s %s\n", color.GreenString("hex[16] :"), color.CyanString(hex))
	return nil
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
