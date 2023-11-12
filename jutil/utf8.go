package jutil
import (
	"unicode"
	"strings"
)


// UnicodeCategory returns the Unicode Character Category of the given rune.
func UnicodeCategory(r rune) string {
    for name, table := range unicode.Categories {
        if len(name) == 2 && unicode.Is(table, r) {
            return name
        }
    }
    return "Cn"
}

const (
    _BYTE1 int = 0b00000000
    _BYTE2 int = 0b11000000
    _BYTE3 int = 0b11100000
    _BYTE4 int = 0b11110000

    _2BYTE_FOLLOW int = 0b00011111
    _3BYTE_FOLLOW int = 0b00001111
    _4BYTE_FOLLOW int = 0b00000111

    _VALID_TRAILING int = 0b10000000
    _MAX_TRAILING   int = 0b00111111
)


func Utf_toCodePoint(b1, b2, b3, b4 int) rune {
	var ord int = 0
	
    switch Utf_sizeOfUtf(b1) {
        case 1:
            return rune(b1)
         
        case 2:
            ord  = ((b1 & _2BYTE_FOLLOW) << 6)
            ord |= ((b2 & _MAX_TRAILING))
            
        case 3:
            ord  = ((b1 & _3BYTE_FOLLOW) << 12)
            ord |= ((b2 & _MAX_TRAILING) <<  6)
            ord |= ((b3 & _MAX_TRAILING))
            
        case 4:
            ord  = ((b1 & _4BYTE_FOLLOW) << 18)
            ord |= ((b2 & _MAX_TRAILING) << 12)
            ord |= ((b3 & _MAX_TRAILING) <<  6)
            ord |= ((b4 & _MAX_TRAILING))
        default:
            
    }

    return rune(ord)
}

func Utf_sizeOfUtf(firstByte int) int {
    if (firstByte & _BYTE4) == _BYTE4 {
		return 4
	} else if (firstByte & _BYTE3) == _BYTE3 {
		return 3
	} else if (firstByte & _BYTE2) == _BYTE2 {
		return 2
	} else if (firstByte & _BYTE1) == 0 {
		return 1
	}
		
    return 0
}

func Utf_isLetter(r rune) bool {
	if unicode.IsLetter(r) || r == '_' {
		return true
	}

	if r < 0x80 {
		if r == '_' {
			return true
		}
		return ((uint8(r) | 0x20) - 0x61) < 26
	}

	switch UnicodeCategory(r) {
        case "Lu":
        case "Ll":
        case "Lt":
        case "Lm":
        case "Lo":
            return true
	}

	return false
}

func Utf_isDigit(r rune) bool {
    if r < 0x80 {
		return (uint8(r) - '0') < 10
	}
	return strings.Compare(UnicodeCategory(r), "Nd") == 0
}

func Utf_isLetterOrDigit(r rune) bool {
    if (r < 0x80) {
		if (r == '_') {
			return true
		}
		if ((uint8(r) | 0x20) - 0x61) < 26 {
			return true
		}
		return (uint8(r) - '0') < 10
	}
	switch UnicodeCategory(r) {
		case "Lu":
		case "Ll":
		case "Lt":
		case "Lm":
		case "Lo":
			return true
        case "Nd":
            return true
	}

	return false
}

func Utf_isWhiteSpace(r rune) bool {
    if r < 0x80 {
        return r == ' ' || (r >= 0x09 && r <= 0x0D)
    } else if (
        r == 0x00 ||
        r == 0x08 ||
        r == 0x09 ||
        r == 0x0a ||
        r == 0x0d ||
        r == 0x20 ) {
        return true
    }
    return strings.Compare(UnicodeCategory(r), "Zs") == 0
}

func Utf_codePointLength(str string) int {
    if len(str) <= 0 {
        return 0
    }

    length := 0

    for index := 0; index < len(str); {
        size := Utf_sizeOfUtf(int(str[index]))
        index += size
        length++
    }

    return length
}