package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.etcd.io/bbolt"
)

// createTempDB creates a temporary database for testing
func createTempDB(t *testing.T) (*Database, string) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "afvikle_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create database directly in temp directory
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create database: %v", err)
	}
	
	database := &Database{db: db}
	
	// Initialize buckets
	if err := database.initBuckets(); err != nil {
		db.Close()
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to initialize buckets: %v", err)
	}
	
	return database, tempDir
}

func TestNewDatabase(t *testing.T) {
	db, tempDir := createTempDB(t)
	defer func() {
		db.Close()
		os.RemoveAll(tempDir)
	}()

	if db == nil {
		t.Fatal("Database should not be nil")
	}

	if db.db == nil {
		t.Fatal("Database connection should not be nil")
	}
}

func TestAddCommand(t *testing.T) {
	db, tempDir := createTempDB(t)
	defer func() {
		db.Close()
		os.RemoveAll(tempDir)
	}()

	tests := []struct {
		name        string
		cmdName     string
		description string
		command     string
		workingDir  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid command",
			cmdName:     "test-cmd",
			description: "Test command",
			command:     "echo hello",
			workingDir:  "",
			expectError: false,
		},
		{
			name:        "Valid command with working directory",
			cmdName:     "test-cmd-dir",
			description: "Test command with dir",
			command:     "ls -la",
			workingDir:  tempDir,
			expectError: false,
		},
		{
			name:        "Empty name",
			cmdName:     "",
			description: "Test",
			command:     "echo test",
			workingDir:  "",
			expectError: true,
			errorMsg:    "command name is required",
		},
		{
			name:        "Empty command",
			cmdName:     "test",
			description: "Test",
			command:     "",
			workingDir:  "",
			expectError: true,
			errorMsg:    "command is required",
		},
		{
			name:        "Invalid working directory",
			cmdName:     "test-invalid-dir",
			description: "Test",
			command:     "echo test",
			workingDir:  "/nonexistent/directory",
			expectError: true,
			errorMsg:    "working directory '/nonexistent/directory' does not exist",
		},
		{
			name:        "Duplicate command name",
			cmdName:     "test-cmd", // Same as first test
			description: "Duplicate",
			command:     "echo duplicate",
			workingDir:  "",
			expectError: true,
			errorMsg:    "command 'test-cmd' already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddCommand(tt.cmdName, tt.description, tt.command, tt.workingDir)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetCommand(t *testing.T) {
	db, tempDir := createTempDB(t)
	defer func() {
		db.Close()
		os.RemoveAll(tempDir)
	}()

	// Add a test command
	err := db.AddCommand("get-test", "Get test command", "echo get-test", tempDir)
	if err != nil {
		t.Fatalf("Failed to add test command: %v", err)
	}

	tests := []struct {
		name        string
		cmdName     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Existing command",
			cmdName:     "get-test",
			expectError: false,
		},
		{
			name:        "Non-existing command",
			cmdName:     "non-existent",
			expectError: true,
			errorMsg:    "command 'non-existent' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := db.GetCommand(tt.cmdName)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if cmd == nil {
					t.Errorf("Command should not be nil")
				} else if cmd.Name != tt.cmdName {
					t.Errorf("Expected command name '%s', got '%s'", tt.cmdName, cmd.Name)
				}
			}
		})
	}
}

func TestGetAllCommands(t *testing.T) {
	db, tempDir := createTempDB(t)
	defer func() {
		db.Close()
		os.RemoveAll(tempDir)
	}()

	// Initially should be empty
	commands, err := db.GetAllCommands()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(commands) != 0 {
		t.Errorf("Expected 0 commands, got %d", len(commands))
	}

	// Add some commands
	testCommands := []struct {
		name        string
		description string
		command     string
		workingDir  string
	}{
		{"cmd1", "Command 1", "echo 1", ""},
		{"cmd2", "Command 2", "echo 2", tempDir},
		{"cmd3", "Command 3", "echo 3", ""},
	}

	for _, tc := range testCommands {
		err := db.AddCommand(tc.name, tc.description, tc.command, tc.workingDir)
		if err != nil {
			t.Fatalf("Failed to add command '%s': %v", tc.name, err)
		}
	}

	// Get all commands
	commands, err = db.GetAllCommands()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(commands) != len(testCommands) {
		t.Errorf("Expected %d commands, got %d", len(testCommands), len(commands))
	}

	// Verify commands
	commandMap := make(map[string]Command)
	for _, cmd := range commands {
		commandMap[cmd.Name] = cmd
	}

	for _, tc := range testCommands {
		cmd, exists := commandMap[tc.name]
		if !exists {
			t.Errorf("Command '%s' not found in results", tc.name)
			continue
		}
		if cmd.Description != tc.description {
			t.Errorf("Expected description '%s', got '%s'", tc.description, cmd.Description)
		}
		if cmd.Command != tc.command {
			t.Errorf("Expected command '%s', got '%s'", tc.command, cmd.Command)
		}
		if cmd.WorkingDir != tc.workingDir {
			t.Errorf("Expected working dir '%s', got '%s'", tc.workingDir, cmd.WorkingDir)
		}
	}
}

