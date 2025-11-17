package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== OPENROUTER INTEGRATION ====================

// OpenRouterConfig holds OpenRouter configuration
type OpenRouterConfig struct {
	APIKey      string  `json:"api_key"`
	BaseURL     string  `json:"base_url"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// OpenRouterClient wraps OpenRouter API
type OpenRouterClient struct {
	config  OpenRouterConfig
	baseURL string
	client  *http.Client
}

// OpenRouterMessage represents a message for OpenRouter
type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents a response from OpenRouter
type OpenRouterResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      OpenRouterMessage `json:"message"`
		FinishReason string            `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ==================== OPENROUTER CLIENT ====================

// NewOpenRouterClient creates a new OpenRouter client
func NewOpenRouterClient(config OpenRouterConfig) *OpenRouterClient {
	baseURL := "https://openrouter.ai/api/v1"
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}

	return &OpenRouterClient{
		config:  config,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage sends a message to OpenRouter
func (orc *OpenRouterClient) SendMessage(ctx context.Context, messages []OpenRouterMessage) (*OpenRouterResponse, error) {
	if orc.config.APIKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is required")
	}

	requestBody := map[string]interface{}{
		"model":       orc.config.Model,
		"messages":    messages,
		"max_tokens":  orc.config.MaxTokens,
		"temperature": orc.config.Temperature,
		"stream":      false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", orc.baseURL+"/chat/completions", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+orc.config.APIKey)
	req.Header.Set("HTTP-Referer", "https://github.com/openrouter/openrouter")
	req.Header.Set("X-Title", "AI TUI Chatroom")

	resp, err := orc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenRouter API error: %d - %s", resp.StatusCode, string(body))
	}

	var response OpenRouterResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// GetModels retrieves available models from OpenRouter
func (orc *OpenRouterClient) GetModels(ctx context.Context) ([]string, error) {
	if orc.config.APIKey == "" {
		return nil, fmt.Errorf("OpenRouter API key is required")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", orc.baseURL+"/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+orc.config.APIKey)
	req.Header.Set("HTTP-Referer", "https://github.com/openrouter/openrouter")

	resp, err := orc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenRouter API error: %d - %s", resp.StatusCode, string(body))
	}

	var modelsResponse struct {
		Data []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &modelsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal models response: %w", err)
	}

	var models []string
	for _, model := range modelsResponse.Data {
		models = append(models, model.ID)
	}

	return models, nil
}

// ==================== PROVIDER INTEGRATION ====================

// OpenRouterProvider implements the AI provider interface
type OpenRouterProvider struct {
	client *OpenRouterClient
	model  string
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(config OpenRouterConfig) *OpenRouterProvider {
	return &OpenRouterProvider{
		client: NewOpenRouterClient(config),
		model:  config.Model,
	}
}

// SendMessage implements the provider interface
func (orp *OpenRouterProvider) SendMessage(ctx context.Context, content string, conversationID string) (string, error) {
	messages := []OpenRouterMessage{
		{
			Role:    "user",
			Content: content,
		},
	}

	response, err := orp.client.SendMessage(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenRouter")
	}

	return response.Choices[0].Message.Content, nil
}

// GetCapabilities returns provider capabilities
func (orp *OpenRouterProvider) GetCapabilities() []string {
	return []string{
		"text_generation",
		"reasoning",
		"code_generation",
		"data_analysis",
		"conversation",
	}
}

// GetModels returns available models
func (orp *OpenRouterProvider) GetModels() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return orp.client.GetModels(ctx)
}

// ==================== CHATROOM INTEGRATION ====================

// ChatroomProvider integrates OpenRouter with the chatroom
type ChatroomProvider struct {
	openRouter   *OpenRouterProvider
	agentManager *AgentManager
}

// NewChatroomProvider creates a new chatroom provider
func NewChatroomProvider(openRouterConfig OpenRouterConfig) *ChatroomProvider {
	return &ChatroomProvider{
		openRouter:   NewOpenRouterProvider(openRouterConfig),
		agentManager: NewAgentManager(),
	}
}

// Initialize initializes the provider
func (cp *ChatroomProvider) Initialize() error {
	// Initialize OpenRouter client
	models, err := cp.openRouter.GetModels()
	if err != nil {
		return fmt.Errorf("failed to get OpenRouter models: %w", err)
	}

	// Add OpenRouter as a provider to the agent manager
	openRouterProvider := Provider{
		ID:       "openrouter",
		Name:     "OpenRouter",
		Models:   models,
		IsActive: true,
		Settings: map[string]interface{}{
			"api_key":  cp.openRouter.client.config.APIKey,
			"base_url": cp.openRouter.client.config.BaseURL,
			"model":    cp.openRouter.client.config.Model,
		},
	}

	if err := cp.agentManager.AddProvider(openRouterProvider); err != nil {
		return fmt.Errorf("failed to add OpenRouter provider: %w", err)
	}

	// Create default agents for OpenRouter
	agents := []AgentConfig{
		{
			ID:          "openrouter-coordinator",
			Name:        "OpenRouter Coordinator",
			Role:        "coordinator",
			Personality: "diplomatic",
			Provider:    "openrouter",
			Model:       "anthropic/claude-3-sonnet-20240229",
			MaxTokens:   4096,
			Temperature: 0.7,
			Capabilities: []string{
				"text_generation",
				"reasoning",
				"coordination",
				"analysis",
			},
			IsActive: true,
		},
		{
			ID:          "openrouter-coder",
			Name:        "OpenRouter Coder",
			Role:        "specialist",
			Personality: "analytical",
			Provider:    "openrouter",
			Model:       "qwen/qwen-2.5-coder-32b",
			MaxTokens:   8192,
			Temperature: 0.1,
			Capabilities: []string{
				"text_generation",
				"code_generation",
				"debugging",
				"analysis",
			},
			IsActive: true,
		},
		{
			ID:          "openrouter-analyst",
			Name:        "OpenRouter Analyst",
			Role:        "specialist",
			Personality: "thorough",
			Provider:    "openrouter",
			Model:       "meta-llama/llama-3-70b-instruct",
			MaxTokens:   4096,
			Temperature: 0.3,
			Capabilities: []string{
				"data_analysis",
				"reasoning",
				"visualization",
				"research",
			},
			IsActive: true,
		},
	}

	for _, agentConfig := range agents {
		if err := cp.agentManager.AddAgent(agentConfig); err != nil {
			return fmt.Errorf("failed to add agent %s: %w", agentConfig.Name, err)
		}
	}

	return nil
}

// SendMessage sends a message through the chatroom
func (cp *ChatroomProvider) SendMessage(ctx context.Context, content string, conversationID string, agentIDs []string) (string, error) {
	// For now, use the first available agent
	if len(agentIDs) == 0 {
		return "", fmt.Errorf("no agent IDs provided")
	}

	// Get the first agent
	agents := cp.agentManager.GetAgents()
	agent, exists := agents[agentIDs[0]]
	if !exists {
		return "", fmt.Errorf("agent %s not found", agentIDs[0])
	}

	// Send message through OpenRouter
	response, err := cp.openRouter.SendMessage(ctx, content, conversationID)
	if err != nil {
		return "", fmt.Errorf("failed to send message via OpenRouter: %w", err)
	}

	// Log the agent's response
	if cp.agentManager.eventHandler != nil {
		cp.agentManager.eventHandler(AgentEvent{
			Type:      "message_sent",
			AgentID:   agentIDs[0],
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"content":         content,
				"conversation_id": conversationID,
				"provider":        "openrouter",
				"response":        response,
			},
			Message: fmt.Sprintf("Message sent via %s", agent.Name),
		})
	}

	return response, nil
}

