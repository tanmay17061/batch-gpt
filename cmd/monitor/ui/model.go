package ui

import (
	"batch-gpt/server/db" // for getting batch data

	"fmt"
	"time"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sashabaranov/go-openai"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Tab      key.Binding
	Help     key.Binding
	Quit     key.Binding
	Refresh  key.Binding
	Enter    key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	GotoTop  key.Binding
	GotoEnd  key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "previous tab"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next tab"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup/g", "go to top"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdown/G", "go to bottom"),
	),
	GotoTop: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("", ""), // Hide from help as it's shown with pgup
	),
	GotoEnd: key.NewBinding(
		key.WithKeys("G"),
		key.WithHelp("", ""), // Hide from help as it's shown with pgdown
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch tab"),
	),
	Help: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "view details"),
	),
}

type tab int

const (
	activeTab tab = iota
	completedTab
	failedTab
)

type Model struct {
	tabs          []string
	currentTab    tab
	batches       []batchItem
	cursor        int // Current selected batch
	offset        int // Scroll offset for viewing batches
	width, height int
	help          bool
	loading       bool
	error         error
	lastUpdate    time.Time
}

func NewModel() Model {
	return Model{
		tabs:       []string{"Active Batches", "Completed Batches", "Failed Batches"},
		currentTab: activeTab,
		help:       true,
		lastUpdate: time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return fetchBatches
}

func fetchBatches() tea.Msg {
	batches, err := db.GetAllBatchStatuses()
	if err != nil {
		return errMsg{err}
	}
	return batchesMsg{batches}
}

type batchesMsg struct {
	batches []openai.BatchResponse
}

type errMsg struct {
	err error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Tab), key.Matches(msg, keys.Right):
			m.currentTab = (m.currentTab + 1) % 3
			m.cursor = 0
			m.offset = 0
		case key.Matches(msg, keys.Left):
			m.currentTab = (m.currentTab - 1 + 3) % 3
			m.cursor = 0
			m.offset = 0
		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
				// Adjust offset if cursor moves above visible area
				if m.cursor < m.offset {
					m.offset = m.cursor
				}
			}
		case key.Matches(msg, keys.Down):
			filteredBatches := m.filterBatches()
			if m.cursor < len(filteredBatches)-1 {
				m.cursor++
				// Adjust offset if cursor moves below visible area
				if m.cursor-m.offset >= m.maxVisibleBatches() {
					m.offset = m.cursor - m.maxVisibleBatches() + 1
				}
			}
		case key.Matches(msg, keys.PageUp), key.Matches(msg, keys.GotoTop):
			m.cursor = 0
			m.offset = 0
		case key.Matches(msg, keys.PageDown), key.Matches(msg, keys.GotoEnd):
			filteredBatches := m.filterBatches()
			m.cursor = len(filteredBatches) - 1
			// Adjust offset to show the last page of batches
			if m.cursor >= m.maxVisibleBatches() {
				m.offset = m.cursor - m.maxVisibleBatches() + 1
			}
		case key.Matches(msg, keys.Help):
			m.help = !m.help
		case key.Matches(msg, keys.Refresh):
			m.loading = true
			return m, fetchBatches
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case batchesMsg:
		m.batches = processBatches(msg.batches)
		m.loading = false
		m.lastUpdate = time.Now()

	case errMsg:
		m.error = msg.err
		m.loading = false
	}

	return m, nil
}

func (m Model) maxVisibleBatches() int {
	// Subtract space for header (1), tabs (3), footer (1), and borders/padding (4)
	availableHeight := m.height - 9
	// Each batch takes 3 lines (id/timestamp line, progress bar, spacing)
	batchHeight := 3
	return max(1, availableHeight/batchHeight)
}

func (m Model) visibleBatches() []batchItem {
	filteredBatches := m.filterBatches()
	maxVisible := m.maxVisibleBatches()

	if len(filteredBatches) <= maxVisible {
		return filteredBatches
	}

	end := min(m.offset+maxVisible, len(filteredBatches))
	return filteredBatches[m.offset:end]
}

func (m Model) View() string {
    headerText := "🖥️  Batch-GPT Monitor"
    timeWithZone := time.Now().Format("15:04:05 MST")
    
    header := lipgloss.NewStyle().
        Width(m.width).
        Align(lipgloss.Center).
        Render(
            titleStyle.Render(headerText) + "\n" + 
            helpStyle.Render(timeWithZone),
        )

	// Build tabs with counts
	var tabs []string
	filteredBatches := m.filterBatches()
	totalBatches := len(filteredBatches)
	for i, t := range m.tabs {
		style := tabStyle
		if tab(i) == m.currentTab {
			style = activeTabStyle
			t = fmt.Sprintf("%s (%d)", t, totalBatches)
		}
		tabs = append(tabs, style.Render(t))
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)

	// Build batch list with scroll indicators
	var content string
	if totalBatches > 0 {
		visibleBatches := m.visibleBatches()
		content = m.renderBatches(visibleBatches)

		// Add scroll indicators if needed
		if m.offset > 0 {
			content = "↑ More batches above\n" + content
		}
		if m.offset+len(visibleBatches) < totalBatches {
			content = content + "\n↓ More batches below"
		}

		// Add batch counter
		content += fmt.Sprintf("\nShowing %d-%d of %d batches",
			m.offset+1,
			m.offset+len(visibleBatches),
			totalBatches)
	} else {
		content = "No batches found"
	}

	// Build footer with navigation help
	footer := helpStyle.Render("←→: switch tabs • ↑↓: navigate • pgup/g: top • pgdown/G: bottom • r: refresh • enter: details • q: quit")
	if m.loading {
		footer = "Loading..."
	}

	// Join all sections
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(primaryColor).
			Render(tabRow),
		content,
		footer,
	)
}
