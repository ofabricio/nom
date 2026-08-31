package nom

import (
	"fmt"
	"strconv"
	"unicode"
)

func ExampleScanner_Expect() {

	s := New("fun{}")
	ok := s.Expect("fun") && s.Expect("{") && s.Expect("}")
	fmt.Println("--- no expected error ---")
	fmt.Println(ok, s.Err)

	s = New("fun()")
	ok = s.Expect("fun") && s.Expect("{") && s.Expect("}")
	fmt.Println("--- default expected error ---")
	fmt.Println(ok)
	fmt.Println(s.Err)

	s = New("fun()")
	ok = s.Expect("fun") && s.ExpectWith("{", "open brackets") && s.Expect("}")
	fmt.Println("--- custom expected error ---")
	fmt.Println(ok)
	fmt.Println(s.Err)

	// Output:
	// --- no expected error ---
	// true <nil>
	// --- default expected error ---
	// false
	// failed to parse: line 1 char 4: expected {
	//       |
	//     1 | fun()
	//       |    ^--
	// --- custom expected error ---
	// false
	// failed to parse: line 1 char 4: expected open brackets
	//       |
	//     1 | fun()
	//       |    ^--
}

func ExampleScanner_MatchFunc() {

	s := New("1a")
	fmt.Println(s.Match(unicode.IsDigit), s.Tail() == "a")
	fmt.Println(s.Match(unicode.IsDigit), s.Tail() == "a")

	// Output:
	// true true
	// false true
}

func Example_expression() {

	s := exprParser{s: New("1+2*(3+4)*5")}

	var e expr
	ok := s.Expr(&e)
	fmt.Println(ok)
	fmt.Println(e.Calc())

	// Output:
	// true
	// 71
}

type exprParser struct {
	s Scanner
}

func (p *exprParser) Expr(out *expr) bool {
	var l, o, r expr
	if p.Term(&l) {
		if p.s.MatchOut("+", &o.V) && p.Expr(&r) {
			out.V = o.V
			out.L = &l
			out.R = &r
			return true
		}
		*out = l
		return true
	}
	return false
}

func (p *exprParser) Term(out *expr) bool {
	var l, o, r expr
	if p.Fact(&l) {
		if p.s.MatchOut("*", &o.V) && p.Term(&r) {
			out.V = o.V
			out.L = &l
			out.R = &r
			return true
		}
		*out = l
		return true
	}
	return false
}

func (p *exprParser) Fact(out *expr) bool {
	return p.s.Match("(") && p.Expr(out) && p.s.Match(")") || p.Number(out)
}

func (p *exprParser) Number(out *expr) bool {
	return p.s.MatchOut(DIGITS, &out.V)
}

type expr struct {
	L *expr
	V Token
	R *expr
}

func (e *expr) Calc() int {
	switch e.V.Text {
	case "+":
		return e.L.Calc() + e.R.Calc()
	case "*":
		return e.L.Calc() * e.R.Calc()
	default:
		v, _ := strconv.Atoi(e.V.Text)
		return v
	}
}
