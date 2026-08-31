package lang

import (
	"fmt"

	. "github.com/ofabricio/nom"
)

// Example_tokenizer shows how to tokenize a string into individual tokens.
func Example_tokenizer() {

	src := `
		fun main () {
			one();
		}

		let x = 1
	`

	p := New(src)
	for p.Optional(WS) && p.More() {
		if m := p.Mark(); p.Match(WORD) || p.Any() {
			t := p.Token(m)
			fmt.Printf("Idx: %2v, Row: %1v, Col: %2v, Text: %v\n", t.Idx, t.Row, t.Col, t.Text)
		}
	}

	// Output:
	// Idx:  3, Row: 2, Col:  3, Text: fun
	// Idx:  7, Row: 2, Col:  7, Text: main
	// Idx: 12, Row: 2, Col: 12, Text: (
	// Idx: 13, Row: 2, Col: 13, Text: )
	// Idx: 15, Row: 2, Col: 15, Text: {
	// Idx: 20, Row: 3, Col:  4, Text: one
	// Idx: 23, Row: 3, Col:  7, Text: (
	// Idx: 24, Row: 3, Col:  8, Text: )
	// Idx: 25, Row: 3, Col:  9, Text: ;
	// Idx: 29, Row: 4, Col:  3, Text: }
	// Idx: 34, Row: 6, Col:  3, Text: let
	// Idx: 38, Row: 6, Col:  7, Text: x
	// Idx: 40, Row: 6, Col:  9, Text: =
	// Idx: 42, Row: 6, Col: 11, Text: 1
}
