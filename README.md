# afvikle

## TL;DR

afvikle (dansk for hold, run or conduct) is a cli tool for running multiple scripts or commands from anywhere on the system.

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
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Check for potential issues
go vet ./...

# Download and tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

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
├── main.go           # Main application entry point
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
├── README.md        # Project documentation
└── afvikle.exe      # Compiled binary (gitignored)
```

### Dependencies

This project uses:

- `github.com/leaanthony/clir` - CLI framework for Go

### Getting Help

- Create an issue for bugs or feature requests
- Check existing issues before creating new ones
- For questions, start a discussion in the repository

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Make sure you understand and agree to the project's license before contributing.
