package rag

import (
	"regexp"
	"strings"
)

// escapeForPrompt neutralises markup so user or document text cannot close or forge the
// <passage>/<question> tags that separate data from instructions (prompt injection, threat A1/A2).
var promptEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func escapeForPrompt(s string) string { return promptEscaper.Replace(s) }

// PII patterns redacted before chat text is stored (threat A5, PDPA).
var piiPatterns = []struct {
	re   *regexp.Regexp
	repl string
}{
	{regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`), "[email]"},
	// Thai national ID: 13 digits, optionally grouped 1-4-5-2-1.
	{regexp.MustCompile(`\b\d[\s-]?\d{4}[\s-]?\d{5}[\s-]?\d{2}[\s-]?\d\b`), "[id]"},
	// Card-like numbers: 13–19 digits with optional separators.
	{regexp.MustCompile(`\b(?:\d[\s-]?){13,19}\b`), "[number]"},
	// Phone numbers: Thai 0x / +66, and Chinese mobile 1xx with 11 digits.
	{regexp.MustCompile(`(?:\+?66[\s-]?|0)\d{1,2}[\s-]?\d{3}[\s-]?\d{4}\b`), "[phone]"},
	{regexp.MustCompile(`\b(?:\+?86[\s-]?)?1[3-9]\d[\s-]?\d{4}[\s-]?\d{4}\b`), "[phone]"},
}

func RedactPII(s string) string {
	for _, p := range piiPatterns {
		s = p.re.ReplaceAllString(s, p.repl)
	}
	return s
}
