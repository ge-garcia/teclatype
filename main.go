package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorText       = lipgloss.Color("#F4E6D2")
	colorIncorrect  = lipgloss.Color("#CA3F3F")
	colorCorrect    = lipgloss.Color("#989F56")
	colorBackground = lipgloss.Color("#181515")
	colorContainer  = lipgloss.Color("#221E1E")
)

type ViewState int

const (
	ViewTitle ViewState = iota
	ViewTest
	ViewResults
)

const (
	statusBarHeight = 1
)

type stats struct {
	wpm int
	cpm int
}

type model struct {
	start        time.Time
	text         string
	typed        string
	lastDuration time.Duration
	success      bool
	statistics   stats
	width        int
	height       int
	state        ViewState
}

type keybind struct {
	key string
	cmd string
}

func initialModel() model {
	return model{
		state: ViewTitle,
		text:  "Sphinx of black quartz, judge my vow.",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC: // quit
			return m, tea.Quit
		case tea.KeyEsc:
			switch m.state {
			case ViewTitle:
				return m, tea.Quit
			case ViewTest:
				m.CalculateStats()
				m.state = ViewResults
			}
		case tea.KeyBackspace:
			if m.state == ViewTest && len(m.typed) > 0 {
				m.typed = m.typed[:len(m.typed)-1]
			}
		case tea.KeyEnter:
			switch m.state {
			case ViewTitle:
				// begin new test
				m.state = ViewTest
				m.start = time.Now()
				m.typed = ""
			case ViewTest:
				// restart test
				m.typed = ""
				m.start = time.Now()
			case ViewResults:
				// begin a new test
				m.state = ViewTest
				m.start = time.Now()
				m.typed = ""
			}
		default:
			if m.state == ViewTest {
				m.typed += string(msg.Runes)
				if m.ShouldEndTest() {
					m.CalculateStats()
					m.state = ViewResults
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m *model) CalculateStats() {
	m.lastDuration = time.Since(m.start)
	durationInMinutes := m.lastDuration.Minutes()
	m.statistics.cpm = int(float64(len(m.typed)) / durationInMinutes)
	m.statistics.wpm = int(float64(len(m.typed)) / 5 / durationInMinutes)
}

func (m model) ShouldEndTest() bool {
	textWords := strings.Fields(m.text)
	typedWords := strings.Fields(m.typed)
	lastTextWord := textWords[len(textWords)-1]
	lastTypedWord := typedWords[len(typedWords)-1]

	// all words typed
	if len(typedWords) == len(textWords) {
		// last word is correct or space after incorrectLastWord
		if lastTypedWord == lastTextWord {
			return true
		}
		// space after incorrect last word
		if strings.HasSuffix(m.typed, " ") {
			return true
		}
	}

	return false
}

func renderFooter(cmds []keybind, width int) string {
	var footerContent string
	for i, cmd := range cmds {
		key := lipgloss.NewStyle().Render(cmd.key)
		command := lipgloss.NewStyle().Render(cmd.cmd)
		footerContent += fmt.Sprintf("%s %s", key, command)
		if i < len(cmds)-1 {
			footerContent += " - "
		}
	}

	return lipgloss.NewStyle().
		Background(colorContainer).
		Foreground(colorText).
		Height(statusBarHeight).
		Width(width).
		Padding(0, 2).
		Render(footerContent)
}

func (m model) View() string {
	// prompt := lipgloss.NewStyle().Foreground(colorTestBlank).Background(colorBackground).Render("Type this sentence:")
	switch m.state {
	case ViewTitle:
		return m.TitleView()
	case ViewTest:
		return m.TestView()
	case ViewResults:
		return m.ResultsView()
	default:
		return "Unknown state"
	}
}

func (m model) TitleView() string {
	header := `
   __                  __           __
  / /_  ___   _____   / /  ____ _  / /_   __  __    ____   ___
 / __/ / _ \ / ___/  / /  / __  / / __/  / / / /   / __ \ / _ \
/ /_  /  __// /__   / /  / /_/ / / /_   / /_/ /   / /_/ //  __/
\__/  \___/ \___/  /_/   \__,_/  \__/   \__, /   / .___/ \___/
                                       /____/   /_/
`
	cmds := []keybind{
		{key: "Enter", cmd: "begin test"},
	}

	title := lipgloss.NewStyle().Foreground(colorText).Render(header)
	footer := renderFooter(cmds, m.width)
	container := lipgloss.NewStyle().Background(colorBackground).Height(m.height-statusBarHeight).Width(m.width).Align(lipgloss.Center, lipgloss.Center).Render(title)
	view := lipgloss.JoinVertical(lipgloss.Center, container, footer)

	// return the view with full window size, background color, and centered
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(colorBackground).
		Render(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view))
}

func (m model) TestView() string {
	var styledText string
	for i, c := range m.text {
		if i < len(m.typed) {
			if m.typed[i] == byte(c) {
				styledText += lipgloss.NewStyle().Foreground(colorCorrect).Render(string(c))
			} else {
				styledText += lipgloss.NewStyle().Foreground(colorIncorrect).Render(string(c))
			}
		} else {
			styledText += lipgloss.NewStyle().Foreground(colorText).Render(string(c))
		}
	}

	cmds := []keybind{
		{key: "Escape", cmd: "stop"},
		{key: "Enter", cmd: "restart"},
		{key: "Control+C", cmd: "quit"},
	}
	container := lipgloss.NewStyle().Background(colorBackground).Height(m.height-statusBarHeight).Width(m.width).Align(lipgloss.Center, lipgloss.Center).Render(styledText)
	footer := renderFooter(cmds, m.width)
	view := lipgloss.JoinVertical(lipgloss.Center, container, footer)

	// return the view with full window size, background color, and centered
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(colorBackground).
		Render(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view))
}

func (m model) ResultsView() string {
	successStyle := lipgloss.NewStyle().Foreground(colorCorrect).Render(fmt.Sprintf("Test completed in %.2f seconds!", m.lastDuration.Seconds()))
	statsStyle := lipgloss.NewStyle().Foreground(colorCorrect).Render(fmt.Sprintf("WPM: %d | CPM: %d", m.statistics.wpm, m.statistics.cpm))
	statsContent := lipgloss.JoinVertical(lipgloss.Center, successStyle, statsStyle)
	statsContainer := lipgloss.NewStyle().
		Background(colorBackground).
		Padding(2, 4).
		Align(lipgloss.Center).
		Render(statsContent)

	container := lipgloss.NewStyle().Background(colorBackground).Height(m.height-statusBarHeight).Width(m.width).Align(lipgloss.Center, lipgloss.Center).Render(statsContainer)

	cmds := []keybind{
		{key: "Enter", cmd: "start"},
		{key: "Control+C", cmd: "quit"},
	}

	footer := renderFooter(cmds, m.width)
	view := lipgloss.JoinVertical(lipgloss.Center, container, footer)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(colorBackground).
		Render(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view))
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Oops! Error: %v", err)
		os.Exit(1)
	}
}
