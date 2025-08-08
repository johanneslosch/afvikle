package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLICommands tests the CLI functionality through subprocess calls
func TestCLICommands(t *testing.T) {
	// Build the binary for testing
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)
	
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "afvikle_cli_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Copy binary to temp directory to test database creation
	testBinary := filepath.Join(tempDir, "afvikle"+filepath.Ext(binaryPath))
	copyFile(t, binaryPath, testBinary)
	
	t.Run("Help Command", func(t *testing.T) {
		testHelpCommand(t, testBinary)
	})
	
	t.Run("Info Command Empty", func(t *testing.T) {
		testInfoCommandEmpty(t, testBinary, tempDir)
	})
	
	t.Run("List Command Empty", func(t *testing.T) {
		testListCommandEmpty(t, testBinary)
	})
	
	t.Run("Add Command", func(t *testing.T) {
		testAddCommand(t, testBinary, tempDir)
	})
	
	t.Run("List Command With Data", func(t *testing.T) {
		testListCommandWithData(t, testBinary)
	})
	
	t.Run("Run Command", func(t *testing.T) {
		testRunCommand(t, testBinary)
	})
	
	t.Run("Delete Command", func(t *testing.T) {
		testDeleteCommand(t, testBinary)
	})
	
	t.Run("Delete All Commands", func(t *testing.T) {
		testDeleteAllCommands(t, testBinary, tempDir)
	})
	
	t.Run("Error Cases", func(t *testing.T) {
		testErrorCases(t, testBinary)
	})
}

func buildTestBinary(t *testing.T) string {
	// Build the binary in a temporary location
	binaryName := "afvikle_test"
	if os.Getenv("GOOS") == "windows" || (os.Getenv("GOOS") == "" && os.PathSeparator == '\\') {
		binaryName += ".exe"
	}
	
	tempDir, err := os.MkdirTemp("", "afvikle_build_*")
	if err != nil {
		t.Fatalf("Failed to create temp build directory: %v", err)
	}
	
	binaryPath := filepath.Join(tempDir, binaryName)
	
	cmd := exec.Command("go", "build", "-o", binaryPath)
	cmd.Dir = "." // Current directory
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
	}
	
	return binaryPath
}

func copyFile(t *testing.T, src, dst string) {
	input, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("Failed to read source file: %v", err)
	}
	
	err = os.WriteFile(dst, input, 0755)
	if err != nil {
		t.Fatalf("Failed to write destination file: %v", err)
	}
}

