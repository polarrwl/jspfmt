package jspfmt

import (
	"regexp"
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

// acceptNot accepts the next rune if it is not a member of invalid.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) < 0 {
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
func (l *lexer) acceptRegexp(valid string) bool {
	match, _ := regexp.MatchString(valid, string(l.next()))
	if !match {
		l.backup()
		return false
	}
	return true
}

func (l *lexer) acceptRunRegexp(valid string) {
	for l.acceptRegexp(valid) {
	}
}

func (l *lexer) ignore() {
	l.start = l.cursor
}
