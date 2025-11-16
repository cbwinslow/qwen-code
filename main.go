package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

// ==================== MODEL ====================

// Secret represents a stored secret
type Secret struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []string  `json:"tags"`
}

// Pane represents a resizable pane
type Pane struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	X           int     `json:"x"`
	Y           int     `json:"y"`
	IsActive    bool    `json:"is_active"`
	IsResizable bool    `json:"is_resizable"`
}

// Progress represents a progress bar
type Progress struct {
	ID          string  `json:"id"`
	Label       string  `json:"label"`
	Percent     float64 `json:"percent"`
	IsActive    bool    `json:"is_active"`
	ShowPercent bool    `json:"show_percent"`
	Color       string  `json:"color"`
}

// Model holds the application state
type Model struct {
	panes       []Pane
	secrets     []Secret
	progress    []Progress
	activePane  int
	focusedPane int
	width       int
	height      int
	loading     bool
	loadingText string
	editingSecret *Secret
	newSecretName string
	newSecretValue string
	showSecrets bool
	showProgress bool
	draggingPane *Pane
	dragStartX   int
	dragStartY   int
}

// ==================== INITIALIZE ====================

func initialModel() Model {
	// Initialize with some default panes
	panes := []Pane{
		{
			ID:          "main",
			Title:       "Main Workspace",
			Content:     "Welcome to your TUI workspace!\n\nThis is a resizable pane system with progress tracking and secret management.\n\n• Drag pane borders to resize\n• Click panes to focus\n• Use Tab to cycle focus",
			Width:       60,
			Height:      20,
			X:           2,
			Y:           2,
			IsActive:    true,
			IsResizable: true,
		},
		{
			ID:          "logs",
			Title:       "Activity Log",
			Content:     "System logs will appear here...\n\n[INFO] Application started\n[INFO] Panes initialized\n[INFO] Ready for user input",
			Width:       40,
			Height:      15,
			X:           65,
			Y:           2,
			IsActive:    false,
			IsResizable: true,
		},
		{
			ID:          "secrets",
			Title:       "Secret Manager",
			Content:     "Click 's' to view and manage secrets",
			Width:       35,
			Height:      12,
			X:           2,
			Y:           23,
			IsActive:    false,
			IsResizable: true,
		},
	}

	// Initialize with sample progress bars
	progress := []Progress{
		{ID: "task1", Label: "Database Migration", Percent: 0.75, IsActive: true, ShowPercent: true, Color: "blue"},
		{ID: "task2", Label: "File Processing", Percent: 0.45, IsActive: true, ShowPercent: true, Color: "green"},
		{ID: "task3", Label: "API Sync", Percent: 0.90, IsActive: true, ShowPercent: true, Color: "yellow"},
	}

	return Model{
		panes:       panes,
		secrets:     loadSecrets(),
		progress:    progress,
		activePane:  0,
		focusedPane: 0,
		width:       100,
		height:      40,
		loading:     false,
		showSecrets: false,
		showProgress: true,
	}
}

// ==================== STYLES ====================

var (
	// Base styles
	baseStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAFAFA"))
	paneStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	activeStyle   = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Padding(0, 1).Foreground(lipgloss.Color("#86E1FC"))
	inactiveStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).Padding(0, 1).Foreground(lipgloss.Color("#6C7086"))
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF79C6"))
	progressStyle = lipgloss.NewStyle().Bold(true)
	secretStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	warningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	successStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
)

// ==================== UPDATE ====================

func (m Model) Init() bubbletea.Cmd {
	return bubbletea.Batch(
		bubbletea.WindowSize(),
		m.startProgressUpdates(),
	)
}

func (m Model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.WindowSizeMsg:
		m.width, m.height = int(msg.Width), int(msg.Height)
		return m, nil

	case bubbletea.KeyMsg:
		return m.handleKey(msg)

	case bubbletea.MouseMsg:
		return m.handleMouse(msg)

	case progressUpdateMsg:
		return m.updateProgress(msg)

	case bubbletea.TickMsg:
		return m, m.startProgressUpdates

	default:
		return m, nil
	}
}

