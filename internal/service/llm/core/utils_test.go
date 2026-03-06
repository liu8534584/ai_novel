package core

import "testing"

func TestThinkTagFilterAcrossChunks(t *testing.T) {
	f := NewThinkTagFilter()
	got := ""
	got += f.Process("ab<th")
	got += f.Process("ink>hidden")
	got += f.Process("</think>cd")
	got += f.Flush()
	if got != "abcd" {
		t.Fatalf("unexpected filtered content: %q", got)
	}
}

func TestRemoveReasoningContent(t *testing.T) {
	in := "x<think>abc</think>y"
	if out := RemoveReasoningContent(in); out != "xy" {
		t.Fatalf("unexpected output: %q", out)
	}
}
