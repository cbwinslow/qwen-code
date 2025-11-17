package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== FILE SHARING TYPES ====================

// FileOperation represents different file operations
type FileOperation string

const (
	OperationUpload   FileOperation = "upload"
	OperationDownload FileOperation = "download"
	OperationShare    FileOperation = "share"
	OperationDelete   FileOperation = "delete"
	OperationRename   FileOperation = "rename"
	OperationMove     FileOperation = "move"
	OperationCopy     FileOperation = "copy"
)

// SharedFile represents a file shared in the chatroom
type SharedFile struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	Type        string    `json:"type"`
	MimeType    string    `json:"mime_type"`
	Owner       string    `json:"owner"`
	Permissions []string  `json:"permissions"`
	Tags        []string  `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	Checksum    string    `json:"checksum"`
	IsPublic    bool      `json:"is_public"`
	Downloads   int       `json:"downloads"`
	Description string    `json:"description,omitempty"`
}

// FilePermission represents file permissions
type FilePermission string

const (
	PermissionRead   FilePermission = "read"
	PermissionWrite  FilePermission = "write"
	PermissionDelete FilePermission = "delete"
	PermissionShare  FilePermission = "share"
	PermissionAdmin  FilePermission = "admin"
)

// FileCategory represents file categories
type FileCategory string

const (
	CategoryDocument FileCategory = "document"
	CategoryImage    FileCategory = "image"
	CategoryVideo    FileCategory = "video"
	CategoryAudio    FileCategory = "audio"
	CategoryCode     FileCategory = "code"
	CategoryData     FileCategory = "data"
	CategoryOther    FileCategory = "other"
)

// ==================== FILE MANAGER ====================

// FileManager manages file sharing and collaboration
type FileManager struct {
	sharedFiles  map[string]*SharedFile
	uploadDir    string
	maxFileSize  int64
	allowedTypes map[string]bool
	mu           sync.RWMutex
	eventHandler func(event FileEvent)
}

// FileEvent represents file-related events
type FileEvent struct {
	Type      string                 `json:"type"`
	FileID    string                 `json:"file_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    string                 `json:"user_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Message   string                 `json:"message,omitempty"`
}

// ==================== FILE MANAGER IMPLEMENTATION ====================

// NewFileManager creates a new file manager
func NewFileManager(uploadDir string) *FileManager {
	return &FileManager{
		sharedFiles: make(map[string]*SharedFile),
		uploadDir:   uploadDir,
		maxFileSize: 100 * 1024 * 1024, // 100MB
		allowedTypes: map[string]bool{
			".txt":  true,
			".md":   true,
			".json": true,
			".yaml": true,
			".yml":  true,
			".go":   true,
			".js":   true,
			".ts":   true,
			".py":   true,
			".html": true,
			".css":  true,
			".png":  true,
			".jpg":  true,
			".jpeg": true,
			".gif":  true,
			".pdf":  true,
			".zip":  true,
			".tar":  true,
			".gz":   true,
		},
		mu: sync.RWMutex{},
	}
}

// UploadFile handles file uploads
func (fm *FileManager) UploadFile(filePath string, owner string, permissions []string, isPublic bool) (*SharedFile, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Validate file
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Check file size
	if fileInfo.Size() > fm.maxFileSize {
		return nil, fmt.Errorf("file size %d exceeds maximum %d", fileInfo.Size(), fm.maxFileSize)
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(filePath))
	if !fm.allowedTypes[ext] {
		return nil, fmt.Errorf("file type %s is not allowed", ext)
	}

	// Calculate checksum
	checksum, err := fm.calculateChecksum(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// Copy file to upload directory
	fileName := filepath.Base(filePath)
	uploadPath := filepath.Join(fm.uploadDir, fileName)

	if err := fm.copyFile(filePath, uploadPath); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Determine MIME type
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	sharedFile := &SharedFile{
		ID:          generateID(),
		Name:        fileName,
		Path:        uploadPath,
		Size:        fileInfo.Size(),
		Type:        fm.getFileCategory(ext),
		MimeType:    mimeType,
		Owner:       owner,
		Permissions: permissions,
		Tags:        []string{},
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		Checksum:    checksum,
		IsPublic:    isPublic,
		Downloads:   0,
	}

	fm.sharedFiles[sharedFile.ID] = sharedFile

	if fm.eventHandler != nil {
		fm.eventHandler(FileEvent{
			Type:      "file_uploaded",
			FileID:    sharedFile.ID,
			Timestamp: time.Now(),
			UserID:    owner,
			Data: map[string]interface{}{
				"file": sharedFile,
			},
			Message: fmt.Sprintf("File %s uploaded by %s", fileName, owner),
		})
	}

	return sharedFile, nil
}

// DownloadFile handles file downloads
func (fm *FileManager) DownloadFile(fileID string) (string, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	file, exists := fm.sharedFiles[fileID]
	if !exists {
		return "", fmt.Errorf("file with ID %s not found", fileID)
	}

	// Increment download count
	file.Downloads++

	if fm.eventHandler != nil {
		fm.eventHandler(FileEvent{
			Type:      "file_downloaded",
			FileID:    fileID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"file": file,
			},
			Message: fmt.Sprintf("File %s downloaded", file.Name),
		})
	}

	return file.Path, nil
}

