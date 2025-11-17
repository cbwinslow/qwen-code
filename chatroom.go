package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== CHATROOM CORE TYPES ====================

// MessageType represents different types of messages
type MessageType string

const (
	MessageTypeUser       MessageType = "user"
	MessageTypeAgent     MessageType = "agent"
	MessageTypeSystem    MessageType = "system"
	MessageTypeTool     MessageType = "tool"
	MessageTypeFile     MessageType = "file"
)

// Message represents a chat message
type Message struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	AgentID   string                 `json:"agent_id,omitempty"`
	Content   string                 `json:"content"`
	Type      MessageType            `json:"type"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	ParentID  string                 `json:"parent_id,omitempty"`
	ThreadID  string                 `json:"thread_id,omitempty"`
}

// Agent represents an AI agent
type Agent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Role        string    `json:"role"`
	Personality string    `json:"personality"`
	Provider    string    `json:"provider"`
	Model       string    `json:"model"`
	Status      string    `json:"status"`
	IsActive    bool      `json:"is_active"`
	LastSeen    time.Time `json:"last_seen"`
	Messages    int       `json:"messages"`
	Avatar      string    `json:"avatar"`
}

// ConversationType represents different conversation modes
type ConversationType string

const (
	ConversationDemocratic   ConversationType = "democratic"
	ConversationEnsemble     ConversationType = "ensemble"
	ConversationHierarchical ConversationType = "hierarchical"
	ConversationCompetitive   ConversationType = "competitive"
	ConversationSpecialist   ConversationType = "specialist"
	ConversationConsensus    ConversationType = "consensus"
)

// Conversation represents a chat conversation
type Conversation struct {
	ID           string         `json:"id"`
	Type         ConversationType `json:"type"`
	Participants []string       `json:"participants"`
	Messages     []Message       `json:"messages"`
	IsActive     bool           `json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	Subject      string         `json:"subject,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
}

// Provider represents an AI provider
type Provider struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	APIKey      string   `json:"api_key,omitempty"`
	BaseURL     string   `json:"base_url,omitempty"`
	Models       []string `json:"models"`
	IsActive    bool     `json:"is_active"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

// ==================== CHATROOM MODEL ====================

// ChatroomModel represents the chatroom application state
type ChatroomModel struct {
	width           int
	height          int
	chatVisible     bool
	inputText       string
	messages        []Message
	agents          []Agent
	conversations    []Conversation
	activeConvID    string
	selectedAgentIDs []string
	providers       []Provider
	activeProvider  string
	scrollOffset     int
	agentStatus     map[string]string
	fileManager     *FileManager
	settings        *Settings
}

// FileManager handles file operations
type FileManager struct {
	sharedFiles []SharedFile
	uploadDir   string
}

// SharedFile represents a file shared in chatroom
type SharedFile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	Type      string    `json:"type"`
	Owner     string    `json:"owner"`
	SharedAt  time.Time `json:"shared_at"`
}

// Settings represents application settings
type Settings struct {
	theme           string
	animationSpeed  float64
	autoSave        bool
	encryptKeys     bool
	maxTokens       int
	temperature     float64
}

// ==================== CHATROOM IMPLEMENTATION ====================

// NewChatroomModel creates a new chatroom model
func NewChatroomModel() ChatroomModel {
	return ChatroomModel{
		chatVisible:     false,
		messages:        []Message{},
		agents:          []Agent{},
		conversations:    []Conversation{},
		activeConvID:    "",
		selectedAgentIDs: []string{},
		providers: []Provider{
			{
				ID:   "openrouter",
				Name: "OpenRouter",
				Models: []string{"qwen/qwen-2.5-coder-32b", "anthropic/claude-3-sonnet-20240229", "meta-llama/llama-3-70b-instruct"},
			},
			{
				ID:   "ollama",
				Name: "Ollama (Local)",
				Models: []string{"llama2", "codellama", "mistral", "vicuna"},
			},
			{
				ID:   "qwen-coder",
				Name: "Qwen Coder",
				Models: []string{"qwen-coder-2.5", "qwen-coder-1.5"},
			},
		},
		activeProvider:  "openrouter",
		scrollOffset:    0,
		agentStatus:     make(map[string]string),
		fileManager:     &FileManager{
			sharedFiles: []SharedFile{},
			uploadDir:   filepath.Join(os.TempDir(), "chatroom-uploads"),
		},
		settings: &Settings{
			theme:          "ocean",
			animationSpeed: 1.0,
			autoSave:       true,
			encryptKeys:    true,
			maxTokens:      4096,
			temperature:    0.7,
		},
	}
}

