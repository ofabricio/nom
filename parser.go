package nom

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// New creates a new parser for the given source string.
func New(src string) Parser {
	return Parser{src: src, Row: 1, Col: 1}
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

// Opt makes cond optional. Ther mark m moves the
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
	r := p.Curr()
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
	return f(p.Curr())
}

// Any matches any characters.
func (p *Parser) Any() bool {
	return p.Next()
}

// Next advances the parser by one character.
func (p *Parser) Next() bool {
	return p.advance(string(p.Curr()))
}

// Curr returns the current character.
func (p *Parser) Curr() rune {
	r, _ := utf8.DecodeRuneInString(p.Tail())
	return r
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
	return len(p.Tail()) > 0
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
	line := e.getErrorLine()
	tabs := strings.Count(line, "\t") * 3
	line = strings.ReplaceAll(line, "\t", "    ")
	a := fmt.Sprintf("failed to parse: line %d char %d: expected %s", e.Row, e.Col, e.Msg)
	aa := fmt.Sprintf("%5s |", "")
	ab := fmt.Sprintf("%5d | %s", e.Row, line)
	ac := fmt.Sprintf("%5s |%s%s", "", strings.Repeat(" ", e.Col+tabs), "^--")
	return fmt.Sprintf("%s\n%s\n%s\n%s", a, aa, ab, ac)
}

func (e *Error) getErrorLine() string {
	ini := strings.LastIndex(e.Head(), "\n") + 1
	end := strings.Index(e.Tail(), "\n")
	if end == -1 {
		end = len(e.Body())
	} else {
		end = len(e.Head()) + end
	}
	return e.Body()[ini:end]
}

var WORD = regexp.MustCompile(`^\w+`)
var DIGITS = regexp.MustCompile(`^\d+`)
var ST = regexp.MustCompile(`^[ \t]+`)
var WS = regexp.MustCompile(`^\s+`)
