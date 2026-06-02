package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Slide struct {
	Raw      string
	Rendered string
}

type model struct {
	slides  []Slide
	current int
	width   int
	height  int
}

type keyMap struct {
	Next key.Binding
	Prev key.Binding
	Quit key.Binding
}

var keys = keyMap{
	Next: key.NewBinding(
		key.WithKeys("right", "l", "j", " ", "enter"),
	),
	Prev: key.NewBinding(
		key.WithKeys("left", "h", "k"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
	),
}

func main() {
	content, err := loadContent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	slides := splitSlides(content)
	if len(slides) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no slides found\n")
		os.Exit(1)
	}

	m := model{slides: slides}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func loadContent() (string, error) {
	if len(os.Args) > 1 {
		path := os.Args[1]
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("cannot read file %s: %w", path, err)
		}
		return string(data), nil
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("cannot read stdin: %w", err)
		}
		return string(data), nil
	}

	return "", fmt.Errorf("usage: clislides <file.md>  or  cat file.md | clislides")
}

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

func renderSlide(raw string, width int) string {
	w := width
	if w <= 0 {
		w = 80
	}
	if w > 120 {
		w = 120
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(w-4),
	)
	if err != nil {
		return raw
	}

	out, err := r.Render(raw)
	if err != nil {
		return raw
	}
	return out
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.rerenderCurrent()
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Next):
			if m.current < len(m.slides)-1 {
				m.current++
				m.rerenderCurrent()
			}
			return m, nil
		case key.Matches(msg, keys.Prev):
			if m.current > 0 {
				m.current--
				m.rerenderCurrent()
			}
			return m, nil
		}
	}
	return m, nil
}

func (m *model) rerenderCurrent() {
	m.slides[m.current].Rendered = renderSlide(m.slides[m.current].Raw, m.width)
}

func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	slide := m.slides[m.current]
	if slide.Rendered == "" {
		m.slides[m.current].Rendered = renderSlide(slide.Raw, m.width)
		slide = m.slides[m.current]
	}

	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render(fmt.Sprintf(" Slide %d/%d │ ←/→ navigate │ q/Esc quit ",
			m.current+1, len(m.slides)))

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

	return body + "\n" + statusBar
}
