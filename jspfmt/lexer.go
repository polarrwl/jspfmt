package jspfmt

import (
	"fmt"
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

// run the state machine.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// emit sends the current scanned token of type t on the tokens channel.
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{typ: t, val: l.input[l.start:l.cursor]}
	l.start = l.cursor
}

// next moves the cursor forward one rune, returning the rune it steps over.
func (l *lexer) next() (r rune) {
	if l.cursor >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.cursor:])
	l.cursor += l.width
	return r
}

// backup moves the cursor back one rune.
func (l *lexer) backup() {
	l.cursor -= l.width
}

// peek returns the next rune in the input without moving the cursor forward.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept accepts the next rune only if it is a member of valid. Returns
// a boolean indicating whether the rune was accepted.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) < 0 {
		l.backup()
		return false
	}
	return true
}

// acceptRun accepts a run of runes that are members of valid.
func (l *lexer) acceptRun(valid string) {
	for l.accept(valid) {
	}
}

// acceptRegexp accepts the next rune only if it matches against the
// regular expression. Returns boolean indicating whether the rune was accepted.
func (l *lexer) acceptRegexp(valid string) bool {
	if match, _ := regexp.MatchString(valid, string(l.next())); !match {
		l.backup()
		return false
	}
	return true
}

// acceptRunRegexp accepts a run of runes which match the regular expression.
// Note that the regular expression is only applied to one input rune at a time.
func (l *lexer) acceptRunRegexp(valid string) {
	for l.acceptRegexp(valid) {
	}
}

// ignore moves start up to the cursor.
func (l *lexer) ignore() {
	l.start = l.cursor
}

// errorf emits an error token whose value is the given message.
func (l *lexer) errorf(format string, a ...interface{}) {
	// Make a copy of the lexer.
	lcp := &lexer{}
	*lcp = *l

	// Modify the input, start, and cursor so that l.emit()
	// sends the error message.
	l.input = fmt.Sprintf(format, a...)
	l.start, l.cursor = 0, len(l.input)

	// Emit the error
	l.emit(tokError)

	// Fix the state of the lexer
	*l = *lcp
	l.start = l.cursor
}
