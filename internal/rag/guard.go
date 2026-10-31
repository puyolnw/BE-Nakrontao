package rag

import (
	"regexp"
	"strings"
)

// escapeForPrompt neutralises markup so user or document text cannot close or forge the
// <passage>/<question> tags that separate data from instructions (prompt injection, threat A1/A2).
var promptEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func escapeForPrompt(s string) string { return promptEscaper.Replace(s) }

// PII patterns redacted before chat text leaves the server or is stored (threat A5, PDPA).
// Regexes cannot catch every personal detail. These cover the common formats, and the privacy
// text says so.
var piiPatterns = []struct {
	re   *regexp.Regexp
	repl string
}{
	{regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`), "[email]"},
	// LINE / WeChat handles, only when labelled with an id keyword or a colon so ordinary words survive.
	{regexp.MustCompile(`(?i)(?:\b(?:line|wechat|weixin)|ไลน์|วีแชท|微信)(?:\s*(?:id|ไอดี)\s*[:：]?|\s*[:：])\s*[a-z0-9_.\-]{3,}`), "[contact]"},
	// Thai national ID: 13 digits, optionally grouped 1-4-5-2-1.
	{regexp.MustCompile(`\b\d[\s-]?\d{4}[\s-]?\d{5}[\s-]?\d{2}[\s-]?\d\b`), "[id]"},
	// Chinese national ID: 17 digits and a check digit that may be X.
	{regexp.MustCompile(`\b\d{17}[\dXx]\b`), "[id]"},
	// Card-like numbers: 13–19 digits with optional separators.
	{regexp.MustCompile(`\b(?:\d[\s-]?){13,19}\b`), "[number]"},
	// Thai mobile and landline, 0x with 9–10 digits in any grouping.
	{regexp.MustCompile(`\b0\d(?:[\s-]?\d){8,9}\b`), "[phone]"},
	// International +66 (Thai) form.
	{regexp.MustCompile(`\+?66[\s-]?\d(?:[\s-]?\d){7,8}\b`), "[phone]"},
	// Chinese mobile 1xx with 11 digits.
	{regexp.MustCompile(`\b(?:\+?86[\s-]?)?1[3-9]\d[\s-]?\d{4}[\s-]?\d{4}\b`), "[phone]"},
	// Bank account and other long numbers that no pattern above caught.
	{regexp.MustCompile(`\b\d{9,12}\b`), "[number]"},
}

// foldDigits maps Thai and full-width digits to ASCII so the patterns see them too.
func foldDigits(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= '๐' && r <= '๙':
			return '0' + (r - '๐')
		case r >= '０' && r <= '９':
			return '0' + (r - '０')
		}
		return r
	}, s)
}

func RedactPII(s string) string {
	s = foldDigits(s)
	for _, p := range piiPatterns {
		s = p.re.ReplaceAllString(s, p.repl)
	}
	return s
}
