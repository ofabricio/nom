package nom

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"unicode"
)

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
		ok := p.Match(tc.GiveM) && p.Expect(tc.GiveE)
		if ok != tc.ThenOk {
			t.Errorf("\nMsg: %s\nGot:\n%v\nExp:\n%v\n", tc.Descr, ok, tc.ThenOk)
		}
		exp := strings.Join(tc.Then, "\n")
		got := fmt.Sprint(p.Err)
		if got != exp {
			t.Errorf("\nMsg: %s\nGot:\n%v\nExp:\n%v\n", tc.Descr, got, exp)
		}
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
