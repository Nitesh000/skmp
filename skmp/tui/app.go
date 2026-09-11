package tui

import (
	"fmt"
	"strings"

	"github.com/Nitesh000/skmp/harness"
	"github.com/Nitesh000/skmp/registry"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tab int

const (
	tabSkills tab = iota
	tabBundles
	tabMySkills
	tabMyBundles
)

// minWidth is the point below which the split pane is dropped for a single list.
const minWidth = 60

type Model struct {
	version         string
	skills          []registry.Skill
	bundles         []registry.Bundle
	filtered        []registry.Skill
	filteredBundles []registry.Bundle
	installed       map[string]bool
	loading         map[string]bool
	indexLoaded     bool
	cursor          int
	activeTab       tab
	width           int
	height          int
	err             error
	searchInput     textinput.Model
	searchFocus     bool
	spinner         spinner.Model
	spinning        bool
	showHelp        bool
	detailFocus     bool
	detailScroll    int
}

// messages
type (
	indexLoadedMsg     struct{ idx *registry.Index }
	indexErrMsg        struct{ err error }
	installedLoadedMsg struct{ skills []string }
	actionCompleteMsg  struct {
		name      string
		isInstall bool
		err       error
	}
)

func New(version string) Model {
	ti := textinput.New()
	ti.Placeholder = "search..."
	ti.Width = 30

	sp := spinner.New(
		spinner.WithSpinner(spinner.MiniDot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(highlight)),
	)

	return Model{
		version:     version,
		installed:   map[string]bool{},
		loading:     map[string]bool{},
		searchInput: ti,
		spinner:     sp,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick, loadIndex(), loadInstalled())
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

	case spinner.TickMsg:
		// Let the animation die out once nothing is in flight.
		if !m.anyLoading() {
			m.spinning = false
			break
		}
		var c tea.Cmd
		m.spinner, c = m.spinner.Update(msg)
		cmds = append(cmds, c)

	case indexLoadedMsg:
		m.indexLoaded = true
		m.skills = msg.idx.Skills
		m.bundles = msg.idx.Bundles
		m.filtered = msg.idx.Skills
		m.filteredBundles = msg.idx.Bundles
		registry.BuildIndex(m.skills, m.bundles)

	case indexErrMsg:
		m.err = msg.err

	case installedLoadedMsg:
		for _, s := range msg.skills {
			m.installed[s] = true
		}

	case actionCompleteMsg:
		delete(m.loading, msg.name)
		if msg.err == nil {
			m.installed[msg.name] = msg.isInstall
		} else {
			m.err = msg.err
		}
		m.clampCursor()

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
			m.detailScroll = 0

			return m, tea.Batch(cmds...)
		}

		if m.showHelp {
			switch msg.String() {
			case "?", "esc", "q":
				m.showHelp = false
			}
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "j", "down":
			if m.detailFocus {
				m.detailScroll++
			} else {
				m.cursor++
				m.detailScroll = 0
			}
		case "k", "up":
			if m.detailFocus {
				if m.detailScroll > 0 {
					m.detailScroll--
				}
			} else if m.cursor > 0 {
				m.cursor--
				m.detailScroll = 0
			}

		case "K", "home", "pgup":
			if m.detailFocus {
				m.detailScroll = 0
			} else {
				m.cursor = 0
				m.detailScroll = 0
			}

		case "J", "end", "pgdown":
			if m.detailFocus {
				m.detailScroll = 999999 // Let clamp/scroll handle max
			} else {
				m.cursor = m.listLen() - 1
				m.detailScroll = 0
			}

		case "1":
			m.switchTab(tabSkills)
		case "2":
			m.switchTab(tabBundles)
		case "3":
			m.switchTab(tabMySkills)
		case "4":
			m.switchTab(tabMyBundles)

		case "tab":
			m.detailFocus = !m.detailFocus

		case "?":
			m.showHelp = true

		case "/":
			m.searchFocus = true
			m.searchInput.Focus()

		case "i":
			return m, m.startAction(true, cmds)

		case "x":
			return m, m.startAction(false, cmds)
		}

		m.clampCursor()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) switchTab(t tab) {
	m.activeTab = t
	m.cursor = 0
	m.detailScroll = 0
}

