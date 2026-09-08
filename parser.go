package nom

import (
	"fmt"
	"regexp"
	"strings"
)

// New creates a new parser for the given source string.
func New(src string) Parser {
	return Parser{src: src, Line: 1, Column: 1}
}

// String matches a string enclosed in the quote character.
// It escapes the quote character with a backslash.
func (p *Parser) String(quote string) bool {
	if p.Match(quote) {
		for p.More() && !p.Equal(quote) {
			if !(p.Match("\\") && p.Match(quote)) {
				p.Next()
			}
		}
		return p.Exp(quote)
	}
	return false
}

// GetLine returns the current line, even if the
// parser is not in the beginning of the line.
func (p *Parser) GetLine() Token {
	ini := strings.LastIndex(p.Head(), "\n") + 1
	end := strings.Index(p.Tail(), "\n")
	if end == -1 {
		end = len(p.Body())
	} else {
		end = len(p.Head()) + end
	}
	return Token{Text: p.src[ini:end], Offset: ini, Line: p.Line, Column: 1}
}

// MatchOut matches the given pattern and outputs
// the corresponding token on success.
func (p *Parser) MatchOut[P Pattern](pattern P, out *Token) bool {
	return p.Out(p.Mark(), p.Match(pattern), out)
}

// ExpOut expects the given pattern and outputs the
// corresponding token on success, or triggers an
// expectation error on failure.
func (p *Parser) ExpOut[P Pattern](pattern P, out *Token) bool {
	return p.Out(p.Mark(), p.Exp(pattern), out)
}

// Out outputs the corresponding token between the mark m
// and the current position of the parser if cond is true.
func (p *Parser) Out(m Marker, cond bool, out *Token) bool {
	if cond {
		*out = p.Token(m)
	}
	return cond
}

// Undo moves the parser back to the mark m if
// cond is false and returns the condition.
func (p *Parser) Undo(m Marker, cond bool) bool {
	if !cond {
		p.Back(m)
	}
	return cond
}

// FindOut is like Find, but outputs the matching token if found.
func (p *Parser) FindOut[P Pattern](pattern P, out *Token) bool {
	return p.Find(pattern) && p.MatchOut(pattern, out)
}

// Find advances through the input until it finds a pattern.
// Returns true if found, and the parser is at the start of
// the pattern. The parser hits the end of the input if not
// found.
func (p *Parser) Find[P Pattern](pattern P) bool {
	for p.More() && !p.Equal(pattern) {
		p.Next()
	}
	return p.More()
}

// Opt optionally matches the given pattern.
func (p *Parser) Opt[P Pattern](pattern P) bool {
	return p.Match(pattern) || true
}

// Exp expects the given pattern and triggers
// an expectation error if it fails.
func (p *Parser) Exp[P Pattern](pattern P) bool {
	switch pattern := any(pattern).(type) {
	case string:
		return p.MatchString(pattern) || p.Expected(pattern)
	case *regexp.Regexp:
		return p.MatchRegex(pattern) || p.Expected(pattern.String())
	case func(rune) bool:
		return p.MatchFunc(pattern) || p.Expected("token")
	default:
		return false
	}
}

// Expected triggers an expectation error for the given pattern.
func (p *Parser) Expected(msg string) bool {
	if p.Err == nil {
		p.Err = &Error{Marker: p.Marker, ErrLine: p.GetLine().Text, Message: msg}
	}
	return false
}

// Match matches the given pattern and advances the parser
// on success. Returns true if it matches.
func (p *Parser) Match[P Pattern](pattern P) bool {
	switch pattern := any(pattern).(type) {
	case string:
		return p.MatchString(pattern)
	case *regexp.Regexp:
		return p.MatchRegex(pattern)
	case func(rune) bool:
		return p.MatchFunc(pattern)
	default:
		return false
	}
}

// MatchString matches the given string and advances the parser
// on success. Returns true if it matches.
func (p *Parser) MatchString(v string) bool {
	return p.EqualString(v) && p.advance(v)
}

// MatchRegex matches the given regular expression and advances
// the parser on success. Returns true if it matches.
func (p *Parser) MatchRegex(v *regexp.Regexp) bool {
	return p.advance(v.FindString(p.Tail()))
}

