package nom

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

func New(src string) Parser {
	return Parser{src: src, Row: 1, Col: 1}
}

func (p *Parser) MatchOut[T MatchType](v T, out *Token) bool {
	return p.Out(p.Mark(), p.Match(v), out)
}

func (p *Parser) ExpectOut[T MatchType](v T, out *Token) bool {
	return p.Out(p.Mark(), p.Expect(v), out)
}

func (p *Parser) Out(m Parser, cond bool, out *Token) bool {
	if cond {
		*out = p.Token(m)
	}
	return cond
}

func (p *Parser) Opt(m Parser, cond bool) bool {
	return p.Undo(m, cond) || true
}

// Undo sends the parser back to the mark m if cond is false.
func (p *Parser) Undo(m Parser, cond bool) bool {
	if !cond {
		p.Back(m)
	}
	return cond
}

func (p *Parser) Optional[T MatchType](v T) bool {
	return p.Match(v) || true
}

func (p *Parser) Expect[T MatchType](v T) bool {
	return p.Match(v) || p.expect(v)
}

func (p *Parser) Expected(msg string) bool {
	return p.expect(msg)
}

func (p *Parser) ExpectWith[T MatchType](v T, msg string) bool {
	return p.Match(v) || p.expect(msg)
}

func (p *Parser) expect[T MatchType](v T) bool {
	if p.Err != nil {
		return false
	}
	switch v := any(v).(type) {
	case string:
		p.Err = &Error{Parser: *p, Msg: v}
	case *regexp.Regexp:
		p.Err = &Error{Parser: *p, Msg: v.String()}
	default:
		p.Err = &Error{Parser: *p, Msg: "token"}
	}
	return false
}

func (p *Parser) Match[T MatchType](v T) bool {
	switch v := any(v).(type) {
	case string:
		return p.MatchString(v)
	case *regexp.Regexp:
		return p.MatchRegex(v)
	case func(rune) bool:
		return p.MatchFunc(v)
	default:
		return false
	}
}

func (p *Parser) MatchString(v string) bool {
	return p.EqualString(v) && p.advance(v)
}

func (p *Parser) MatchRegex(v *regexp.Regexp) bool {
	return p.advance(v.FindString(p.Tail()))
}

func (p *Parser) MatchFunc(f func(rune) bool) bool {
	r := p.Curr()
	return f(r) && p.advance(string(r))
}

func (p *Parser) Equal[T MatchType](v T) bool {
	switch v := any(v).(type) {
	case string:
		return p.EqualString(v)
	case *regexp.Regexp:
		return p.EqualRegex(v)
	case func(rune) bool:
		return p.EqualFunc(v)
	default:
		return false
	}
}

func (p *Parser) EqualString(v string) bool {
	return strings.HasPrefix(p.Tail(), v)
}

func (p *Parser) EqualRegex(v *regexp.Regexp) bool {
	return v.MatchString(p.Tail())
}

func (p *Parser) EqualFunc(f func(rune) bool) bool {
	return f(p.Curr())
}

func (p *Parser) Any() bool {
	return p.Next()
}

func (p *Parser) Next() bool {
	return p.advance(string(p.Curr()))
}

func (p *Parser) Curr() rune {
	r, _ := utf8.DecodeRuneInString(p.Tail())
	return r
}

func (p *Parser) Head() string {
	return p.src[:p.Idx]
}

func (p *Parser) Tail() string {
	return p.src[p.Idx:]
}

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

func (p Parser) Mark() Parser {
	return p
}

func (p *Parser) Back(m Parser) {
	*p = m
}

func (p *Parser) Token(m Parser) Token {
	return Token{Text: p.src[m.Idx:p.Idx], Idx: m.Idx, Row: m.Row, Col: m.Col}
}

func (p Parser) More() bool {
	return len(p.Tail()) > 0
}

type Parser struct {
	src string
	Idx int
	Row int
	Col int
	Err error
}

type Token struct {
	Text string
	Idx  int
	Row  int
	Col  int
}

type MatchType interface {
	int | string | *regexp.Regexp | func(rune) bool
}

// Error represents a parsing error.
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
	ini := strings.LastIndex(e.Head(), "\n")
	if ini == -1 {
		ini = 0
	} else {
		ini += 1 // Exclude the \n.
	}
	end := strings.Index(e.Tail(), "\n")
	switch end {
	case -1: // No newline found after the error position.
		end = len(e.Body())
	case 0: // Newline immediately after the error position.
		end = len(e.Head())
	default:
		end = len(e.Head()) + end
	}
	return e.Body()[ini:end]
}

var WORD = regexp.MustCompile(`^\w+`)
var DIGITS = regexp.MustCompile(`^\d+`)
var ST = regexp.MustCompile(`^[ \t]+`)
var WS = regexp.MustCompile(`^\s+`)
