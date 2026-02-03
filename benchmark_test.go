package rut

import (
	"testing"
)

func BenchmarkValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Validate("12.345.678-5")
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("123456785")
	}
}

func BenchmarkFormat_Complete(b *testing.B) {
	r := RUT{Number: 12345678, DV: '5'}
	for i := 0; i < b.N; i++ {
		r.Format(FormatComplete)
	}
}

func BenchmarkFormat_Escaped(b *testing.B) {
	r := RUT{Number: 12345678, DV: '5'}
	for i := 0; i < b.N; i++ {
		r.Format(FormatEscaped)
	}
}

func BenchmarkFormat_WithDash(b *testing.B) {
	r := RUT{Number: 12345678, DV: '5'}
	for i := 0; i < b.N; i++ {
		r.Format(FormatWithDash)
	}
}

func BenchmarkCalculateDV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateDV(12345678)
	}
}
