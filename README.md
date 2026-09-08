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
See [playground](https://go.dev/play/p/3Cz5GXCJVjc).

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
        fmt.Printf("Kind: %v, Idx: %2v, Ln: %1v, Col: %2v, Text: %v\n", kind, t.Offset, t.Line, t.Column, t.Text)
    }

    // Output:
    // Kind: Word, Idx:  3, Ln: 2, Col:  3, Text: fun
    // Kind: Word, Idx:  7, Ln: 2, Col:  7, Text: main
    // Kind: Symb, Idx: 12, Ln: 2, Col: 12, Text: (
    // Kind: Symb, Idx: 13, Ln: 2, Col: 13, Text: )
    // Kind: Symb, Idx: 15, Ln: 2, Col: 15, Text: {
    // Kind: Word, Idx: 20, Ln: 3, Col:  4, Text: one
    // Kind: Symb, Idx: 23, Ln: 3, Col:  7, Text: (
    // Kind: Symb, Idx: 24, Ln: 3, Col:  8, Text: )
    // Kind: Symb, Idx: 25, Ln: 3, Col:  9, Text: ;
    // Kind: Symb, Idx: 29, Ln: 4, Col:  3, Text: }
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

| Function | Description |
| :---     | :---        |
| Match | Matches and advances the parser on success. |
| MatchOut | Same as `Match`, but outputs the matching token on success. |
| Exp | Matches and advances the parser on success, but triggers an expectation error on failure. |
| ExpOut | Same as `Exp`, but outputs the matching token on success. |
| Opt | Matches optionally. |
| Any | Matches any character and advances the parser on success. |
| Equal | Tests a pattern without advancing the parser on success. |
| More | Tells if there are more characters to parse. |
| Find | Advances through the input until it finds a pattern. |
| FindOut | Same as `Find`, but outputs the matching token on success. |
| GetLine | Returns the current line. |
| String | Matches a string. |
| Out | Outputs a token on success. |
| Mark | Sets a mark at the current position. |
| Back | Sends the parser back to a mark. |
| Token | Returns the token between two marks. |
| Rune | Returns the current character as a rune. |
| Char | Returns the current character as a string. |
| Next | Advances the parser by one character. |
| Head | Returns the portion of the input before the current position. |
| Tail | Returns the portion of the input from the current position onwards. |
| Body | Returns the entire input. |
| Undo | Moves the parser back to a mark. |
| Tokenize | Tokenizes the input based on the provided tokenizers `T` |
| T | A tokenizer used in `Tokenize` |
