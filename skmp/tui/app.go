package tui

import (
	"fmt"
	"strings"

	"github.com/Nitesh000/skmp/registry"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	tabSkills tab = iota
	tabbundles
)

type Model struct {
	skills    []registry.Skill
	bundles   []registry.Bundle
	filtered  []registry.Skill
	installed map[string]bool
	cursor    int
	activeTab tab
	width     int
	height    int
	err       error
}

// messages
type (
	indexLoadedMsg struct{ idx *registry.Index }
	indexErrMsg    struct{ err error }
)

func New() Model {
	return Model{
		installed: map[string]bool{},
	}
}

func (m Model) Init() tea.Cmd {
	return loadIndex()
}

func loadIndex() tea.Cmd {
	return func() tea.Msg {
		idx, err := registry.Load()
		if err != nil {
			return indexErrMsg{err}
		}
		return indexLoadedMsg{idx}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case indexLoadedMsg:
		m.skills = msg.idx.Skills
		m.bundles = msg.idx.Bundles
		m.filtered = msg.idx.Skills
		registry.BuildIndex(m.skills)

	case indexErrMsg:
		m.err = msg.err

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "j", "down":
			m.cursor++
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "1":
			m.activeTab = tabSkills
			m.cursor = 0

		case "2":
			m.activeTab = tabbundles
			m.cursor = 0
		}

		// clamp cursor
		max := m.listLen() - 1
		if m.cursor > max {
			m.cursor = max
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
	}

	return m, nil
}

func (m Model) listLen() int {
	if m.activeTab == tabbundles {
		return len(m.bundles)
	}

	return len(m.filtered)
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %s\n\npress q to quit", m.err)
	}
	if m.width == 0 {
		return "loading..."
	}

	tabBar := m.tabBarView()
	body := m.bodyView()
	help := mutedStyle.Render("  j/k move · 1/2 tabs · q quit")

	return lipgloss.JoinVertical(lipgloss.Left, tabBar, body, help)
}

func (m Model) tabBarView() string {
	s := ""
	if m.activeTab == tabSkills {
		s += "[1] Skills   "
		s += " 2  Bundles"
	} else {
		s += " 1  Skills   "
		s += "[2] Bundles"
	}
	return s
}

func (m Model) listView() string {
	if m.activeTab == tabSkills {
		return m.skillsListView()
	}
	return m.bundlesListView()
}

func (m Model) skillsListView() string {
	if len(m.filtered) == 0 {
		return "  no skills found"
	}
	var s strings.Builder
	for i, skill := range m.filtered {
		prefix := "  "
		if i == m.cursor {
			prefix = "▶ "
		}
		fmt.Fprintf(&s, "%s%s\n", prefix, skill.Name)
	}
	return s.String()
}

func (m Model) bundlesListView() string {
	if len(m.bundles) == 0 {
		return "  no bundles found"
	}
	var s strings.Builder
	for i, bun := range m.bundles {
		prefix := "  "
		if i == m.cursor {
			prefix = "▶ "
		}
		fmt.Fprintf(&s, "%s%s (%d skills)\n", prefix, bun.Name, len(bun.Skills))
	}
	return s.String()
}

func (m Model) bodyView() string {
	listW := m.width / 3
	detailW := m.width - listW - 3 // 3 = borders + gap
	innerH := m.height - 4

	list := m.paneList(listW, innerH)
	detail := m.paneDetail(detailW, innerH)

	return lipgloss.JoinHorizontal(lipgloss.Top, list, detail)
}

func (m Model) paneList(w, h int) string {
	content := m.listView()

	style := borderStyle.Width(w).Height(h)
	return style.Render(content)
}

func (m Model) paneDetail(w, h int) string {
	var content string
	if m.activeTab == tabSkills {
		content = m.skillDetailView(w)
	} else {
		content = m.bundleDetailsView(w)
	}

	style := borderStyle.Width(w).Height(h)
	return style.Render(content)
}

func (m Model) skillDetailView(w int) string {
	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
		return mutedStyle.Render("select a skill")
	}

	s := m.filtered[m.cursor]

	tags := ""
	for _, t := range s.Tags {
		tags += tagStyle.Render(t)
	}

	installed := mutedStyle.Render("○ not installed")
	if m.installed[s.Name] {
		installed = installedStyle.Render("● installed")
	}

	return lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render(s.Name), "", wrapText(s.Description, w-4),
		"",
		labelStyle.Render("Version")+valueStyle.Render(s.Version),
		labelStyle.Render("Author")+valueStyle.Render(s.Author),
		labelStyle.Render("Harnesses")+valueStyle.Render(strings.Join(s.Harnesses, ",")),
		labelStyle.Render("Tags")+tags,
		"",
		installed,
	)
}

func (m Model) bundleDetailsView(w int) string {
	if len(m.bundles) == 0 || m.cursor >= len(m.bundles) {
		return mutedStyle.Render("select a bundle")
	}

	b := m.bundles[m.cursor]

	// count installed skills in bundle
	installedCount := 0
	for _, name := range b.Skills {
		if m.installed[name] {
			installedCount++
		}
	}

	total := len(b.Skills)

	var status string
	switch {
	case installedCount == total:
		status = installedStyle.Render("● installed")
	case installedCount > 0:
		status = installedStyle.Render(fmt.Sprintf("◐ partial (%d/%d)", installedCount, total))
	default:
		status = mutedStyle.Render("○ not installed")
	}

	skillList := ""
	for _, name := range b.Skills {
		badge := "  "
		if m.installed[name] {
			badge = installedStyle.Render("✓ ")
		}
		skillList += badge + name + "\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(b.Name),
		"",
		wrapText(b.Description, w-4),
		"",
		labelStyle.Render("Author")+valueStyle.Render(b.Author),
		labelStyle.Render("Skills")+valueStyle.Render(fmt.Sprintf("%d", total)),
		"",
		status,
		"",
		skillList,
	)
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)

	var lines []string
	line := ""
	for _, w := range words {
		if len(line)+len(w)+1 > width && line != "" {
			lines = append(lines, line)
			line = w
		} else {
			if line != "" {
				line += " "
			}
			line += w
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
