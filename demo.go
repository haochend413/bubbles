package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/haochend413/bubbles/textarea"
)

type model struct {
	textarea textarea.Model
}

func NewModel() model {
	ta := textarea.New()
	ta.Placeholder = "Start typing... (Press 'i' for INSERT mode, 'esc' for VIEW mode)"
	ta.SetWidth(80)
	ta.SetHeight(15)
	ta.Focus()
	ta.SetValue("Welcome to the Textarea Demo!\n\nNew features:\n- Press 'i' to enter INSERT mode (editing enabled)\n- Press 'esc' to enter VIEW mode (navigation only)\n- Mode and word count shown in status bar\n\nTry it out!")

	return model{
		textarea: ta,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	header := headerStyle.Render("Textarea Demo - VIM-like Mode Switching") + "\n"

	help := infoStyle.Render(
		"Mode: Press 'i' (INSERT) | 'esc' (VIEW) | ctrl+c/ctrl+q (quit)\n" +
			"INSERT: Full editing | VIEW: Navigation only\n",
	)

	stats := infoStyle.Render(
		fmt.Sprintf(
			"Lines: %d | Cursor: Line %d, Col %d\n",
			m.textarea.LineCount(),
			m.textarea.Line()+1,
			m.textarea.Column(),
		),
	)

	return tea.NewView(
		header + help + stats + "\n" + m.textarea.View(),
	)
}

func main() {
	p := tea.NewProgram(NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
