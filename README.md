# afvikle

## TL;DR

afvikle (dansk for hold, run or conduct) is a CLI tool for running multiple scripts or commands from anywhere on the system. It includes a built-in database to store and manage your frequently used commands with working directory support.

## Features

- **Command Storage**: Store commands with names, descriptions, and working directories in a local database
- **Working Directory Support**: Commands can store and use specific working directories
- **Directory Shortcuts**: Use `.` (current), `~` (home), and `~/path` (home subdirectory) shortcuts
- **Runtime Directory Override**: Override working directories when running commands
- **Cross-Platform**: Works on Windows, Linux, and macOS with proper path handling
- **Portable Database**: Database file is stored alongside the executable - no installation required
- **Simple CLI**: Intuitive command-line interface with comprehensive help
- **Bulk Operations**: Delete all commands at once with confirmation
- **Path Resolution**: Automatic resolution of relative paths to absolute paths
- **Safe Operations**: Confirmation prompts for destructive operations

## Commands

### Core Commands

| Command      | Description               | Example                                             |
| ------------ | ------------------------- | --------------------------------------------------- |
| `afv add`    | Store a new command       | `afv add --name "build" --cmd "go build" --dir "."` |
| `afv list`   | Show all stored commands  | `afv list`                                          |
| `afv run`    | Execute a stored command  | `afv run --name "build"`                            |
| `afv delete` | Remove command(s)         | `afv delete --name "old-cmd"` or `afv delete --all` |
| `afv info`   | Show database information | `afv info`                                          |

### Command Flags

#### `afv add` - Add Command

- `--name` (required): Unique command name
- `--cmd` (required): Command to execute
- `--desc` (optional): Command description
- `--dir` (optional): Working directory (supports `.`, `~`, `~/path`)

#### `afv run` - Run Command

- `--name` (required): Command name to execute
- `--dir` (optional): Override working directory for this run

#### `afv delete` - Delete Command(s)

- `--name`: Delete specific command
- `--all`: Delete all commands (with confirmation)

## Usage

### Adding Commands

Store commands with optional working directories:

```bash
# Basic command
afv add --name "hello" --cmd "echo Hello World"

# Command with current directory
afv add --name "build" --desc "Build the project" --cmd "go build" --dir "."

# Command with home directory
afv add --name "backup" --desc "Backup files" --cmd "backup.sh" --dir "~"

# Command with specific path
afv add --name "deploy" --desc "Deploy app" --cmd "./scripts/deploy.sh" --dir "~/projects/myapp"
```

### Listing Commands

See all stored commands with their working directories:

```bash
afv list
```

Output:

```
Available commands:
  build           Build the project (dir: /home/user/project)
  deploy          Deploy app (dir: /home/user/projects/myapp)
  hello           Hello World
  backup          Backup files (dir: /home/user)
```

### Running Commands

Execute stored commands:

```bash
# Run with stored working directory
afv run --name "build"

# Override working directory
afv run --name "build" --dir "/different/path"

# Use directory shortcuts for override
afv run --name "build" --dir "."          # Current directory
afv run --name "build" --dir "~"          # Home directory
afv run --name "build" --dir "~/Desktop"  # Home subdirectory
```

### Managing Commands

Delete commands individually or all at once:

```bash
# Delete specific command
afv delete --name "old-command"

# Delete all commands (with confirmation)
afv delete --all
```

### Database Information

View database location and statistics:

```bash
afv info
```

Output:

```
Database location: /path/to/afvikle.db
Total commands: 5
```

## Working Directory Features

### Directory Shortcuts

- **`.`** - Current directory (resolved to absolute path when storing)
- **`~`** - User's home directory
- **`~/path`** - Subdirectory under home directory

### Directory Priority (when running commands)

1. **Runtime `--dir` flag** (highest priority)
2. **Stored working directory**
3. **Current directory** (lowest priority)

### Cross-Platform Path Handling

- Windows: `C:\Users\username`, `C:\path\to\project`
- Linux/macOS: `/home/username`, `/path/to/project`
- Home directory (`~`) resolves correctly on all platforms

## Database

### Location

The database (`afvikle.db`) is automatically created in the same directory as the executable, making the tool completely portable.

### Portability

