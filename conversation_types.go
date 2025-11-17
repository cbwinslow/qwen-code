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

// ==================== CONVERSATION TYPES ====================

// ConversationType represents different conversation modes
type ConversationType string

const (
	ConversationDemocratic   ConversationType = "democratic"
	ConversationEnsemble     ConversationType = "ensemble"
	ConversationHierarchical ConversationType = "hierarchical"
	ConversationCompetitive  ConversationType = "competitive"
	ConversationSpecialist   ConversationType = "specialist"
	ConversationConsensus    ConversationType = "consensus"
	ConversationBrainstorm   ConversationType = "brainstorm"
	ConversationDebate       ConversationType = "debate"
	ConversationPeerReview   ConversationType = "peer_review"
	ConversationSocratic     ConversationType = "socratic"
)

// ConversationConfig holds configuration for conversation types
type ConversationConfig struct {
	Type              ConversationType       `json:"type"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Icon              string                 `json:"icon"`
	MaxParticipants   int                    `json:"max_participants"`
	MinParticipants   int                    `json:"min_participants"`
	RequiresModerator bool                   `json:"requires_moderator"`
	Settings          map[string]interface{} `json:"settings,omitempty"`
	Enabled           bool                   `json:"enabled"`
}

// ConversationRule defines rules for conversations
type ConversationRule struct {
	ID          string                 `json:"id"`
	Type        ConversationType       `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
}

// ConversationState represents the state of a conversation
type ConversationState struct {
	ID           string                 `json:"id"`
	Type         ConversationType       `json:"type"`
	Participants []string               `json:"participants"`
	IsActive     bool                   `json:"is_active"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Subject      string                 `json:"subject,omitempty"`
	Messages     []ConversationMessage  `json:"messages"`
	CurrentTurn  int                    `json:"current_turn"`
	TurnOrder    []string               `json:"turn_order"`
	Moderator    string                 `json:"moderator,omitempty"`
	Settings     map[string]interface{} `json:"settings"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ConversationMessage represents a message in a conversation
type ConversationMessage struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	AgentID     string                 `json:"agent_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	Content     string                 `json:"content"`
	Type        string                 `json:"type"` // "user", "agent", "system"
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	ParentID    string                 `json:"parent_id,omitempty"`
	ThreadID    string                 `json:"thread_id,omitempty"`
	Votes       map[string]int         `json:"votes,omitempty"`
	Reactions   map[string][]string    `json:"reactions,omitempty"`
	Edited      bool                   `json:"edited,omitempty"`
	EditHistory []EditHistory          `json:"edit_history,omitempty"`
}

// EditHistory tracks changes to messages
type EditHistory struct {
	ID         string    `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	UserID     string    `json:"user_id"`
	OldContent string    `json:"old_content"`
	NewContent string    `json:"new_content"`
	Reason     string    `json:"reason"`
}

// ==================== CONVERSATION MANAGER ====================

// ConversationManager manages different conversation types
type ConversationManager struct {
	configs      map[string]ConversationConfig
	rules        map[string]ConversationRule
	states       map[string]*ConversationState
	activeConv   string
	eventHandler func(event ConversationEvent)
	mu           sync.RWMutex
}