func TestUpdateCommand(t *testing.T) {
	db, tempDir := createTempDB(t)
	defer func() {
		db.Close()
		os.RemoveAll(tempDir)
	}()

	// Add a command to update
	err := db.AddCommand("update-test", "Original description", "echo original", "")
	if err != nil {
		t.Fatalf("Failed to add test command: %v", err)
	}

	tests := []struct {
		name        string
		cmdName     string
		description string
		command     string
		workingDir  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid update",
			cmdName:     "update-test",
			description: "Updated description",
			command:     "echo updated",
			workingDir:  tempDir,
			expectError: false,
		},
		{
			name:        "Update non-existing command",
			cmdName:     "non-existent",
			description: "Test",
			command:     "echo test",
			workingDir:  "",
			expectError: true,
			errorMsg:    "command 'non-existent' not found",
		},
		{
			name:        "Empty command name",
			cmdName:     "",
			description: "Test",
			command:     "echo test",
			workingDir:  "",
			expectError: true,
			errorMsg:    "command name is required",
		},
		{
			name:        "Empty command",
			cmdName:     "update-test",
			description: "Test",
			command:     "",
			workingDir:  "",
			expectError: true,
			errorMsg:    "command is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.UpdateCommand(tt.cmdName, tt.description, tt.command, tt.workingDir)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				
				// Verify the update
				cmd, err := db.GetCommand(tt.cmdName)
				if err != nil {
					t.Errorf("Failed to get updated command: %v", err)
				} else {
					if cmd.Description != tt.description {
						t.Errorf("Expected description '%s', got '%s'", tt.description, cmd.Description)
					}
					if cmd.Command != tt.command {
						t.Errorf("Expected command '%s', got '%s'", tt.command, cmd.Command)
					}
					if cmd.WorkingDir != tt.workingDir {
						t.Errorf("Expected working dir '%s', got '%s'", tt.workingDir, cmd.WorkingDir)
					}
				}
			}
		})
	}
}

func TestDeleteCommand(t *testing.T) {
	db, tempDir := createTempDB(t)
	defer func() {
		db.Close()
		os.RemoveAll(tempDir)
	}()

	// Add a command to delete
	err := db.AddCommand("delete-test", "Delete test command", "echo delete", "")
	if err != nil {
		t.Fatalf("Failed to add test command: %v", err)
	}

	tests := []struct {
		name        string
		cmdName     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Delete existing command",
			cmdName:     "delete-test",
			expectError: false,
		},
		{
			name:        "Delete non-existing command",
			cmdName:     "non-existent",
			expectError: true,
			errorMsg:    "command 'non-existent' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.DeleteCommand(tt.cmdName)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				
				// Verify deletion
				_, err := db.GetCommand(tt.cmdName)
				if err == nil {
					t.Errorf("Command should have been deleted")
				}
			}
		})
	}
}

func TestCommandFields(t *testing.T) {
	db, tempDir := createTempDB(t)
	defer func() {
		db.Close()
		os.RemoveAll(tempDir)
	}()

	// Test default description
	err := db.AddCommand("test-default", "", "echo test", "")
	if err != nil {
		t.Fatalf("Failed to add command: %v", err)
	}

	cmd, err := db.GetCommand("test-default")
	if err != nil {
		t.Fatalf("Failed to get command: %v", err)
	}

	if cmd.Description != "No description provided" {
		t.Errorf("Expected default description 'No description provided', got '%s'", cmd.Description)
	}

	// Test CreatedAt field
	if cmd.CreatedAt == "" {
		t.Errorf("CreatedAt should not be empty")
	}

	// Parse time to verify format
	_, err = time.Parse("2006-01-02 15:04:05", cmd.CreatedAt)
	if err != nil {
		t.Errorf("CreatedAt has invalid format: %v", err)
	}

	// Test whitespace trimming
	err = db.AddCommand("  trim-test  ", "  trim description  ", "  echo trim  ", "")
	if err != nil {
		t.Fatalf("Failed to add command: %v", err)
	}

	cmd, err = db.GetCommand("trim-test")
	if err != nil {
		t.Fatalf("Failed to get command: %v", err)
	}

	if cmd.Name != "trim-test" {
		t.Errorf("Expected trimmed name 'trim-test', got '%s'", cmd.Name)
	}
	if cmd.Description != "trim description" {
		t.Errorf("Expected trimmed description 'trim description', got '%s'", cmd.Description)
	}
	if cmd.Command != "echo trim" {
		t.Errorf("Expected trimmed command 'echo trim', got '%s'", cmd.Command)
	}
}

func TestGetDatabasePath(t *testing.T) {
	db, tempDir := createTempDB(t)
	defer func() {
		db.Close()
		os.RemoveAll(tempDir)
	}()

	// Since we can't mock os.Executable, we'll test that the function works
	// by calling it and verifying it returns a valid path
	path, err := db.GetDatabasePath()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if path == "" {
		t.Errorf("Database path should not be empty")
	}

	// Verify it ends with afvikle.db
	if !filepath.IsAbs(path) {
		t.Errorf("Database path should be absolute, got: %s", path)
	}

	if filepath.Base(path) != "afvikle.db" {
		t.Errorf("Database path should end with 'afvikle.db', got: %s", path)
	}
}
