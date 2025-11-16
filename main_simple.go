package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	panes      []Pane
	activePane int
	width      int
	height     int
}

type Pane struct {
	ID       string
	Title    string
	Content  string
	Width    int
	Height   int
	X        int
	Y        int
	IsActive bool
}

func initialModel() Model {
	panes := []Pane{
		{
			ID:       "main",
			Title:    "Main Workspace",
			Content:  "Welcome to the TUI system!\n\nFeatures:\n• Resizable panes\n• Progress tracking\n• Secret management\n\nPress 'q' to quit",
			Width:    50,
			Height:   15,
			X:        2,
			Y:        2,
			IsActive: true,
		},
		{
			ID:       "progress",
			Title:    "Progress",
			Content:  "📊 Progress Tracking\n\nTask 1: 75%\nTask 2: 45%\nTask 3: 90%",
			Width:    40,
			Height:   10,
			X:        55,
			Y:        2,
			IsActive: false,
		},
		{
			ID:       "secrets",
			Title:    "Secrets",
			Content:  "🔐 Secret Manager\n\nStored secrets:\n• API Keys\n• Database credentials\n• Configuration tokens",
			Width:    35,
			Height:   8,
			X:        2,
			Y:        18,
			IsActive: false,
		},
	}

	return Model{
		panes:      panes,
		activePane: 0,
		width:      100,
		height:     30,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = int(msg.Width), int(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyTab:
			m.activePane = (m.activePane + 1) % len(m.panes)
			for i := range m.panes {
				m.panes[i].IsActive = (i == m.activePane)
			}
			return m, nil

		case tea.KeyRunes:
			if len(msg.Runes) > 0 && msg.Runes[0] == 's' {
				return m, tea.Printf("Secret management would open here")
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width < 80 || m.height < 24 {
		return "Terminal too small! Please resize to at least 80x24"
	}

	var content strings.Builder

	// Header
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#86E1FC")).
		Background(lipgloss.Color("#1E1E2E")).
		Padding(0, 2).
		Render("🚀 Multi-Pane TUI System")

	content.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Left, title))
	content.WriteString("\n\n")

	// Render panes
	for _, pane := range m.panes {
		style := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		if pane.IsActive {
			style = style.Border(lipgloss.ThickBorder()).
				Foreground(lipgloss.Color("#86E1FC"))
		}

		// Truncate content to fit pane
		lines := strings.Split(pane.Content, "\n")
		maxLines := pane.Height - 2 // Account for border
		if len(lines) > maxLines {
			lines = lines[:maxLines]
		}
		paneContent := strings.Join(lines, "\n")

		renderedPane := style.Width(pane.Width).Height(pane.Height).Render(paneContent)

		content.WriteString(lipgloss.Place(m.height, m.width,
			lipgloss.Left, lipgloss.Top,
			renderedPane,
		))
	}

	// Footer
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6C7086")).
		Background(lipgloss.Color("#1E1E2E")).
		Padding(0, 1).
		Render("[Tab] Switch Pane | [s] Secrets | [q] Quit")

	content.WriteString(lipgloss.Place(m.height, m.width,
		lipgloss.Left, lipgloss.Bottom,
		footer,
	))

	return content.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
