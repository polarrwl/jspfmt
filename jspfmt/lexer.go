package jspfmt

import (
	"strings"
	"unicode/utf8"
)

type lexer struct {
	name  string // used for error reports.
	input string // the string being scanned.

	start  int // start position of current token.
	cursor int // current position of the input.
	width  int // width of last rune read.

	tokens chan token // channel of scanned tokens.
}

func (l *lexer) run() {
	for state := lexHTML; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{typ: t, val: l.input[l.start:l.cursor]}
	l.start = l.cursor
}

func (l *lexer) next() (r rune) {
	if l.cursor >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.cursor:])
	l.cursor += l.width
	return r
}

func (l *lexer) backup() {
	l.cursor -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) hasPrefix(pfx string) bool {
	return strings.HasPrefix(l.input[l.cursor:], pfx)
}

// acceptNot accepts the next rune if it is not a member of invalid.
func (l *lexer) accept(valid string) bool {
	r := l.next()
	if strings.IndexRune(valid, r) < 0 {
		l.backup()
		return false
	}
	return true
}

func (l *lexer) acceptRun(valid string) {
	for l.accept(valid) {
	}
}

// acceptNot accepts the next rune if it is not a member of invalid.
func (l *lexer) acceptNot(invalid string) bool {
	if strings.IndexRune(invalid, l.next()) >= 0 {
		l.backup()
		return false
	}
	return true
}

func (l *lexer) acceptRunNot(invalid string) {
	for l.acceptNot(invalid) {
	}
}
