package jspfmt

import (
	"fmt"
)

const debug = false

// Format returns a pretty formatted version of the input JSP string.
func Format(name, input string) {
	tokens := make(chan token)
	l := lexer{
		name:   name,
		input:  input,
		tokens: tokens,
	}
	go l.run()
	depth := 0
	tabSize := 4
	if debug {
		printDebug(tokens)
		return
	}
	for t := range tokens {
		switch t.typ {
		case tokOpenTag:
			fmt.Printf("%*s%s\n", depth*tabSize, "", t.val)
			depth++
		case tokCloseTag:
			depth--
			fmt.Printf("%*s%s\n", depth*tabSize, "", t.val)
		case tokSelfClosingTag, tokText:
			fmt.Printf("%*s%s\n", depth*tabSize, "", t.val)
		case tokError:
			fmt.Println(t.val)
		}
	}
}

func printDebug(tokens chan token) {
	for t := range tokens {
		switch t.typ {
		case tokOpenTag:
			fmt.Println("open tag: " + t.val)
		case tokCloseTag:
			fmt.Println("close tag: " + t.val)
		case tokSelfClosingTag:
			fmt.Println("self-close tag: " + t.val)
		case tokText:
			fmt.Println("text: " + t.val)
		case tokError:
			fmt.Println("error: " + t.val)
		}
	}
}