// ==================== UI COMPONENTS ====================

// Styles for the chatroom interface
var (
	chatroomStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#0a0e27")).
			Foreground(lipgloss.Color("#ffffff"))

	agentStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1e3a8a")).
			Foreground(lipgloss.Color("#ffffff"))

	userStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2d3748")).
			Foreground(lipgloss.Color("#ffffff"))

	systemStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#6b7280")).
			Foreground(lipgloss.Color("#ffffff"))

	inputStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1a1a2e")).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#4a90e2")).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 1).
			Bold(true)

	activeButtonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#61dafb")).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4a90e2"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#61dafb")).
			Bold(true).
			Padding(0, 1)
)

	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
)

	accentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f39c12")).
			Bold(true)
)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff4444")).
			Bold(true)
)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffa500")).
			Bold(true)
)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#17a2b8")).
			Bold(true)
)

	// Agent status colors
	agentStatusColors = map[string]lipgloss.Style{
		"active":    successStyle,
		"idle":     subtleStyle,
		"busy":     warningStyle,
		"error":     errorStyle,
		"offline":   errorStyle,
	}
)

// ==================== CHATROOM MODEL METHODS ====================

// Init initializes the chatroom model
func (m ChatroomModel) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("AI TUI - Multi-Agent Chatroom"),
		tea.EnterAltScreen,
		tea.EnableMouseCellMotion,
	)
}

// Update handles messages and updates
func (m ChatroomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case time.Time:
		// Update agent statuses
		for i := range m.agents {
			if m.agents[i].IsActive {
				m.agents[i].LastSeen = time.Now()
			}
		}
		return m, nil

	default:
		return m, nil
	}
}

// View renders the chatroom interface
func (m ChatroomModel) View() string {
	if !m.chatVisible {
		return m.renderSplash()
	}

	return m.renderChatroom()
}

// ==================== KEY HANDLING ====================

// handleKey processes keyboard input
func (m ChatroomModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		return m, tea.Quit

	case tea.KeyCtrlT:
		m.chatVisible = !m.chatVisible
		return m, nil

	case tea.KeyCtrlN:
		return m.createNewConversation()

	case tea.KeyEnter:
		if m.inputText != "" {
			return m.sendMessage()
		}

	case tea.KeyBackspace:
		if len(m.inputText) > 0 {
			m.inputText = m.inputText[:len(m.inputText)-1]
		}

	case tea.KeyUp, tea.KeyDown:
		return m.navigateHistory(msg)

	case tea.KeyPgUp, tea.KeyPgDown:
		return m.navigateMessages(msg)

	case tea.KeyTab:
		return m.cycleFocus()

	default:
		if msg.Type == tea.KeyRunes {
			m.inputText += string(msg.Runes)
		}
		return m, nil
	}
}

// ==================== MOUSE HANDLING ====================

// handleMouse processes mouse input
func (m ChatroomModel) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Handle mouse clicks on agent panels, messages, etc.
	return m, nil
}

// ==================== RENDERING METHODS ====================

// renderSplash displays the splash screen
func (m ChatroomModel) renderSplash() string {
	splash := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#61dafb")).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height).
		Render("🤖 Multi-Agent Chatroom\n\nPress Ctrl+T to begin")

	return splash
}

// renderChatroom renders the main chatroom interface
func (m ChatroomModel) renderChatroom() string {
	// Create layout
	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		MaxWidth(m.width).
		MaxHeight(m.height).
		Render(m.renderMainContent())

	return content
}

