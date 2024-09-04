package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TitleView struct {
	width  int
	height int
}

func NewTitleView(width int, height int) *TitleView {
	return &TitleView{
		width:  width,
		height: height,
	}
}

func (tv TitleView) Init() tea.Cmd {
	return nil
}

func (tv TitleView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return tv, tea.Quit
		case tea.KeyEnter:
			return NewTestView(tv.width, tv.height), nil
		}
	case tea.WindowSizeMsg:
		tv.width = msg.Width
		tv.height = msg.Height
	}

	return tv, nil
}

func (tv TitleView) View() string {
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
	footer := renderFooter(cmds, tv.width)
	container := lipgloss.NewStyle().Background(colorBackground).Height(tv.height-statusBarHeight).Width(tv.width).Align(lipgloss.Center, lipgloss.Center).Render(title)
	view := lipgloss.JoinVertical(lipgloss.Center, container, footer)

	// return the view with full window size, background color, and centered
	return lipgloss.NewStyle().
		Width(tv.width).
		Height(tv.height).
		Background(colorBackground).
		Render(lipgloss.Place(tv.width, tv.height, lipgloss.Center, lipgloss.Center, view))
}
