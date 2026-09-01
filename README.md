# nom

A text parser.

## Install

```sh
go get github.com/ofabricio/nom
```

## Examples

See more examples in the [examples](/examples) folder.

## Example: Tokenizer

This example shows how to tokenize some text.
See [playground](https://go.dev/play/p/WPk_5aysaS5).

```go
package main

import . "github.com/ofabricio/nom"

func main() {

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
```

## Example: Parsing a math expression

This example shows how to parse and evaluate a mathematical expression.
See [playground](https://go.dev/play/p/Qdpjl8VP1ZQ).

https://github.com/ofabricio/nom/blob/b37e363829856b707c64201dfb3f16148ce5b2ef/examples/expr_test.go#L11-L84

## Documentation

This parser works by matching a pattern against the current position and advancing as it matches.

### Functions

There are a few functions to help handling the parsing.

**Matchers:** functions that advance the parser as they match.

- Match
- Expect
- Optional
- Any

**Testers:** functions that test for a pattern without advancing the parser.

- Equal
- More

**Capturers:** functions that capture tokens on a match.

- MatchOut
- ExpectOut
- Out
- Mark, Token
- Curr

**Movement:** functions that move the parser back and forth.

- Mark, Back
- Next

**Other:** other functions.

- Head
- Tail
- Body

**Utils:** utility functions.

- Undo
- Opt
- Tokenize, T
