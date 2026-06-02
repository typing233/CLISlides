package main

import (
	"regexp"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

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

	case execResultMsg:
		m.execOutput = msg.output
		m.mode = ModeExecOutput
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case ModeSearch:
			return m.handleKeySearch(msg)
		case ModeExecConfirm:
			return m.handleKeyExecConfirm(msg)
		case ModeExecOutput:
			return m.handleKeyExecOutput(msg)
		default:
			return m.handleKeyNormal(msg)
		}
	}
	return m, nil
}

func (m model) handleKeyNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	if key == "ctrl+c" {
		return m, tea.Quit
	}

	if m.pendingG {
		m.pendingG = false
		if key == "g" {
			n := m.consumeNumeric()
			if n > 0 {
				m.gotoSlide(n - 1)
			} else {
				m.gotoSlide(0)
			}
			return m, nil
		}
	}

	if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
		m.numericBuffer += key
		return m, nil
	}

	switch key {
	case "q", "esc":
		return m, tea.Quit
	case "g":
		m.pendingG = true
		return m, nil
	case "G":
		n := m.consumeNumeric()
		if n > 0 {
			m.gotoSlide(n - 1)
		} else {
			m.gotoSlide(len(m.slides) - 1)
		}
	case "/":
		m.mode = ModeSearch
		m.searchInput.Focus()
		m.searchInput.SetValue("")
		return m, m.searchInput.Cursor.BlinkCmd()
	case "n":
		if len(m.searchResults) > 0 {
			m.nextSearchMatch()
		} else {
			m.advance(m.consumeNumericOr(1))
		}
	case "N":
		if len(m.searchResults) > 0 {
			m.prevSearchMatch()
		} else {
			m.advance(m.consumeNumericOr(1))
		}
	case "p", "left", "h", "k":
		m.retreat(m.consumeNumericOr(1))
	case "right", "l", "j", " ", "enter":
		m.advance(m.consumeNumericOr(1))
	case "e", "x":
		if !m.noExec {
			blocks := detectCodeBlocks(m.slides[m.current].Raw)
			if len(blocks) > 0 {
				m.codeBlocks = blocks
				m.selectedBlock = 0
				m.mode = ModeExecConfirm
			}
		}
	}
	return m, nil
}

func (m model) handleKeySearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "enter":
		m.performSearch()
		return m, nil
	case "esc", "ctrl+c":
		m.mode = ModeNormal
		m.searchInput.Blur()
		return m, nil
	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}
}

func (m model) handleKeyExecConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "esc", "ctrl+c", "q":
		m.mode = ModeNormal
		return m, nil
	case "up", "k":
		if m.selectedBlock > 0 {
			m.selectedBlock--
		}
		return m, nil
	case "down", "j":
		if m.selectedBlock < len(m.codeBlocks)-1 {
			m.selectedBlock++
		}
		return m, nil
	case "enter", "y":
		block := m.codeBlocks[m.selectedBlock]
		return m, executeBlock(block)
	}
	return m, nil
}

func (m model) handleKeyExecOutput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.mode = ModeNormal
	m.execOutput = ""
	return m, nil
}

func (m *model) gotoSlide(n int) {
	if n < 0 {
		n = 0
	}
	if n >= len(m.slides) {
		n = len(m.slides) - 1
	}
	m.current = n
	m.rerenderCurrent()
}

func (m *model) advance(count int) {
	target := m.current + count
	if target >= len(m.slides) {
		target = len(m.slides) - 1
	}
	if target != m.current {
		m.current = target
		m.rerenderCurrent()
	}
}

func (m *model) retreat(count int) {
	target := m.current - count
	if target < 0 {
		target = 0
	}
	if target != m.current {
		m.current = target
		m.rerenderCurrent()
	}
}

func (m *model) consumeNumeric() int {
	if m.numericBuffer == "" {
		return 0
	}
	n, _ := strconv.Atoi(m.numericBuffer)
	m.numericBuffer = ""
	return n
}

func (m *model) consumeNumericOr(def int) int {
	n := m.consumeNumeric()
	if n == 0 {
		return def
	}
	return n
}

func (m *model) performSearch() {
	pattern := m.searchInput.Value()
	m.mode = ModeNormal
	m.searchInput.Blur()

	if pattern == "" {
		return
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}

	m.searchRegex = re
	m.searchQuery = pattern
	m.searchResults = nil

	for i, s := range m.slides {
		if re.MatchString(s.Raw) {
			m.searchResults = append(m.searchResults, i)
		}
	}

	if len(m.searchResults) > 0 {
		m.searchIndex = 0
		for i, idx := range m.searchResults {
			if idx >= m.current {
				m.searchIndex = i
				break
			}
		}
		m.gotoSlide(m.searchResults[m.searchIndex])
	}
}

func (m *model) nextSearchMatch() {
	if len(m.searchResults) == 0 {
		return
	}
	m.searchIndex = (m.searchIndex + 1) % len(m.searchResults)
	m.gotoSlide(m.searchResults[m.searchIndex])
}

func (m *model) prevSearchMatch() {
	if len(m.searchResults) == 0 {
		return
	}
	m.searchIndex--
	if m.searchIndex < 0 {
		m.searchIndex = len(m.searchResults) - 1
	}
	m.gotoSlide(m.searchResults[m.searchIndex])
}
