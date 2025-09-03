package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/bubbles/statusbar"
	// "github.com/haochending/bubbles/statusbar" // Adjust this import path as needed
)

// StatusbarDemo represents our application state
type StatusbarDemo struct {
	statusbar statusbar.Model
	width     int
	height    int
	tick      int
	mode      string
	showHelp  bool
	showTags  bool
}

// tickMsg is sent when the timer ticks
type tickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() StatusbarDemo {
	// Initialize ElemsMap
	sb := statusbar.New(
		statusbar.WithWidth(100),
		statusbar.WithHeight(1),
	)
	sb.ElemsMap = make(map[string]*statusbar.Elem) // Initialize the map

	// Add elements with tags to left side
	statusElem := sb.AddLeft(20, "Status: Ready").SetColors("0", "46")
	sb.SetTag(statusElem, "status")

	cpuElem := sb.AddLeft(15, "CPU: 42%").SetColors("0", "33")
	sb.SetTag(cpuElem, "cpu")

	memElem := sb.AddLeft(15, "Mem: 1.2GB").SetColors("0", "177")
	sb.SetTag(memElem, "memory")

	// Add elements with tags to right side
	helpElem := sb.AddRight(22, "Press ? for help").SetColors("252", "236")
	sb.SetTag(helpElem, "help")

	timeElem := sb.AddRight(10, time.Now().Format("15:04:05")).SetColors("0", "226")
	sb.SetTag(timeElem, "time")

	return StatusbarDemo{
		statusbar: sb,
		width:     100,
		height:    30,
		tick:      0,
		mode:      "normal",
		showHelp:  false,
		showTags:  false,
	}
}

func (m StatusbarDemo) Init() tea.Cmd {
	return tickEvery()
}

func (m StatusbarDemo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "?":
			// Toggle help screen
			m.showHelp = !m.showHelp
			helpText := "Press ? for help"
			if m.showHelp {
				helpText = "Press ? to exit help"
			}
			m.statusbar.GetTag("help").SetValue(helpText)

		case "t":
			// Toggle tag view
			m.showTags = !m.showTags

		case "1":
			// Normal mode
			m.mode = "normal"
			statusElem := m.statusbar.GetTag("status")
			if statusElem != nil {
				statusElem.SetValue("Status: Ready").SetColors("0", "46")
			}

		case "2":
			// Warning mode
			m.mode = "warning"
			statusElem := m.statusbar.GetTag("status")
			if statusElem != nil {
				statusElem.SetValue("Status: Warning").SetColors("0", "208")
			}

		case "3":
			// Error mode
			m.mode = "error"
			statusElem := m.statusbar.GetTag("status")
			if statusElem != nil {
				statusElem.SetValue("Status: Error").SetColors("255", "196")
			}

		case "c":
			// Clear CPU/Mem stats
			cpuElem := m.statusbar.GetTag("cpu")
			if cpuElem != nil {
				cpuElem.SetValue("CPU: 0%")
			}

			memElem := m.statusbar.GetTag("memory")
			if memElem != nil {
				memElem.SetValue("Mem: 0.0GB")
			}

		case "r":
			// Remove an element (as a test)
			if len(m.statusbar.LeftElems) > 0 {
				m.statusbar.RemoveLeft(len(m.statusbar.LeftElems) - 1)
			}

		case "a":
			// Add a new element (as a test)
			newElem := m.statusbar.AddLeft(15, "Added: "+fmt.Sprint(m.tick)).SetColors("0", "99")
			m.statusbar.SetTag(newElem, fmt.Sprintf("added-%d", m.tick))
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusbar.SetWidth(msg.Width)

	case tickMsg:
		// Update the time
		timeElem := m.statusbar.GetTag("time")
		if timeElem != nil {
			timeElem.SetValue(time.Now().Format("15:04:05"))
		}

		// Update tick counter
		m.tick++

		// Simulate changing CPU usage
		cpuValue := 40 + ((m.tick * 5) % 50)
		cpuElem := m.statusbar.GetTag("cpu")
		if cpuElem != nil {
			// Change CPU color based on usage
			cpuColor := "33" // Default blue
			if cpuValue > 70 {
				cpuColor = "208" // Orange for high
			}
			if cpuValue > 90 {
				cpuColor = "196" // Red for very high
			}
			cpuElem.SetValue(fmt.Sprintf("CPU: %d%%", cpuValue)).SetColors("0", cpuColor)
		}

		// Simulate changing memory usage
		memValue := 0.5 + float64((m.tick*7)%70)/20
		memElem := m.statusbar.GetTag("memory")
		if memElem != nil {
			memElem.SetValue(fmt.Sprintf("Mem: %.1fGB", memValue))
		}

		return m, tickEvery()
	}

	return m, nil
}

