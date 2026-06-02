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
	out := renderSlide("# Hello\n\n- one\n- two", 80)
	if out == "" {
		t.Error("render returned empty string")
	}
	if out == "# Hello\n\n- one\n- two" {
		t.Error("render returned raw input (glamour did not process)")
	}
}