// ==================== KEY HANDLING ====================

func (m Model) handleKey(msg bubbletea.KeyMsg) (Model, bubbletea.Cmd) {
	switch msg.Type {
	case bubbletea.KeyCtrlC, bubbletea.KeyEsc:
		return m, bubbletea.Quit

	case bubbletea.KeyTab:
		m.focusedPane = (m.focusedPane + 1) % len(m.panes)
		m.activePane = m.focusedPane
		return m, nil

	case bubbletea.KeyShiftTab:
		m.focusedPane = (m.focusedPane - 1 + len(m.panes)) % len(m.panes)
		m.activePane = m.focusedPane
		return m, nil

	case bubbletea.KeyEnter:
		if m.showSecrets && m.editingSecret != nil {
			m.saveSecret()
		}
		return m, nil

	case bubbletea.KeyBackspace:
		if m.showSecrets && len(m.newSecretValue) > 0 {
			m.newSecretValue = m.newSecretValue[:len(m.newSecretValue)-1]
		}
		return m, nil

	case bubbletea.KeyRunes:
		if m.showSecrets {
			m.newSecretValue += string(msg.Runes)
		}
		return m, nil

	case 's':
		m.showSecrets = !m.showSecrets
		m.showProgress = !m.showProgress
		return m, nil

	case 'p':
		m.showProgress = !m.showProgress
		m.showSecrets = !m.showSecrets
		return m, nil

	case 'n':
		if m.showSecrets {
			m.editingSecret = &Secret{ID: generateID(), CreatedAt: time.Now()}
			m.newSecretName = ""
			m.newSecretValue = ""
		}
		return m, nil

	case 'q':
		if m.showSecrets && m.editingSecret != nil {
			m.editingSecret = nil
			m.newSecretName = ""
			m.newSecretValue = ""
		}
		return m, nil

	case 'r':
		if m.showSecrets {
			m.refreshSecrets()
		}
		return m, nil
	}

	return m, nil
}

// ==================== MOUSE HANDLING ====================

func (m Model) handleMouse(msg bubbletea.MouseMsg) (Model, bubbletea.Cmd) {
	if msg.Type == bubbletea.MouseLeft {
		x, y := msg.X, msg.Y
		
		// Check if clicking on pane border for resizing
		for i := range m.panes {
			pane := &m.panes[i]
			if m.isOnPaneBorder(x, y, pane) {
				m.draggingPane = pane
				m.dragStartX = x
				m.dragStartY = y
				return m, nil
			}
			
			// Check if clicking inside pane
			if m.isInsidePane(x, y, pane) {
				m.focusedPane = i
				m.activePane = i
				return m, nil
			}
		}
	} else if msg.Type == bubbletea.MouseMotion && m.draggingPane != nil {
		// Handle pane resizing
		dx := msg.X - m.dragStartX
		dy := msg.Y - m.dragStartY
		
		m.draggingPane.Width = max(20, m.draggingPane.Width+dx)
		m.draggingPane.Height = max(5, m.draggingPane.Height+dy)
		
		m.dragStartX = msg.X
		m.dragStartY = msg.Y
		return m, nil
	} else if msg.Type == bubbletea.MouseRelease {
		m.draggingPane = nil
		return m, nil
	}

	return m, nil
}

// ==================== PROGRESS SYSTEM ====================

type progressUpdateMsg struct {
	ID      string
	Percent float64
}

func (m Model) startProgressUpdates() bubbletea.Cmd {
	return bubbletea.Tick(time.Second, func(t time.Time) bubbletea.Msg {
		updates := []progressUpdateMsg{
			{ID: "task1", Percent: min(1.0, m.progress[0].Percent+0.05)},
			{ID: "task2", Percent: min(1.0, m.progress[1].Percent+0.03)},
			{ID: "task3", Percent: min(1.0, m.progress[2].Percent+0.02)},
		}
		return updates
	})
}

func (m Model) updateProgress(msg progressUpdateMsg) (Model, bubbletea.Cmd) {
	for i := range m.progress {
		if m.progress[i].ID == msg.ID {
			m.progress[i].Percent = msg.Percent
			if msg.Percent >= 1.0 {
				m.progress[i].IsActive = false
			}
			break
		}
	}
	return m, nil
}

