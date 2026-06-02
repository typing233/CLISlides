package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type CodeBlock struct {
	Language string
	Code     string
}

type execResultMsg struct {
	output string
}

var codeBlockRe = regexp.MustCompile("(?m)^```(\\w*)\\n([\\s\\S]*?)^```$")

func detectCodeBlocks(raw string) []CodeBlock {
	matches := codeBlockRe.FindAllStringSubmatch(raw, -1)
	var blocks []CodeBlock
	for _, m := range matches {
		lang := m[1]
		code := m[2]
		blocks = append(blocks, CodeBlock{Language: lang, Code: code})
	}
	return blocks
}

func executeBlock(block CodeBlock) tea.Cmd {
	return func() tea.Msg {
		runner := resolveRunner(block.Language)
		cmd := exec.Command("sh", "-c", runner)
		cmd.Stdin = strings.NewReader(block.Code)
		out, err := cmd.CombinedOutput()
		output := string(out)
		if err != nil {
			output += fmt.Sprintf("\n[exit: %v]", err)
		}
		return execResultMsg{output: output}
	}
}

func resolveRunner(lang string) string {
	switch lang {
	case "go":
		return "go run /dev/stdin"
	case "python", "python3":
		return "python3"
	case "bash", "sh", "":
		return "sh"
	case "node", "javascript", "js":
		return "node"
	case "ruby":
		return "ruby"
	case "perl":
		return "perl"
	default:
		return lang
	}
}
