package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.etcd.io/bbolt"
)

type Database struct {
	db *bbolt.DB
}

type Command struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
	WorkingDir  string `json:"working_dir"`
	CreatedAt   string `json:"created_at"`
}

var commandsBucket = []byte("commands")

// NewDatabase creates a new database connection and initializes buckets
func NewDatabase() (*Database, error) {
	// Get the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}
	
	execDir := filepath.Dir(execPath)
	dbPath := filepath.Join(execDir, "afvikle.db")
	
	// Create or open the database
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	
	database := &Database{db: db}
	
	// Initialize buckets
	if err := database.initBuckets(); err != nil {
		return nil, fmt.Errorf("failed to initialize buckets: %v", err)
	}
	
	return database, nil
}

// initBuckets creates the necessary buckets if they don't exist
func (d *Database) initBuckets() error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(commandsBucket)
		return err
	})
}

// AddCommand adds a new command to the database
func (d *Database) AddCommand(name, description, command, workingDir string) error {
	// Validate required fields
	if name == "" {
		return fmt.Errorf("command name is required")
	}
	if command == "" {
		return fmt.Errorf("command is required")
	}
	
	// Trim whitespace
	name = strings.TrimSpace(name)
	command = strings.TrimSpace(command)
	description = strings.TrimSpace(description)
	workingDir = strings.TrimSpace(workingDir)
	
	// Set default description if empty
	if description == "" {
		description = "No description provided"
	}
	
	// Validate working directory if provided
	if workingDir != "" {
		if _, err := os.Stat(workingDir); os.IsNotExist(err) {
			return fmt.Errorf("working directory '%s' does not exist", workingDir)
		}
	}
	
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(commandsBucket)
		
		// Check if command already exists
		if b.Get([]byte(name)) != nil {
			return fmt.Errorf("command '%s' already exists", name)
		}
		
		cmd := Command{
			Name:        name,
			Description: description,
			Command:     command,
			WorkingDir:  workingDir,
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		}
		
		data, err := json.Marshal(cmd)
		if err != nil {
			return err
		}
		
		return b.Put([]byte(name), data)
	})
}

// GetCommand retrieves a command by name
func (d *Database) GetCommand(name string) (*Command, error) {
	var cmd Command
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(commandsBucket)
		data := b.Get([]byte(name))
		if data == nil {
			return fmt.Errorf("command '%s' not found", name)
		}
		
		return json.Unmarshal(data, &cmd)
	})
	
	if err != nil {
		return nil, err
	}
	
	return &cmd, nil
}

// GetAllCommands retrieves all commands from the database
func (d *Database) GetAllCommands() ([]Command, error) {
	var commands []Command
	
	err := d.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(commandsBucket)
		
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var cmd Command
			if err := json.Unmarshal(v, &cmd); err != nil {
				return err
			}
			commands = append(commands, cmd)
		}
		
		return nil
	})
	
	return commands, err
}

// UpdateCommand updates an existing command
func (d *Database) UpdateCommand(name, description, command, workingDir string) error {
	// Validate required fields
	if name == "" {
		return fmt.Errorf("command name is required")
	}
	if command == "" {
		return fmt.Errorf("command is required")
	}
	
	// Trim whitespace
	name = strings.TrimSpace(name)
	command = strings.TrimSpace(command)
	description = strings.TrimSpace(description)
	workingDir = strings.TrimSpace(workingDir)
	
	// Set default description if empty
	if description == "" {
		description = "No description provided"
	}
	
	// Validate working directory if provided
	if workingDir != "" {
		if _, err := os.Stat(workingDir); os.IsNotExist(err) {
			return fmt.Errorf("working directory '%s' does not exist", workingDir)
		}
	}
	
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(commandsBucket)
		
		// Check if command exists
		data := b.Get([]byte(name))
		if data == nil {
			return fmt.Errorf("command '%s' not found", name)
		}
		
		var cmd Command
		if err := json.Unmarshal(data, &cmd); err != nil {
			return err
		}
		
		// Update fields
		cmd.Description = description
		cmd.Command = command
		cmd.WorkingDir = workingDir
		
		data, err := json.Marshal(cmd)
		if err != nil {
			return err
		}
		
		return b.Put([]byte(name), data)
	})
}

// DeleteCommand removes a command from the database
func (d *Database) DeleteCommand(name string) error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(commandsBucket)
		
		// Check if command exists
		if b.Get([]byte(name)) == nil {
			return fmt.Errorf("command '%s' not found", name)
		}
		
		return b.Delete([]byte(name))
	})
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// GetDatabasePath returns the path to the database file
func (d *Database) GetDatabasePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	
	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, "afvikle.db"), nil
}
