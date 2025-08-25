package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/haochend413/bubbles/statusbar"
	// "github.com/haochending/bubbles/statusbar" // Adjust this path as needed
)

// StatusbarDemoModel represents our application state
type StatusbarDemoModel struct {
	statusbar statusbar.Model
	width     int
	height    int
	tick      int
	cpuUsage  int
	memUsage  int
	mode      string
	showHelp  bool
}

// tickMsg is sent when the timer ticks
type tickMsg time.Time

func tickEvery() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() StatusbarDemoModel {
	// Create our statusbar
	sb := statusbar.New(
		statusbar.WithWidth(100),
		statusbar.WithHeight(1),
	)

	// Add some elements to the left side
	sb.AddLeft(20, "Status: Ready").SetColors("0", "46") // Black on green
	sb.AddLeft(15, "CPU: 42%").SetColors("0", "33")      // Black on blue
	sb.AddLeft(15, "Mem: 1.2GB").SetColors("0", "177")   // Black on purple

	// Add some elements to the right side
	sb.AddRight(22, "Press ? for help").SetColors("252", "236")          // Light text on dark gray
	sb.AddRight(10, time.Now().Format("15:04:05")).SetColors("0", "226") // Black on yellow

	return StatusbarDemoModel{
		statusbar: sb,
		width:     100,
		height:    30,
		tick:      0,
		cpuUsage:  42,
		memUsage:  30,
		mode:      "normal",
		showHelp:  false,
	}
}

func (m StatusbarDemoModel) Init() tea.Cmd {
	return tickEvery()
}

func (m StatusbarDemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.statusbar.GetRight(0).SetValue(helpText)

		case "1":
			// Normal mode
			m.mode = "normal"
			m.statusbar.GetLeft(0).SetValue("Status: Ready")

		case "2":
			// Warning mode
			m.mode = "warning"
			m.statusbar.GetLeft(0).SetValue("Status: Warning")

		case "3":
			// Error mode
			m.mode = "error"
			m.statusbar.GetLeft(0).SetValue("Status: Error")

		case "c":
			// Clear CPU/Mem stats
			m.cpuUsage = 0
			m.memUsage = 0
			m.statusbar.GetLeft(1).SetValue("CPU: 0%")
			m.statusbar.GetLeft(2).SetValue("Mem: 0.0GB")
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.statusbar.SetWidth(msg.Width)

	case tickMsg:
		// Update the time
		m.statusbar.GetRight(1).SetValue(time.Now().Format("15:04:05"))

		// Update tick counter
		m.tick++

		// Simulate changing CPU usage
		cpuDelta := m.tick % 5
		if m.tick%10 >= 5 {
			cpuDelta = -cpuDelta
		}
		m.cpuUsage += cpuDelta
		if m.cpuUsage < 0 {
			m.cpuUsage = 0
		} else if m.cpuUsage > 100 {
			m.cpuUsage = 100
		}

		// // Change CPU color based on usage
		// cpuColor := "33" // Default blue
		// if m.cpuUsage > 70 {
		// 	cpuColor = "208" // Orange for high
		// }
		// if m.cpuUsage > 90 {
		// 	cpuColor = "196" // Red for very high
		// }
		m.statusbar.GetLeft(1).
			SetValue(fmt.Sprintf("CPU: %d%%", m.cpuUsage))
		// 	SetColors("0", cpuColor)

		// Simulate changing memory usage
		memDelta := (m.tick % 3) * 10
		if m.tick%6 >= 3 {
			memDelta = -memDelta
		}
		m.memUsage += memDelta
		if m.memUsage < 0 {
			m.memUsage = 0
		} else if m.memUsage > 100 {
			m.memUsage = 100
		}

		memGB := float64(m.memUsage) / 100 * 4 // Simulate 4GB max
		m.statusbar.GetLeft(2).SetValue(fmt.Sprintf("Mem: %.1fGB", memGB))

		return m, tickEvery()
	}

	return m, nil
}

func (m StatusbarDemoModel) View() string {
	if m.showHelp {
		return helpView(m) + "\n" + m.statusbar.Render()
	}
	return mainView(m) + "\n" + m.statusbar.Render()
}

func mainView(m StatusbarDemoModel) string {
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
1-3 to change modes
c to clear stats
q to quit
`)

	stats := fmt.Sprintf("CPU: %d%%\nMemory: %.1f GB\nTick: %d",
		m.cpuUsage,
		float64(m.memUsage)/100*4,
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

func helpView(m StatusbarDemoModel) string {
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
1 - Switch to Normal mode
2 - Switch to Warning mode
3 - Switch to Error mode
c - Clear CPU/Memory stats
q - Quit application

The statusbar shows:
- Status indicator (changes with modes)
- CPU usage (color changes with load)
- Memory usage
- Help hint
- Current time (updates automatically)
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

// // This function prints a warning about the for loop issue
// func fixRangeLoop() {
// 	fmt.Println("Before running this demo, please fix the following issue in statusbar.go:")
// 	fmt.Println("In WithLeftLen and WithRightLen functions, change:")
// 	fmt.Println("for i := range n { ... }")
// 	fmt.Println("to:")
// 	fmt.Println("for i := 0; i < n; i++ { ... }")
// 	fmt.Println("\nPress Enter to continue...")
// 	fmt.Scanln()
// }
