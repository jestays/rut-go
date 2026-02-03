// Package rut provides functions to validate, parse, and format
// Chilean RUTs (Rol Único Tributario - Chilean tax ID).
//
// Supported formats:
//   - With dots and dash: "12.345.678-9"
//   - Only dash: "12345678-9"
//   - No separators: "123456789"
package rut

import (
	"errors"
	"strconv"
	"strings"
)

// Package errors
var (
	ErrInvalidFormat = errors.New("rut: invalid format")
	ErrEmptyRUT      = errors.New("rut: empty string")
	ErrTooShort      = errors.New("rut: too short (minimum 5 characters)")
	ErrTooLong       = errors.New("rut: too long (maximum 10 characters)")
)

// FormatStyle defines the formatting style for the RUT.
type FormatStyle int

const (
	// FormatComplete formats as "12.345.678-9" (dots and dash)
	FormatComplete FormatStyle = iota
	// FormatEscaped formats as "123456789" (no separators)
	FormatEscaped
	// FormatWithDash formats as "12345678-9" (dash only)
	FormatWithDash
)

// multipliers is a lookup table for the check digit calculation
var multipliers = [6]int{2, 3, 4, 5, 6, 7}

// isValidRUTChar checks if a character is valid for a RUT and normalizes it.
// Returns the normalized character and true if valid, 0 and false otherwise.
func isValidRUTChar(c byte) (byte, bool) {
	if c >= '0' && c <= '9' {
		return c, true
	}
	if c == 'k' || c == 'K' {
		return 'K', true
	}
	return 0, false
}

// RUT represents a parsed Chilean RUT.
type RUT struct {
	Number int  // RUT number without check digit
	DV     byte // Check digit ('0'-'9' or 'K')
}

// Validate checks if a RUT string is valid.
// It accepts formats with or without dots and with or without dash.
// Case insensitive for 'K'.
func Validate(rut string) bool {
	r, err := Parse(rut)
	if err != nil {
		return false
	}
	return r.Validate()
}

// Parse extracts the number and check digit from a RUT string.
// It returns an error if the format is invalid or the length is out of bounds.
func Parse(s string) (RUT, error) {
	if s == "" {
		return RUT{}, ErrEmptyRUT
	}

	// Clean separators and validate characters
	var (
		raw [12]byte
		n   int
	)

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '.' || c == '-' {
			continue
		}
		if n >= 12 {
			return RUT{}, ErrTooLong
		}

		// Validate and normalize character
		char, ok := isValidRUTChar(c)
		if !ok {
			return RUT{}, ErrInvalidFormat
		}

		raw[n] = char
		n++
	}

	// Length validation (5 to 10 characters as requested)
	// We count the digits + DV
	if n < 5 {
		return RUT{}, ErrTooShort
	}
	if n > 10 {
		return RUT{}, ErrTooLong
	}

	// DV is the last character
	dv := raw[n-1]

	// Check if 'K' is in the wrong place
	for i := 0; i < n-1; i++ {
		if raw[i] == 'K' {
			return RUT{}, ErrInvalidFormat
		}
	}

	// Parse number
	numStr := string(raw[:n-1])
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return RUT{}, ErrInvalidFormat
	}

	return RUT{
		Number: num,
		DV:     dv,
	}, nil
}

// Format normalizes and formats a RUT string according to the specified style.
func Format(s string, style FormatStyle) (string, error) {
	r, err := Parse(s)
	if err != nil {
		return "", err
	}
	return r.Format(style), nil
}

// CalculateDV computes the check digit for a given RUT number.
func CalculateDV(number int) byte {
	if number == 0 {
		return '0'
	}

	sum := 0
	multiplierIdx := 0

	for number > 0 {
		digit := number % 10
		sum += digit * multipliers[multiplierIdx]

		number /= 10
		multiplierIdx = (multiplierIdx + 1) % 6
	}

	remainder := sum % 11
	checkResult := 11 - remainder

	switch checkResult {
	case 11:
		return '0'
	case 10:
		return 'K'
	default:
		return byte(checkResult + '0')
	}
}

// String implements fmt.Stringer using FormatComplete style.
func (r RUT) String() string {
	return r.Format(FormatComplete)
}

// Format returns the RUT formatted according to the specified style.
func (r RUT) Format(style FormatStyle) string {
	numStr := strconv.Itoa(r.Number)

	switch style {
	case FormatEscaped:
		var b strings.Builder
		b.Grow(len(numStr) + 1)
		b.WriteString(numStr)
		b.WriteByte(r.DV)
		return b.String()

	case FormatWithDash:
		var b strings.Builder
		b.Grow(len(numStr) + 2)
		b.WriteString(numStr)
		b.WriteByte('-')
		b.WriteByte(r.DV)
		return b.String()

	case FormatComplete:
		fallthrough
	default:
		// Format: XX.XXX.XXX-X
		var b strings.Builder
		// Max length is 12: 12.345.678-K
		b.Grow(12)

		n := len(numStr)
		for i, c := range numStr {
			b.WriteRune(c)
			// Add dots from right to left every 3 digits
			distFromEnd := n - i - 1
			if distFromEnd > 0 && distFromEnd%3 == 0 {
				b.WriteByte('.')
			}
		}

		b.WriteByte('-')
		b.WriteByte(r.DV)
		return b.String()
	}
}

// Validate checks if the RUT's check digit matches the calculated one.
func (r RUT) Validate() bool {
	if r.Number <= 0 {
		return false
	}
	return r.DV == CalculateDV(r.Number)
}
