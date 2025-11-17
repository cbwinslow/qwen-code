package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== AGENT MANAGEMENT TYPES ====================

// AgentRole defines different roles for agents
type AgentRole string

const (
	RoleCoordinator AgentRole = "coordinator"
	RoleSpecialist  AgentRole = "specialist"
	RoleCritic      AgentRole = "critic"
	RoleSynthesizer AgentRole = "synthesizer"
	RoleResearcher  AgentRole = "researcher"
	RoleImplementer AgentRole = "implementer"
	RoleReviewer    AgentRole = "reviewer"
)

// AgentPersonality defines different personality types
type AgentPersonality string

const (
	PersonalityAnalytical    AgentPersonality = "analytical"
	PersonalityCreative      AgentPersonality = "creative"
	PersonalitySkeptical     AgentPersonality = "skeptical"
	PersonalityCollaborative AgentPersonality = "collaborative"
	PersonalityEfficient     AgentPersonality = "efficient"
	PersonalityThorough      AgentPersonality = "thorough"
	PersonalityDiplomatic    AgentPersonality = "diplomatic"
)

// AgentCapability defines what an agent can do
type AgentCapability string

const (
	CapabilityTextGeneration  AgentCapability = "text_generation"
	CapabilityCodeGeneration  AgentCapability = "code_generation"
	CapabilityDataAnalysis    AgentCapability = "data_analysis"
	CapabilityWebSearch       AgentCapability = "web_search"
	CapabilityFileProcessing  AgentCapability = "file_processing"
	CapabilityImageGeneration AgentCapability = "image_generation"
	CapabilityReasoning       AgentCapability = "reasoning"
	CapabilityTranslation     AgentCapability = "translation"
	CapabilitySummarization   AgentCapability = "summarization"
)

// AgentConfig defines configuration for an agent
type AgentConfig struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Role         AgentRole              `json:"role"`
	Personality  AgentPersonality       `json:"personality"`
	Provider     string                 `json:"provider"`
	Model        string                 `json:"model"`
	APIKey       string                 `json:"api_key,omitempty"`
	MaxTokens    int                    `json:"max_tokens"`
	Temperature  float64                `json:"temperature"`
	Capabilities []AgentCapability      `json:"capabilities"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
	IsActive     bool                   `json:"is_active"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// AgentTask represents a task assigned to an agent
