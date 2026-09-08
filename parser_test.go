package nom

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

func ExampleParser_String() {

	p := New(`
		''
		'\''
		'one'
		'\'one\''
		'one\'two\'three'
		'\'one\'two\'three\''
	`)

	for p.Opt(WS) && p.More() {
		if m := p.Mark(); p.String("'") {
			fmt.Println(p.Token(m).Text)
			continue
		}
		break
	}

	p = New(`"abc`)
	fmt.Println("---")
	fmt.Println(p.String("\""), p.Err)

	// Output:
	// ''
	// '\''
	// 'one'
	// '\'one\''
	// 'one\'two\'three'
	// '\'one\'two\'three\''
	// ---
	// false failed to parse: line 1 char 5: expected "
	//       |
	//     1 | "abc
	//       |     ^--
}

func ExampleParser_Line() {

	p := New("\naa\nb\n\n")

	var t Token
	for p.Out(p.Mark(), p.Line(), &t) && p.Match("\n") {
		fmt.Printf("Idx=%1d Row=%1d Col=%1d Text='%s'\n", t.Idx, t.Row, t.Col, t.Text)
	}

	fmt.Println("---")
	p = New("\naa\nb\n\n")

	for ; p.Line(); p.Match("\n") {
		t := p.GetLine()
		fmt.Printf("Idx=%1d Row=%1d Col=%1d Text='%s'\n", t.Idx, t.Row, t.Col, t.Text)
	}

	// Output:
	// Idx=0 Row=1 Col=1 Text=''
	// Idx=1 Row=2 Col=1 Text='aa'
	// Idx=4 Row=3 Col=1 Text='b'
	// Idx=6 Row=4 Col=1 Text=''
	// ---
	// Idx=0 Row=1 Col=1 Text=''
	// Idx=1 Row=2 Col=1 Text='aa'
	// Idx=4 Row=3 Col=1 Text='b'
	// Idx=6 Row=4 Col=1 Text=''
}

func ExampleParser_Find() {

	p := New(`Coffee is $5, but he sold for $4.`)

	money := regexp.MustCompile(`^\$\d+`)

	var out Token
	for p.Find(money) && p.MatchOut(money, &out) {
		fmt.Println(out.Text)
	}

	fmt.Println(p.Err)

	// Output:
	// $5
	// $4
	// <nil>
}

func ExampleParser_FindOut() {

	p := New(`Coffee is $5, but he sold for $4.`)

	money := regexp.MustCompile(`^\$\d+`)

	var out Token
	for p.FindOut(money, &out) {
		fmt.Println(out.Text)
	}

	fmt.Println(p.Err)

	// Output:
	// $5
	// $4
	// <nil>
}

func TestParserExpectedErrorMsg(t *testing.T) {

	tt := []struct {
		Descr  string
		GiveI  string // Input
		GiveM  string // Match
		GiveE  string // Expect
		Then   []string
		ThenOk bool
	}{
		{
			Descr:  "just one line",
			GiveI:  "111",
			GiveM:  "11",
			GiveE:  "x",
			ThenOk: false,
			Then: []string{
				"failed to parse: line 1 char 3: expected x",
				"      |",
				"    1 | 111",
				"      |   ^--",
			},
		},
		{
			Descr:  "first line on the end edge",
			GiveI:  "111\n",
			GiveM:  "111",
			GiveE:  "x",
			ThenOk: false,
			Then: []string{
				"failed to parse: line 1 char 4: expected x",
				"      |",
				"    1 | 111",
				"      |    ^--",
			},
		},
		{
			Descr:  "second line with start edge",
			GiveI:  "\n2222",
			GiveM:  "\n22",
			GiveE:  "x",
			ThenOk: false,
			Then: []string{
				"failed to parse: line 2 char 3: expected x",
				"      |",
				"    2 | 2222",
				"      |   ^--",
			},
		},
		{
			Descr:  "second line of three",
			GiveI:  "111\n2222\n33",
			GiveM:  "111\n22",
			GiveE:  "x",
			ThenOk: false,
			Then: []string{
				"failed to parse: line 2 char 3: expected x",
				"      |",
				"    2 | 2222",
				"      |   ^--",
			},
		},
		{
			Descr:  "three empty lines",
			GiveI:  "\n\n\n",
			GiveM:  "\n",
			GiveE:  "x",
			ThenOk: false,
			Then: []string{
				"failed to parse: line 2 char 1: expected x",
				"      |",
				"    2 | ",
				"      | ^--",
			},
		},
		{
			Descr:  "should format \t accordingly",
			GiveI:  "\n\t\t2222\n",
			GiveM:  "\n\t\t22",
			GiveE:  "x",
			ThenOk: false,
			Then: []string{
				"failed to parse: line 2 char 5: expected x",
				"      |",
				"    2 |         2222",
				"      |           ^--",
			},
		},
	}

	for _, tc := range tt {
		p := New(tc.GiveI)
		ok := p.Match(tc.GiveM) && p.Exp(tc.GiveE)
		assert(t, ok == tc.ThenOk, fmt.Sprint(ok), fmt.Sprint(tc.ThenOk), tc.Descr)
		exp := strings.Join(tc.Then, "\n")
		got := fmt.Sprint(p.Err)
		assert(t, got == exp, got, exp, tc.Descr)
	}
}

func ExampleParser_MatchFunc() {

	p := New("1a")
	fmt.Println(p.Match(unicode.IsDigit), p.Tail() == "a")
	fmt.Println(p.Match(unicode.IsDigit), p.Tail() == "a")

	// Output:
	// true true
	// false true
}

func ExampleParser_Rune() {

	p := New("d😊b")
	fmt.Println(p.Rune() == 'd')
	p.Next()
	fmt.Println(p.Rune() == '😊')
	p.Next()
	fmt.Println(p.Rune() == 'b')
	p.Next()
	fmt.Println(p.Rune() == 0)

	// Output:
	// true
	// true
	// true
	// true
}

func ExampleParser_Char() {

	p := New("d😊b")
	fmt.Println(p.Char() == "d")
	p.Next()
	fmt.Println(p.Char() == "😊")
	p.Next()
	fmt.Println(p.Char() == "b")
	p.Next()
	fmt.Println(p.Char() == "")

	// Output:
	// true
	// true
	// true
	// true
}

func assert(t *testing.T, cond bool, got, exp string, msg string) {
	if !cond {
		t.Errorf("\nMsg: %s\nGot:\n%v\nExp:\n%v\n", msg, got, exp)
	}
}