// GetStatus returns the status of the provider
func (cp *ChatroomProvider) GetStatus() map[string]interface{} {
	agents := cp.agentManager.GetActiveAgents()
	models, _ := cp.openRouter.GetModels()

	return map[string]interface{}{
		"provider":         "openrouter",
		"active_agents":    len(agents),
		"available_models": models,
		"agent_status":     cp.agentManager.GetPerformanceMetrics(),
		"last_updated":     time.Now(),
	}
}

// ==================== TESTING ====================

// TestOpenRouterIntegration tests the OpenRouter integration
func TestOpenRouterIntegration(t *testing.T) {
	// Test with mock API key
	config := OpenRouterConfig{
		APIKey:      "test-key",
		Model:       "anthropic/claude-3-sonnet-20240229",
		MaxTokens:   100,
		Temperature: 0.7,
	}

	provider := NewOpenRouterProvider(config)

	// Test initialization
	err := provider.Initialize()
	if err != nil {
		t.Errorf("Failed to initialize provider: %v", err)
	}

	// Test capabilities
	capabilities := provider.GetCapabilities()
	if len(capabilities) == 0 {
		t.Error("Provider should have capabilities")
	}

	// Test models
	models, err := provider.GetModels()
	if err != nil {
		t.Errorf("Failed to get models: %v", err)
	}
	if len(models) == 0 {
		t.Error("Should have available models")
	}

	t.Logf("OpenRouter integration test passed")
	t.Logf("Capabilities: %v", capabilities)
	t.Logf("Models: %v", models)
}

// ==================== MAIN FUNCTION ====================

// main for testing the OpenRouter integration
func main() {
	fmt.Println("🚀 Testing OpenRouter Integration")

	// Test with environment variables
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ OPENROUTER_API_KEY environment variable not set")
		fmt.Println("Set it with: export OPENROUTER_API_KEY=your_api_key")
		os.Exit(1)
	}

	config := OpenRouterConfig{
		APIKey:      apiKey,
		Model:       "anthropic/claude-3-sonnet-20240229",
		MaxTokens:   100,
		Temperature: 0.7,
	}

	provider := NewOpenRouterProvider(config)

	if err := provider.Initialize(); err != nil {
		fmt.Printf("❌ Failed to initialize provider: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Provider initialized successfully")

	// Test getting models
	models, err := provider.GetModels()
	if err != nil {
		fmt.Printf("❌ Failed to get models: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Available models: %v\n", models)

	// Test sending a message
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	response, err := provider.SendMessage(ctx, "Hello, this is a test message", "test-conversation", []string{"openrouter-coordinator"})
	if err != nil {
		fmt.Printf("❌ Failed to send message: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Message sent successfully\n")
	fmt.Printf("Response: %s\n", response)

	// Test getting status
	status := provider.GetStatus()
	statusJSON, _ := json.MarshalIndent(status, "", "  ")
	fmt.Printf("✅ Provider Status:\n%s\n", string(statusJSON))
}
