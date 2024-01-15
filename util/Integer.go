package util

import (
	"fmt"
	"math"
	"strings"
)

func ParseInt(number string) int {
	if len(number) <= 0 {
		return 0
	}

	for number[0] == '0' && len(number) > 1 {
		number = number[1:]
	}

	sig := 0          // 0 := + , 1 := -
	var value int = 0 // max 64 bit signed

	if number[0] == '+' {
		sig = 0
		number = number[1:]
	} else if number[0] == '-' {
		sig = 1
		number = number[1:]
	}

	for number[0] == '0' && len(number) > 1 {
		number = number[1:]
	}

	for i := 0; i < len(number); i++ {
		charint := int(number[i] - '0')
		value = (value * 10) + charint
	}

	if sig == 0 {
		return (value & math.MaxInt)
	} else {
		return (-value | math.MinInt)
	}
}

func ParseHex(number string) int {
	sig := 0          // 0 := + , 1 := -
	var value int = 0 // max 64 bit signed

	if number[0] == '+' {
		sig = 0
		number = number[1:]
	} else if number[0] == '-' {
		sig = 1
		number = number[1:]
	}

	if (strings.Compare(number[0:2], "0x") == 0) || (strings.Compare(number[0:2], "0X") == 0) {
		number = number[2:]
	} else {
		panic(fmt.Sprintf("invalid number format %s!!!", number))
	}

	for i := 0; i < len(number); i++ {
		charint := valueOf(number[i])
		value = (value * 16) + charint

		if FixedBitCount(value) >= 63 {
			break
		}
	}

	if sig == 0 {
		return (value & math.MaxInt)
	} else {
		return (-value | math.MinInt)
	}
}

func ParseOct(number string) int {
	sig := 0          // 0 := + , 1 := -
	var value int = 0 // max 64 bit signed

	if number[0] == '+' {
		sig = 0
		number = number[1:]
	} else if number[0] == '-' {
		sig = 1
		number = number[1:]
	}

	if (strings.Compare(number[0:2], "0o") == 0) || (strings.Compare(number[0:2], "0O") == 0) {
		number = number[2:]
	} else {
		panic(fmt.Sprintf("invalid number format %s!!!", number))
	}

	for i := 0; i < len(number); i++ {
		charint := valueOf(number[i])
		value = (value * 8) + charint

		if FixedBitCount(value) >= 63 {
			break
		}
	}

	if sig == 0 {
		return (value & math.MaxInt)
	} else {
		return (-value | math.MinInt)
	}
}

func ParseBin(number string) int {
	sig := 0          // 0 := + , 1 := -
	var value int = 0 // max 64 bit signed

	if number[0] == '+' {
		sig = 0
		number = number[1:]
	} else if number[0] == '-' {
		sig = 1
		number = number[1:]
	}

	if (strings.Compare(number[0:2], "0b") == 0) || (strings.Compare(number[0:2], "0B") == 0) {
		number = number[2:]
	} else {
		panic(fmt.Sprintf("invalid number format %s!!!", number))
	}

	for i := 0; i < len(number); i++ {
		charint := valueOf(number[i])
		value = (value * 2) + charint

		if FixedBitCount(value) >= 63 {
			break
		}
	}

	if sig == 0 {
		return (value & math.MaxInt)
	} else {
		return (-value | math.MinInt)
	}
}

func Parse(number string, base int) int {
	switch base {
	case 10:
		return ParseInt(number)
	case 16:
		return ParseHex(number)
	case 8:
		return ParseOct(number)
	case 2:
		return ParseBin(number)
	}
	panic(fmt.Sprintf("invalid integer base %d!!!", base))
}
