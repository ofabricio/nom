package lang

import (
	"fmt"
	"strconv"

	. "github.com/ofabricio/nom"
)

// ExampleExpr shows parsing and evaluating a mathematical expression.
func ExampleExpr() {

	src := "1+2*(3+4)*5"

	var out Expr
	var exp ExprParser
	err := exp.Parse(src, &out)

	fmt.Println(err)
	fmt.Println(out.Calc())

	// Output:
	// <nil>
	// 71
}

type ExprParser struct {
	Parser
}

func (p *ExprParser) Parse(src string, out *Expr) error {
	p.Parser = New(src)
	p.Expr(out)
	return p.Err
}

func (p *ExprParser) Expr(out *Expr) bool {
	var l, r Expr
	if p.Term(&l) {
		if p.MatchOut("+", &out.V) && p.Expr(&r) {
			out.L = &l
			out.R = &r
			return true
		}
		*out = l
		return true
	}
	return p.Expect("expression")
}

func (p *ExprParser) Term(out *Expr) bool {
	var l, r Expr
	if p.Fact(&l) {
		if p.MatchOut("*", &out.V) && p.Term(&r) {
			out.L = &l
			out.R = &r
			return true
		}
		*out = l
		return true
	}
	return p.Expect("expression")
}

func (p *ExprParser) Fact(out *Expr) bool {
	return p.Match("(") && p.Expr(out) && p.Expect(")") || p.MatchOut(DIGITS, &out.V)
}

type Expr struct {
	V    Token
	L, R *Expr
}

func (e *Expr) Calc() int {
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
