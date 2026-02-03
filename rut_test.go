package rut

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"12.345.678-5", true},
		{"12345678-5", true},
		{"123456785", true},
		{"1.009-k", true},
		{"1.009-K", true},
		{"7.654.321-6", true},
		{"11.111.111-1", true},
		{"12.345.678-0", false}, // Invalid DV
		{"1.234.567-4", true},
		{"5.555.555-5", false},
		{"", false},
		{"123", false},              // Too short
		{"12.345.678.901-2", false}, // Too long
		{"abc-d", false},            // Invalid chars
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := Validate(tt.input); got != tt.expected {
				t.Errorf("Validate(%q) = %v; want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		wantNum int
		wantDV  byte
		wantErr bool
		errType error
	}{
		{"12.345.678-5", 12345678, '5', false, nil},
		{"1.009-K", 1009, 'K', false, nil},
		{"1-9", 0, 0, true, ErrTooShort},
		{"1234-5", 1234, '5', false, nil}, // Minimum valid (5 chars)
		{"12345678901", 0, 0, true, ErrTooLong},
		{"12.34K.678-5", 0, 0, true, ErrInvalidFormat},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Number != tt.wantNum {
					t.Errorf("Parse(%q) Number = %v, want %v", tt.input, got.Number, tt.wantNum)
				}
				if got.DV != tt.wantDV {
					t.Errorf("Parse(%q) DV = %c, want %c", tt.input, got.DV, tt.wantDV)
				}
			}
		})
	}
}

func TestFormat(t *testing.T) {
	input := "123456785"

	tests := []struct {
		name     string
		style    FormatStyle
		expected string
	}{
		{"Complete", FormatComplete, "12.345.678-5"},
		{"Escaped", FormatEscaped, "123456785"},
		{"WithDash", FormatWithDash, "12345678-5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(input, tt.style)
			if err != nil {
				t.Fatalf("Format() error = %v", err)
			}
			if got != tt.expected {
				t.Errorf("Format() = %q; want %q", got, tt.expected)
			}
		})
	}
}

func TestRUT_String(t *testing.T) {
	r := RUT{Number: 12345678, DV: '5'}
	expected := "12.345.678-5"
	if got := r.String(); got != expected {
		t.Errorf("RUT.String() = %q; want %q", got, expected)
	}
}

func TestCalculateDV(t *testing.T) {
	tests := []struct {
		num      int
		expected byte
	}{
		{12345678, '5'},
		{7654321, '6'},
		{11111111, '1'},
		{1009, 'K'},
		{14555848, '4'},
		{0, '0'},
	}

	for _, tt := range tests {
		if got := CalculateDV(tt.num); got != tt.expected {
			t.Errorf("CalculateDV(%d) = %c; want %c", tt.num, got, tt.expected)
		}
	}
}
