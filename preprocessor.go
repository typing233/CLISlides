package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var tildeBlockRe = regexp.MustCompile(`(?m)^~~~(\S+)\n([\s\S]*?)^~~~$`)

func preprocess(content string) string {
	return tildeBlockRe.ReplaceAllStringFunc(content, func(match string) string {
		groups := tildeBlockRe.FindStringSubmatch(match)
		if len(groups) < 3 {
			return match
		}
		command := groups[1]
		body := groups[2]

		cmd := exec.Command("sh", "-c", command)
		cmd.Stdin = strings.NewReader(body)
		out, err := cmd.Output()
		if err != nil {
			return fmt.Sprintf("```\n[preprocessor error: %v]\n```", err)
		}
		return strings.TrimRight(string(out), "\n")
	})
}
