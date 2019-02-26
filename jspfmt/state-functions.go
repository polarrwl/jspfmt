package jspfmt

import (
	"strings"
)

type stateFn func(*lexer) stateFn

// cf. https://www.w3.org/TR/html/infrastructure.html#space-characters
const spaceChars = " \f\t\n\r\v"

// lexText is the default state. It looks for the next tag.
func lexText(l *lexer) stateFn {
	l.acceptRun(spaceChars)
	l.ignore()

	for {
		if strings.HasPrefix(l.input[l.cursor:], lessSlash) {
			if l.cursor > l.start {
				l.emit(tokText)
			}
			return lexTagClose
		}
		if strings.HasPrefix(l.input[l.cursor:], lessThan) {
			if l.cursor > l.start {
				l.emit(tokText)
			}
			return lexTagOpen
		}
		// After the checks, advance the cursor. If we've reached the
		// end of input, break out of the loop.
		if l.next() == eof {
			break
		}
	}
	// Reached EOF.
	if l.cursor > l.start {
		l.emit(tokText)
	}
	l.emit(tokEOF)
	return nil
}

// lexTagOpen assumes l.cursor is sitting on a literal "<" and scans until after
// the tag name.
func lexTagOpen(l *lexer) stateFn {
	// step inside
	l.cursor += len(lessThan)
	l.emit(tokLessThan)

	// next token must be tagname
	if l.peek() == eof {
		l.errorf("unexpected EOF: missing tag name")
		return nil
	}

	// bang included for doctype tag
	l.acceptRunRegexp("[0-9a-zA-Z!]")
	if l.cursor == l.start {
		l.errorf("%d: bad character %q after %q", l.start, l.peek(), lessThan)
		return nil
	}
	l.emit(tokTagName)
	return lexAttributes
}

// lexTagClose assumes that l.cursor is sitting on a literal "</" and scans until
// the end of a closing tag.
func lexTagClose(l *lexer) stateFn {
	// step inside
	l.cursor += len(lessSlash)
	l.emit(tokLessSlash)

	// next token must be a tagname
	if l.peek() == eof {
		l.errorf("unexpected EOF: missing tag name")
		return nil
	}

	// cf. https://www.w3.org/TR/html/syntax.html#tag-name
	l.acceptRunRegexp("[0-9a-zA-Z]")
	if l.cursor == l.start {
		l.errorf("%d: bad character %q after %q", l.start, l.peek(), lessSlash)
		return nil
	}
	l.emit(tokTagName)

	// gobble any whitespace
	l.acceptRun(spaceChars)
	l.ignore()

	// the next part of the input must be tagCloseRight
	if !strings.HasPrefix(l.input[l.cursor:], greaterThan) {
		l.errorf("unexpected token %s", string(l.peek()))
		return nil
	}
	l.cursor += len(greaterThan)
	l.emit(tokSlashGreater)
	return lexText
}

// lexAttributes is called when the cursor is inside an open tag after the
// tag name or previous attributes. It scans until a literal "/>", literal ">",
// or an attribute.
func lexAttributes(l *lexer) stateFn {
	l.acceptRun(spaceChars)

	if strings.HasPrefix(l.input[l.cursor:], slashGreater) {
		// ignore the whitespace
		l.ignore()
		l.cursor += len(slashGreater)
		l.emit(tokSlashGreater)
		return lexText
	}

	if strings.HasPrefix(l.input[l.cursor:], greaterThan) {
		// ignore the whitespace
		l.ignore()
		l.cursor += len(greaterThan)
		l.emit(tokGreaterThan)
		return lexText
	}

	if l.peek() == eof {
		l.errorf("unexpected EOF: tag did not end")
		return nil
	}

	if l.start == l.cursor {
		l.errorf("at least one space character required before attribute key")
		return nil
	}
	l.ignore() // no token for the space characters

	// cf. https://www.w3.org/TR/html/syntax.html#elements-attributes
	l.acceptRunRegexp("[^ \f\t\n\r\v\"'/>=[:cntrl:]\ufdd0-\ufddf]")
	if l.cursor == l.start {
		l.errorf("%d: bad character %q in attribute name", l.start, l.peek())
		return nil
	}
	l.emit(tokAttrKey)

	l.acceptRun(spaceChars)
	l.ignore()

	if !l.accept(equalsSign) { // empty attribute
		return lexAttributes
	}
	l.emit(tokEquals)

	l.acceptRun(spaceChars)
	l.ignore()

	// There are three valid formats for attribute values whether enclosed
	// within double quotes, single quotes, or no quotes. Each has separate
	// rules about allowed characters. The quotes are preserved in the
	// token value.
	switch {
	case l.accept("\""):
		l.acceptRunRegexp("[^\"[:cntrl:]]")
		if l.start == l.cursor {
			l.errorf("missing attribute value")
			return nil
		}
		if !l.accept("\"") {
			l.errorf("missing closing quote in attribute value")
			return nil
		}
	case l.accept("'"):
		l.acceptRunRegexp("[^'[:cntrl:]]")
		if l.start == l.cursor {
			l.errorf("missing attribute value")
			return nil
		}
		if !l.accept("'") {
			l.errorf("missing closing apostrophe in attribute value")
			return nil
		}
	default:
		l.acceptRunRegexp("[^ \f\t\n\r\v\"'=<>`[:cntrl:]]")
		if l.start == l.cursor {
			l.errorf("missing attribute value")
			return nil
		}
	}
	l.emit(tokAttrVal)
	return lexAttributes
}
