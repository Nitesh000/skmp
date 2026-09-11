package tui

import (
	"fmt"

	"github.com/Nitesh000/skmp/registry"
	tea "github.com/charmbracelet/bubbletea"
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
	return Model{}
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

	tabs := m.tabBarView()
	list := m.listView()
	help := "\n  j/k move · 1/2 tabs · q quit"

	return tabs + "\n" + list + help
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
	s := ""
	for i, skill := range m.filtered {
		prefix := "  "
		if i == m.cursor {
			prefix = "▶ "
		}
		s += fmt.Sprintf("%s%s\n", prefix, skill.Name)
	}
	return s
}

func (m Model) bundlesListView() string {
	if len(m.bundles) == 0 {
		return "  no bundles found"
	}
	s := ""
	for i, bun := range m.bundles {
		prefix := "  "
		if i == m.cursor {
			prefix = "▶ "
		}
		s += fmt.Sprintf("%s%s (%d skills)\n", prefix, bun.Name, len(bun.Skills))
	}
	return s
}
