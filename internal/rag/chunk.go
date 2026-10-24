package rag

import "strings"

const (
	chunkRunes   = 700
	chunkOverlap = 120
)

// Chunk splits text on paragraph boundaries into pieces of about chunkRunes runes, with overlap
// so a fact cut at a boundary still appears whole in one chunk. Counts runes, not bytes, because
// Thai and Chinese characters are multi-byte.
func Chunk(text string) []string {
	text = strings.TrimSpace(strings.ReplaceAll(text, "\r\n", "\n"))
	if text == "" {
		return nil
	}
	var paras []string
	for _, p := range strings.Split(text, "\n") {
		if p = strings.TrimSpace(p); p != "" {
			paras = append(paras, p)
		}
	}
	var chunks []string
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			chunks = append(chunks, strings.TrimSpace(string(cur)))
			if len(cur) > chunkOverlap {
				cur = append([]rune(nil), cur[len(cur)-chunkOverlap:]...)
			} else {
				cur = nil
			}
		}
	}
	for _, p := range paras {
		pr := []rune(p)
		for len(pr) > chunkRunes { // a single huge paragraph
			if len(cur) > chunkOverlap {
				flush()
			}
			take := chunkRunes - len(cur)
			cur = append(cur, pr[:take]...)
			pr = pr[take:]
			flush()
		}
		if len(cur)+len(pr)+1 > chunkRunes && len(cur) > chunkOverlap {
			flush()
		}
		if len(cur) > 0 {
			cur = append(cur, '\n')
		}
		cur = append(cur, pr...)
	}
	if len(cur) > chunkOverlap || len(chunks) == 0 {
		chunks = append(chunks, strings.TrimSpace(string(cur)))
	}
	return chunks
}
