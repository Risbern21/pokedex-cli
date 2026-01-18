package styles

import "github.com/charmbracelet/lipgloss"

var (
	ListFocusedStyle   = lipgloss.NewStyle().Width(30).Border(lipgloss.ThickBorder(), true, true).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("9"))
	ListUnFocusedStyle = lipgloss.NewStyle().Width(30).Border(lipgloss.ThickBorder(), true, true).BorderStyle(lipgloss.RoundedBorder())

	ListItemSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)

	TitleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 3)
	}()

	InfoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return TitleStyle.BorderStyle(b)
	}()

	SpinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("200"))

	PokedexListTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("206")).Background(lipgloss.Color("190"))
)
