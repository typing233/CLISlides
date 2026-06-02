package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func splitSlides(content string) []Slide {
	parts := strings.Split(content, "\n---\n")
	slides := make([]Slide, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		slides = append(slides, Slide{Raw: trimmed})
	}
	return slides
}

func renderSlide(raw string, width int, theme string) string {
	w := width
	if w <= 0 {
		w = 80
	}
	if w > 120 {
		w = 120
	}

	opts := []glamour.TermRendererOption{
		glamour.WithWordWrap(w - 4),
	}

	switch theme {
	case "dark":
		opts = append(opts, glamour.WithAutoStyle())
	case "light":
		opts = append(opts, glamour.WithAutoStyle())
	case "notty":
		opts = append(opts, glamour.WithAutoStyle())
	default:
		opts = append(opts, glamour.WithAutoStyle())
	}

	r, err := glamour.NewTermRenderer(opts...)
	if err != nil {
		return raw
	}

	out, err := r.Render(raw)
	if err != nil {
		return raw
	}
	return out
}

func (m *model) rerenderCurrent() {
	m.slides[m.current].Rendered = renderSlide(m.slides[m.current].Raw, m.width, m.metadata.Theme)
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	slide := m.slides[m.current]
	if slide.Rendered == "" {
		m.slides[m.current].Rendered = renderSlide(slide.Raw, m.width, m.metadata.Theme)
		slide = m.slides[m.current]
	}

	switch m.mode {
	case ModeExecOutput:
		return m.viewExecOutput()
	case ModeExecConfirm:
		return m.viewExecConfirm()
	default:
		return m.viewNormal(slide)
	}
}

func (m model) viewNormal(slide Slide) string {
	contentHeight := m.height - 2
	content := slide.Rendered

	lines := strings.Split(content, "\n")
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}
	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	body := strings.Join(lines, "\n")

	var statusBar string
	switch m.mode {
	case ModeSearch:
		statusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Render("/" + m.searchInput.View())
	default:
		statusBar = m.buildStatusBar()
	}

	return body + "\n" + statusBar
}

func (m model) buildStatusBar() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	pagination := m.paginationStr()

	var extra string
	if m.numericBuffer != "" {
		extra = fmt.Sprintf(" │ %s…", m.numericBuffer)
	}
	if m.pendingG {
		extra += " │ g…"
	}
	if len(m.searchResults) > 0 && m.searchQuery != "" {
		extra += fmt.Sprintf(" │ /%s [%d/%d]", m.searchQuery, m.searchIndex+1, len(m.searchResults))
	}

	return style.Render(fmt.Sprintf(" %s%s │ ←/→ nav │ /search │ e exec │ q quit ", pagination, extra))
}

func (m model) paginationStr() string {
	switch m.metadata.Pagination {
	case "false", "none":
		return ""
	case "count":
		return fmt.Sprintf("Slide %d of %d", m.current+1, len(m.slides))
	default:
		return fmt.Sprintf("%d/%d", m.current+1, len(m.slides))
	}
}

func (m model) viewExecConfirm() string {
	contentHeight := m.height - 4
	var lines []string

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("  Select code block to execute:"))
	lines = append(lines, "")

	for i, block := range m.codeBlocks {
		prefix := "  "
		if i == m.selectedBlock {
			prefix = "▶ "
		}
		preview := block.Code
		if len(preview) > 60 {
			preview = preview[:60] + "…"
		}
		preview = strings.Split(preview, "\n")[0]
		lang := block.Language
		if lang == "" {
			lang = "sh"
		}
		lines = append(lines, fmt.Sprintf("  %s[%s] %s", prefix, lang, preview))
	}

	for len(lines) < contentHeight {
		lines = append(lines, "")
	}
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}

	body := strings.Join(lines, "\n")
	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Render(" ↑/↓ select │ enter/y execute │ esc cancel ")

	return body + "\n" + statusBar
}

func (m model) viewExecOutput() string {
	contentHeight := m.height - 2
	lines := strings.Split(m.execOutput, "\n")
	if len(lines) > contentHeight {
		lines = lines[:contentHeight]
	}
	for len(lines) < contentHeight {
		lines = append(lines, "")
	}

	body := strings.Join(lines, "\n")
	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Render(" Press any key to return ")

	return body + "\n" + statusBar
}
