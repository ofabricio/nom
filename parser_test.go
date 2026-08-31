package nom

import (
	"fmt"
	"unicode"
)

func ExampleParser_Expect() {

	p := New("fun{}")
	ok := p.Expect("fun") && p.Expect("{") && p.Expect("}")
	fmt.Println("--- no expected error ---")
	fmt.Println(ok, p.Err)

	p = New("fun()")
	ok = p.Expect("fun") && p.Expect("{") && p.Expect("}")
	fmt.Println("--- default expected error ---")
	fmt.Println(ok)
	fmt.Println(p.Err)

	p = New("fun()")
	ok = p.Expect("fun") && p.ExpectWith("{", "open brackets") && p.Expect("}")
	fmt.Println("--- custom expected error ---")
	fmt.Println(ok)
	fmt.Println(p.Err)

	p = New("arg(abc)")
	ok = p.Expect("arg") && p.Expect("(") && p.Expect(DIGITS) && p.Expect(")")
	fmt.Println("--- arg example ---")
	fmt.Println(ok)
	fmt.Println(p.Err)

	p = New("\n\t\tfunc\n")
	ok = p.Match(WS) && p.Match("func") && p.Expect(ST)
	fmt.Println("--- fun example ---")
	fmt.Println(ok)
	fmt.Println(p.Err)

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
	// --- arg example ---
	// false
	// failed to parse: line 1 char 5: expected ^\d+
	//       |
	//     1 | arg(abc)
	//       |     ^--
	// --- fun example ---
	// false
	// failed to parse: line 2 char 7: expected ^[ \t]+
	//       |
	//     2 |         func
	//       |             ^--
}

func ExampleParser_MatchFunc() {

	p := New("1a")
	fmt.Println(p.Match(unicode.IsDigit), p.Tail() == "a")
	fmt.Println(p.Match(unicode.IsDigit), p.Tail() == "a")

	// Output:
	// true true
	// false true
}