func (m *Model) clampCursor() {
	if max := m.listLen() - 1; m.cursor > max {
		m.cursor = max
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// startAction queues installs or removals for the current selection. A bundle
// selection fans out to every skill it contains.
func (m *Model) startAction(install bool, cmds []tea.Cmd) tea.Cmd {
	var names []string
	var sources []string

	if m.showingBundles() {
		bundles := m.visibleBundles()
		if len(bundles) == 0 || m.cursor >= len(bundles) {
			return tea.Batch(cmds...)
		}
		b := bundles[m.cursor]
		for _, name := range b.Skills {
			names = append(names, name)
			sources = append(sources, registry.BundleSkillSource(&registry.BundleFile{
				Repo: b.Repo, Branch: b.Branch, SkillsPath: b.SkillsPath,
			}, name))
		}
	} else {
		skills := m.visibleSkills()
		if len(skills) == 0 || m.cursor >= len(skills) {
			return tea.Batch(cmds...)
		}
		names = append(names, skills[m.cursor].Name)
		sources = append(sources, skills[m.cursor].Source)
	}

	started := false
	for i, name := range names {
		if m.installed[name] == install || m.loading[name] {
			continue
		}
		m.loading[name] = true
		started = true
		if install {
			cmds = append(cmds, doInstall(name, sources[i]))
		} else {
			cmds = append(cmds, doRemove(name))
		}
	}

	if started && !m.spinning {
		m.spinning = true
		cmds = append(cmds, m.spinner.Tick)
	}

	return tea.Batch(cmds...)
}

func (m Model) anyLoading() bool {
	return !m.indexLoaded || len(m.loading) > 0
}

func (m Model) showingBundles() bool {
	return m.activeTab == tabBundles || m.activeTab == tabMyBundles
}

// visibleSkills is the search-filtered list, narrowed to installed skills on
// the "My Skills" tab.
func (m Model) visibleSkills() []registry.Skill {
	if m.activeTab != tabMySkills {
		return m.filtered
	}
	var out []registry.Skill
	for _, s := range m.filtered {
		if m.installed[s.Name] {
			out = append(out, s)
		}
	}
	return out
}

// visibleBundles narrows to bundles with at least one installed skill on the
// "My Bundles" tab, so partial installs stay visible.
func (m Model) visibleBundles() []registry.Bundle {
	if m.activeTab != tabMyBundles {
		return m.filteredBundles
	}
	var out []registry.Bundle
	for _, b := range m.filteredBundles {
		if m.installedIn(b) > 0 {
			out = append(out, b)
		}
	}
	return out
}

func (m Model) installedIn(b registry.Bundle) int {
	n := 0
	for _, name := range b.Skills {
		if m.installed[name] {
			n++
		}
	}
	return n
}

func (m Model) loadingIn(b registry.Bundle) bool {
	for _, name := range b.Skills {
		if m.loading[name] {
			return true
		}
	}
	return false
}

func (m Model) listLen() int {
	if m.showingBundles() {
		return len(m.visibleBundles())
	}
	return len(m.visibleSkills())
}

// installedCount counts registry skills only, so it matches what the
// "My Skills" tab can actually list.
func (m Model) installedCount() int {
	n := 0
	for _, s := range m.skills {
		if m.installed[s.Name] {
			n++
		}
	}
	return n
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %s\n\npress q to quit", m.err)
	}
	if m.width == 0 {
		return "loading..."
	}

	help := "  j/k move · K/J top/bottom · 1-4 tabs · / search · i install · x remove · ? help · q quit"
	if m.width < minWidth {
		help = "  j/k move · K/J top/bottom · 1-4 tabs · ? help · q quit"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		m.tabBarView(),
		m.bodyView(),
		m.statusBarView(),
		mutedStyle.MaxWidth(m.width).Render(help),
	)
}

func (m Model) tabBarView() string {
	labels := []string{
		"Skills",
		"Bundles",
		fmt.Sprintf("My Skills (%d)", m.installedCount()),
		fmt.Sprintf("My Bundles (%d)", len(m.installedBundles())),
	}

	var parts []string
	for i, label := range labels {
		if tab(i) == m.activeTab {
			parts = append(parts, activeTabStyle.Render(fmt.Sprintf(" %d %s ", i+1, label)))
			continue
		}
		// Narrow terminals only get room for the inactive tab numbers.
		if m.width < minWidth {
			parts = append(parts, inactiveTabStyle.Render(fmt.Sprintf(" %d ", i+1)))
		} else {
			parts = append(parts, inactiveTabStyle.Render(fmt.Sprintf(" %d %s ", i+1, label)))
		}
	}
	tabs := strings.Join(parts, " ")

	search := m.searchInput.View()
	padding := m.width - lipgloss.Width(tabs) - lipgloss.Width(search) - 2
	if padding < 1 {
		// No room for both — the tabs matter more.
		return lipgloss.NewStyle().MaxWidth(m.width).Render(tabs)
	}

	return tabs + strings.Repeat(" ", padding) + search
}

func (m Model) statusBarView() string {
	left := fmt.Sprintf("  %d installed · %d skills · %d bundles",
		m.installedCount(), len(m.skills), len(m.bundles))

	if m.anyLoading() {
		if !m.indexLoaded {
			left += fmt.Sprintf(" · %s fetching registry", m.spinner.View())
		} else {
			left += fmt.Sprintf(" · %s %d in progress", m.spinner.View(), len(m.loading))
		}
	}

	right := "skmp v" + m.version + "  "
	padding := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if padding < 1 {
		return mutedStyle.MaxWidth(m.width).Render(left)
	}

	return mutedStyle.Render(left) + strings.Repeat(" ", padding) + mutedStyle.Render(right)
}

func (m Model) installedBundles() []registry.Bundle {
	var out []registry.Bundle
	for _, b := range m.bundles {
		if m.installedIn(b) > 0 {
			out = append(out, b)
		}
	}
	return out
}

func (m Model) bodyView() string {
	tabH := lipgloss.Height(m.tabBarView())
	statusH := lipgloss.Height(m.statusBarView())
	helpStr := "  j/k move · K/J top/bottom · 1-4 tabs · / search · i install · x remove · ? help · q quit"
	if m.width < minWidth {
		helpStr = "  j/k move · K/J top/bottom · 1-4 tabs · ? help · q quit"
	}
	helpH := lipgloss.Height(mutedStyle.MaxWidth(m.width).Render(helpStr))
	outerH := m.height - tabH - statusH - helpH
	if outerH < 3 {
		outerH = 3
	}
	innerH := outerH - 2 // account for top and bottom borders

	if m.showHelp {
		return lipgloss.Place(m.width, outerH, lipgloss.Center, lipgloss.Center, helpView())
	}

	if !m.indexLoaded {
		loadingMsg := fmt.Sprintf("%s Fetching registry...", m.spinner.View())
		return lipgloss.Place(m.width, outerH, lipgloss.Center, lipgloss.Center, mutedStyle.Render(loadingMsg))
	}

	if m.width < minWidth {
		return m.paneList(m.width-2, outerH, innerH)
	}

	listW := m.width / 3
	detailW := m.width - listW - 4 // 2 borders per pane

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.paneList(listW, outerH, innerH),
		m.paneDetail(detailW, outerH, innerH),
	)
}

