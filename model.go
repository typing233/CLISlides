package main

import (
	"regexp"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

type InputMode int

const (
	ModeNormal      InputMode = iota
	ModeSearch                // "/" active, typing regex
	ModeExecConfirm           // prompting to confirm code execution
	ModeExecOutput            // showing execution output
)

type Slide struct {
	Raw      string
	Rendered string
}

type Metadata struct {
	Theme      string `yaml:"theme"`
	Author     string `yaml:"author"`
	Date       string `yaml:"date"`
	Title      string `yaml:"title"`
	Pagination string `yaml:"pagination"`
	SSHPort    int    `yaml:"ssh_port"`
}

type model struct {
	slides   []Slide
	current  int
	width    int
	height   int
	metadata Metadata

	mode          InputMode
	numericBuffer string
	pendingG      bool

	searchInput   textinput.Model
	searchQuery   string
	searchRegex   *regexp.Regexp
	searchResults []int
	searchIndex   int

	codeBlocks    []CodeBlock
	selectedBlock int
	execOutput    string
	execViewport  viewport.Model

	noExec bool
}

func newModel(slides []Slide, meta Metadata, noExec bool) model {
	ti := textinput.New()
	ti.Placeholder = "regex pattern"
	ti.CharLimit = 256

	vp := viewport.New(80, 20)

	return model{
		slides:       slides,
		metadata:     meta,
		searchInput:  ti,
		execViewport: vp,
		noExec:       noExec,
	}
}
