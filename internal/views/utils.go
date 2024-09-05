package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type keybind struct {
	key string
	cmd string
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
		Background(ColorContainer).
		Foreground(ColorText).
		Height(StatusBarHeight).
		Width(width).
		Padding(0, 2).
		Render(footerContent)
}
