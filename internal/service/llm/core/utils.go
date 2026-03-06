package core

import (
	"regexp"
	"strings"
)

// RemoveReasoningContent removes content within <think>...</think> tags or "Thinking Process:" blocks
func RemoveReasoningContent(content string) string {
	// 1. Match <think>...</think> (including newlines)
	reTag := regexp.MustCompile(`(?s)<think>.*?</think>`)
	content = reTag.ReplaceAllString(content, "")

	// 2. Match "Thinking Process:", "Thought:", or "Reasoning:" followed by text before a JSON block
	// Go's regexp doesn't support lookahead (?=), so we use a simpler regex and manual check.
	rePrefix := regexp.MustCompile(`(?si)(Thinking Process|Thought|Reasoning):`)
	loc := rePrefix.FindStringIndex(content)
	if loc != nil {
		prefixEnd := loc[1]
		rest := content[prefixEnd:]
		// Find first JSON marker: {, [, or ```
		firstBrace := strings.Index(rest, "{")
		firstBracket := strings.Index(rest, "[")
		firstFence := strings.Index(rest, "```")

		markerIdx := -1
		for _, idx := range []int{firstBrace, firstBracket, firstFence} {
			if idx != -1 && (markerIdx == -1 || idx < markerIdx) {
				markerIdx = idx
			}
		}

		if markerIdx != -1 {
			// Remove from prefix start to marker start
			content = content[:loc[0]] + rest[markerIdx:]
		} else {
			// If no JSON marker found, remove everything from the prefix to the end
			content = content[:loc[0]]
		}
	}

	return content
}

// ThinkTagFilter incrementally strips <think>...</think> in streaming text.
// It handles tags split across chunks.
type ThinkTagFilter struct {
	inThink bool
	pending string
}

func NewThinkTagFilter() *ThinkTagFilter {
	return &ThinkTagFilter{}
}

func (f *ThinkTagFilter) Process(chunk string) string {
	data := f.pending + chunk
	f.pending = ""
	var out strings.Builder

	// Check for "Thinking Process:" if not already in think mode
	if !f.inThink {
		rePrefix := regexp.MustCompile(`(?si)(Thinking Process|Thought|Reasoning):`)
		loc := rePrefix.FindStringIndex(data)
		if loc != nil {
			// Found reasoning prefix. 
			// Check if there is a JSON marker or code block after it in this chunk
			rest := data[loc[1]:]
			firstMarker := strings.IndexAny(rest, "{[")
			firstFence := strings.Index(rest, "```")
			
			markerIdx := -1
			if firstMarker != -1 && (firstFence == -1 || firstMarker < firstFence) {
				markerIdx = firstMarker
			} else if firstFence != -1 {
				markerIdx = firstFence
			}

			if markerIdx != -1 {
				// Marker found, skip everything between prefix and marker
				out.WriteString(data[:loc[0]])
				data = rest[markerIdx:]
			} else {
				// No marker yet, enter think mode and wait for marker
				out.WriteString(data[:loc[0]])
				f.inThink = true
				f.pending = "" // Everything after prefix is skipped
				return out.String()
			}
		}
	}

	for len(data) > 0 {
		if f.inThink {
			// If we entered think mode via <think> tag
			endIdx := strings.Index(data, "</think>")
			if endIdx != -1 {
				data = data[endIdx+len("</think>"):]
				f.inThink = false
				continue
			}

			// If we entered think mode via "Thinking Process:" prefix, we look for JSON markers
			firstMarker := strings.IndexAny(data, "{[")
			firstFence := strings.Index(data, "```")
			markerIdx := -1
			if firstMarker != -1 && (firstFence == -1 || firstMarker < firstFence) {
				markerIdx = firstMarker
			} else if firstFence != -1 {
				markerIdx = firstFence
			}

			if markerIdx != -1 {
				data = data[markerIdx:]
				f.inThink = false
				continue
			}

			// Still in think mode, keep only potential partial markers
			keep := 8 // max(len("</think>"), len("```"))
			if len(data) > keep {
				f.pending = data[len(data)-keep:]
			} else {
				f.pending = data
			}
			return out.String()
		}

		startIdx := strings.Index(data, "<think>")
		if startIdx == -1 {
			// Keep only the possible suffix for a split opening tag.
			keep := len("<think>") - 1
			if len(data) > keep {
				out.WriteString(data[:len(data)-keep])
				f.pending = data[len(data)-keep:]
			} else {
				f.pending = data
			}
			return out.String()
		}

		out.WriteString(data[:startIdx])
		data = data[startIdx+len("<think>"):]
		f.inThink = true
	}

	return out.String()
}

func (f *ThinkTagFilter) Flush() string {
	if f.inThink {
		f.pending = ""
		return ""
	}
	out := f.pending
	f.pending = ""
	return out
}

// ParseJSON cleans the LLM response to extract the JSON content
func ParseJSON(content string) string {
	// First, remove reasoning content if present
	content = RemoveReasoningContent(content)

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

	trimmed := strings.TrimSpace(content)
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return trimmed
	}

	return ""
}