// renderMainContent renders the main chatroom content
func (m ChatroomModel) renderMainContent() string {
	// Split screen into sections
	leftWidth := m.width / 3
	centerWidth := m.width / 3
	rightWidth := m.width - leftWidth - centerWidth

	leftPanel := m.renderAgentPanel(leftWidth)
	centerPanel := m.renderConversationPanel(centerWidth)
	rightPanel := m.renderControlPanel(rightWidth)

	return lipgloss.JoinHorizontal(leftPanel, centerPanel, rightPanel)
}

// renderAgentPanel renders the agent management panel
func (m ChatroomModel) renderAgentPanel(width int) string {
	title := titleStyle.Render("🤖 Agents")
	
	agentsList := strings.Builder{}
	for _, agent := range m.agents {
		statusColor := agentStatusColors[agent.Status]
		status := statusColor.Render(fmt.Sprintf("[%s]", agent.Status))
		
		agentInfo := fmt.Sprintf("%s %s %s", 
			agent.Avatar, agent.Name, status)
		
		if contains(m.selectedAgentIDs, agent.ID) {
			agentInfo = accentStyle.Render(agentInfo)
		}
		
		agentsList.WriteString(lipgloss.NewStyle().
			Width(width-4).
			Render(agentInfo))
		agentsList.WriteString("\n")
	}

	panel := lipgloss.NewStyle().
		Width(width).
		Height(m.height - 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4a90e2")).
		Render(lipgloss.JoinVertical(
			title,
			"",
			agentsList.String(),
		))

	return panel
}

// renderConversationPanel renders the conversation area
func (m ChatroomModel) renderConversationPanel(width int) string {
	title := titleStyle.Render("💬 Conversation")
	
	// Render messages
	messagesView := strings.Builder{}
	visibleMessages := m.getVisibleMessages()
	
	for i, msg := range visibleMessages {
		messageStyle := m.getMessageStyle(msg)
		formattedMsg := messageStyle.Render(fmt.Sprintf("[%s] %s: %s", 
			msg.Timestamp.Format("15:04"), 
			m.getAgentName(msg.AgentID), 
			msg.Content))
		
		messagesView.WriteString(formattedMsg)
		messagesView.WriteString("\n")
	}

	// Render input area
	inputArea := lipgloss.JoinHorizontal(
		inputStyle.Render("📝 "),
		lipgloss.NewStyle().
			Width(width-10).
			Background(lipgloss.Color("#1a1a2e")).
			Foreground(lipgloss.Color("#ffffff")).
			Render(m.inputText),
	)

	panel := lipgloss.NewStyle().
		Width(width).
		Height(m.height - 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4a90e2")).
		Render(lipgloss.JoinVertical(
			title,
			"",
			messagesView.String(),
			"",
			inputArea,
		))

	return panel
}

// renderControlPanel renders the control panel
func (m ChatroomModel) renderControlPanel(width int) string {
	title := titleStyle.Render("⚙️ Controls")
	
	// Conversation type selector
	convTypes := []ConversationType{
		ConversationDemocratic,
		ConversationEnsemble,
		ConversationHierarchical,
		ConversationCompetitive,
		ConversationSpecialist,
		ConversationConsensus,
	}
	
	convSelector := strings.Builder{}
	for i, convType := range convTypes {
		prefix := "  "
		if i == 0 {
			prefix = "▶ "
		}
		convSelector.WriteString(fmt.Sprintf("%s%s %s\n", 
			prefix, convType, m.getConversationDescription(convType)))
	}

	// Provider selector
	providerSelector := strings.Builder{}
	for i, provider := range m.providers {
		prefix := "  "
		if provider.ID == m.activeProvider {
			prefix = "▶ "
		}
		providerSelector.WriteString(fmt.Sprintf("%s%s (%s)\n", 
			prefix, provider.Name, strings.Join(provider.Models, ", ")))
	}

	// Action buttons
	buttons := lipgloss.JoinHorizontal(
		buttonStyle.Render("👥 Add Agent"),
		buttonStyle.Render("📁 Files"),
		buttonStyle.Render("💾 Save"),
		buttonStyle.Render("⚙️ Settings"),
	)

	panel := lipgloss.NewStyle().
		Width(width).
		Height(m.height - 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4a90e2")).
		Render(lipgloss.JoinVertical(
			title,
			"",
			lipgloss.NewStyle().Render("Conversation Type:"),
			convSelector.String(),
			"",
			lipgloss.NewStyle().Render("AI Provider:"),
			providerSelector.String(),
			"",
			buttons,
		))

	return panel
}

// ==================== HELPER FUNCTIONS ====================

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

// contains checks if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ==================== HELPER METHODS ====================

// getMessageStyle returns appropriate style for message type
func (m ChatroomModel) getMessageStyle(msg Message) lipgloss.Style {
	switch msg.Type {
	case MessageTypeUser:
		return userStyle
	case MessageTypeAgent:
		return agentStyle
	case MessageTypeSystem:
		return systemStyle
	case MessageTypeTool:
		return warningStyle
	case MessageTypeFile:
		return infoStyle
	default:
		return chatroomStyle
	}
}

// getAgentName returns agent name by ID
func (m ChatroomModel) getAgentName(agentID string) string {
	for _, agent := range m.agents {
		if agent.ID == agentID {
			return agent.Name
		}
	}
	return "Unknown"
}

// getVisibleMessages returns messages for current view
func (m ChatroomModel) getVisibleMessages() []Message {
	maxVisible := m.height - 10 // Reserve space for input and UI
	start := m.scrollOffset
	end := start + maxVisible
	
	if end > len(m.messages) {
		end = len(m.messages)
	}
	
	if start >= len(m.messages) {
		return []Message{}
	}
	
	return m.messages[start:end]
}

// getConversationDescription returns description for conversation type
func (m ChatroomModel) getConversationDescription(convType ConversationType) string {
	switch convType {
	case ConversationDemocratic:
		return "🗳️ Democratic - All agents discuss equally"
	case ConversationEnsemble:
		return "🎭 Ensemble - Agents build on each other's responses"
	case ConversationHierarchical:
		return "🏛️ Hierarchical - Lead agent directs others"
	case ConversationCompetitive:
		return "⚔️ Competitive - Agents compete for best solution"
	case ConversationSpecialist:
		return "🎯 Specialist - Each agent has domain expertise"
	case ConversationConsensus:
		return "🤝 Consensus - Agents work toward agreement"
	default:
		return "❓ Unknown conversation type"
	}
}

// ==================== ACTION METHODS ====================

// sendMessage sends a user message
func (m ChatroomModel) sendMessage() tea.Cmd {
	msg := Message{
		ID:        generateID(),
		Timestamp: time.Now(),
		Content:   m.inputText,
		Type:      MessageTypeUser,
		Metadata:  map[string]interface{}{
			"conversation_id": m.activeConvID,
		},
	}
	
	m.messages = append(m.messages, msg)
	m.inputText = ""
	
	// In a real implementation, this would send to active agents
	return func() tea.Msg {
		return msg
	}
}

// createNewConversation creates a new conversation
func (m ChatroomModel) createNewConversation() tea.Cmd {
	conv := Conversation{
		ID:       generateID(),
		Type:     ConversationDemocratic,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive: true,
	}
	
	m.conversations = append(m.conversations, conv)
	m.activeConvID = conv.ID
	
	return func() tea.Msg {
		return conv
	}
}

// navigateHistory handles message history navigation
func (m ChatroomModel) navigateHistory(msg tea.KeyMsg) tea.Cmd {
	// Implement message history navigation
	return nil
}

// navigateMessages handles message list navigation
func (m ChatroomModel) navigateMessages(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyUp:
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
	case tea.KeyDown:
		if m.scrollOffset < len(m.messages)-10 {
			m.scrollOffset++
		}
	case tea.KeyPgUp:
		m.scrollOffset -= 10
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
	case tea.KeyPgDown:
		m.scrollOffset += 10
		if m.scrollOffset > len(m.messages)-10 {
			m.scrollOffset = len(m.messages) - 10
		}
	}
	return nil
}

// cycleFocus cycles through UI sections
func (m ChatroomModel) cycleFocus() tea.Cmd {
	// Implement focus cycling between agent panel, conversation, and controls
	return nil
}

// ==================== MAIN FUNCTION ====================

// main runs the chatroom application
func main() {
	p := tea.NewProgram(NewChatroomModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running chatroom: %v", err)
		os.Exit(1)
	}
}