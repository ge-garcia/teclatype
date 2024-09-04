package main

import (
	"bufio"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	colorText        = lipgloss.Color("#636363")
	colorIncorrect   = lipgloss.Color("#CA3F3F")
	colorCorrect     = lipgloss.Color("#FFFFFF")
	colorBackground  = lipgloss.Color("#181515")
	colorContainer   = lipgloss.Color("#221E1E")
	colorCurrentWord = lipgloss.Color("#a8a8a8")
	colorNextWord    = lipgloss.Color("#808080")
	colorHighlight   = lipgloss.Color("#d4d4d4")
	colorBgHighlight = lipgloss.Color("#5e5e5e")
)

type ViewState int

const (
	statusBarHeight = 1
)

type model struct {
	view tea.Model

	width  int
	height int
}

type keybind struct {
	key string
	cmd string
}

func readWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func initialModel() model {
	return model{
		view: NewTitleView(0, 0),
	}
}

func (m model) Init() tea.Cmd {
	return m.view.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	view, cmd := m.view.Update(msg)
	m.view = view
	return m, cmd
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
	return m.view.View()
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Oops! Error: %v", err)
		os.Exit(1)
	}
}
