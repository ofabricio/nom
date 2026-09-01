package lang

import (
	"fmt"
	"regexp"

	. "github.com/ofabricio/nom"
)

// ExampleTokenize shows how to split a string into individual tokens.
func ExampleTokenize() {

	src := `
		fun main () {
			one();
		}
	`

	SYMB := regexp.MustCompile(`^[^\s]`)

	for kind, t := range Tokenize(src, T("Word", WORD), T("Symb", SYMB)) {
		fmt.Printf("Kind: %v, Idx: %2v, Row: %1v, Col: %2v, Text: %v\n", kind, t.Idx, t.Row, t.Col, t.Text)
	}

	// Output:
	// Kind: Word, Idx:  3, Row: 2, Col:  3, Text: fun
	// Kind: Word, Idx:  7, Row: 2, Col:  7, Text: main
	// Kind: Symb, Idx: 12, Row: 2, Col: 12, Text: (
	// Kind: Symb, Idx: 13, Row: 2, Col: 13, Text: )
	// Kind: Symb, Idx: 15, Row: 2, Col: 15, Text: {
	// Kind: Word, Idx: 20, Row: 3, Col:  4, Text: one
	// Kind: Symb, Idx: 23, Row: 3, Col:  7, Text: (
	// Kind: Symb, Idx: 24, Row: 3, Col:  8, Text: )
	// Kind: Symb, Idx: 25, Row: 3, Col:  9, Text: ;
	// Kind: Symb, Idx: 29, Row: 4, Col:  3, Text: }
}
