package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		var output string
		var err error

		if block.Language == "go" {
			output, err = executeGo(block.Code)
		} else {
			runner := resolveRunner(block.Language)
			cmd := exec.Command("sh", "-c", runner)
			cmd.Stdin = strings.NewReader(block.Code)
			out, e := cmd.CombinedOutput()
			output = string(out)
			err = e
		}

		if err != nil {
			output += fmt.Sprintf("\n[exit: %v]", err)
		}
		return execResultMsg{output: output}
	}
}

func executeGo(code string) (string, error) {
	dir, err := os.MkdirTemp("", "clislides-exec-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	file := filepath.Join(dir, "main.go")
	if err := os.WriteFile(file, []byte(code), 0644); err != nil {
		return "", err
	}

	cmd := exec.Command("go", "run", file)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func resolveRunner(lang string) string {
	switch lang {
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
