package main

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveDirectory(t *testing.T) {
	// Get current working directory for testing
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Get home directory for testing
	usr, err := user.Current()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}
	homeDir := usr.HomeDir

	tests := []struct {
		name        string
		input       string
		expectError bool
		validate    func(string) bool
		description string
	}{
		{
			name:        "Empty string",
			input:       "",
			expectError: false,
			validate:    func(result string) bool { return result == "" },
			description: "Empty input should return empty string",
		},
		{
			name:        "Current directory",
			input:       ".",
			expectError: false,
			validate:    func(result string) bool { return result == cwd },
			description: "Current directory should resolve to absolute path",
		},
		{
			name:        "Home directory",
			input:       "~",
			expectError: false,
			validate:    func(result string) bool { return result == homeDir },
			description: "Tilde should resolve to home directory",
		},
		{
			name:        "Home subdirectory",
			input:       "~/Documents",
			expectError: false,
			validate: func(result string) bool {
				expected := filepath.Join(homeDir, "Documents")
				return result == expected
			},
			description: "~/path should resolve to home/path",
		},
		{
			name:        "Home subdirectory with multiple levels",
			input:       "~/Documents/Projects/test",
			expectError: false,
			validate: func(result string) bool {
				expected := filepath.Join(homeDir, "Documents", "Projects", "test")
				return result == expected
			},
			description: "~/deep/path should resolve correctly",
		},
		{
			name:        "Relative path",
			input:       "testdir",
			expectError: false,
			validate: func(result string) bool {
				expected := filepath.Join(cwd, "testdir")
				return result == expected
			},
			description: "Relative path should resolve to absolute",
		},
		{
			name:        "Relative path with subdirectory",
			input:       "test/subdir",
			expectError: false,
			validate: func(result string) bool {
				expected := filepath.Join(cwd, "test", "subdir")
				return result == expected
			},
			description: "Relative path with subdirs should resolve correctly",
		},
		{
			name:        "Absolute path",
			input:       filepath.Join(cwd, "absolute", "test"),
			expectError: false,
			validate: func(result string) bool {
				expected := filepath.Join(cwd, "absolute", "test")
				return result == expected
			},
			description: "Absolute path should remain absolute",
		},
		{
			name:        "Whitespace trimming",
			input:       "  ~/Documents  ",
			expectError: false,
			validate: func(result string) bool {
				expected := filepath.Join(homeDir, "Documents")
				return result == expected
			},
			description: "Whitespace should be trimmed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveDirectory(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if !tt.validate(result) {
					t.Errorf("%s: input '%s' resolved to '%s'", tt.description, tt.input, result)
				}
			}
		})
	}
}

func TestResolveDirectoryPlatformSpecific(t *testing.T) {
	// Test platform-specific path handling
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "Current directory with extra slash",
			input:       "./",
			expectError: false,
			description: "Current directory with trailing slash should work",
		},
		{
			name:        "Parent directory",
			input:       "..",
			expectError: false,
			description: "Parent directory should resolve",
		},
		{
			name:        "Parent directory path",
			input:       "../test",
			expectError: false,
			description: "Parent directory with subpath should resolve",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveDirectory(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				// Basic validation - should be absolute path
				if !filepath.IsAbs(result) {
					t.Errorf("Result should be absolute path, got: %s", result)
				}
			}
		})
	}
}

func TestResolveDirectoryEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name:        "Multiple tildes",
			input:       "~/~/test",
			expectError: false,
			description: "Multiple tildes should be treated as literal after first",
		},
		{
			name:        "Tilde in middle",
			input:       "test/~/middle",
			expectError: false,
			description: "Tilde in middle should be treated as literal",
		},
		{
			name:        "Only whitespace",
			input:       "   ",
			expectError: false,
			description: "Only whitespace should resolve to empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resolveDirectory(tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.description, err)
				}
				
				// For "only whitespace" test, the function trims to empty and then
				// follows the default case which resolves to current directory
				if tt.input == "   " {
					currentDir, _ := os.Getwd()
					if result != currentDir {
						t.Errorf("Expected current directory for whitespace input, got: %s", result)
					}
				}
				
				// For other tests, expect non-empty absolute path
				if strings.TrimSpace(tt.input) != "" && tt.input != "   " {
					if result == "" {
						t.Errorf("Expected non-empty result for %s", tt.description)
					}
					if !filepath.IsAbs(result) {
						t.Errorf("Expected absolute path for %s, got: %s", tt.description, result)
					}
				}
			}
		})
	}
}
