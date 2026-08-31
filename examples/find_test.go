package lang

import (
	"fmt"

	. "github.com/ofabricio/nom"
)

// ExampleFind shows the usage of the Find parser to collect
// either the section numbers or the years in the input text.
func ExampleFind() {

	src := `
		The standard chunk of Lorem Ipsum used since 1966 is reproduced below
		for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus
		Bonorum et Malorum" by Cicero are also reproduced in their exact
		original form, accompanied by English versions from the 1914
		translation by H. Rackham.
	`

	var out []string
	s := New(src)
	for s.More() {
		if m := s.Mark(); s.Match(DIGITS) && s.Opt(s.Mark(), s.Match(".") && s.Match(DIGITS) && s.Match(".") && s.Match(DIGITS)) {
			out = append(out, s.Token(m).Text)
			continue
		}
		s.Next()
	}

	fmt.Println(s.Err, out)

	// Output:
	// <nil> [1966 1.10.32 1.10.33 1914]
}

// ExampleFind shows the usage of the Find parser to collect
// either the section numbers or the years in the input text.
func ExampleParser_find() {

	src := `
		The standard chunk of Lorem Ipsum used since 1966 is reproduced below
		for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus
		Bonorum et Malorum" by Cicero are also reproduced in their exact
		original form, accompanied by English versions from the 1914
		translation by H. Rackham.
	`

	var out []string

	dgts := M(DIGITS)
	year := dgts
	sect := And(dgts, M("."), dgts, M("."), dgts)
	eith := Or(sect, year).On(Grabs(&out))
	root := Find(eith).ZeroToMany()

	err := root.Parse(src)

	fmt.Println(err, out)

	// Output:
	// <nil> [1966 1.10.32 1.10.33 1914]
}
