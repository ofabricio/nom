package nom

import (
	"iter"
)

// Tokenize tokenizes the input string based on the provided tokenizers.
func Tokenize[K any](src string, match ...Tokenizer[K]) iter.Seq2[K, Token] {
	return func(yield func(K, Token) bool) {
		p := New(src)
	outer:
		for p.More() {
			mrk := p.Mark()
			for _, m := range match {
				if k, ok := m(&p); ok {
					if !yield(k, p.Token(mrk)) {
						return
					}
					continue outer
				}
			}
			p.Next()
		}
	}
}

// T matches and identifies a token.
func T[Kind any, P Pattern](kind Kind, pattern P) Tokenizer[Kind] {
	return func(s *Parser) (Kind, bool) {
		return kind, s.Match(pattern)
	}
}

// Tokenizer is a function that attempts to match a
// token and returns its kind and success status.
type Tokenizer[Kind any] func(s *Parser) (Kind, bool)
