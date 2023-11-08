package main
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


func utf_toCodePoint(b1, b2, b3, b4 int) rune {
	var ord int = 0
	
    switch utf_sizeOfUtf(b1) {
        case 1:
            return rune(b1)
         
        case 2:
            ord  = ((b1 & _2BYTE_FOLLOW) << 6)
            ord |= ((b2 & _MAX_TRAILING))
            break
        case 3:
            ord  = ((b1 & _3BYTE_FOLLOW) << 12)
            ord |= ((b2 & _MAX_TRAILING) <<  6)
            ord |= ((b3 & _MAX_TRAILING))
            break
        case 4:
            ord  = ((b1 & _4BYTE_FOLLOW) << 18)
            ord |= ((b2 & _MAX_TRAILING) << 12)
            ord |= ((b3 & _MAX_TRAILING) <<  6)
            ord |= ((b4 & _MAX_TRAILING))
        default:
            break
    }

    return rune(ord)
}

func utf_sizeOfUtf(firstByte int) int {
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

func utf_isLetter(r rune) bool {
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

func utf_isDigit(r rune) bool {
    if r < 0x80 {
		return (uint8(r) - '0') < 10
	}
	return strings.Compare(UnicodeCategory(r), "Nd") == 0
}

func utf_isLetterOrDigit(r rune) bool {
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

func utf_isWhiteSpace(r rune) bool {
    if r < 0x80 {
        return r == ' ' || (r >= 0x09 && r <= 0x0D)
    }
    return strings.Compare(UnicodeCategory(r), "Zs") == 0
}