// ShareFile generates a shareable link for a file
func (fm *FileManager) ShareFile(fileID string, expires time.Duration) (string, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	file, exists := fm.sharedFiles[fileID]
	if !exists {
		return "", fmt.Errorf("file with ID %s not found", fileID)
	}

	// Generate share link (in a real implementation, this would be a URL)
	shareLink := fmt.Sprintf("https://chatroom.local/share/%s?expires=%d", fileID, expires.Seconds())

	if fm.eventHandler != nil {
		fm.eventHandler(FileEvent{
			Type:      "file_shared",
			FileID:    fileID,
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"file":       file,
				"share_link": shareLink,
				"expires":    expires,
			},
			Message: fmt.Sprintf("File %s shared", file.Name),
		})
	}

	return shareLink, nil
}

// DeleteFile removes a file
func (fm *FileManager) DeleteFile(fileID string, userID string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	file, exists := fm.sharedFiles[fileID]
	if !exists {
		return fmt.Errorf("file with ID %s not found", fileID)
	}

	// Check permissions
	if !fm.hasPermission(file.Permissions, userID, PermissionDelete) {
		return fmt.Errorf("user %s does not have delete permission", userID)
	}

	// Remove file from filesystem
	if err := os.Remove(file.Path); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Remove from shared files
	delete(fm.sharedFiles, fileID)

	if fm.eventHandler != nil {
		fm.eventHandler(FileEvent{
			Type:      "file_deleted",
			FileID:    fileID,
			Timestamp: time.Now(),
			UserID:    userID,
			Data: map[string]interface{}{
				"file": file,
			},
			Message: fmt.Sprintf("File %s deleted by %s", file.Name, userID),
		})
	}

	return nil
}

// ListFiles returns all shared files
func (fm *FileManager) ListFiles(filter map[string]interface{}) ([]*SharedFile, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	var files []*SharedFile
	for _, file := range fm.sharedFiles {
		// Apply filters
		if filter != nil {
			if fileType, ok := filter["type"]; ok && fileType != "" {
				if file.Type != fileType {
					continue
				}
			}
			if owner, ok := filter["owner"]; ok && owner != "" {
				if file.Owner != owner {
					continue
				}
			}
			if isPublic, ok := filter["is_public"]; ok {
				if file.IsPublic != isPublic {
					continue
				}
			}
		}
		files = append(files, file)
	}

	return files, nil
}

// GetFile returns a specific file
func (fm *FileManager) GetFile(fileID string) (*SharedFile, error) {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	file, exists := fm.sharedFiles[fileID]
	if !exists {
		return nil, fmt.Errorf("file with ID %s not found", fileID)
	}

	return file, nil
}

// ==================== HELPER METHODS ====================

// calculateChecksum calculates MD5 checksum of a file
func (fm *FileManager) calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// copyFile copies a file with progress tracking
func (fm *FileManager) copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

// getFileCategory determines file category from extension
func (fm *FileManager) getFileCategory(ext string) FileCategory {
	switch ext {
	case ".txt", ".md", ".doc", ".docx", ".pdf":
		return CategoryDocument
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg":
		return CategoryImage
	case ".mp4", ".avi", ".mov", ".wmv", ".flv":
		return CategoryVideo
	case ".mp3", ".wav", ".flac", ".ogg":
		return CategoryAudio
	case ".go", ".js", ".ts", ".py", ".java", ".cpp", ".c", ".h":
		return CategoryCode
	case ".json", ".yaml", ".yml", ".csv", ".xml":
		return CategoryData
	default:
		return CategoryOther
	}
}

