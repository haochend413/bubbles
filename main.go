package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/bubbles/statusbar"
)

type model struct {
	statusbar statusbar.Model
	width     int
	height    int
}

func initialModel() model {
	// Create a new statusbar
	var sb statusbar.Model
	sb = sb.New()

	// Set dimensions
	sb.SetWidth(80)
	sb.SetHeight(1)

	// Add left elements
	sb.AddLeft(10, "Status: OK")
	sb.AddLeft(15, "CPU: 23%")

	// Add right elements
	sb.AddRight(20, "Memory: 512MB/4GB")
	sb.AddRight(12, "11:45 AM")
	// sb.RemoveRight(0)

	// Customize elements
	cpuElem := sb.GetLeft(0)
	cpuElem.SetColors("15", "27") // White text on blue background

	timeElem := sb.GetRight(0)
	timeElem.SetColors("0", "226") // Black text on yellow background

	return model{
		statusbar: sb,
		width:     80,
		height:    24,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			// Update a status element when "r" is pressed
			cpuElem := m.statusbar.GetLeft(1)
			cpuElem.SetValue("CPU: 78%")
			cpuElem.SetColors("0", "196") // Black text on red background
		case "t":
			// Update time when "t" is pressed
			timeElem := m.statusbar.GetRight(1)
			timeElem.SetValue("12:00 PM")

		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusbar.SetWidth(msg.Width)
	}
	return m, nil
}

func (m model) View() string {
	// Create a content area
	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-1). // Leave space for statusbar
		Align(lipgloss.Center, lipgloss.Center).
		Render("Press:\n" +
			"q to quit\n" +
			"r to update CPU usage\n" +
			"t to update time\n" +
			"b to toggle border on status")

	// Render the statusbar at the bottom
	return lipgloss.JoinVertical(
		lipgloss.Bottom,
		content,
		m.statusbar.Render(),
	)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
