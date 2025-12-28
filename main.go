package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/bubbles/textarea_vim"
)

type model struct {
	textarea    textarea_vim.Model
	width       int
	height      int
	showHelp    bool
	keypressIdx int
}

type tickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	ta := textarea_vim.New()
	ta.SetWidth(60)  // Slightly larger width
	ta.SetHeight(10) // Slightly larger height
	ta.Focus()       // Ensure textarea is focused

	return model{
		textarea:    ta,
		width:       0,
		height:      0,
		showHelp:    false,
		keypressIdx: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.textarea.Focus(),
		tickEvery(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		fmt.Printf("\r\033[K")
		fmt.Printf(strconv.Itoa(m.keypressIdx) + ":" + msg.String())
		m.keypressIdx += 1
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp
		}
		var taCmd tea.Cmd
		m.textarea, taCmd = m.textarea.Update(msg)
		cmds = append(cmds, taCmd)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		cmds = append(cmds, tickEvery())
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	boxWidth := 62  // Slightly larger than textarea
	boxHeight := 12 // Slightly larger than textarea
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Width(boxWidth).
		Height(boxHeight)

	content := m.textarea.View()

	centered := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		borderStyle.Render(content),
	)

	if m.showHelp {
		return helpView(m) + "\n" + centered
	}
	return centered
}

func helpView(m model) string {
	s := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center)

	helpTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		MarginBottom(1).
		Render("TEXTAREA_VIM DEMO HELP")

	helpContent := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(40).
		Render(`
This demo shows a centered, bordered textarea_vim component.

Keyboard commands:
-----------------
? - Toggle this help screen
q - Quit application

All other keys are passed to the textarea_vim component.
`)

	return s.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			helpTitle,
			helpContent,
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("Press ? to return"),
		),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