type AgentTask struct {
	ID          string                 `json:"id"`
	AgentID     string                 `json:"agent_id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Priority    int                    `json:"priority"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Result      interface{}            `json:"result,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AgentStatus represents the current status of an agent
type AgentStatus struct {
	AgentID     string    `json:"agent_id"`
	Status      string    `json:"status"`
	LastSeen    time.Time `json:"last_seen"`
	CurrentTask string    `json:"current_task,omitempty"`
	TasksTotal  int       `json:"tasks_total"`
	TasksDone   int       `json:"tasks_done"`
	Performance float64   `json:"performance"`
	ErrorRate   float64   `json:"error_rate"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ==================== AGENT MANAGER ====================

// AgentManager manages multiple AI agents
type AgentManager struct {
	agents       map[string]*ManagedAgent
	tasks        chan AgentTask
	configs      map[string]AgentConfig
	statuses     map[string]*AgentStatus
	mu           sync.RWMutex
	eventHandler func(event AgentEvent)
}

// ManagedAgent represents an agent being managed
type ManagedAgent struct {
	Config      AgentConfig
	Status      AgentStatus
	Tasks       []AgentTask
	Performance AgentPerformance
	mu          sync.RWMutex
}

// AgentPerformance tracks agent performance metrics
type AgentPerformance struct {
	TasksCompleted      int       `json:"tasks_completed"`
	TasksTotal          int       `json:"tasks_total"`
	AverageResponseTime float64   `json:"average_response_time"`
	SuccessRate         float64   `json:"success_rate"`
	ErrorRate           float64   `json:"error_rate"`
	QualityScore        float64   `json:"quality_score"`
	LastUpdated         time.Time `json:"last_updated"`
}

// AgentEvent represents events from agents
type AgentEvent struct {
	Type      string                 `json:"type"`
	AgentID   string                 `json:"agent_id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Message   string                 `json:"message,omitempty"`
}

// ==================== AGENT MANAGER IMPLEMENTATION ====================

// NewAgentManager creates a new agent manager
func NewAgentManager() *AgentManager {
	return &AgentManager{
		agents:       make(map[string]*ManagedAgent),
		tasks:        make(chan AgentTask, 100),
		configs:      make(map[string]AgentConfig),
		statuses:     make(map[string]*AgentStatus),
		eventHandler: func(event AgentEvent) {},
	}
}

// LoadConfigs loads agent configurations from file
func (am *AgentManager) LoadConfigs(configPath string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default configs if file doesn't exist
			return am.createDefaultConfigs(configPath)
		}
		return err
	}

	var configs map[string]AgentConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("failed to unmarshal agent configs: %w", err)
	}

	am.configs = configs
	return am.initializeAgents()
}

// SaveConfigs saves agent configurations to file
func (am *AgentManager) SaveConfigs(configPath string) error {
	am.mu.RLock()
	defer am.mu.RUnlock()

	data, err := json.MarshalIndent(am.configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal agent configs: %w", err)
	}

	return os.WriteFile(configPath, data, 0644)
}

// createDefaultConfigs creates default agent configurations
func (am *AgentManager) createDefaultConfigs(configPath string) error {
	defaultConfigs := map[string]AgentConfig{
		"coordinator": {
			ID:          "coordinator-1",
			Name:        "Coordinator",
			Role:        RoleCoordinator,
			Personality: PersonalityDiplomatic,
			Provider:    "openrouter",
			Model:       "anthropic/claude-3-sonnet-20240229",
			MaxTokens:   4096,
			Temperature: 0.7,
			Capabilities: []AgentCapability{
				CapabilityTextGeneration,
				CapabilityReasoning,
				CapabilityCoordination,
			},
			IsActive:  true,
			CreatedAt: time.Now(),
		},
		"specialist-coder": {
			ID:          "specialist-coder-1",
			Name:        "Code Specialist",
			Role:        RoleSpecialist,
			Personality: PersonalityAnalytical,
			Provider:    "openrouter",
			Model:       "qwen/qwen-2.5-coder-32b",
			MaxTokens:   8192,
			Temperature: 0.1,
			Capabilities: []AgentCapability{
				CapabilityCodeGeneration,
				CapabilityCodeAnalysis,
				CapabilityDebugging,
			},
			IsActive:  true,
			CreatedAt: time.Now(),
		},
		"specialist-analyst": {
			ID:          "specialist-analyst-1",
			Name:        "Data Analyst",
			Role:        RoleSpecialist,
			Personality: PersonalityThorough,
			Provider:    "openrouter",
			Model:       "meta-llama/llama-3-70b-instruct",
			MaxTokens:   4096,
			Temperature: 0.3,
			Capabilities: []AgentCapability{
				CapabilityDataAnalysis,
				CapabilityReasoning,
				CapabilityVisualization,
			},
			IsActive:  true,
			CreatedAt: time.Now(),
		},
		"critic": {
			ID:          "critic-1",
			Name:        "Quality Critic",
			Role:        RoleCritic,
			Personality: PersonalitySkeptical,
			Provider:    "openrouter",
			Model:       "anthropic/claude-3-sonnet-20240229",
			MaxTokens:   4096,
			Temperature: 0.5,
			Capabilities: []AgentCapability{
				CapabilityTextGeneration,
				CapabilityReasoning,
				CapabilityQualityAssessment,
			},
			IsActive:  true,
			CreatedAt: time.Now(),
		},
		"researcher": {
			ID:          "researcher-1",
			Name:        "Researcher",
			Role:        RoleResearcher,
			Personality: PersonalityThorough,
			Provider:    "openrouter",
			Model:       "meta-llama/llama-3-70b-instruct",
			MaxTokens:   4096,
			Temperature: 0.2,
			Capabilities: []AgentCapability{
				CapabilityWebSearch,
				CapabilityDataAnalysis,
				CapabilityFactChecking,
			},
			IsActive:  true,
			CreatedAt: time.Now(),
		},
	}

	am.configs = defaultConfigs

	data, err := json.MarshalIndent(defaultConfigs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// initializeAgents initializes agents from configs
func (am *AgentManager) initializeAgents() error {
	for id, config := range am.configs {
		agent := &ManagedAgent{
			Config: config,
			Status: AgentStatus{
				AgentID:    id,
				Status:     "idle",
				LastSeen:   time.Now(),
				TasksTotal: 0,
				TasksDone:  0,
				Performance: AgentPerformance{
					LastUpdated: time.Now(),
				},
				UpdatedAt: time.Now(),
			},
			Tasks: []AgentTask{},
			Performance: AgentPerformance{
				LastUpdated: time.Now(),
			},
		}

		am.agents[id] = agent
		am.statuses[id] = &agent.Status
	}

	return nil
}

// AddAgent adds a new agent
func (am *AgentManager) AddAgent(config AgentConfig) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.configs[config.ID]; exists {
		return fmt.Errorf("agent with ID %s already exists", config.ID)
	}

	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()
	am.configs[config.ID] = config

	agent := &ManagedAgent{
		Config: config,
		Status: AgentStatus{
			AgentID:    config.ID,
			Status:     "idle",
			LastSeen:   time.Now(),
			TasksTotal: 0,
			TasksDone:  0,
			Performance: AgentPerformance{
				LastUpdated: time.Now(),
			},
			UpdatedAt: time.Now(),
		},
		Tasks: []AgentTask{},
		Performance: AgentPerformance{
			LastUpdated: time.Now(),
		},
	}

	am.agents[config.ID] = agent
	am.statuses[config.ID] = &agent.Status

	if am.eventHandler != nil {
		am.eventHandler(AgentEvent{
			Type:      "agent_added",
			AgentID:   config.ID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"config": config,
			},
			Message: fmt.Sprintf("Agent %s added", config.Name),
		})
	}

	return nil
}

// RemoveAgent removes an agent
func (am *AgentManager) RemoveAgent(agentID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.agents[agentID]; !exists {
		return fmt.Errorf("agent with ID %s not found", agentID)
	}

	delete(am.agents, agentID)
	delete(am.configs, agentID)
	delete(am.statuses, agentID)

	if am.eventHandler != nil {
		am.eventHandler(AgentEvent{
			Type:      "agent_removed",
			AgentID:   agentID,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("Agent %s removed", agentID),
		})
	}

	return nil
}

// UpdateAgent updates an agent configuration
func (am *AgentManager) UpdateAgent(config AgentConfig) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.agents[config.ID]; !exists {
		return fmt.Errorf("agent with ID %s not found", config.ID)
	}

	config.UpdatedAt = time.Now()
	am.configs[config.ID] = config

	if agent := am.agents[config.ID]; agent != nil {
		agent.Config = config
		agent.Status.UpdatedAt = time.Now()
	}

	if am.eventHandler != nil {
		am.eventHandler(AgentEvent{
			Type:      "agent_updated",
			AgentID:   config.ID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"config": config,
			},
			Message: fmt.Sprintf("Agent %s updated", config.Name),
		})
	}

	return nil
}

// AssignTask assigns a task to an agent
func (am *AgentManager) AssignTask(task AgentTask) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if agent := am.agents[task.AgentID]; agent == nil {
		return fmt.Errorf("agent with ID %s not found", task.AgentID)
	}

	agent.mu.Lock()
	defer agent.mu.Unlock()

	task.CreatedAt = time.Now()
	agent.Tasks = append(agent.Tasks, task)
	agent.Status.TasksTotal++
	agent.Status.CurrentTask = task.ID
	agent.Status.UpdatedAt = time.Now()

	// Send task to agent
	select {
	case am.tasks <- task:
	default:
		return fmt.Errorf("task queue is full")
	}

	if am.eventHandler != nil {
		am.eventHandler(AgentEvent{
			Type:      "task_assigned",
			AgentID:   task.AgentID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"task": task,
			},
			Message: fmt.Sprintf("Task assigned to agent %s", task.AgentID),
		})
	}

	return nil
}

// GetAgents returns all agents
func (am *AgentManager) GetAgents() map[string]*ManagedAgent {
	am.mu.RLock()
	defer am.mu.RUnlock()

	result := make(map[string]*ManagedAgent)
	for id, agent := range am.agents {
		result[id] = agent
	}

	return result
}

// GetAgentStatus returns status of a specific agent
func (am *AgentManager) GetAgentStatus(agentID string) (*AgentStatus, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	status, exists := am.statuses[agentID]
	if !exists {
		return nil, fmt.Errorf("agent with ID %s not found", agentID)
	}

	return status, nil
}

// UpdateAgentStatus updates agent status
func (am *AgentManager) UpdateAgentStatus(agentID string, status AgentStatus) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if agent := am.agents[agentID]; agent == nil {
		return fmt.Errorf("agent with ID %s not found", agentID)
	}

	agent.mu.Lock()
	defer agent.mu.Unlock()

	agent.Status = status
	agent.Status.UpdatedAt = time.Now()
	am.statuses[agentID] = &agent.Status

	if am.eventHandler != nil {
		am.eventHandler(AgentEvent{
			Type:      "status_updated",
			AgentID:   agentID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"status": status,
			},
			Message: fmt.Sprintf("Agent %s status updated to %s", agentID, status.Status),
		})
	}

	return nil
}

// SetEventHandler sets the event handler
func (am *AgentManager) SetEventHandler(handler func(event AgentEvent)) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.eventHandler = handler
}

// GetActiveAgents returns currently active agents
func (am *AgentManager) GetActiveAgents() []*ManagedAgent {
	am.mu.RLock()
	defer am.mu.RUnlock()

	var active []*ManagedAgent
	for _, agent := range am.agents {
		if agent.Config.IsActive {
			active = append(active, agent)
		}
	}

	return active
}

// GetPerformanceMetrics returns performance metrics for all agents
func (am *AgentManager) GetPerformanceMetrics() map[string]AgentPerformance {
	am.mu.RLock()
	defer am.mu.RUnlock()

	metrics := make(map[string]AgentPerformance)
	for id, agent := range am.agents {
		metrics[id] = agent.Performance
	}

	return metrics
}

// ==================== AGENT COORDINATION ====================

// TaskDistributor handles task distribution among agents
type TaskDistributor struct {
	agentManager *AgentManager
	strategy     CoordinationStrategy
}

type CoordinationStrategy interface {
	DistributeTask(task AgentTask, agents []*ManagedAgent) (*ManagedAgent, error)
	ShouldIntervene(task AgentTask, agents []*ManagedAgent) bool
}

// RoundRobinStrategy distributes tasks in round-robin fashion
type RoundRobinStrategy struct {
	lastAssigned map[string]int
}

func (rr *RoundRobinStrategy) DistributeTask(task AgentTask, agents []*ManagedAgent) (*ManagedAgent, error) {
	if len(agents) == 0 {
		return nil, fmt.Errorf("no agents available")
	}

	// Find next agent in round-robin fashion
	var selectedAgent *ManagedAgent
	minTasks := int(^uint(0) >> 1) // Max int

	for _, agent := range agents {
		if agent.Status.TasksTotal < minTasks {
			minTasks = agent.Status.TasksTotal
			selectedAgent = agent
		}
	}

	if selectedAgent == nil {
		selectedAgent = agents[0]
	}

	rr.lastAssigned[selectedAgent.Config.ID] = rr.lastAssigned[selectedAgent.Config.ID] + 1

	return selectedAgent, nil
}

func (rr *RoundRobinStrategy) ShouldIntervene(task AgentTask, agents []*ManagedAgent) bool {
	// Simple intervention logic
	for _, agent := range agents {
		if agent.Status.ErrorRate > 0.5 { // 50% error rate
			return true
		}
	}
	return false
}

// NewTaskDistributor creates a new task distributor
func NewTaskDistributor(am *AgentManager) *TaskDistributor {
	return &TaskDistributor{
		agentManager: am,
		strategy: &RoundRobinStrategy{
			lastAssigned: make(map[string]int),
		},
	}
}

// DistributeTask distributes a task using the configured strategy
func (td *TaskDistributor) DistributeTask(task AgentTask) (*ManagedAgent, error) {
	agents := td.agentManager.GetActiveAgents()
	if len(agents) == 0 {
		return nil, fmt.Errorf("no active agents available")
	}

	// Check if intervention is needed
	if td.strategy.ShouldIntervene(task, agents) {
		// Use a different strategy or escalate
		return td.distributeWithFallback(task, agents)
	}

	return td.strategy.DistributeTask(task, agents)
}

// distributeWithFallback handles fallback distribution
func (td *TaskDistributor) distributeWithFallback(task AgentTask, agents []*ManagedAgent) (*ManagedAgent, error) {
	// Simple fallback: choose agent with best performance
	var bestAgent *ManagedAgent
	bestScore := 0.0

	for _, agent := range agents {
		score := agent.Performance.SuccessRate - agent.Performance.ErrorRate
		if score > bestScore {
			bestScore = score
			bestAgent = agent
		}
	}

	if bestAgent == nil {
		bestAgent = agents[0]
	}

	return bestAgent, nil
}

// ==================== MAIN FUNCTION ====================

// main for testing the agent manager
func main() {
	am := NewAgentManager()

	// Load configurations
	configPath := filepath.Join(os.TempDir(), "agent-configs.json")
	if err := am.LoadConfigs(configPath); err != nil {
		fmt.Printf("Error loading configs: %v\n", err)
		os.Exit(1)
	}

	// Set event handler to see events
	am.SetEventHandler(func(event AgentEvent) {
		fmt.Printf("Event: %s - %s\n", event.Type, event.Message)
		if event.Data != nil {
			data, _ := json.MarshalIndent(event.Data, "", "  ")
			fmt.Printf("Data: %s\n", string(data))
		}
	})

	// Test adding a task
	task := AgentTask{
		ID:          "test-task-1",
		Type:        "code_review",
		Description: "Review the chatroom code for quality",
		Priority:    1,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Get active agents and assign task
	agents := am.GetActiveAgents()
	if len(agents) > 0 {
		agent, err := am.AssignTask(task)
		if err != nil {
			fmt.Printf("Error assigning task: %v\n", err)
		} else {
			fmt.Printf("Task assigned to agent: %s\n", agent.Config.Name)
		}
	}

	// Print agent status
	for id, status := range am.GetPerformanceMetrics() {
		fmt.Printf("Agent %s: Tasks=%d/%d, Success=%.2f, Error=%.2f\n",
			id, status.TasksDone, status.TasksTotal, status.SuccessRate, status.ErrorRate)
	}

	fmt.Printf("Agent management system initialized with %d agents\n", len(am.agents))
}