// ==================== SECRET MANAGEMENT ====================

func loadSecrets() []Secret {
	home, _ := os.UserHomeDir()
	secretFile := filepath.Join(home, ".go-tui-secrets.json")
	
	data, err := os.ReadFile(secretFile)
	if err != nil {
		return []Secret{}
	}
	
	var secrets []Secret
	if err := json.Unmarshal(data, &secrets); err != nil {
		return []Secret{}
	}
	
	return secrets
}

func saveSecrets(secrets []Secret) error {
	home, _ := os.UserHomeDir()
	secretFile := filepath.Join(home, ".go-tui-secrets.json")
	
	data, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(secretFile, data, 0600)
}

func (m *Model) saveSecret() {
	if m.editingSecret == nil || m.newSecretName == "" || m.newSecretValue == "" {
		return
	}
	
	m.editingSecret.Name = m.newSecretName
	m.editingSecret.Value = m.newSecretValue
	m.editingSecret.UpdatedAt = time.Now()
	
	m.secrets = append(m.secrets, *m.editingSecret)
	saveSecrets(m.secrets)
	
	m.editingSecret = nil
	m.newSecretName = ""
	m.newSecretValue = ""
}

func (m *Model) refreshSecrets() {
	m.secrets = loadSecrets()
}

// ==================== RENDER ====================

func (m Model) View() string {
	if m.width < 80 || m.height < 24 {
		return m.renderTooSmall()
	}
	
	var content strings.Builder
	
	// Header
	content.WriteString(m.renderHeader())
	
	// Main content area
	if m.showSecrets {
		content.WriteString(m.renderSecrets())
	} else if m.showProgress {
		content.WriteString(m.renderProgress())
	} else {
		content.WriteString(m.renderPanes())
	}
	
	// Footer
	content.WriteString(m.renderFooter())
	
	return content.String()
}

func (m Model) renderHeader() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#86E1FC")).
		Background(lipgloss.Color("#1E1E2E")).
		Padding(0, 2).
		Render("🚀 Multi-Pane TUI System")
	
	return lipgloss.PlaceHorizontal(m.width, lipgloss.Left, title)
}

func (m Model) renderPanes() string {
	var content strings.Builder
	
	// Render each pane
	for i, pane := range m.panes {
		style := inactiveStyle
		if i == m.activePane {
			style = activeStyle
		}
		
		// Create pane content
		paneContent := fmt.Sprintf("%s\n\n%s", 
			titleStyle.Render(pane.Title),
			pane.Content)
		
		// Truncate content to fit pane
		lines := strings.Split(paneContent, "\n")
		maxLines := pane.Height - 2 // Account for border
		if len(lines) > maxLines {
			lines = lines[:maxLines]
		}
		paneContent = strings.Join(lines, "\n")
		
		// Render pane with border
		renderedPane := style.Width(pane.Width).Height(pane.Height).Render(paneContent)
		
		// Position pane
		content.WriteString(lipgloss.Place(m.height, m.width, 
			lipgloss.Left, lipgloss.Top,
			renderedPane,
			lipgloss.WithWhitespaceChars(" ", " "),
			lipgloss.WithWhitespacePane(lipgloss.NewStyle().Background(lipgloss.Color("#1E1E2E")),
		))
	}
	
	return content.String()
}

func (m Model) renderProgress() string {
	var content strings.Builder
	
	title := titleStyle.Render("📊 Progress Tracking")
	content.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title))
	content.WriteString("\n\n")
	
	// Create progress table
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NormalBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return titleStyle
			}
			return baseStyle
		}).
		Headers("TASK", "PROGRESS", "STATUS").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return titleStyle
			}
			return baseStyle
		})
	
	for _, p := range m.progress {
		status := "🔄 Active"
		if !p.IsActive {
			status = "✅ Complete"
		}
		
		progressBar := m.renderProgressBar(p.Percent, p.Color)
		t.Row(p.Label, progressBar, status)
	}
	
	content.WriteString(t.String())
	return content.String()
}

