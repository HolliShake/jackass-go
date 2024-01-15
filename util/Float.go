package util

import (
	"math"
)

func ParseFloat(number string) float64 {
	if len(number) <= 0 {
		return 0
	}

	for number[0] == '0' && len(number) > 1 {
		number = number[1:]
	}

	sig, other_exponent_sig := 0, 0                      // 0 := + , 1 := -
	var exponent, other_exponent, mantissa int = 0, 0, 0 // 11 max 52

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

	decimalPlaces := 0
	truncated, decimal := false, false

	for index := 0; index < len(number); index++ {
		char := number[index]

		if !truncated {
			if !decimal {
				if char == '.' {
					decimal = true
				} else if char == 'e' || char == 'E' {
					truncated = true
				} else {
					exponent = (exponent * 10) + valueOf(char)

					if FixedBitCount(exponent) > 11 {
						decimal = true
					}
				}
			} else {
				if char == 'e' || char == 'E' {
					truncated = true
				} else {
					mantissa = (mantissa * 10) + valueOf(char)
					decimalPlaces++

					if FixedBitCount(mantissa) >= 52 {
						break
					}
				}
			}
		} else {
			if char == '+' {
				other_exponent_sig = 0
			} else if char == '-' {
				other_exponent_sig = 1
			} else {
				other_exponent = (other_exponent * 10) + valueOf(char)

				if FixedBitCount(other_exponent) > 11 {
					break
				}
			}
		}
	}

	value := 0.0

	if mantissa <= 0 {
		value = float64(exponent)
	} else {
		value = float64(exponent) + (float64(mantissa) / math.Pow(10, float64(decimalPlaces)))
	}

	finalValue := 0.0

	if other_exponent_sig == 0 {
		finalValue = value * (math.Pow(10, float64(other_exponent)))
	} else if other_exponent_sig == 1 {
		finalValue = value * (math.Pow(10, float64(-other_exponent)))
	}

	if sig == 0 {
		return finalValue
	} else {
		return -finalValue
	}
}
