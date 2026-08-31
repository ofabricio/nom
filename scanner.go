package nom

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

func New(src string) Scanner {
	return Scanner{src: src, Row: 1, Col: 1}
}

func (s *Scanner) MatchOut[T MatchType](v T, out *Token) bool {
	m := s.Mark()
	return s.Match(v) && s.Out(m, out)
}

func (s *Scanner) ExpectOut[T MatchType](v T, out *Token) bool {
	m := s.Mark()
	return s.Expect(v) && s.Out(m, out)
}

func (s *Scanner) Out(m Scanner, out *Token) bool {
	*out = s.Token(m)
	return true
}

func (s *Scanner) Opt(m Scanner, cond bool) bool {
	return s.Undo(m, cond) || true
}

// Undo sends the parser back to the mark m if the condition is false.
func (s *Scanner) Undo(m Scanner, cond bool) bool {
	if !cond {
		s.Back(m)
	}
	return cond
}

func (s *Scanner) Optional[T MatchType](v T) bool {
	return s.Match(v) || true
}

func (s *Scanner) Expect[T MatchType](v T) bool {
	return s.Match(v) || s.expect(v)
}

func (s *Scanner) Expected(msg string) bool {
	return s.expect(msg)
}

func (s *Scanner) ExpectWith[T MatchType](v T, msg string) bool {
	return s.Match(v) || s.expect(msg)
}

func (s *Scanner) expect[T MatchType](v T) bool {
	switch v := any(v).(type) {
	case string:
		s.Err = &Error{Scanner: *s, Msg: v}
	case *regexp.Regexp:
		s.Err = &Error{Scanner: *s, Msg: v.String()}
	default:
		s.Err = &Error{Scanner: *s, Msg: "token"}
	}
	return false
}

func (s *Scanner) Match[T MatchType](v T) bool {
	switch v := any(v).(type) {
	case string:
		return s.MatchString(v)
	case *regexp.Regexp:
		return s.MatchRegex(v)
	case func(rune) bool:
		return s.MatchFunc(v)
	default:
		return false
	}
}

func (s *Scanner) MatchString(v string) bool {
	return s.EqualString(v) && s.advance(v)
}

func (s *Scanner) MatchRegex(v *regexp.Regexp) bool {
	return s.advance(v.FindString(s.Tail()))
}

func (s *Scanner) MatchFunc(f func(rune) bool) bool {
	r := s.Curr()
	return f(r) && s.advance(string(r))
}

func (s *Scanner) Equal[T MatchType](v T) bool {
	switch v := any(v).(type) {
	case string:
		return s.EqualString(v)
	case *regexp.Regexp:
		return s.EqualRegex(v)
	case func(rune) bool:
		return s.EqualFunc(v)
	default:
		return false
	}
}

func (s *Scanner) EqualString(v string) bool {
	return strings.HasPrefix(s.Tail(), v)
}

func (s *Scanner) EqualRegex(v *regexp.Regexp) bool {
	return v.MatchString(s.Tail())
}

func (s *Scanner) EqualFunc(f func(rune) bool) bool {
	return f(s.Curr())
}

func (s *Scanner) Next() bool {
	return s.advance(string(s.Curr()))
}

func (s *Scanner) Curr() rune {
	r, _ := utf8.DecodeRuneInString(s.Tail())
	return r
}

func (s *Scanner) Head() string {
	return s.src[:s.Idx]
}

func (s *Scanner) Tail() string {
	return s.src[s.Idx:]
}

func (s *Scanner) Body() string {
	return s.src
}

func (s *Scanner) advance(v string) bool {
	s.Idx += len(v)
	s.coln(v)
	return len(v) > 0
}

func (s *Scanner) coln(v string) {
	for _, r := range v {
		s.Col++
		if r == '\n' {
			s.Row++
			s.Col = 1
		}
	}
}

func (s Scanner) Mark() Scanner {
	return s
}

func (s *Scanner) Back(m Scanner) {
	*s = m
}

func (s *Scanner) Token(m Scanner) Token {
	return Token{Text: s.src[m.Idx:s.Idx], Idx: m.Idx, Row: m.Row, Col: m.Col}
}

func (s Scanner) More() bool {
	return len(s.Tail()) > 0
}

type Scanner struct {
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

var WORD = regexp.MustCompile(`^\w+`)
var DIGITS = regexp.MustCompile(`^\d+`)
var ST = regexp.MustCompile(`^[ \t]+`)
var WS = regexp.MustCompile(`^\s+`)

// Error represents a parsing error.
type Error struct {
	Scanner
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
