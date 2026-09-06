package nom

import (
	"fmt"
	"regexp"
	"strings"
)

// New creates a new parser for the given source string.
func New(src string) Parser {
	return Parser{src: src, Row: 1, Col: 1}
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
		return p.Expect(quote)
	}
	return false
}

// Line matches the rest of a line.
func (p *Parser) Line() bool {
	return p.Equal("\n") || p.Find("\n")
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
	return Token{Text: p.src[ini:end], Idx: ini, Row: p.Row, Col: 1}
}

// MatchOut matches the given pattern and outputs
// the corresponding token on success.
func (p *Parser) MatchOut[P Pattern](pattern P, out *Token) bool {
	return p.Out(p.Mark(), p.Match(pattern), out)
}

// ExpectOut expects the given pattern and outputs the
// corresponding token on success, or triggers an
// expectation error on failure.
func (p *Parser) ExpectOut[P Pattern](pattern P, out *Token) bool {
	return p.Out(p.Mark(), p.Expect(pattern), out)
}

// Out outputs the corresponding token between the mark m
// and the current position of the parser if cond is true.
func (p *Parser) Out(m Parser, cond bool, out *Token) bool {
	if cond {
		*out = p.Token(m)
	}
	return cond
}

// Opt makes cond optional. The mark m moves the
// parser back to it on failure.
func (p *Parser) Opt(m Parser, cond bool) bool {
	return p.Undo(m, cond) || true
}

// Undo moves the parser back to the mark m if
// cond is false and returns the condition.
func (p *Parser) Undo(m Parser, cond bool) bool {
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

// Optional optionally parses the given pattern.
func (p *Parser) Optional[P Pattern](pattern P) bool {
	return p.Match(pattern) || true
}

// Expect expects the given pattern and triggers
// an expectation error if it fails.
func (p *Parser) Expect[P Pattern](pattern P) bool {
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

// Expects expects the given pattern and triggers an expectation
// error with the given message if it fails.
func (p *Parser) Expects[P Pattern](pattern P, msg string) bool {
	return p.Match(pattern) || p.Expected(msg)
}

// Expected triggers an expectation error for the given pattern.
func (p *Parser) Expected(msg string) bool {
	if p.Err == nil {
		p.Err = &Error{Parser: *p, Msg: msg}
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
	return p.src[:p.Idx]
}

// Tail returns the portion of the source from
// the current parser position onwards.
func (p *Parser) Tail() string {
	return p.src[p.Idx:]
}

// Body returns the entire source string.
func (p *Parser) Body() string {
	return p.src
}

func (p *Parser) advance(v string) bool {
	p.Idx += len(v)
	p.coln(v)
	return len(v) > 0
}

func (p *Parser) coln(v string) {
	for _, r := range v {
		p.Col++
		if r == '\n' {
			p.Row++
			p.Col = 1
		}
	}
}

// Mark returns a mark of the current parser state.
func (p Parser) Mark() Parser {
	return p
}

// Back restores the parser state to the given mark.
func (p *Parser) Back(m Parser) {
	*p = m
}

// Token returns a token representing the text between
// the given mark and the current parser position.
func (p *Parser) Token(m Parser) Token {
	return Token{Text: p.src[m.Idx:p.Idx], Idx: m.Idx, Row: m.Row, Col: m.Col}
}

// More checks if there are more characters to parse.
func (p Parser) More() bool {
	return p.Idx < len(p.src)
}

// Parser represents a parser.
type Parser struct {
	src string
	Idx int
	Row int
	Col int
	Err error
}

// Token represents a token extracted from the source string.
type Token struct {
	Text string
	Idx  int
	Row  int
	Col  int
}

// Pattern represents a pattern that can be matched by the parser.
type Pattern interface {
	string | *regexp.Regexp | func(rune) bool
}

// Error represents a parsing error that occurred during parsing.
type Error struct {
	Parser
	Msg string
}

func (e *Error) Error() string {
	line := e.GetLine().Text
	tabs := strings.Count(line, "\t") * 3
	line = strings.ReplaceAll(line, "\t", "    ")
	a := fmt.Sprintf("failed to parse: line %d char %d: expected %s", e.Row, e.Col, e.Msg)
	b := fmt.Sprintf("%5s |", "")
	c := fmt.Sprintf("%5d | %s", e.Row, line)
	d := fmt.Sprintf("%5s |%s%s", "", strings.Repeat(" ", e.Col+tabs), "^--")
	return fmt.Sprintf("%s\n%s\n%s\n%s", a, b, c, d)
}

var WORD = regexp.MustCompile(`^\w+`)
var DIGITS = regexp.MustCompile(`^\d+`)
var HS = regexp.MustCompile(`^[ \t]+`)
var WS = regexp.MustCompile(`^\s+`)