- Copy the executable and `.db` file together
- No external dependencies or installation required
- Works immediately on any compatible system

## Installation

1. Download the binary for your platform from the releases page
2. Place it in a directory that's in your PATH
3. Start using it immediately - no setup required!

## Examples

### Project Development Workflow

```bash
# Set up project commands from project directory
cd /path/to/my/project
afv add --name "build" --desc "Build project" --cmd "make" --dir "."
afv add --name "test" --desc "Run tests" --cmd "npm test" --dir "."
afv add --name "start" --desc "Start dev server" --cmd "npm start" --dir "."

# Run from anywhere - commands execute in correct directory
cd /tmp
afv run --name "build"  # Runs make in /path/to/my/project
afv run --name "test"   # Runs npm test in /path/to/my/project
```

### System Administration

```bash
# System-wide commands (no working directory needed)
afv add --name "sysinfo" --desc "System info" --cmd "uname -a"
afv add --name "diskspace" --desc "Disk usage" --cmd "df -h"

# Home directory maintenance
afv add --name "cleanup" --desc "Clean downloads" --cmd "cleanup.sh" --dir "~/Downloads"
```

### Multi-Project Management

```bash
# Different projects with same command patterns
afv add --name "frontend-build" --cmd "npm run build" --dir "~/projects/frontend"
afv add --name "backend-build" --cmd "go build" --dir "~/projects/backend"
afv add --name "mobile-build" --cmd "flutter build" --dir "~/projects/mobile"

# Run any project's build from anywhere
afv run --name "frontend-build"
```

## Contribution

We welcome contributions to afvikle! This guide will help you get started with development and understand the project workflow.

### Prerequisites

- Go 1.24.5 or later
- Git
- A text editor or IDE with Go support
- Or the GitHub Codespace

### Development Setup

1. **Fork and Clone the Repository**

   ```bash
   git clone https://github.com/johanneslosch/afvikle.git
   cd afvikle
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   ```

3. **Verify Setup**
   ```bash
   go run main.go --help
   ```

### Go Commands Reference

Here are the essential Go commands for working with this project:

#### Building and Running

```bash
# Run the application directly
go run main.go

# Build for current platform
go build -o afvikle.exe main.go

# Build for different platforms
# Windows
GOOS=windows GOARCH=amd64 go build -o afvikle.exe main.go
# Linux
GOOS=linux GOARCH=amd64 go build -o afvikle main.go
# macOS
GOOS=darwin GOARCH=amd64 go build -o afvikle main.go
```

#### Testing and Quality

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage report
go test -cover ./...

# Run tests with detailed coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test file
go test -v ./database_test.go
go test -v ./cli_test.go
go test -v ./resolve_test.go

# Run specific test function
go test -v -run TestAddCommand

# Format code
go fmt ./...

# Check for potential issues
go vet ./...

# Download and tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

### Test Suite

The project includes comprehensive test coverage:

#### Database Tests (`database_test.go`)

- **CRUD Operations**: Add, Get, Update, Delete commands
- **Validation**: Required fields, duplicate names, invalid directories
- **Data Integrity**: Field trimming, default values, timestamps
- **Error Handling**: Non-existent commands, invalid input

#### CLI Integration Tests (`cli_test.go`)

- **Command Execution**: Full CLI workflow testing
- **Help System**: Command help and usage information
- **Error Scenarios**: Invalid input, missing arguments
- **User Interaction**: Confirmation prompts, bulk operations
- **Working Directory**: Directory resolution and execution

#### Directory Resolution Tests (`resolve_test.go`)

- **Path Shortcuts**: `.` (current), `~` (home), `~/path` (subdirectory)
- **Cross-Platform**: Windows, Linux, macOS path handling
- **Edge Cases**: Whitespace, invalid paths, complex scenarios
- **Validation**: Absolute path conversion, error conditions

#### Running Tests

```bash
# Quick test run
go test ./...

# Comprehensive test with coverage
go test -v -cover ./...

# Test specific functionality
go test -v -run TestDatabase     # Database tests only
go test -v -run TestCLI          # CLI tests only
go test -v -run TestResolve      # Directory resolution tests only
```

#### Test Coverage

The test suite covers:

