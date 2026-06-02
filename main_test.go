package main

import (
	"testing"
)

func TestSplitSlides(t *testing.T) {
	content := "# Slide 1\n\nHello\n\n---\n\n# Slide 2\n\nWorld\n\n---\n\n# Slide 3"
	slides := splitSlides(content)
	if len(slides) != 3 {
		t.Fatalf("expected 3 slides, got %d", len(slides))
	}
	if slides[0].Raw != "# Slide 1\n\nHello" {
		t.Errorf("slide 0 raw = %q", slides[0].Raw)
	}
	if slides[2].Raw != "# Slide 3" {
		t.Errorf("slide 2 raw = %q", slides[2].Raw)
	}
}

func TestSplitSlidesEmpty(t *testing.T) {
	slides := splitSlides("")
	if len(slides) != 0 {
		t.Fatalf("expected 0 slides from empty input, got %d", len(slides))
	}
}

func TestSplitSlidesSingle(t *testing.T) {
	slides := splitSlides("# Only One Slide")
	if len(slides) != 1 {
		t.Fatalf("expected 1 slide, got %d", len(slides))
	}
}

func TestRenderSlide(t *testing.T) {
	out := renderSlide("# Hello\n\n- one\n- two", 80, "")
	if out == "" {
		t.Error("render returned empty string")
	}
	if out == "# Hello\n\n- one\n- two" {
		t.Error("render returned raw input (glamour did not process)")
	}
}

func TestParseMetadata(t *testing.T) {
	input := "---\ntitle: Test\nauthor: Alice\ntheme: dark\n---\n# First Slide"
	meta, remaining := parseMetadata(input)
	if meta.Title != "Test" {
		t.Errorf("expected title 'Test', got %q", meta.Title)
	}
	if meta.Author != "Alice" {
		t.Errorf("expected author 'Alice', got %q", meta.Author)
	}
	if meta.Theme != "dark" {
		t.Errorf("expected theme 'dark', got %q", meta.Theme)
	}
	if remaining != "# First Slide" {
		t.Errorf("unexpected remaining: %q", remaining)
	}
}

func TestParseMetadataNoFrontmatter(t *testing.T) {
	input := "# Just a slide"
	meta, remaining := parseMetadata(input)
	if meta.Title != "" {
		t.Errorf("expected empty title, got %q", meta.Title)
	}
	if remaining != input {
		t.Errorf("content should be unchanged")
	}
}

func TestPreprocess(t *testing.T) {
	input := "# Title\n\n~~~cat\nhello world\n~~~\n\nMore text"
	result := preprocess(input)
	if result == input {
		t.Error("preprocess did not modify content")
	}
	if !contains(result, "hello world") {
		t.Errorf("expected 'hello world' in output, got: %q", result)
	}
}

func TestPreprocessNoBlocks(t *testing.T) {
	input := "# No preprocessor blocks here\n\n```go\nfmt.Println(\"hi\")\n```"
	result := preprocess(input)
	if result != input {
		t.Error("preprocess should not modify content without tilde blocks")
	}
}

func TestDetectCodeBlocks(t *testing.T) {
	raw := "# Test\n\n```go\nfmt.Println(\"hello\")\n```\n\n```python\nprint('hi')\n```"
	blocks := detectCodeBlocks(raw)
	if len(blocks) != 2 {
		t.Fatalf("expected 2 code blocks, got %d", len(blocks))
	}
	if blocks[0].Language != "go" {
		t.Errorf("expected language 'go', got %q", blocks[0].Language)
	}
	if blocks[1].Language != "python" {
		t.Errorf("expected language 'python', got %q", blocks[1].Language)
	}
}

func TestDetectCodeBlocksNone(t *testing.T) {
	raw := "# No code here\n\nJust text"
	blocks := detectCodeBlocks(raw)
	if len(blocks) != 0 {
		t.Fatalf("expected 0 code blocks, got %d", len(blocks))
	}
}

func TestNavigationGoto(t *testing.T) {
	slides := []Slide{
		{Raw: "# One"},
		{Raw: "# Two"},
		{Raw: "# Three"},
		{Raw: "# Four"},
		{Raw: "# Five"},
	}
	m := newModel(slides, Metadata{}, false)
	m.width = 80
	m.height = 24

	m.gotoSlide(3)
	if m.current != 3 {
		t.Errorf("expected current=3, got %d", m.current)
	}

	m.gotoSlide(-1)
	if m.current != 0 {
		t.Errorf("expected current=0 after goto -1, got %d", m.current)
	}

	m.gotoSlide(100)
	if m.current != 4 {
		t.Errorf("expected current=4 after goto 100, got %d", m.current)
	}
}

func TestNavigationAdvanceRetreat(t *testing.T) {
	slides := []Slide{
		{Raw: "# One"},
		{Raw: "# Two"},
		{Raw: "# Three"},
	}
	m := newModel(slides, Metadata{}, false)
	m.width = 80
	m.height = 24

	m.advance(2)
	if m.current != 2 {
		t.Errorf("expected current=2, got %d", m.current)
	}

	m.retreat(1)
	if m.current != 1 {
		t.Errorf("expected current=1, got %d", m.current)
	}

	m.advance(100)
	if m.current != 2 {
		t.Errorf("expected current=2 (clamped), got %d", m.current)
	}
}

func TestNumericBuffer(t *testing.T) {
	slides := []Slide{{Raw: "a"}, {Raw: "b"}, {Raw: "c"}}
	m := newModel(slides, Metadata{}, false)
	m.width = 80
	m.height = 24
	m.numericBuffer = "2"

	n := m.consumeNumericOr(1)
	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
	if m.numericBuffer != "" {
		t.Errorf("buffer should be empty after consume")
	}

	n = m.consumeNumericOr(1)
	if n != 1 {
		t.Errorf("expected default 1, got %d", n)
	}
}

func contains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
