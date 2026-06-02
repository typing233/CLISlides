package main

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func parseMetadata(content string) (Metadata, string) {
	var meta Metadata

	if !strings.HasPrefix(content, "---\n") {
		return meta, content
	}

	end := strings.Index(content[4:], "\n---")
	if end == -1 {
		return meta, content
	}
	end += 4

	yamlBlock := content[4:end]
	remaining := content[end+4:]
	remaining = strings.TrimPrefix(remaining, "\n")

	_ = yaml.Unmarshal([]byte(yamlBlock), &meta)
	return meta, remaining
}