func (m StatusbarDemo) View() string {
	if m.showHelp {
		return helpView(m) + "\n" + m.statusbar.Render()
	}

	if m.showTags {
		return tagView(m) + "\n" + m.statusbar.Render()
	}

	return mainView(m) + "\n" + m.statusbar.Render()
}

func mainView(m StatusbarDemo) string {
	s := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-1). // Subtract 1 for statusbar
		Align(lipgloss.Center, lipgloss.Center)

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		MarginBottom(1).
		Render("STATUSBAR COMPONENT DEMO")

	// Create a simple animation
	anim := ""
	for i := 0; i < 20; i++ {
		char := "·"
		if (i+m.tick)%5 == 0 {
			char = "○"
		}
		color := fmt.Sprintf("%d", 30+(i+m.tick)%7)
		anim += lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Render(char)
	}

	modeStyle := lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
	modeText := ""

	switch m.mode {
	case "normal":
		modeText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")). // Green
			Render("Normal Mode")
	case "warning":
		modeText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")). // Orange
			Render("Warning Mode")
	case "error":
		modeText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // Red
			Render("Error Mode")
	}

	keysHelp := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")). // Dark gray
		Render(`
Press ? for help screen
t to view tag information
1-3 to change modes
c to clear stats
r to remove an element
a to add an element
q to quit
`)

	stats := fmt.Sprintf("CPU: %s\nMemory: %s\nTime: %s\nTick: %d",
		m.statusbar.GetTag("cpu").Content,
		m.statusbar.GetTag("memory").Content,
		m.statusbar.GetTag("time").Content,
		m.tick)

	return s.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			anim,
			modeStyle.Render(modeText),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Render(stats),
			"",
			keysHelp,
		),
	)
}

func tagView(m StatusbarDemo) string {
	s := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-1). // Subtract 1 for statusbar
		Align(lipgloss.Center, lipgloss.Center)

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		MarginBottom(1).
		Render("STATUSBAR TAG INFORMATION")

	// Display information about tags
	tagsInfo := "Registered Tags:\n\n"

	// Check if map is initialized
	if m.statusbar.ElemsMap == nil {
		tagsInfo += "No tags found (ElemsMap is nil)"
	} else {
		for tag, elem := range m.statusbar.ElemsMap {
			tagsInfo += fmt.Sprintf("Tag: %s\n", tag)
			tagsInfo += fmt.Sprintf("  Content: %s\n", elem.Content)
			tagsInfo += fmt.Sprintf("  Colors: fg=%s, bg=%s\n", elem.FgColor, elem.BgColor)
			tagsInfo += fmt.Sprintf("  Width: %d\n\n", elem.Width)
		}

		if len(m.statusbar.ElemsMap) == 0 {
			tagsInfo += "No tags registered."
		}
	}

	return s.Render(
		lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Render(tagsInfo),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("Press t to return to main view"),
		),
	)
}

func helpView(m StatusbarDemo) string {
	s := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-1). // Subtract 1 for statusbar
		Align(lipgloss.Center, lipgloss.Center)

	helpTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		MarginBottom(1).
		Render("STATUSBAR HELP")

	helpContent := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Width(60).
		Render(`
This demo shows the statusbar component in action.

Keyboard commands:
-----------------
? - Toggle this help screen
t - View tag information
1 - Switch to Normal mode
2 - Switch to Warning mode
3 - Switch to Error mode
c - Clear CPU/Memory stats
r - Remove last left element (test)
a - Add a new element (test)
q - Quit application

The statusbar shows:
- Status indicator (changes with modes)
- CPU usage (color changes with load)
- Memory usage
- Help hint
- Current time (updates automatically)

This demo showcases:
- Element tagging system
- Dynamic styling
- Real-time updates
- Element addition/removal
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
	// Fix the for loop in WithLeftLen and WithRightLen before running
	// fixRangeLoop()

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