func (m Model) paneList(w, outerH, innerH int) string {
	content := m.listView(w, innerH)
	return m.paneStyle(!m.detailFocus).Width(w).Height(innerH).Render(content)
}

func (m Model) paneDetail(w, outerH, innerH int) string {
	var content string
	if m.showingBundles() {
		content = m.bundleDetailsView(w)
	} else {
		content = m.skillDetailView(w)
	}
	return m.paneStyle(m.detailFocus).Width(w).Height(innerH).Render(scroll(content, m.detailScroll, innerH))
}

func (m Model) paneStyle(focused bool) lipgloss.Style {
	if focused {
		return activeBorderStyle
	}
	return borderStyle
}

func (m Model) viewport(total, height int) (offset, end int) {
	if m.cursor >= height {
		offset = m.cursor - height + 1
	}
	end = offset + height
	if end > total {
		end = total
	}
	return
}

// scroll drops the first offset lines, keeping at least one screen of content.
func scroll(content string, offset, height int) string {
	lines := strings.Split(content, "\n")
	if max := len(lines) - height; offset > max {
		offset = max
	}
	if offset < 1 {
		return content
	}
	return strings.Join(lines[offset:], "\n")
}

func (m Model) listView(w, h int) string {
	if m.showingBundles() {
		return m.bundlesListView(w, h)
	}
	return m.skillsListView(w, h)
}

func (m Model) skillsListView(w, h int) string {
	skills := m.visibleSkills()
	if len(skills) == 0 {
		if m.activeTab == tabMySkills {
			return mutedStyle.Render("  no skills installed")
		}
		return mutedStyle.Render("  no skills found")
	}

	offset, end := m.viewport(len(skills), h)
	var s strings.Builder
	for i := offset; i < end; i++ {
		s.WriteString(m.row(i, m.skillBadge(skills[i].Name), skills[i].Name, "", w))
	}
	return strings.TrimRight(s.String(), "\n")
}

func (m Model) bundlesListView(w, h int) string {
	bundles := m.visibleBundles()
	if len(bundles) == 0 {
		if m.activeTab == tabMyBundles {
			return mutedStyle.Render("  no bundles installed")
		}
		return mutedStyle.Render("  no bundles found")
	}

	offset, end := m.viewport(len(bundles), h)
	var s strings.Builder
	for i := offset; i < end; i++ {
		count := fmt.Sprintf(" (%d/%d)", m.installedIn(bundles[i]), len(bundles[i].Skills))
		s.WriteString(m.row(i, m.bundleBadge(bundles[i]), bundles[i].Name, count, w))
	}
	return strings.TrimRight(s.String(), "\n")
}

