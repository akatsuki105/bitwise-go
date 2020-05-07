package main

import "strconv"

type bitArray struct {
	buf [32]byte
	str string
}

func newBitArray() *bitArray {
	b := &bitArray{}
	b.update(0)
	return b
}

func (b *bitArray) decimal() int64 {
	decimal := int64(0)

	for i := 0; i < 32; i++ {
		decimal += int64(b.buf[i]) << i
	}

	return decimal
}

func (b *bitArray) update(decimal int64) {
	bin := strconv.FormatInt(decimal, 2)
	for i := len(bin); i < 32; i++ {
		bin = "0" + bin
	}
	for i := 0; i < 32; i++ {
		switch bin[31-i] {
		case '0':
			b.buf[i] = 0
		case '1':
			b.buf[i] = 1
		}
	}

	result := ""
	for i := 0; i < 32; i++ {
		result += strconv.Itoa(int(b.buf[31-i]))

		if i == 31 {
			break
		}
		result += " "

		if i > 0 && i%8 == 7 {
			result += "| "
		}
	}
	b.str = result
}

func (b *bitArray) toString() string {
	return b.str
}
