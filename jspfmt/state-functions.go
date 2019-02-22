package jspfmt

type stateFn func(*lexer) stateFn

// lexHtml accepts text until a leftMeta is found
func lexHTML(l *lexer) stateFn {
	for {
		if l.hasPrefix("<") {
			if l.cursor > l.start {
				l.emit(tokText)
			}
			if l.hasPrefix("</") {
				return lexCloseTag
			}
			return lexOpenTag
		}
		if l.next() == eof {
			break
		}
	}
	// We've correctly reached EOF.
	if l.cursor > l.start {
		l.emit(tokText)
	}
	l.emit(tokEOF)
	return nil
}

func lexOpenTag(l *lexer) stateFn {
	l.cursor += len("<") // step inside
	l.acceptRunNot("</>")
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
	l.acceptRunNot("</>")
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
