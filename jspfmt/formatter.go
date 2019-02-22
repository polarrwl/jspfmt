package jspfmt

import "fmt"

// Format returns a pretty formatted version of the input JSP string.
func Format(name, input string) {
	tokens := make(chan token)
	l := lexer{
		name:   name,
		input:  input,
		tokens: tokens,
	}
	go l.run()
	for t := range tokens {
		fmt.Println(t)
	}
}
