package lang

import (
	"fmt"
	"strings"

	. "github.com/ofabricio/nom"
)

// ExampleLang shows how to parse a programming language.
func ExampleLang() {

	src := `
		fun main () {
			aaa()
			bbb()
		}

		fun ccc () {
			ddd()
			fun eee () {
				fff()
			}
		}
		
		let x = 1
		let y = 2
	`

	var out Lang
	var lng LangParserLow
	err := lng.Parse(src, &out)

	fmt.Println(err)
	LangPrint(out, 0)

	// Output:
	// <nil>
	// Function: main [
	//     Call: aaa
	//     Call: bbb
	// ]
	// Function: ccc [
	//     Call: ddd
	//     Function: eee [
	//         Call: fff
	//     ]
	// ]
	// Assignment: x = 1
	// Assignment: y = 2
}

type LangParserLow struct {
	Parser
}

func (p *LangParserLow) Parse(src string, out *Lang) error {
	p.Parser = New(src)
	p.Program(out)
	return p.Err
}

func (p *LangParserLow) Program(out *Lang) bool {
	for p.Optional(WS) {
		if v := (LangFunctionDef{}); p.FunctionDef(&v) {
			out.D = append(out.D, v)
			continue
		}
		if v := (LangAssignment{}); p.Assignment(&v) {
			out.A = append(out.A, v)
			continue
		}
		break
	}
	return p.More() && p.Expected("EOF")
}

func (p *LangParserLow) FunctionDef(out *LangFunctionDef) bool {
	return p.Match("fun") &&
		p.Expect(ST) &&
		p.ExpectOut(WORD, &out.Name) &&
		p.Optional(ST) &&
		p.Expect("(") &&
		p.Expect(")") &&
		p.Optional(ST) &&
		p.Expect("{") &&
		p.FunctionBody(&out.Body) &&
		p.Expect("}")
}

func (p *LangParserLow) FunctionBody(out *[]LangStatement) bool {
	var s LangStatement
	for p.Optional(WS) && p.Statement(&s) {
		*out = append(*out, s)
		s = LangStatement{}
	}
	return true
}

func (p *LangParserLow) Statement(out *LangStatement) bool {
	if v := (LangFunctionDef{}); p.FunctionDef(&v) {
		out.D = &v
		return true
	}
	if v := (LangFunctionCall{}); p.FunctionCall(&v) {
		out.C = &v
		return true
	}
	return false
}

func (p *LangParserLow) FunctionCall(out *LangFunctionCall) bool {
	return p.MatchOut(WORD, &out.Name) &&
		p.Expect("(") &&
		p.Expect(")")
}

func (p *LangParserLow) Assignment(out *LangAssignment) bool {
	return p.Match("let") &&
		p.Expect(ST) &&
		p.ExpectOut(WORD, &out.Name) &&
		p.Optional(ST) &&
		p.Expect("=") &&
		p.Optional(ST) &&
		p.ExpectOut(DIGITS, &out.Value)
}

type Lang struct {
	D []LangFunctionDef
	A []LangAssignment
}

type LangFunctionDef struct {
	Name Token
	Body []LangStatement
}

type LangStatement struct {
	C *LangFunctionCall
	D *LangFunctionDef
}

type LangFunctionCall struct {
	Name Token
}

type LangAssignment struct {
	Name  Token
	Value Token
}

func LangPrint(v any, depth int) {
	switch v := v.(type) {
	case Lang:
		for _, f := range v.D {
			LangPrint(f, depth)
		}
		for _, a := range v.A {
			LangPrint(a, depth)
		}
	case LangFunctionDef:
		LangPrint(fmt.Sprintf("Function: %s [", v.Name.Text), depth)
		for _, s := range v.Body {
			LangPrint(s, depth+1)
		}
		LangPrint("]", depth)
	case LangAssignment:
		LangPrint(fmt.Sprintf("Assignment: %s = %s", v.Name.Text, v.Value.Text), depth)
	case LangStatement:
		if v.C != nil {
			LangPrint(*v.C, depth)
		}
		if v.D != nil {
			LangPrint(*v.D, depth)
		}
	case LangFunctionCall:
		LangPrint(fmt.Sprintf("Call: %s", v.Name.Text), depth)
	default:
		fmt.Print(strings.Repeat("    ", depth))
		fmt.Printf("%v\n", v)
	}
}
