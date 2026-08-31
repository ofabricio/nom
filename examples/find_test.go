package lang

import (
	"fmt"

	. "github.com/ofabricio/nom"
)

// Example_find shows the usage of the parser to collect tokens.
func Example_find() {

	src := `
		The standard chunk of Lorem Ipsum used since 1966 is reproduced below
		for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus
		Bonorum et Malorum" by Cicero are also reproduced in their exact
		original form, accompanied by English versions from the 1914
		translation by H. Rackham.
	`

	var out []string

	p := New(src)
	for p.More() {
		if m := p.Mark(); p.Match(DIGITS) && p.Opt(p.Mark(), p.Match(".") && p.Match(DIGITS) && p.Match(".") && p.Match(DIGITS)) {
			out = append(out, p.Token(m).Text)
			continue
		}
		p.Next()
	}

	fmt.Println(p.Err, out)

	// Output:
	// <nil> [1966 1.10.32 1.10.33 1914]
}
