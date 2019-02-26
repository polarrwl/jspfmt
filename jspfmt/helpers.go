package jspfmt

import (
	"strings"
)

// voidTagNames are names of void tags.
var voidTagNames = []string{
	"!doctype",
	"area",
	"base",
	"br",
	"col",
	"command",
	"embed",
	"hr",
	"img",
	"input",
	"link",
	"meta",
	"param",
	"source",
	"track",
}

// isVoidTag returns true if the tagname is one of the void HTML tags, and false otherwise.
func isVoidTagname(tagname string) bool {
	for _, name := range voidTagNames {
		// perform case-insensitive comparison
		if strings.HasPrefix(strings.ToLower(tagname), strings.ToLower(name)) {
			return true
		}
	}
	return false
}
