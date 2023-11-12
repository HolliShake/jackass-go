package jutil
import (
	"math"
	"fmt"
)

const (
	MAX_SAFE_INDEX int32 = math.MaxInt32
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
