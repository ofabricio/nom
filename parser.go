package nom

func Grab(out *string) Outer {
	return func(t Token) {
		*out = t.Text
	}
}

func Grabs(out *[]string) Outer {
	return func(t Token) {
		*out = append(*out, t.Text)
	}
}

// M matches the input against the given value.
func M[T MatchType](v T) Parser {
	return func(c *Context) bool {
		return c.s.Match(v)
	}
}

// T tests if the input equals the given value.
func T[T MatchType](v T) Parser {
	return func(c *Context) bool {
		return c.s.Equal(v)
	}
}

// O makes a parser optional.
func O[T MatchType](v T) Parser {
	return func(c *Context) bool {
		return c.s.Match(v) || true
	}
}

// E expects the input to match the given value.
func E[T MatchType](v T) Parser {
	return func(c *Context) bool {
		return c.s.Expect(v)
	}
}

// P is a pointer to a parser, useful for recursive calls.
func P(m *Parser) Parser {
	return func(c *Context) bool {
		return (*m)(c)
	}
}

// Any matches any single token from the input.
func Any() Parser {
	return func(c *Context) bool {
		return c.s.Next()
	}
}

func Find(f Parser) Parser {
	return func(c *Context) bool {
		for c.s.More() && !f(c) && c.s.Next() {
		}
		return c.s.More()
	}
}

// And matches all of the given parsers in sequence.
func And(and ...Parser) Parser {
	return func(c *Context) bool {
		m := c.s.Mark()
		for _, fn := range and {
			if !fn(c) {
				m.Err = c.s.Err
				c.s.Back(m)
				return false
			}
		}
		return true
	}
}

// Or matches one of the given parsers.
func Or(or ...Parser) Parser {
	return func(c *Context) bool {
		for _, fn := range or {
			if fn(c) {
				return true
			}
		}
		return false
	}
}

// EOF tests for the end of the input.
func EOF() Parser {
	return func(c *Context) bool {
		return !c.s.More()
	}
}

// Parser is the entry point for parsing an input.
func (f Parser) Parse(src string) error {
	c := &Context{s: New(src)}
	f(c)
	return c.s.Err
}

// Out outputs the matched token.
func (f Parser) Out(out *Token) Parser {
	return f.On(func(tk Token) {
		*out = tk
	})
}

// Set assigns v to out, then resets v to its zero value.
func (f Parser) Set[T any](v *T, out **T) Parser {
	return f.On(func(Token) {
		var zero T
		*out = new(*v)
		*v = zero
	})
}

// Add appends v to out, then resets v to its zero value.
func (f Parser) Add[T any](v *T, out *[]T) Parser {
	return f.On(func(Token) {
		var zero T
		*out = append(*out, *v)
		*v = zero
	})
}

// On executes the given function when a parser successfully matches.
func (f Parser) On(fn Outer) Parser {
	return func(c *Context) bool {
		if m := c.s.Mark(); f(c) {
			fn(c.s.Token(m))
			return true
		}
		return false
	}
}

// ZeroToMany applies a parser zero or more times.
func (f Parser) ZeroToMany() Parser {
	return func(c *Context) bool {
		v := 0
		for f(c) {
			v++
		}
		return v >= 0 && c.s.Err == nil
	}
}

// OneToMany applies a parser one or more times.
func (f Parser) OneToMany() Parser {
	return func(c *Context) bool {
		v := 0
		for f(c) {
			v++
		}
		return v >= 1
	}
}

// True makes a parser always succeed.
func (f Parser) True() Parser {
	return func(c *Context) bool {
		return f(c) || true
	}
}

// False makes a parser always fail.
func (f Parser) False() Parser {
	return func(c *Context) bool {
		return f(c) && false
	}
}

// Expect makes a parser triggers an error
// with the expected message on failure.
//
//	M(DIGITS).Expect("number").Parse("hello")
func (f Parser) Expect(expected string) Parser {
	return func(c *Context) bool {
		return f(c) || c.s.expect(expected)
	}
}

// Parser is a parsing function that operates on a context
// and returns a boolean indicating success.
type Parser func(*Context) bool

type Outer func(Token)

// Context holds the current state of a parser.
type Context struct {
	s Scanner
}
