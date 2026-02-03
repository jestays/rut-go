![rut-go](https://plvvvfsabjjcaadqfygn.supabase.co/storage/v1/object/public/assets/banner-rut-go.png)

![Go](https://github.com/jestays/rut-go/actions/workflows/go.yml/badge.svg?branch=main)
![Go Version](https://img.shields.io/github/go-mod/go-version/jestays/rut-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/jestays/rut-go.svg)](https://pkg.go.dev/github.com/jestays/rut-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jestays/rut-go)](https://goreportcard.com/report/github.com/jestays/rut-go)
![License](https://img.shields.io/github/license/jestays/rut-go)

Validate, parse, and format Chilean RUTs (Rol Unico Tributario) in Go.

## Features
- Validate RUT strings with or without separators
- Parse into a structured `RUT` type
- Format in three styles: dots+dash, dash only, or no separators
- Fast, allocation-light implementation

## Requirements
- Go 1.21+

## Install
```bash
go get github.com/jestays/rut-go
```

## License
MIT. See `LICENSE`.

## Usage
```go
package main

import (
	"fmt"

	"github.com/jestays/rut-go"
)

func main() {
	// Validate a string
	fmt.Println(rut.Validate("12.345.678-5")) // true

	// Parse a RUT
	r, err := rut.Parse("123456785")
	if err != nil {
		panic(err)
	}

	// Format in different styles
	fmt.Println(r.Format(rut.FormatComplete)) // "12.345.678-5"
	fmt.Println(r.Format(rut.FormatWithDash)) // "12345678-5"
	fmt.Println(r.Format(rut.FormatEscaped))  // "123456785"

	// Compute a check digit
	fmt.Printf("%c\n", rut.CalculateDV(12345678)) // '5'
}
```

## Format styles
```go
const (
	FormatComplete FormatStyle = iota // "12.345.678-9"
	FormatEscaped                    // "123456789"
	FormatWithDash                   // "12345678-9"
)
```

`Format` accepts any supported input format and returns the selected style:
```go
out, err := rut.Format("12.345.678-5", rut.FormatEscaped)
// out == "123456785"
```

## Validation rules
- Separators are optional. Dots and dashes are ignored during parsing.
- The check digit can be numeric or `K` (case-insensitive).
- After removing separators, the length must be **5 to 10 characters**
  (digits + check digit).

## API summary
- `Validate(string) bool`
- `Parse(string) (RUT, error)`
- `Format(string, FormatStyle) (string, error)`
- `CalculateDV(int) byte`
- `type RUT struct { Number int; DV byte }`
  - `func (RUT) Validate() bool`
  - `func (RUT) Format(FormatStyle) string`
  - `func (RUT) String() string` (uses `FormatComplete`)

## Errors
`Parse` and `Format` can return:
- `ErrEmptyRUT`
- `ErrTooShort`
- `ErrTooLong`
- `ErrInvalidFormat`
- `ErrInvalidDV`

## Tests and benchmarks
```bash
go test -v .
go test -bench=. -benchmem
```
