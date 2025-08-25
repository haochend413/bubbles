package statusbar

import (
	"github.com/charmbracelet/lipgloss"
)

type Partition int

const (
	Left Partition = iota
	Right
)

type Elem struct {
	Content string
	Width   int
	// Partition   Partition
	BgColor string
	FgColor string
}

func (e *Elem) SetValue(content string) {
	e.Content = content
}

func (e *Elem) SetColors(fg, bg string) {
	e.FgColor = fg
	e.BgColor = bg
}

func (e *Elem) SetWidth(width int) {
	e.Width = width
}

//	func (e *Elem) SetPartition(p Partition) {
//		e.Partition = p
//	}
func (e Elem) Render(h int) string {
	// Custom border: block characters that touch the content directly

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(e.FgColor)).
		Background(lipgloss.Color(e.BgColor)).
		Width(e.Width).
		Height(h)

	return style.Render(" " + e.Content)
}

type Model struct {
	LeftELems  []*Elem
	RightElems []*Elem
	Height     int
	Width      int
}

func (m *Model) New() Model {
	return Model{
		LeftELems:  []*Elem{},
		RightElems: []*Elem{},
		Height:     1,
		Width:      100,
	}
}

func (m *Model) AddLeft(w int, c string) {
	newElem := &Elem{
		Content: c,
		Width:   w,
		BgColor: "236",
		FgColor: "252",
	}
	m.LeftELems = append(m.LeftELems, newElem)
}

func (m *Model) RemoveLeft(i int) {
	if i < len(m.LeftELems) {
		m.LeftELems = append(m.LeftELems[:i], m.LeftELems[i+1:]...)
	}
}

func (m *Model) RemoveRight(i int) {
	if i < len(m.RightElems) {
		m.RightElems = append(m.RightElems[:i], m.RightElems[i+1:]...)
	}
}

func (m *Model) AddRight(w int, c string) {
	newElem := &Elem{
		Content: c,
		Width:   w,
		BgColor: "236",
		FgColor: "252",
	}
	m.RightElems = append(m.RightElems, newElem)
}

func (m *Model) GetLeft(index int) *Elem {
	return m.LeftELems[index]
}

func (m *Model) GetRight(index int) *Elem {
	return m.RightElems[index]
}

func (m *Model) SetWidth(w int) {
	m.Width = w
}

func (m *Model) SetHeight(h int) {
	m.Height = h
}

func (m Model) Render() string {
	// Render left elements
	leftContent := ""
	for _, elem := range m.LeftELems {
		leftContent += elem.Render(m.Height)
	}

	// Render right elements
	rightContent := ""
	for _, elem := range m.RightElems {

		rightContent += elem.Render(m.Height)
	}

	// Calculate space between left and right elements
	leftWidth := lipgloss.Width(leftContent)
	rightWidth := lipgloss.Width(rightContent)
	middleWidth := max(0, m.Width-leftWidth-rightWidth)

	// Create empty middle space with appropriate width
	middleSpace := lipgloss.NewStyle().
		Width(middleWidth).
		Height(m.Height).
		Render("")

	// Join left content, middle space, and right content
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		leftContent,
		middleSpace,
		rightContent,
	)
}
