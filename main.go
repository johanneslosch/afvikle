package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/leaanthony/clir"
)

// resolveDirectory resolves special directory shortcuts like "." and "~"
func resolveDirectory(dir string) (string, error) {
	if dir == "" {
		return "", nil
	}
	
	dir = strings.TrimSpace(dir)
	
	switch dir {
	case ".":
		// Current directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %v", err)
		}
		return cwd, nil
	case "~":
		// Home directory
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %v", err)
		}
		return usr.HomeDir, nil
	default:
		// Handle paths starting with ~/ (home directory expansion)
		if strings.HasPrefix(dir, "~/") {
			usr, err := user.Current()
			if err != nil {
				return "", fmt.Errorf("failed to get user home directory: %v", err)
			}
			return filepath.Join(usr.HomeDir, dir[2:]), nil
		}
		// Regular path - convert to absolute if relative
		absPath, err := filepath.Abs(dir)
		if err != nil {
			return "", fmt.Errorf("failed to resolve path: %v", err)
		}
		return absPath, nil
	}
}

func main() {
	cli := clir.NewCli("afv", "Short for afvikle. CLI to speed up the process of running multiple scripts without creating another script. Run from anywhere.", "v1.0.0")

	// Initialize database
	db, err := NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// List command - show all stored commands
	cli.NewSubCommand("list", "Returns a list of commands runnable with afvikle").
		Action(func() error {
			commands, err := db.GetAllCommands()
			if err != nil {
				return fmt.Errorf("failed to get commands: %v", err)
			}

			if len(commands) == 0 {
				fmt.Println("No commands found. Use 'afv add' to add commands.")
				return nil
			}

			fmt.Println("Available commands:")
			for _, cmd := range commands {
				fmt.Printf("  %-15s %s", cmd.Name, cmd.Description)
				if cmd.WorkingDir != "" {
					fmt.Printf(" (dir: %s)", cmd.WorkingDir)
				}
				fmt.Println()
			}
			return nil
		})

	// Add command - store a new command
	addCmd := cli.NewSubCommand("add", "Add a new command to the database")
	var addName, addDesc, addCommand, addWorkingDir string
	addCmd.StringFlag("name", "Command name", &addName)
	addCmd.StringFlag("desc", "Command description", &addDesc)
	addCmd.StringFlag("cmd", "Command to execute", &addCommand)
	addCmd.StringFlag("dir", "Working directory for the command (optional)", &addWorkingDir)
	addCmd.Action(func() error {
		if addName == "" {
			return fmt.Errorf("name is required")
		}
		if addCommand == "" {
			return fmt.Errorf("cmd is required")
		}

		if addDesc == "" {
			addDesc = "No description provided"
		}

		// Handle special directory shortcuts
		resolvedDir, err := resolveDirectory(addWorkingDir)
		if err != nil {
			return fmt.Errorf("failed to resolve directory: %v", err)
		}

		err = db.AddCommand(addName, addDesc, addCommand, resolvedDir)
		if err != nil {
			return fmt.Errorf("failed to add command: %v", err)
		}

		fmt.Printf("Command '%s' added successfully.\n", addName)
		if resolvedDir != "" {
			fmt.Printf("Working directory: %s\n", resolvedDir)
		}
		return nil
	})

	// Run command - execute a stored command
	runCmd := cli.NewSubCommand("run", "Run a stored command")
	var runName string
	var workingDir string
	runCmd.StringFlag("name", "Command name to run", &runName)
	runCmd.StringFlag("dir", "Working directory to run the command in (optional)", &workingDir)
	runCmd.Action(func() error {
		if runName == "" {
			return fmt.Errorf("name is required")
		}

		command, err := db.GetCommand(runName)
		if err != nil {
			return fmt.Errorf("failed to get command: %v", err)
		}

		// Determine working directory with resolution
		var cmdDir string
		if workingDir != "" {
			// Use specified working directory (resolve shortcuts)
			resolvedDir, err := resolveDirectory(workingDir)
			if err != nil {
				return fmt.Errorf("failed to resolve working directory: %v", err)
			}
			cmdDir = resolvedDir
		} else if command.WorkingDir != "" {
			// Use stored working directory
			cmdDir = command.WorkingDir
		} else {
			// Use current directory
			cmdDir, _ = os.Getwd()
		}

		fmt.Printf("Executing: %s\n", command.Command)
		if cmdDir != "" {
			fmt.Printf("Working directory: %s\n", cmdDir)
		}

		// Parse and execute the command
		parts := strings.Fields(command.Command)
		if len(parts) == 0 {
			return fmt.Errorf("empty command")
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
		// Set working directory if specified
		if cmdDir != "" {
			cmd.Dir = cmdDir
		}

		return cmd.Run()
	})

	// Delete command - remove a stored command
	deleteCmd := cli.NewSubCommand("delete", "Delete a stored command")
	var deleteName string
	var deleteAll bool
	deleteCmd.StringFlag("name", "Command name to delete", &deleteName)
	deleteCmd.BoolFlag("all", "Delete all commands", &deleteAll)
	deleteCmd.Action(func() error {
		if deleteAll {
			// Delete all commands
			commands, err := db.GetAllCommands()
			if err != nil {
				return fmt.Errorf("failed to get commands: %v", err)
			}

			if len(commands) == 0 {
				fmt.Println("No commands to delete.")
				return nil
			}

			fmt.Printf("This will delete %d command(s). Are you sure? (y/N): ", len(commands))
			var response string
			_, _ = fmt.Scanln(&response) // Ignore error - user input handling
			
			if strings.ToLower(strings.TrimSpace(response)) != "y" && strings.ToLower(strings.TrimSpace(response)) != "yes" {
				fmt.Println("Operation cancelled.")
				return nil
			}

			// Delete all commands
			for _, cmd := range commands {
				err := db.DeleteCommand(cmd.Name)
				if err != nil {
					return fmt.Errorf("failed to delete command '%s': %v", cmd.Name, err)
				}
			}

			fmt.Printf("Successfully deleted %d command(s).\n", len(commands))
			return nil
		}

		if deleteName == "" {
			return fmt.Errorf("either --name or --all is required")
		}

		err := db.DeleteCommand(deleteName)
		if err != nil {
			return fmt.Errorf("failed to delete command: %v", err)
		}

		fmt.Printf("Command '%s' deleted successfully.\n", deleteName)
		return nil
	})

	// Info command - show database information
	cli.NewSubCommand("info", "Show database information").
		Action(func() error {
			dbPath, err := db.GetDatabasePath()
			if err != nil {
				return fmt.Errorf("failed to get database path: %v", err)
			}

			commands, err := db.GetAllCommands()
			if err != nil {
				return fmt.Errorf("failed to get commands: %v", err)
			}

			fmt.Printf("Database location: %s\n", dbPath)
			fmt.Printf("Total commands: %d\n", len(commands))
			return nil
		})

	// Starte the CLI
	if err := cli.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