// MatchFunc matches the given rune function and advances the
// parser on success. Returns true if it matches.
func (p *Parser) MatchFunc(f func(rune) bool) bool {
	r := p.Rune()
	return f(r) && p.advance(string(r))
}

// Equal checks if the given pattern matches the
// current parser position without advancing.
func (p *Parser) Equal[P Pattern](pattern P) bool {
	switch pattern := any(pattern).(type) {
	case string:
		return p.EqualString(pattern)
	case *regexp.Regexp:
		return p.EqualRegex(pattern)
	case func(rune) bool:
		return p.EqualFunc(pattern)
	default:
		return false
	}
}

// EqualString checks if the given string matches
// the current parser position without advancing.
func (p *Parser) EqualString(v string) bool {
	return strings.HasPrefix(p.Tail(), v)
}

// EqualRegex checks if the given regular expression matches
// the current parser position without advancing.
func (p *Parser) EqualRegex(v *regexp.Regexp) bool {
	return v.MatchString(p.Tail())
}

// EqualFunc checks if the given rune function matches
// the current parser position without advancing.
func (p *Parser) EqualFunc(f func(rune) bool) bool {
	return f(p.Rune())
}

// Any matches any characters.
func (p *Parser) Any() bool {
	return p.Next()
}

// Next advances the parser by one character.
func (p *Parser) Next() bool {
	return p.advance(p.Char())
}

// Rune returns the current character as a rune.
func (p *Parser) Rune() rune {
	for _, v := range p.Tail() {
		return v
	}
	return 0
}

// Char returns the current character as a string.
func (p *Parser) Char() string {
	for _, v := range p.Tail() {
		return string(v)
	}
	return ""
}

// Head returns the portion of the source
// before the current parser position.
func (p *Parser) Head() string {
	return p.src[:p.Offset]
}

// Tail returns the portion of the source from
// the current parser position onwards.
func (p *Parser) Tail() string {
	return p.src[p.Offset:]
}

// Body returns the entire source string.
func (p *Parser) Body() string {
	return p.src
}

func (p *Parser) advance(v string) bool {
	p.Offset += len(v)
	p.coln(v)
	return len(v) > 0
}

func (p *Parser) coln(v string) {
	for _, r := range v {
		p.Column++
		if r == '\n' {
			p.Line++
			p.Column = 1
		}
	}
}

// Mark returns a mark of the current parser state.
func (p Parser) Mark() Marker {
	return p.Marker
}

// Back restores the parser state to the given mark.
func (p *Parser) Back(m Marker) {
	p.Marker = m
}

// Token returns a token representing the text between
// the given mark and the current parser position.
func (p *Parser) Token(m Marker) Token {
	return Token{Text: p.src[m.Offset:p.Offset], Marker: m}
}

// More checks if there are more characters to parse.
func (p Parser) More() bool {
	return p.Offset < len(p.src)
}

// Parser represents a parser.
type Parser struct {
	Marker
	src string
	Err error
}

// Token represents a token extracted from the source string.
type Token struct {
	Marker
	Text string
}

// Error represents a parsing error that occurred during parsing.
type Error struct {
	Marker
	ErrLine string
	Message string
}

type Marker struct {
	Offset int
	Line   int
	Column int
}

func (e *Error) Error() string {
	tabs := strings.Count(e.ErrLine, "\t") * 3
	line := strings.ReplaceAll(e.ErrLine, "\t", "    ")
	a := fmt.Sprintf("failed to parse: line %d char %d: expected %s", e.Line, e.Column, e.Message)
	b := fmt.Sprintf("%5s |", "")
	c := fmt.Sprintf("%5d | %s", e.Line, line)
	d := fmt.Sprintf("%5s |%s%s", "", strings.Repeat(" ", e.Column+tabs), "^--")
	return fmt.Sprintf("%s\n%s\n%s\n%s", a, b, c, d)
}

// Pattern represents a pattern that can be matched by the parser.
type Pattern interface {
	string | *regexp.Regexp | func(rune) bool
}

var WORD = regexp.MustCompile(`^\w+`)
var DIGITS = regexp.MustCompile(`^\d+`)
var HS = regexp.MustCompile(`^[ \t]+`)
var WS = regexp.MustCompile(`^\s+`)