// ConversationEvent represents conversation-related events
type ConversationEvent struct {
	Type      string                 `json:"type"`
	ConvID    string                 `json:"conv_id,omitempty"`
	AgentID   string                 `json:"agent_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Message   string                 `json:"message,omitempty"`
}

// ==================== CONVERSATION MANAGER IMPLEMENTATION ====================

// NewConversationManager creates a new conversation manager
func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		configs:      make(map[string]ConversationConfig),
		rules:        make(map[string]ConversationRule),
		states:       make(map[string]*ConversationState),
		eventHandler: func(event ConversationEvent) {},
		mu:           sync.RWMutex{},
	}
}

// LoadConfigs loads conversation configurations
func (cm *ConversationManager) LoadConfigs(configPath string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cm.createDefaultConfigs(configPath)
		}
		return err
	}

	var configs map[string]ConversationConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("failed to unmarshal conversation configs: %w", err)
	}

	cm.configs = configs
	return cm.initializeConversations()
}

// SaveConfigs saves conversation configurations
func (cm *ConversationManager) SaveConfigs(configPath string) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.MarshalIndent(cm.configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversation configs: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// createDefaultConfigs creates default conversation configurations
func (cm *ConversationManager) createDefaultConfigs(configPath string) error {
	defaultConfigs := map[string]ConversationConfig{
		"democratic": {
			Type:              ConversationDemocratic,
			Name:              "Democratic",
			Description:       "All participants discuss equally, decisions made by majority vote",
			Icon:              "🗳️",
			MaxParticipants:   20,
			MinParticipants:   2,
			RequiresModerator: false,
			Settings: map[string]interface{}{
				"voting_timeout":  30,
				"allow_anonymous": false,
			},
			Enabled: true,
		},
		"ensemble": {
			Type:              ConversationEnsemble,
			Name:              "Ensemble",
			Description:       "Agents build on each other's responses, creating a refined final output",
			Icon:              "🎭",
			MaxParticipants:   10,
			MinParticipants:   2,
			RequiresModerator: false,
			Settings: map[string]interface{}{
				"refinement_rounds": 3,
				"synthesis_method":  "weighted_average",
			},
			Enabled: true,
		},
		"hierarchical": {
			Type:              ConversationHierarchical,
			Name:              "Hierarchical",
			Description:       "Lead agent directs conversation, others follow chain of command",
			Icon:              "🏛️",
			MaxParticipants:   50,
			MinParticipants:   2,
			RequiresModerator: true,
			Settings: map[string]interface{}{
				"allow_delegation":   true,
				"escalation_timeout": 60,
			},
			Enabled: true,
		},
		"competitive": {
			Type:              ConversationCompetitive,
			Name:              "Competitive",
			Description:       "Agents compete to provide the best solution, winner takes all",
			Icon:              "⚔️",
			MaxParticipants:   8,
			MinParticipants:   2,
			RequiresModerator: true,
			Settings: map[string]interface{}{
				"scoring_criteria": "quality_speed",
				"time_limit":       300,
			},
			Enabled: true,
		},
		"specialist": {
			Type:              ConversationSpecialist,
			Name:              "Specialist",
			Description:       "Each agent has domain expertise, provides specialized insights",
			Icon:              "🎯",
			MaxParticipants:   15,
			MinParticipants:   1,
			RequiresModerator: false,
			Settings: map[string]interface{}{
				"domain_focus":    true,
				"expertise_level": "senior",
			},
			Enabled: true,
		},
		"consensus": {
			Type:              ConversationConsensus,
			Name:              "Consensus",
			Description:       "Agents work together until agreement is reached",
			Icon:              "🤝",
			MaxParticipants:   12,
			MinParticipants:   3,
			RequiresModerator: true,
			Settings: map[string]interface{}{
				"consensus_threshold": 0.8,
				"max_discussion_time": 600,
			},
			Enabled: true,
		},
		"brainstorm": {
			Type:              ConversationBrainstorm,
			Name:              "Brainstorm",
			Description:       "Free-flow idea generation without criticism",
			Icon:              "💡",
			MaxParticipants:   25,
			MinParticipants:   3,
			RequiresModerator: false,
			Settings: map[string]interface{}{
				"idea_generation_phase": "divergent",
				"convergence_phase":     "convergent",
			},
			Enabled: true,
		},
		"debate": {
			Type:              ConversationDebate,
			Name:              "Debate",
			Description:       "Structured debate on opposing viewpoints with formal rules",
			Icon:              "🎤�",
			MaxParticipants:   6,
			MinParticipants:   2,
			RequiresModerator: true,
			Settings: map[string]interface{}{
				"speaking_time": 120,
				"rebuttal_time": 60,
				"moderation":    true,
			},
			Enabled: true,
		},
		"peer_review": {
			Type:              ConversationPeerReview,
			Name:              "Peer Review",
			Description:       "Agents review and critique each other's work constructively",
			Icon:              "👥",
			MaxParticipants:   8,
			MinParticipants:   2,
			RequiresModerator: false,
			Settings: map[string]interface{}{
				"review_criteria": "completeness_accuracy",
				"anonymity":       true,
			},
			Enabled: true,
		},
		"socratic": {
			Type:              ConversationSocratic,
			Name:              "Socratic",
			Description:       "Question-based dialogue to uncover truth through inquiry",
			Icon:              "🔍",
			MaxParticipants:   12,
			MinParticipants:   2,
			RequiresModerator: true,
			Settings: map[string]interface{}{
				"questioning_method":  "socratic",
				"follow_up_questions": true,
			},
			Enabled: true,
		},
	}

	cm.configs = defaultConfigs

	data, err := json.MarshalIndent(defaultConfigs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// initializeConversations initializes conversation states
func (cm *ConversationManager) initializeConversations() error {
	for id, config := range cm.configs {
		if config.Enabled {
			state := &ConversationState{
				ID:           generateID(),
				Type:         config.Type,
				Participants: []string{},
				IsActive:     false,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Settings:     config.Settings,
				Metadata: map[string]interface{}{
					"config": config,
				},
			}
			cm.states[id] = state
		}
	}

	return nil
}

// CreateConversation creates a new conversation
func (cm *ConversationManager) CreateConversation(convType string, participants []string, subject string, userID string) (*ConversationState, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	config, exists := cm.configs[convType]
	if !exists || !config.Enabled {
		return nil, fmt.Errorf("conversation type %s is not available or not enabled", convType)
	}

	// Validate participants
	if len(participants) < config.MinParticipants {
		return nil, fmt.Errorf("conversation type %s requires at least %d participants", convType, config.MinParticipants)
	}

	if len(participants) > config.MaxParticipants {
		return nil, fmt.Errorf("conversation type %s supports maximum %d participants", convType, config.MaxParticipants)
	}

	// Check moderator requirement
	if config.RequiresModerator {
		hasModerator := false
		for _, participant := range participants {
			if cm.isModerator(participant) {
				hasModerator = true
				break
			}
		}

		if !hasModerator {
			return nil, fmt.Errorf("conversation type %s requires a moderator", convType)
		}
	}

	state := &ConversationState{
		ID:           generateID(),
		Type:         config.Type,
		Participants: participants,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Subject:      subject,
		CurrentTurn:  0,
		TurnOrder:    participants,
		Moderator:    "",
		Settings:     config.Settings,
		Metadata: map[string]interface{}{
			"creator_id": userID,
			"config":     config,
		},
	}

	cm.states[state.ID] = state
	cm.activeConv = state.ID

	if cm.eventHandler != nil {
		cm.eventHandler(ConversationEvent{
			Type:      "conversation_created",
			ConvID:    state.ID,
			UserID:    userID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"conversation": state,
				"config":       config,
			},
			Message: fmt.Sprintf("Conversation %s created with %d participants", config.Name, len(participants)),
		})
	}

	return state, nil
}

// AddMessage adds a message to a conversation
func (cm *ConversationManager) AddMessage(convID string, message ConversationMessage) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	state, exists := cm.states[convID]
	if !exists {
		return fmt.Errorf("conversation %s not found", convID)
	}

	message.ID = generateID()
	message.Timestamp = time.Now()
	state.Messages = append(state.Messages, message)
	state.UpdatedAt = time.Now()

	if cm.eventHandler != nil {
		cm.eventHandler(ConversationEvent{
			Type:      "message_added",
			ConvID:    convID,
			AgentID:   message.AgentID,
			UserID:    message.UserID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"message": message,
			},
			Message: fmt.Sprintf("Message added to conversation %s", convID),
		})
	}

	return nil
}

// GetConversation returns a conversation state
func (cm *ConversationManager) GetConversation(convID string) (*ConversationState, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	state, exists := cm.states[convID]
	if !exists {
		return nil, fmt.Errorf("conversation %s not found", convID)
	}

	return state, nil
}

// GetActiveConversation returns the currently active conversation
func (cm *ConversationManager) GetActiveConversation() (*ConversationState, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if cm.activeConv == "" {
		return nil, fmt.Errorf("no active conversation")
	}

	return cm.GetConversation(cm.activeConv)
}

// EndConversation ends a conversation
func (cm *ConversationManager) EndConversation(convID string, summary string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	state, exists := cm.states[convID]
	if !exists {
		return fmt.Errorf("conversation %s not found", convID)
	}

	state.IsActive = false
	state.UpdatedAt = time.Now()

	// Add summary as final system message if provided
	if summary != "" {
		summaryMessage := ConversationMessage{
			ID:        generateID(),
			Timestamp: time.Now(),
			Type:      "system",
			Content:   fmt.Sprintf("Conversation ended: %s", summary),
			Metadata: map[string]interface{}{
				"summary":  summary,
				"duration": time.Since(state.CreatedAt).String(),
			},
		}
		state.Messages = append(state.Messages, summaryMessage)
	}

	cm.activeConv = ""

	if cm.eventHandler != nil {
		cm.eventHandler(ConversationEvent{
			Type:      "conversation_ended",
			ConvID:    convID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"conversation": state,
				"summary":      summary,
			},
			Message: fmt.Sprintf("Conversation %s ended", convID),
		})
	}

	return nil
}

// isModerator checks if a user has moderator privileges
func (cm *ConversationManager) isModerator(userID string) bool {
	// In a real implementation, this would check against a user database
	// For now, we'll use a simple check
	moderators := []string{"admin", "moderator", "owner"}
	for _, moderator := range moderators {
		if userID == moderator {
			return true
		}
	}
	return false
}

// SetEventHandler sets the event handler
func (cm *ConversationManager) SetEventHandler(handler func(event ConversationEvent)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.eventHandler = handler
}

// GetAvailableTypes returns all available conversation types
func (cm *ConversationManager) GetAvailableTypes() []ConversationConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var types []ConversationConfig
	for _, config := range cm.configs {
		if config.Enabled {
			types = append(types, config)
		}
	}

	return types
}

// GetConversationRules returns rules for a conversation type
func (cm *ConversationManager) GetConversationRules(convType string) ([]ConversationRule, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var rules []ConversationRule
	for _, rule := range cm.rules {
		if rule.Type == convType && rule.Enabled {
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

// ==================== MAIN FUNCTION ====================

// main for testing conversation type selection
func main() {
	fmt.Println("🗣️ Testing Conversation Type Selection")

	// Create conversation manager
	convoManager := NewConversationManager()

	// Load configurations
	configPath := "/tmp/conversation-configs.json"
	if err := convoManager.LoadConfigs(configPath); err != nil {
		fmt.Printf("❌ Failed to load configs: %v\n", err)
		return
	}

	// Set event handler
	convoManager.SetEventHandler(func(event ConversationEvent) {
		fmt.Printf("🗣️ Conversation Event: %s - %s\n", event.Type, event.Message)
		if event.Data != nil {
			data, _ := json.MarshalIndent(event.Data, "", "  ")
			fmt.Printf("   Data: %s\n", string(data))
		}
	})

	// Test creating different conversation types
	participants := []string{"agent-1", "agent-2", "agent-3"}

	// Test democratic conversation
	democraticConv, err := convoManager.CreateConversation("democratic", participants, "Test Democratic Conversation", "test-user")
	if err != nil {
		fmt.Printf("❌ Failed to create democratic conversation: %v\n", err)
	} else {
		fmt.Printf("✅ Democratic conversation created: %s\n", democraticConv.ID)
	}

	// Test ensemble conversation
	ensembleConv, err := convoManager.CreateConversation("ensemble", participants[:2], "Test Ensemble Conversation", "test-user")
	if err != nil {
		fmt.Printf("❌ Failed to create ensemble conversation: %v\n", err)
	} else {
		fmt.Printf("✅ Ensemble conversation created: %s\n", ensembleConv.ID)
	}

	// Test hierarchical conversation
	hierarchicalConv, err := convoManager.CreateConversation("hierarchical", participants, "Test Hierarchical Conversation", "test-user")
	if err != nil {
		fmt.Printf("❌ Failed to create hierarchical conversation: %v\n", err)
	} else {
		fmt.Printf("✅ Hierarchical conversation created: %s\n", hierarchicalConv.ID)
	}

	// Display available types
	types := convoManager.GetAvailableTypes()
	fmt.Printf("✅ Available conversation types:\n")
	for _, convType := range types {
		fmt.Printf("  %s %s - %s\n", convType.Icon, convType.Name, convType.Description)
	}

	// Test adding messages
	testMessage := ConversationMessage{
		ID:        generateID(),
		Timestamp: time.Now(),
		Content:   "This is a test message for the conversation system.",
		Type:      "user",
		UserID:    "test-user",
	}

	if err := convoManager.AddMessage(democraticConv.ID, testMessage); err != nil {
		fmt.Printf("❌ Failed to add message: %v\n", err)
	} else {
		fmt.Printf("✅ Message added to democratic conversation\n")
	}

	// Test ending conversation
	if err := convoManager.EndConversation(democraticConv.ID, "Test completed successfully"); err != nil {
		fmt.Printf("❌ Failed to end conversation: %v\n", err)
	} else {
		fmt.Printf("✅ Democratic conversation ended\n")
	}

	// Get conversation details
	conv, err := convoManager.GetConversation(democraticConv.ID)
	if err != nil {
		fmt.Printf("❌ Failed to get conversation: %v\n", err)
		return
	}

	fmt.Printf("✅ Retrieved conversation with %d messages\n", len(conv.Messages))

	// Save configurations
	if err := convoManager.SaveConfigs(configPath); err != nil {
		fmt.Printf("❌ Failed to save configs: %v\n", err)
		return
	}

	fmt.Printf("✅ Conversation configurations saved\n")

	fmt.Println("🎉 Conversation type selection test completed successfully!")
}
