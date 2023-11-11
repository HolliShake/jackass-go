package jutil
import (
	"fmt"
	"math"
	"strings"
)

func valueOf(r byte) int {
	switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return int(r - '0')
		case 'a', 'A':
			return 10
		case 'b', 'B':
			return 11
		case 'c', 'C':
			return 12
		case 'd', 'D':
			return 13
		case 'e', 'E':
			return 14
		case 'f', 'F':
			return 15
		default:
			panic(fmt.Sprintf("invalid character %c for number!!!", r))
		
	}
}

func BitCount(i int) int {
	i = i - ((i >> 1) & 0x5555555555555555);
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333);
	i = (i + (i >> 4)) & 0x0f0f0f0f0f0f0f0f;
	i = i + (i >> 8);
	i = i + (i >> 16);
	i = i + (i >> 32);
	return i & 0x7f;
}

func FixedBitCount(i int) int {
	bit := BitCount(i)

	if bit <= 7 && (i >= math.MinInt8 && i <= math.MaxInt8) {
		return 7
	} else if bit <= 15 && (i >= math.MinInt16 && i <= math.MaxInt16) {
		return 15
	} else if bit <= 31 && (i >= math.MinInt32 && i <= math.MaxInt32) {
		return 31
	//lint:ignore SA4003 max is int64
	} else if bit <= 63 && (i >= math.MinInt64 && i <= math.MaxInt64) {
		return 63
	} else {
		return 0
	}
}

func ParseInt(number string) int {
	if len(number) <= 0 {
		return 0
	}

	for number[0] == '0' && len(number) > 1 {
		number = number[1:]
	}

	sig := 0 // 0 := + , 1 := -
	var value int = 0 // max 64 bit signed

	if number[0] == '+' {
		sig = 0
		number = number[1:]
	} else if number[0] == '-' {
		sig = 1
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
	sig := 0 // 0 := + , 1 := -
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
	sig := 0 // 0 := + , 1 := -
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
	sig := 0 // 0 := + , 1 := -
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