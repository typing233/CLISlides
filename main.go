package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	serveFlag := flag.Bool("serve", false, "Start SSH presentation server")
	portFlag := flag.Int("port", 2222, "SSH server port")
	noExecFlag := flag.Bool("no-exec", false, "Disable preprocessor and code execution")
	flag.Parse()

	content, err := loadContent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	meta, content := parseMetadata(content)

	if !*noExecFlag {
		content = preprocess(content)
	}

	slides := splitSlides(content)
	if len(slides) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no slides found\n")
		os.Exit(1)
	}

	if *serveFlag {
		port := meta.SSHPort
		if port == 0 {
			port = *portFlag
		}
		if err := serveSSH(slides, meta, port, *noExecFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	m := newModel(slides, meta, *noExecFlag)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func loadContent() (string, error) {
	args := flag.Args()
	if len(args) > 0 {
		path := args[0]
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

	return "", fmt.Errorf("usage: clislides [flags] <file.md>  or  cat file.md | clislides")
}
