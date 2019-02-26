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

	tokLessThan     // Literal "<"
	tokGreaterThan  // Literal ">"
	tokLessSlash    // Literal "</"
	tokSlashGreater // Literal "/>"
	tokTagName      // Name of the tag.
	tokAttrKey      // Name of the attribute in the tag.
	tokEquals       // Literal "="
	tokAttrVal      // Value of the attribute, including quotes.
	tokText         // Uninterpreted text inside tags, e.g. <tagname attr="val">[text]</tagname>
)

const (
	lessThan     string = "<"
	greaterThan  string = ">"
	lessSlash    string = "</"
	slashGreater string = "/>"
	equalsSign   string = "="
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
	}
	// if len(t.val) > 10 {
	// 	return fmt.Sprintf("%v: %.10q...", t.typ, t.val)
	// }
	return fmt.Sprintf("%s: %q", t.name(), t.val)
}

func (t token) name() string {
	switch t.typ {
	case tokError:
		return "error"
	case tokEOF:
		return "EOF"
	case tokLessThan:
		return "less than"
	case tokGreaterThan:
		return "greater than"
	case tokLessSlash:
		return "less slash"
	case tokSlashGreater:
		return "slash greater"
	case tokTagName:
		return "tag name"
	case tokAttrKey:
		return "attr key"
	case tokEquals:
		return "equals sign"
	case tokAttrVal:
		return "attr value"
	case tokText:
		return "text"
	default:
		return "UNKNOWN"
	}
}
