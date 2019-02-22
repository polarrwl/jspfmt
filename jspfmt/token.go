package jspfmt

import "fmt"

// tokenType identifies the type of lexer tokens.
type tokenType int32

type token struct {
	typ tokenType
	val string
}

const (
	tokError tokenType = iota // error occurred; value is the text of the error.
	tokEOF                    // end of file

	tokOpenTag
	tokCloseTag
	tokSelfClosingTag
	tokText
)

const (
	eof rune = 0
)

func (t token) String() string {
	switch t.typ {
	case tokEOF:
		return "EOF"
	case tokError:
		return "Error: " + t.val
	case tokOpenTag:
		return fmt.Sprintf("Open Tag: %q", t.val)
	case tokCloseTag:
		return fmt.Sprintf("Close Tag: %q", t.val)
	case tokSelfClosingTag:
		return fmt.Sprintf("Self-closing Tag: %q", t.val)
	case tokText:
		return fmt.Sprintf("Text: %q", t.val)
	}
	// if len(t.val) > 10 {
	// 	return fmt.Sprintf("%v: %.10q...", t.typ, t.val)
	// }
	return fmt.Sprintf("UNKNOWN: %q", t.val)
}
