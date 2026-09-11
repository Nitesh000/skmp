package tui

import (
	"fmt"
	"strings"

	"github.com/Nitesh000/skmp/registry"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	tabSkills tab = iota
	tabbundles
)

type Model struct {
	skills      []registry.Skill
	bundles     []registry.Bundle
	filtered        []registry.Skill
	filteredBundles []registry.Bundle
	installed   map[string]bool
	cursor      int
	activeTab   tab
	width       int
	height      int
	err         error
	searchInput textinput.Model
	searchFocus bool
}

// messages
type (
	indexLoadedMsg struct{ idx *registry.Index }
	indexErrMsg    struct{ err error }
)

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "search..."
	ti.Width = 30

	return Model{
		installed:   map[string]bool{},
		searchInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, loadIndex())
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
	var cmds []tea.Cmd

	// Always update the text input (handles blinking, etc)
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case indexLoadedMsg:
		m.skills = msg.idx.Skills
		m.bundles = msg.idx.Bundles
		m.filtered = msg.idx.Skills
		m.filteredBundles = msg.idx.Bundles
		registry.BuildIndex(m.skills, m.bundles)

	case indexErrMsg:
		m.err = msg.err

	case tea.KeyMsg:
		// type inside search box
		if m.searchFocus {
			switch msg.String() {
			case "esc", "enter":
				m.searchFocus = false
				m.searchInput.Blur()
				return m, tea.Batch(cmds...)
			}

			resSkills, _ := registry.SearchSkills(m.searchInput.Value())
			resBundles, _ := registry.SearchBundles(m.searchInput.Value())
			m.filtered = resSkills
			m.filteredBundles = resBundles
			m.cursor = 0

			return m, tea.Batch(cmds...)
		}

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

		case "/":
			m.searchFocus = true
			m.searchInput.Focus()
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

	return m, tea.Batch(cmds...)
}

func (m Model) listLen() int {
	if m.activeTab == tabbundles {
		return len(m.filteredBundles)
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
	tabs := ""
	if m.activeTab == tabSkills {
		tabs += "[1] Skills   "
		tabs += " 2  Bundles"
	} else {
		tabs += " 1  Skills   "
		tabs += "[2] Bundles"
	}

	search := m.searchInput.View()

	padding := m.width - lipgloss.Width(tabs) - lipgloss.Width(search) - 2
	if padding < 0 {
		padding = 0
	}

	return tabs + strings.Repeat(" ", padding) + search
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
	if len(m.filteredBundles) == 0 {
		return "  no bundles found"
	}
	var s strings.Builder
	for i, bun := range m.filteredBundles {
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
	if len(m.filteredBundles) == 0 || m.cursor >= len(m.filteredBundles) {
		return mutedStyle.Render("select a bundle")
	}

	b := m.filteredBundles[m.cursor]

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