func (m Model) renderSecrets() string {
	var content strings.Builder
	
	title := titleStyle.Render("🔐 Secret Manager")
	content.WriteString(lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title))
	content.WriteString("\n\n")
	
	if m.editingSecret != nil {
		// Editing form
		content.WriteString(secretStyle.Render("📝 Add New Secret\n\n"))
		content.WriteString(fmt.Sprintf("Name: %s\n", m.newSecretName))
		content.WriteString(fmt.Sprintf("Value: %s\n", m.newSecretValue))
		content.WriteString("\n")
		content.WriteString(baseStyle.Render("Commands: [Enter] Save [q] Cancel"))
	} else {
		// Secrets list
		if len(m.secrets) == 0 {
			content.WriteString(warningStyle.Render("No secrets stored yet."))
		} else {
			t := table.New().
				Border(lipgloss.RoundedBorder()).
				StyleFunc(func(row, col int) lipgloss.Style {
					if row == 0 {
						return titleStyle
					}
					return secretStyle
				}).
				Headers("NAME", "CREATED", "TAGS")
			
			for _, secret := range m.secrets {
				tags := strings.Join(secret.Tags, ", ")
				if tags == "" {
					tags = "-"
				}
				t.Row(secret.Name, secret.CreatedAt.Format("2006-01-02"), tags)
			}
			
			content.WriteString(t.String())
		}
		
		content.WriteString("\n\n")
		content.WriteString(baseStyle.Render("Commands: [n] New [r] Refresh [s] Switch to Panes"))
	}
	
	return content.String()
}

func (m Model) renderFooter() string {
	commands := "[Tab] Switch Pane | [s] Secrets | [p] Progress | [Ctrl+C] Quit"
	if m.showSecrets {
		commands = "[n] New | [r] Refresh | [q] Cancel | [s] Switch"
	}
	
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6C7086")).
		Background(lipgloss.Color("#1E1E2E")).
		Padding(0, 1)
	
	return lipgloss.Place(m.height, m.width,
		lipgloss.Left, lipgloss.Bottom,
		footerStyle.Render(commands),
	)
}

func (m Model) renderTooSmall() string {
	return errorStyle.Render("Terminal too small! Please resize to at least 80x24")
}

func (m Model) renderProgressBar(percent float64, color string) string {
	width := 20
	filled := int(percent * float64(width))
	empty := width - filled
	
	barColor := lipgloss.Color("#50FA7B")
	switch color {
	case "blue":
		barColor = lipgloss.Color("#86E1FC")
	case "green":
		barColor = lipgloss.Color("#50FA7B")
	case "yellow":
		barColor = lipgloss.Color("#FFB86C")
	}
	
	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", empty)
	
	percentText := fmt.Sprintf("%.0f%%", percent*100)
	
	return lipgloss.NewStyle().
		Foreground(barColor).
		Render(fmt.Sprintf("[%s%s] %s", filledBar, emptyBar, percentText))
}

// ==================== HELPERS ====================

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func (m Model) isInsidePane(x, y int, pane *Pane) bool {
	return x >= pane.X && x < pane.X+pane.Width &&
		y >= pane.Y && y < pane.Y+pane.Height
}

func (m Model) isOnPaneBorder(x, y int, pane *Pane) bool {
	// Check if on right or bottom border (for resizing)
	onRightBorder := x == pane.X+pane.Width-1 && y >= pane.Y && y < pane.Y+pane.Height
	onBottomBorder := y == pane.Y+pane.Height-1 && x >= pane.X && x < pane.X+pane.Width
	
	return onRightBorder || onBottomBorder
}

// ==================== MAIN ====================

func main() {
	// Check terminal capabilities
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println("This application must be run in a terminal")
		os.Exit(1)
	}
	
	// Create and start the application
	p := bubbletea.NewProgram(
		initialModel(),
		bubbletea.WithAltScreen(),
		bubbletea.WithMouseCellMotion(),
		bubbletea.WithMouseAllMotion(),
	)
	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error starting application: %v", err)
		os.Exit(1)
	}
}