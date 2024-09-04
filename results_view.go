package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type stats struct {
	wpm int
	cpm int
}

type ResultsView struct {
	duration   time.Duration
	statistics stats

	width  int
	height int
}

func NewResultsView(width int, height int, duration time.Duration, chars_typed int) *ResultsView {
	durationInMinutes := duration.Minutes()

	statistics := stats{
		cpm: int(float64(chars_typed) / durationInMinutes),
		wpm: int(float64(chars_typed) / 5 / durationInMinutes),
	}

	return &ResultsView{
		duration:   duration,
		statistics: statistics,
		width:      width,
		height:     height,
	}
}

func (rv ResultsView) Init() tea.Cmd {
	return nil
}

func (rv ResultsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return NewTestView(rv.width, rv.height), nil
		}
	case tea.WindowSizeMsg:
		rv.width = msg.Width
		rv.height = msg.Height
	}

	return rv, nil
}

func (rv ResultsView) View() string {
	successStyle := lipgloss.NewStyle().Foreground(colorCorrect).Render(fmt.Sprintf("Test completed in %.2f seconds!", rv.duration.Seconds()))
	statsStyle := lipgloss.NewStyle().Foreground(colorCorrect).Render(fmt.Sprintf("WPM: %d | CPM: %d", rv.statistics.wpm, rv.statistics.cpm))
	statsContent := lipgloss.JoinVertical(lipgloss.Center, successStyle, statsStyle)
	statsContainer := lipgloss.NewStyle().
		Background(colorBackground).
		Padding(2, 4).
		Align(lipgloss.Center).
		Render(statsContent)

	container := lipgloss.NewStyle().Background(colorBackground).Height(rv.height-statusBarHeight).Width(rv.width).Align(lipgloss.Center, lipgloss.Center).Render(statsContainer)

	cmds := []keybind{
		{key: "Enter", cmd: "start"},
		{key: "Control+C", cmd: "quit"},
	}

	footer := renderFooter(cmds, rv.width)
	view := lipgloss.JoinVertical(lipgloss.Center, container, footer)

	return lipgloss.NewStyle().
		Width(rv.width).
		Height(rv.height).
		Background(colorBackground).
		Render(lipgloss.Place(rv.width, rv.height, lipgloss.Center, lipgloss.Center, view))
}