- ✅ **Database operations** (100 % of CRUD methods)
- ✅ **CLI commands** (All subcommands and flags)
- ✅ **Directory resolution** (All path types and shortcuts)
- ✅ **Error handling** (Invalid input and edge cases)
- ✅ **Cross-platform compatibility** (Path handling differences)

#### Development Workflow

```bash
# Install development tools (optional, but recommended)
go install golang.org/x/tools/cmd/goimports@latest

# Check imports and format (if goimports is installed)
goimports -w .

# Alternative: Use built-in Go tools only
go fmt ./...
go mod tidy
```

**Note for Windows users:** If `goimports` is not recognized, you can either:

1. Install it with `go install golang.org/x/tools/cmd/goimports@latest` and ensure your Go bin directory is in PATH
2. Use the built-in `go fmt ./...` command instead, which provides basic formatting

### Continuous Integration

The project uses GitHub Actions for automated testing and releases:

#### **Testing Workflow** (`.github/workflows/test.yml`)

- **Triggers**: Pull requests and pushes to `main` and `develop` branches
- **Matrix Testing**:
  - Operating Systems: Ubuntu, Windows, macOS
  - Go Versions: 1.24.x, 1.23.x
- **Quality Checks**:
  - `go test` with race detection and coverage
  - `go vet` for static analysis
  - `golangci-lint` for comprehensive linting
  - `gosec` security scanner
- **Build Verification**: Ensures the binary builds successfully
- **Coverage Reports**: Uploads to Codecov (Ubuntu + Go 1.24.x only)

#### **Release Workflow** (`.github/workflows/release.yml`)

- **Trigger**: Git tags matching `v*` pattern (e.g., `v1.0.0`)
- **Cross-Platform Builds**: Linux, Windows, macOS (amd64 and arm64)
- **Automated Releases**: Creates GitHub releases with binaries
- **Optimized Binaries**: Built with `-ldflags="-s -w"` for smaller size

#### Running CI Locally

```bash
# Run the same tests as CI
go test -v -race -coverprofile=coverage.out ./...
go vet ./...
go build -v .

# Install and run golangci-lint (optional)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run
```

### Making Changes

1. **Create a Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Your Changes**

   - Follow Go conventions and best practices
   - Write clear, descriptive commit messages
   - Add tests for new functionality
   - Update documentation as needed

3. **Test Your Changes**

   ```bash
   go test ./...
   go build -o afvikle.exe main.go
   ./afvikle.exe --help
   ```

4. **Format and Lint**

   ```bash
   # Format code (always available)
   go fmt ./...

   # Check for potential issues
   go vet ./...

   # Optional: Advanced import formatting (if goimports is installed)
   goimports -w .
   ```

### Code Style Guidelines

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Write documentation for exported functions and types
- Keep functions small and focused
- Handle errors appropriately
- Use Go modules for dependency management

### Submitting Changes

1. **Commit Your Changes**

   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

2. **Push to Your Fork**

   ```bash
   git push origin feature/your-feature-name
   ```

3. **Create a Pull Request**
   - Provide a clear description of your changes
   - Reference any related issues
   - Ensure all tests pass
   - Wait for code review

### Project Structure

```
afvikle/
├── main.go              # Main application entry point
├── database.go          # Database operations and models
├── database_test.go     # Database unit tests
├── cli_test.go          # CLI integration tests
├── resolve_test.go      # Directory resolution tests
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── README.md           # Project documentation
├── COMMAND_REFERENCE.md # Complete command reference
├── WORKING_DIRECTORY_GUIDE.md # Working directory guide
├── FIELD_REQUIREMENTS.md # Database field documentation
├── LICENSE             # MIT license
├── .gitignore          # Git ignore rules
├── .devcontainer/      # Development container config
├── .github/            # GitHub Actions workflows
│   └── workflows/
│       ├── test.yml    # CI testing workflow
│       └── release.yml # Release automation
└── afvikle.exe         # Compiled binary (gitignored)
```

### Dependencies

This project uses:

- `github.com/leaanthony/clir` - CLI framework for Go
- `go.etcd.io/bbolt` - Pure Go key/value database

### Getting Help

- Create an issue for bugs or feature requests
- Check existing issues before creating new ones
- For questions, start a discussion in the repository

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Make sure you understand and agree to the project's license before contributing.
