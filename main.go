package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	start        time.Time
	text         string
	typed        string
	lastDuration time.Duration
	typing       bool
	success      bool
}

func (m *model) Start() {
	if !m.typing { // restart if not already typing
		m.typed = ""
		m.start = time.Now()
	}

	m.typing = true
}

func (m *model) Stop() {
	m.typing = false
}

func (m *model) Finish() {
	m.lastDuration = time.Since(m.start)
	m.success = true
	m.Stop()
}

func initialModel() model {
	return model{
		text: "hello world this is some text to type",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc: // leave typing mode or quit if out
			if m.typing {
				m.Stop()
			} else {
				return m, tea.Quit
			}
		case tea.KeyBackspace:
			if m.typing && len(m.typed) > 0 {
				m.typed = m.typed[:len(m.typed)-1]
			}
		case tea.KeyEnter:
			m = initialModel()
		default:
			m.Start()
			m.typed += string(msg.Runes)

			if m.typed == m.text {
				m.Finish()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Tecla\ntype this sentence:\n\n"

	if m.success {
		s += fmt.Sprintf("success! finished in %v seconds!\n", m.lastDuration)
	} else {
		s += fmt.Sprintf("%s\n%s\n", m.text, m.typed)
	}

	if m.typing {
		s += "\nesc to stop"
		s += "\tenter to restart"
	} else {
		s += "\nesc to quit"
	}

	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Oops! Error: %v", err)
		os.Exit(1)
	}
}
