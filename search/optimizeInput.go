package search

import (
	"regexp"
	"strings"
)

var htmlScriptRemover = regexp.MustCompile(`<script[^>]*>[\s\S]*?</script>`)
var whiteSpaceCompactor = regexp.MustCompile(`\s+`)

func optimizeString(s string) string {
	s = strings.ToLower(s)
	s = htmlScriptRemover.ReplaceAllString(s, "")
	s = whiteSpaceCompactor.ReplaceAllString(s, " ")
	return s
}
