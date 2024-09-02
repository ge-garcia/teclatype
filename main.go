package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBgRed  = "\033[41m"
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
	typing       bool
	success      bool
	statistics   stats
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

func (m *model) CalculateStats() {
	durationInMinutes := m.lastDuration.Minutes()
	m.statistics.cpm = int(float64(len(m.typed)) / durationInMinutes)
	m.statistics.wpm = int(float64(len(m.typed)) / 5 / durationInMinutes)
}

func (m *model) Finish() {
	m.lastDuration = time.Since(m.start)
	m.success = true
	m.CalculateStats()
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

func colorText(text, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, colorReset)
}

func (m model) View() string {
	s := "Tecla\ntype this sentence:\n\n"

	if m.success {
		s += fmt.Sprintf("success! finished in %v seconds!\n", m.lastDuration)
		s += fmt.Sprintf("WPM: %d CPM: %d\n", m.statistics.wpm, m.statistics.cpm)
	} else {
		for i, c := range m.text {
			if i < len(m.typed) {
				if m.typed[i] == byte(c) {
					s += colorText(string(c), colorGreen)
				} else {
					if string(c) != " " {
						s += colorText(string(c), colorRed)
					} else {
						s += colorText(" ", colorBgRed)
					}
				}
			} else {
				s += string(c)
			}
		}
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
