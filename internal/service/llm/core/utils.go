package core

import (
	"regexp"
	"strings"
)

// ParseJSON cleans the LLM response to extract the JSON content
func ParseJSON(content string) string {
	// Regular expression to find JSON blocks enclosed in markdown code fences
	re := regexp.MustCompile(`(?s)\` + "```" + `(?:json)?\s*(.*?)\s*\` + "```")
	
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// Improved strategy: find first '{' or '[' and last '}' or ']'
	firstBrace := strings.IndexAny(content, "{[")
	lastBrace := strings.LastIndexAny(content, "}]")

	if firstBrace != -1 && lastBrace != -1 && lastBrace > firstBrace {
		return content[firstBrace : lastBrace+1]
	}

	return strings.TrimSpace(content)
}
