package jspfmt

import (
	"strings"
)

type stateFn func(*lexer) stateFn

// lexHtml accepts text until a leftMeta is found
func lexHTML(l *lexer) stateFn {
	l.acceptRun(" \n")
	l.ignore()
	for {
		if strings.HasPrefix(l.input[l.cursor:], "</") {
			if l.cursor > l.start {
				l.emit(tokText)
			}
			return lexCloseTag
		}
		if strings.HasPrefix(l.input[l.cursor:], "<") {
			if l.cursor > l.start {
				l.emit(tokText)
			}
			return lexOpenTag
		}
		if l.next() == eof {
			break
		}
	}
	if l.cursor > l.start {
		l.emit(tokText)
	}
	l.emit(tokEOF)
	return nil
}

func lexOpenTag(l *lexer) stateFn {
	l.cursor += len("<") // step inside
	l.acceptRunRegexp("[^</>]")
	// Cannot open a tag inside the tag definition.
	if l.accept("<") {
		l.emit(tokError)
		return nil
	}

	// Could be a self-closing tag.
	if l.accept("/") {
		if !l.accept(">") {
			l.emit(tokError)
			return nil
		}
		l.emit(tokSelfClosingTag)
		return lexHTML
	}

	// The tag definition must now end.
	if !l.accept(">") {
		l.emit(tokError)
		return nil
	}

	l.emit(tokOpenTag)
	return lexHTML
}

func lexCloseTag(l *lexer) stateFn {
	l.cursor += len("</") // step inside
	l.acceptRunRegexp("[^</>]")
	// Cannot open a tag inside the tag definition or look self-closing.
	if l.accept("</") {
		l.emit(tokError)
		return nil
	}

	// The tag definition must now end.
	if !l.accept(">") {
		l.emit(tokError)
		return nil
	}

	l.emit(tokCloseTag)
	return lexHTML
}
