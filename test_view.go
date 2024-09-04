package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TestView struct {
	start time.Time
	text  string
	typed string

	// TODO: abstract into some TestSource interface which generates text, to
	// differentiate between random words vs specific text (Rust enums pls)
	words      []string // words to be randomized
	testLength int

	width  int
	height int
}

func NewTestView(width int, height int) *TestView {
	words, err := readWords("common-words-en.list")
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	tv := TestView{
		words:      words,
		testLength: 20,
		width:      width,
		height:     height,
	}
	tv.GenerateTest()

	return &tv
}

func (tv *TestView) GenerateTest() {
	selectedWords := make([]string, tv.testLength)

	for i := range selectedWords {
		selectedWords[i] = tv.words[rand.Intn(len(tv.words))]
	}

	tv.text = strings.Join(selectedWords, " ")
}

func (tv TestView) Init() tea.Cmd {
	return nil
}

func (tv TestView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	os.WriteFile("/tmp/log", []byte(fmt.Sprintf("%v", tv)), 0644)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewTitleView(tv.width, tv.height), nil
		case tea.KeyEnter:
			tv.GenerateTest()
			tv.start = time.Time{}
			tv.typed = ""
		case tea.KeyBackspace:
			if len(tv.typed) > 0 {
				tv.typed = tv.typed[:len(tv.typed)-1]
			}
		default:
			if tv.start.IsZero() {
				tv.start = time.Now()
			}

			if len(tv.typed) < len(tv.text) {
				tv.typed += string(msg.Runes)
			}

			if tv.ShouldEndTest() {
				return NewResultsView(tv.width, tv.height, time.Since(tv.start), len(tv.typed)), nil
			}
		}
	case tea.WindowSizeMsg:
		tv.width = msg.Width
		tv.height = msg.Height
	}

	return tv, nil
}

func (tv TestView) View() string {
	currWord := true
	nextWord := false
	var styledText string
	for i, c := range tv.text {
		charStyle := lipgloss.NewStyle().Foreground(colorText)
		if i < len(tv.typed) {
			if tv.typed[i] == byte(c) {
				charStyle = lipgloss.NewStyle().Foreground(colorCorrect) // correct letter
			} else {
				charStyle = lipgloss.NewStyle().Foreground(colorIncorrect) // incorrect letter
				if byte(c) == ' ' {
					charStyle = lipgloss.NewStyle().Background(colorIncorrect)
				}
			}
		} else if i == len(tv.typed) {
			charStyle = lipgloss.NewStyle().Foreground(colorHighlight).Background(colorBgHighlight) // current letter user is on
		} else {
			if currWord {
				if tv.text[i-1] != ' ' {
					charStyle = lipgloss.NewStyle().Foreground(colorCurrentWord)
				} else {
					currWord = false
					nextWord = true
					i++
				}
			}
			if nextWord {
				if tv.text[i] != ' ' {
					charStyle = lipgloss.NewStyle().Foreground(colorNextWord)
				} else {
					nextWord = false
				}
			}
		}
		styledText += charStyle.Render(string(c))
	}

	cmds := []keybind{
		{key: "Escape", cmd: "stop"},
		{key: "Enter", cmd: "restart"},
		{key: "Control+C", cmd: "quit"},
	}
	container := lipgloss.NewStyle().Background(colorBackground).Height(tv.height-statusBarHeight).Width(tv.width).Align(lipgloss.Center, lipgloss.Center).Render(styledText)
	footer := renderFooter(cmds, tv.width)
	view := lipgloss.JoinVertical(lipgloss.Center, container, footer)

	// return the view with full window size, background color, and centered
	return lipgloss.NewStyle().
		Width(tv.width).
		Height(tv.height).
		Background(colorBackground).
		Render(lipgloss.Place(tv.width, tv.height, lipgloss.Center, lipgloss.Center, view))
}

func (tv TestView) ShouldEndTest() bool {
	textWords := strings.Fields(tv.text)
	typedWords := strings.Fields(tv.typed)
	lastTextWord := textWords[len(textWords)-1]
	lastTypedWord := ""
	if len(typedWords) > 0 {
		lastTypedWord = typedWords[len(typedWords)-1]
	}

	// all words typed
	if len(typedWords) == len(textWords) {
		// last word is correct or space after incorrectLastWord
		if lastTypedWord == lastTextWord {
			return true
		}
		// space after incorrect last word
		if strings.HasSuffix(tv.typed, " ") {
			return true
		}
	}

	return false
}