// hasPermission checks if a user has a specific permission
func (fm *FileManager) hasPermission(permissions []string, userID string, permission FilePermission) bool {
	// Check if user is owner
	for _, file := range fm.sharedFiles {
		if file.Owner == userID {
			return true // Owner has all permissions
		}
	}

	// Check specific permission
	for _, perm := range permissions {
		if perm == permission {
			return true
		}
	}

	return false
}

// SetEventHandler sets the event handler
func (fm *FileManager) SetEventHandler(handler func(event FileEvent)) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.eventHandler = handler
}

// GetStats returns file manager statistics
func (fm *FileManager) GetStats() map[string]interface{} {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	totalSize := int64(0)
	totalFiles := len(fm.sharedFiles)
	typeCounts := make(map[string]int)
	downloads := 0

	for _, file := range fm.sharedFiles {
		totalSize += file.Size
		downloads += file.Downloads
		typeCounts[file.Type]++
	}

	return map[string]interface{}{
		"total_files":   totalFiles,
		"total_size":    totalSize,
		"downloads":     downloads,
		"file_types":    typeCounts,
		"max_file_size": fm.maxFileSize,
		"upload_dir":    fm.uploadDir,
	}
}

// ==================== COLLABORATION FEATURES ====================

// CollabSession represents a collaborative editing session
type CollabSession struct {
	ID           string                 `json:"id"`
	FileID       string                 `json:"file_id"`
	Participants []string               `json:"participants"`
	IsActive     bool                   `json:"is_active"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Content      string                 `json:"content"`
	Version      int                    `json:"version"`
	Cursor       map[string]interface{} `json:"cursor,omitempty"`
	Changes      []CollabChange         `json:"changes"`
}

// CollabChange represents a change in collaborative editing
type CollabChange struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Type       string                 `json:"type"` // "insert", "delete", "replace", "format"
	Position   map[string]interface{} `json:"position"`
	Content    string                 `json:"content,omitempty"`
	OldContent string                 `json:"old_content,omitempty"`
	NewContent string                 `json:"new_content,omitempty"`
}

// CollabManager manages collaborative editing sessions
type CollabManager struct {
	sessions     map[string]*CollabSession
	mu           sync.RWMutex
	eventHandler func(event FileEvent)
}

// NewCollabManager creates a new collaboration manager
func NewCollabManager() *CollabManager {
	return &CollabManager{
		sessions: make(map[string]*CollabSession),
		mu:       sync.RWMutex{},
	}
}

// CreateSession creates a new collaborative editing session
func (cm *CollabManager) CreateSession(fileID string, participants []string, initialContent string) (*CollabSession, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	session := &CollabSession{
		ID:           generateID(),
		FileID:       fileID,
		Participants: participants,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Content:      initialContent,
		Version:      1,
		Cursor:       make(map[string]interface{}),
		Changes:      []CollabChange{},
	}

	cm.sessions[session.ID] = session

	if cm.eventHandler != nil {
		cm.eventHandler(FileEvent{
			Type:      "collab_session_created",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"session": session,
			},
			Message: fmt.Sprintf("Collaborative session created for file %s", fileID),
		})
	}

	return session, nil
}

// JoinSession adds a participant to a collaborative session
func (cm *CollabManager) JoinSession(sessionID string, userID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	session, exists := cm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	// Check if user is already a participant
	for _, participant := range session.Participants {
		if participant == userID {
			return fmt.Errorf("user %s is already in session %s", userID, sessionID)
		}
	}

	session.Participants = append(session.Participants, userID)
	session.UpdatedAt = time.Now()

	if cm.eventHandler != nil {
		cm.eventHandler(FileEvent{
			Type:      "user_joined_session",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"session_id": sessionID,
				"user_id":    userID,
			},
			Message: fmt.Sprintf("User %s joined session %s", userID, sessionID),
		})
	}

	return nil
}

// ApplyChange applies a change to a collaborative session
func (cm *CollabManager) ApplyChange(sessionID string, change CollabChange) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	session, exists := cm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session %s not found", sessionID)
	}

	change.ID = generateID()
	change.Timestamp = time.Now()
	session.Changes = append(session.Changes, change)
	session.Version++
	session.UpdatedAt = time.Now()

	if cm.eventHandler != nil {
		cm.eventHandler(FileEvent{
			Type:      "change_applied",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"session_id": sessionID,
				"change":     change,
			},
			Message: fmt.Sprintf("Change applied to session %s", sessionID),
		})
	}

	return nil
}

// GetSession returns a collaborative session
func (cm *CollabManager) GetSession(sessionID string) (*CollabSession, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	session, exists := cm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session %s not found", sessionID)
	}

	return session, nil
}

// GetActiveSessions returns all active collaborative sessions
func (cm *CollabManager) GetActiveSessions() []*CollabSession {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var active []*CollabSession
	for _, session := range cm.sessions {
		if session.IsActive {
			active = append(active, session)
		}
	}

	return active
}

// SetEventHandler sets the event handler for collaboration events
func (cm *CollabManager) SetEventHandler(handler func(event FileEvent)) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.eventHandler = handler
}

// ==================== MAIN FUNCTION ====================

// main for testing file sharing and collaboration
func main() {
	fmt.Println("📁 Testing File Sharing and Collaboration")

	// Create upload directory
	uploadDir := "/tmp/chatroom-uploads"
	os.MkdirAll(uploadDir, 0755)

	// Test file manager
	fileManager := NewFileManager(uploadDir)

	// Set event handler
	fileManager.SetEventHandler(func(event FileEvent) {
		fmt.Printf("📁 File Event: %s - %s\n", event.Type, event.Message)
		if event.Data != nil {
			data, _ := json.MarshalIndent(event.Data, "", "  ")
			fmt.Printf("   Data: %s\n", string(data))
		}
	})

	// Test file upload
	testFile := "/tmp/test-document.txt"
	testContent := "This is a test document for file sharing."

	err := os.WriteFile(testFile, testContent)
	if err != nil {
		fmt.Printf("❌ Failed to create test file: %v\n", err)
		return
	}

	// Upload file
	sharedFile, err := fileManager.UploadFile(testFile, "test-user", []string{PermissionRead, PermissionWrite}, false)
	if err != nil {
		fmt.Printf("❌ Failed to upload file: %v\n", err)
		return
	}

	fmt.Printf("✅ File uploaded successfully: %s\n", sharedFile.Name)

	// Test file listing
	files, err := fileManager.ListFiles(nil)
	if err != nil {
		fmt.Printf("❌ Failed to list files: %v\n", err)
		return
	}

	fmt.Printf("✅ Found %d shared files\n", len(files))
	for _, file := range files {
		fmt.Printf("  - %s (%s, %d bytes)\n", file.Name, file.Type, file.Size)
	}

	// Test file download
	downloadPath, err := fileManager.DownloadFile(sharedFile.ID)
	if err != nil {
		fmt.Printf("❌ Failed to download file: %v\n", err)
		return
	}

	fmt.Printf("✅ File downloaded to: %s\n", downloadPath)

	// Test file sharing
	shareLink, err := fileManager.ShareFile(sharedFile.ID, 24*time.Hour)
	if err != nil {
		fmt.Printf("❌ Failed to share file: %v\n", err)
		return
	}

	fmt.Printf("✅ File share link: %s\n", shareLink)

	// Test collaboration
	collabManager := NewCollabManager()
	collabManager.SetEventHandler(func(event FileEvent) {
		fmt.Printf("🤝 Collab Event: %s - %s\n", event.Type, event.Message)
	})

	session, err := collabManager.CreateSession(sharedFile.ID, []string{"test-user"}, testContent)
	if err != nil {
		fmt.Printf("❌ Failed to create collab session: %v\n", err)
		return
	}

	fmt.Printf("✅ Collab session created: %s\n", session.ID)

	// Test joining session
	err = collabManager.JoinSession(session.ID, "test-user-2")
	if err != nil {
		fmt.Printf("❌ Failed to join session: %v\n", err)
		return
	}

	fmt.Printf("✅ User joined session\n")

	// Test applying change
	change := CollabChange{
		Type:    "insert",
		Content: "New content added",
		Position: map[string]interface{}{
			"line":   1,
			"column": 10,
		},
	}

	err = collabManager.ApplyChange(session.ID, change)
	if err != nil {
		fmt.Printf("❌ Failed to apply change: %v\n", err)
		return
	}

	fmt.Printf("✅ Change applied to session\n")

	// Get session info
	retrievedSession, err := collabManager.GetSession(session.ID)
	if err != nil {
		fmt.Printf("❌ Failed to get session: %v\n", err)
		return
	}

	fmt.Printf("✅ Session retrieved with %d changes\n", len(retrievedSession.Changes))

	// Get stats
	stats := fileManager.GetStats()
	statsJSON, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Printf("📊 File Manager Stats:\n%s\n", string(statsJSON))

	fmt.Println("🎉 File sharing and collaboration test completed successfully!")
}