func (m Model) row(i int, badge, name, suffix string, w int) string {
	cursor := "  "
	label := normalStyle.Render(name)
	if i == m.cursor {
		cursor = "▶ "
		label = selectedStyle.Render(name)
	}
	line := cursor + badge + label + mutedStyle.Render(suffix)
	return lipgloss.NewStyle().MaxWidth(w).Render(line) + "\n"
}

// skillBadge shows install state in the list: spinner while working, filled
// dot when installed, hollow dot otherwise.
func (m Model) skillBadge(name string) string {
	switch {
	case m.loading[name]:
		return m.spinner.View() + " "
	case m.installed[name]:
		return installedStyle.Render("● ")
	default:
		return mutedStyle.Render("○ ")
	}
}

func (m Model) bundleBadge(b registry.Bundle) string {
	installed := m.installedIn(b)
	switch {
	case m.loadingIn(b):
		return m.spinner.View() + " "
	case len(b.Skills) > 0 && installed == len(b.Skills):
		return installedStyle.Render("● ")
	case installed > 0:
		return partialStyle.Render("◐ ")
	default:
		return mutedStyle.Render("○ ")
	}
}

func (m Model) skillDetailView(w int) string {
	skills := m.visibleSkills()
	if len(skills) == 0 || m.cursor >= len(skills) {
		return mutedStyle.Render("select a skill")
	}

	s := skills[m.cursor]

	tags := ""
	for _, t := range s.Tags {
		tags += tagStyle.Render(t)
	}

	var status string
	switch {
	case m.loading[s.Name]:
		status = m.spinner.View() + mutedStyle.Render(" working...")
	case m.installed[s.Name]:
		status = installedStyle.Render("● installed")
	default:
		status = mutedStyle.Render("○ not installed")
	}

	return lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render(s.Name), "", wrapText(s.Description, w-4),
		"",
		labelStyle.Render("Version")+valueStyle.Render(s.Version),
		labelStyle.Render("Author")+valueStyle.Render(s.Author),
		labelStyle.Render("Harnesses")+valueStyle.Render(strings.Join(s.Harnesses, ",")),
		labelStyle.Render("Tags")+tags,
		"",
		status,
	)
}

func (m Model) bundleDetailsView(w int) string {
	bundles := m.visibleBundles()
	if len(bundles) == 0 || m.cursor >= len(bundles) {
		return mutedStyle.Render("select a bundle")
	}

	b := bundles[m.cursor]
	installedCount := m.installedIn(b)
	total := len(b.Skills)

	var status string
	switch {
	case total > 0 && installedCount == total:
		status = installedStyle.Render("● installed")
	case installedCount > 0:
		status = partialStyle.Render(fmt.Sprintf("◐ partial (%d/%d)", installedCount, total))
	default:
		status = mutedStyle.Render("○ not installed")
	}

	skillList := ""
	for _, name := range b.Skills {
		skillList += m.skillBadge(name) + name + "\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(b.Name),
		"",
		wrapText(b.Description, w-4),
		"",
		labelStyle.Render("Author")+valueStyle.Render(b.Author),
		labelStyle.Render("Repo")+valueStyle.Render(b.Repo),
		labelStyle.Render("Skills")+valueStyle.Render(fmt.Sprintf("%d", total)),
		"",
		status,
		"",
		skillList,
	)
}

func helpView() string {
	bindings := [][2]string{
		{"j / ↓", "move down"},
		{"k / ↑", "move up"},
		{"1", "all skills"},
		{"2", "all bundles"},
		{"3", "installed skills"},
		{"4", "installed bundles"},
		{"tab", "switch list / detail focus"},
		{"/", "search (esc to leave)"},
		{"i", "install selection"},
		{"x", "remove selection"},
		{"?", "close this help"},
		{"q", "quit"},
	}

	var rows []string
	rows = append(rows, titleStyle.Render("Keybindings"), "")
	for _, b := range bindings {
		rows = append(rows, labelStyle.Render(b[0])+valueStyle.Render(b[1]))
	}

	return activeBorderStyle.Padding(1, 3).Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
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

func loadInstalled() tea.Cmd {
	return func() tea.Msg {
		skills, _ := harness.InstalledSkills()
		return installedLoadedMsg{skills}
	}
}

func doInstall(name, source string) tea.Cmd {
	return func() tea.Msg {
		err := harness.Install(name, source)
		return actionCompleteMsg{name: name, isInstall: true, err: err}
	}
}

func doRemove(name string) tea.Cmd {
	return func() tea.Msg {
		err := harness.Remove(name)
		return actionCompleteMsg{name: name, isInstall: false, err: err}
	}
}
