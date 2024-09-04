package main

import (
	"fmt"
	"os"

	views "github.com/ge-garcia/tecla/internal/views"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	view tea.Model

	width  int
	height int
}

func initialModel() model {
	return model{
		view: views.NewTitleView(0, 0),
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