func runCommand(t *testing.T, binary string, args ...string) (string, string, error) {
	cmd := exec.Command(binary, args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func runCommandWithInput(t *testing.T, binary string, input string, args ...string) (string, string, error) {
	cmd := exec.Command(binary, args...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = strings.NewReader(input)
	
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func testHelpCommand(t *testing.T, binary string) {
	stdout, stderr, err := runCommand(t, binary, "--help")
	if err != nil {
		t.Errorf("Help command failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "afv") {
		t.Errorf("Help output should contain 'afv', got: %s", stdout)
	}
	
	if !strings.Contains(stdout, "Available commands") {
		t.Errorf("Help output should contain 'Available commands', got: %s", stdout)
	}
}

func testInfoCommandEmpty(t *testing.T, binary string, tempDir string) {
	stdout, stderr, err := runCommand(t, binary, "info")
	if err != nil {
		t.Errorf("Info command failed: %v\nStderr: %s", err, stderr)
	}
	
	expectedPath := filepath.Join(tempDir, "afvikle.db")
	if !strings.Contains(stdout, expectedPath) {
		t.Errorf("Info output should contain database path '%s', got: %s", expectedPath, stdout)
	}
	
	if !strings.Contains(stdout, "Total commands: 0") {
		t.Errorf("Info output should show 0 commands, got: %s", stdout)
	}
}

func testListCommandEmpty(t *testing.T, binary string) {
	stdout, stderr, err := runCommand(t, binary, "list")
	if err != nil {
		t.Errorf("List command failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "No commands found") {
		t.Errorf("List output should indicate no commands found, got: %s", stdout)
	}
}

func testAddCommand(t *testing.T, binary string, tempDir string) {
	// Test basic add command
	stdout, stderr, err := runCommand(t, binary, "add", "--name", "test-cmd", "--desc", "Test command", "--cmd", "echo hello")
	if err != nil {
		t.Errorf("Add command failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "Command 'test-cmd' added successfully") {
		t.Errorf("Add output should confirm success, got: %s", stdout)
	}
	
	// Test add command with working directory
	stdout, stderr, err = runCommand(t, binary, "add", "--name", "test-cmd-dir", "--desc", "Test with dir", "--cmd", "echo hello", "--dir", tempDir)
	if err != nil {
		t.Errorf("Add command with dir failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "Command 'test-cmd-dir' added successfully") {
		t.Errorf("Add with dir output should confirm success, got: %s", stdout)
	}
	
	if !strings.Contains(stdout, fmt.Sprintf("Working directory: %s", tempDir)) {
		t.Errorf("Add with dir should show working directory, got: %s", stdout)
	}
	
	// Test add command with current directory shortcut
	stdout, stderr, err = runCommand(t, binary, "add", "--name", "test-cmd-current", "--desc", "Test current dir", "--cmd", "echo current", "--dir", ".")
	if err != nil {
		t.Errorf("Add command with current dir failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "Command 'test-cmd-current' added successfully") {
		t.Errorf("Add with current dir should confirm success, got: %s", stdout)
	}
}

func testListCommandWithData(t *testing.T, binary string) {
	stdout, stderr, err := runCommand(t, binary, "list")
	if err != nil {
		t.Errorf("List command failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "Available commands:") {
		t.Errorf("List output should show available commands, got: %s", stdout)
	}
	
	if !strings.Contains(stdout, "test-cmd") {
		t.Errorf("List output should contain test-cmd, got: %s", stdout)
	}
	
	if !strings.Contains(stdout, "test-cmd-dir") {
		t.Errorf("List output should contain test-cmd-dir, got: %s", stdout)
	}
	
	// Check that working directory is shown
	if !strings.Contains(stdout, "(dir:") {
		t.Errorf("List output should show working directory for some commands, got: %s", stdout)
	}
}

func testRunCommand(t *testing.T, binary string) {
	// Test running a simple command
	stdout, stderr, err := runCommand(t, binary, "run", "--name", "test-cmd")
	if err != nil {
		t.Errorf("Run command failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "Executing: echo hello") {
		t.Errorf("Run output should show executing command, got: %s", stdout)
	}
	
	if !strings.Contains(stdout, "hello") {
		t.Errorf("Run output should contain command output, got: %s", stdout)
	}
}

func testDeleteCommand(t *testing.T, binary string) {
	// Test deleting a specific command
	stdout, stderr, err := runCommand(t, binary, "delete", "--name", "test-cmd")
	if err != nil {
		t.Errorf("Delete command failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "Command 'test-cmd' deleted successfully") {
		t.Errorf("Delete output should confirm success, got: %s", stdout)
	}
	
	// Verify command is gone - should show error message but exit with code 0
	stdout, _, err = runCommand(t, binary, "run", "--name", "test-cmd")
	// The command actually succeeds (exit code 0) but prints error message
	if err != nil {
		t.Errorf("Run deleted command should print error but exit 0, got error: %v", err)
	}
	
	// Check that the error message indicates the command wasn't found
	if !strings.Contains(stdout, "command 'test-cmd' not found") {
		t.Errorf("Run deleted command should indicate command not found, got: %s", stdout)
	}
}

func testDeleteAllCommands(t *testing.T, binary string, tempDir string) {
	// First add some commands to delete
	_, _, err := runCommand(t, binary, "add", "--name", "delete-test-1", "--cmd", "echo 1")
	if err != nil {
		t.Fatalf("Failed to add test command 1: %v", err)
	}
	_, _, err = runCommand(t, binary, "add", "--name", "delete-test-2", "--cmd", "echo 2")
	if err != nil {
		t.Fatalf("Failed to add test command 2: %v", err)
	}
	
	// Test delete all with "no" response
	stdout, stderr, err := runCommandWithInput(t, binary, "n\n", "delete", "--all")
	if err != nil {
		t.Errorf("Delete all with 'n' failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "Operation cancelled") {
		t.Errorf("Delete all with 'n' should be cancelled, got: %s", stdout)
	}
	
	// Test delete all with "yes" response
	stdout, stderr, err = runCommandWithInput(t, binary, "y\n", "delete", "--all")
	if err != nil {
		t.Errorf("Delete all with 'y' failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "Successfully deleted") {
		t.Errorf("Delete all with 'y' should confirm deletion, got: %s", stdout)
	}
	
	// Verify all commands are gone
	stdout, stderr, err = runCommand(t, binary, "list")
	if err != nil {
		t.Errorf("List after delete all failed: %v\nStderr: %s", err, stderr)
	}
	
	if !strings.Contains(stdout, "No commands found") {
		t.Errorf("List after delete all should show no commands, got: %s", stdout)
	}
}

func testErrorCases(t *testing.T, binary string) {
	// Test add without required fields - clir prints error but exits with code 0
	stdout, _, err := runCommand(t, binary, "add")
	if err != nil {
		t.Errorf("Add without arguments should print error but exit 0, got error: %v", err)
	}
	
	// Check error message
	if !strings.Contains(stdout, "name is required") {
		t.Errorf("Add without name should indicate name is required, got: %s", stdout)
	}
	
	// Test run non-existent command
	stdout, _, err = runCommand(t, binary, "run", "--name", "non-existent")
	if err != nil {
		t.Errorf("Run non-existent command should print error but exit 0, got error: %v", err)
	}
	
	if !strings.Contains(stdout, "command 'non-existent' not found") {
		t.Errorf("Run non-existent command should indicate command not found, got: %s", stdout)
	}
	
	// Test delete non-existent command
	stdout, _, err = runCommand(t, binary, "delete", "--name", "non-existent")
	if err != nil {
		t.Errorf("Delete non-existent command should print error but exit 0, got error: %v", err)
	}
	
	if !strings.Contains(stdout, "command 'non-existent' not found") {
		t.Errorf("Delete non-existent command should indicate command not found, got: %s", stdout)
	}
	
	// Test delete without arguments
	stdout, _, err = runCommand(t, binary, "delete")
	if err != nil {
		t.Errorf("Delete without arguments should print error but exit 0, got error: %v", err)
	}
	
	if !strings.Contains(stdout, "either --name or --all is required") {
		t.Errorf("Delete without arguments should indicate name or all is required, got: %s", stdout)
	}
}
