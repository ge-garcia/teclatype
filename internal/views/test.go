package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ge-garcia/tecla/internal/source"
)

type TestView struct {
	text      string
	typed     string
	source    source.TestSource
	stopwatch stopwatch.Model
	width     int
	height    int
}

func NewTestView(width int, height int) *TestView {
	source := source.NewWordsSource("common-words-en.list", 20)
	tv := TestView{
		text:      source.Generate(),
		source:    source,
		stopwatch: stopwatch.New(),
		width:     width,
		height:    height,
	}

	return &tv
}

func (tv TestView) Init() tea.Cmd {
	return nil
}

func (tv TestView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewTitleView(tv.width, tv.height), nil
		case tea.KeyEnter:
			return tv.GenerateTest()
		case tea.KeyBackspace:
			if len(tv.typed) > 0 {
				tv.typed = tv.typed[:len(tv.typed)-1]
			}
		case tea.KeyCtrlW:
			if tv.stopwatch.Running() {
				break
			}

			if ws, ok := tv.source.(*source.WordsSource); ok {
				ws.Count *= 2

				if ws.Count > 80 {
					ws.Count = 10
				}

				return tv.GenerateTest()
			}
		default:
			if len(tv.typed) < len(tv.text) {
				tv.typed += string(msg.Runes)
			}

			if tv.ShouldEndTest() {
				rv := NewResultsView(tv.width, tv.height, tv.stopwatch.Elapsed(), len(tv.typed))
				return rv, tv.stopwatch.Stop()
			}

			if !tv.stopwatch.Running() {
				return tv, tv.stopwatch.Start()
			}
		}
	case tea.WindowSizeMsg:
		tv.width = msg.Width
		tv.height = msg.Height
	}

	// update stopwatch otherwise
	var cmd tea.Cmd
	tv.stopwatch, cmd = tv.stopwatch.Update(msg)
	return tv, cmd
}

func (tv TestView) View() string {
	currWord := true
	nextWord := false

	styleDefault := lipgloss.NewStyle().Foreground(ColorText).Background(ColorBackground)
	styleCorrect := styleDefault.Foreground(ColorCorrect)
	styleIncorrect := styleDefault.Foreground(ColorIncorrect)
	styleIncorrectBlank := styleDefault.Background(ColorIncorrect)

	styleCursor := lipgloss.NewStyle().Foreground(ColorHighlight).Background(ColorBgHighlight)
	styleWord := styleDefault.Foreground(ColorCurrentWord)
	styleAhead := styleDefault.Foreground(ColorNextWord)

	var styledText string
	for i, c := range tv.text {
		charStyle := styleDefault
		if i < len(tv.typed) {
			if tv.typed[i] == byte(c) {
				charStyle = styleCorrect
			} else {
				charStyle = styleIncorrect
				if byte(c) == ' ' {
					charStyle = styleIncorrectBlank
				}
			}
		} else if i == len(tv.typed) {
			charStyle = styleCursor
		} else {
			if currWord {
				if tv.text[i-1] != ' ' {
					charStyle = styleWord
				} else {
					currWord = false
					nextWord = true
					i++
				}
			}
			if nextWord {
				if tv.text[i] != ' ' {
					charStyle = styleAhead
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

	// commands only available while not in a test
	if !tv.stopwatch.Running() {
		if ws, ok := tv.source.(*source.WordsSource); ok {
			cmds = append(cmds, keybind{key: "Control+W", cmd: fmt.Sprintf("word count (%d)", ws.Count)})
		}
	}

	container := styleDefault.Height(tv.height-StatusBarHeight).Width(tv.width).Align(lipgloss.Center, lipgloss.Center).Render(styledText)
	footer := renderFooter(cmds, tv.width)
	view := lipgloss.JoinVertical(lipgloss.Center, container, footer)

	// return the view with full window size, background color, and centered
	return styleDefault.
		Width(tv.width).
		Height(tv.height).
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

func (tv TestView) GenerateTest() (TestView, tea.Cmd) {
	tv.text = tv.source.Generate()
	tv.typed = ""

	cmds := tea.Batch(tv.stopwatch.Stop(), tv.stopwatch.Reset())
	return tv, cmds
}
