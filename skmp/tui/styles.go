package tui

import "github.com/charmbracelet/lipgloss"

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#2B8A8A", Dark: "#5CC8C8"}
	tabBg     = lipgloss.AdaptiveColor{Light: "#135B5B", Dark: "#1A7373"}
	muted     = lipgloss.AdaptiveColor{Light: "#888888", Dark: "#626262"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	warn      = lipgloss.AdaptiveColor{Light: "#B8860B", Dark: "#E6C07B"}

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(subtle)

	activeBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(highlight)

	selectedStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"})

	mutedStyle = lipgloss.NewStyle().
			Foreground(muted)

	labelStyle = lipgloss.NewStyle().
			Foreground(muted).
			Width(12)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"})

	tagStyle = lipgloss.NewStyle().
			Foreground(special).
			Background(lipgloss.AdaptiveColor{Light: "#E8FFE8", Dark: "#1A3A1A"}).
			Padding(0, 1).
			MarginRight(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight)

	installedStyle = lipgloss.NewStyle().
			Foreground(special).
			Bold(true)

	partialStyle = lipgloss.NewStyle().
			Foreground(warn).
			Bold(true)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(tabBg).
			Bold(true)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(muted)
)
