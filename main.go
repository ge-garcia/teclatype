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
	if !m.typing { // start timer if we weren't already typing
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
		switch msg.String() {
		case "ctrl+c", "esc": // leave typing mode or quit if out
			if m.typing {
				m.Stop()
			} else {
				return m, tea.Quit
			}
		case "backspace":
			if m.typing && len(m.typed) > 0 {
				m.typed = m.typed[:len(m.typed)-1]
			}
		case "enter": // restart
			m = initialModel()
		default:
			m.Start()
			m.typed += msg.String()

			if m.typed == m.text {
				m.Finish()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Tecla\n\n"

	if m.success {
		s += fmt.Sprintf("success! finished in %v seconds!\n", m.lastDuration)
	} else {
		s += fmt.Sprintf("%s\n%s\n", m.text, m.typed)

		if m.typing {
			s += "typing!\n"
		} else {
			s += "paused\n"
		}
	}

	s += "\nesc to quit"
	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Oops! Error: %v", err)
		os.Exit(1)
	}
}
