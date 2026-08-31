package nom

import (
	"fmt"
)

func ExampleParser_Expect() {

	fmt.Println("--- arg example ---")
	fmt.Println(And(M("arg"), M("("), M(DIGITS).Expect("number")).Parse("arg(abc)"))

	fmt.Println("--- fun example ---")
	fmt.Println(And(O(WS), M("func"), E(ST)).Parse("\n\t\tfunc\n"))

	// Output:
	// --- arg example ---
	// failed to parse: line 1 char 5: expected number
	//       |
	//     1 | arg(abc)
	//       |     ^--
	// --- fun example ---
	// failed to parse: line 2 char 7: expected ^[ \t]+
	//       |
	//     2 |         func
	//       |             ^--
}

func ExampleParser_expression() {

	src := "1+(2*(3+4))"

	var expr, term, fact Parser
	expr = And(P(&term), And(M("+"), P(&expr)).ZeroToMany())
	term = And(P(&fact), And(M("*"), P(&term)).ZeroToMany())
	fact = Or(And(M("("), expr, M(")")), M(DIGITS))

	err := expr.Parse(src)

	fmt.Println(err)

	// Output:
	// <nil>
}
