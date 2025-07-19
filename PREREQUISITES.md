# Prerequisites and Setup Guide for AWL Learning

This guide helps novice programmers set up their development environment and acquire the necessary knowledge to understand and recreate the AWL mesh VPN project.

## Programming Prerequisites

### Essential Knowledge (Must Have)

#### 1. Go Programming Fundamentals
You need to be comfortable with:

```go
// Variables and basic types
var name string = "Alice"
port := 8080
connected := true

// Functions
func connectToPeer(address string) error {
    // Implementation
    return nil
}

// Structs and methods
type Peer struct {
    ID   string
    Name string
}

func (p *Peer) Connect() error {
    fmt.Printf("Connecting to %s\n", p.Name)
    return nil
}

// Interfaces
type Connector interface {
    Connect() error
    Disconnect() error
}

// Goroutines and channels
func startServer() {
    messages := make(chan string, 10)
    
    go func() {
        for msg := range messages {
            fmt.Println("Received:", msg)
        }
    }()
    
    messages <- "Hello"
    close(messages)
}

// Error handling
func riskyOperation() error {
    if someCondition {
        return fmt.Errorf("something went wrong")
    }
    return nil
}

func main() {
    if err := riskyOperation(); err != nil {
        log.Fatal(err)
    }
}
```

**Learning Resources for Go**:
- [Tour of Go](https://tour.golang.org/) - Interactive tutorial
- [Go by Example](https://gobyexample.com/) - Practical examples
- [Effective Go](https://golang.org/doc/effective_go.html) - Best practices

#### 2. Basic Networking Concepts
Understanding of:

- **IP Addresses**: What 192.168.1.100 and 10.66.0.1 mean
- **Ports**: How applications use port numbers (like :8080)
- **TCP vs UDP**: Reliable vs unreliable protocols
- **Client/Server model**: How programs connect to each other

```go
// Basic TCP server example you should understand
func basicServer() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleConnection(conn) // Handle each client
    }
}
```

#### 3. Command Line Basics
Comfortable with terminal/command prompt:

```bash
# Navigate directories
cd /path/to/project
ls -la                    # List files
pwd                       # Current directory

# Run Go programs
go run main.go
go build .
go mod init project-name
go mod tidy

# Basic Git
git clone <repository>
git status
git add .
git commit -m "message"
```

### Helpful Knowledge (Good to Have)

#### JSON and HTTP APIs
Understanding how programs communicate:

```go
// JSON data format
type Config struct {
    Name    string `json:"name"`
    Address string `json:"address"`
    Port    int    `json:"port"`
}

// HTTP client
resp, err := http.Get("http://api.example.com/peers")
if err != nil {
    return err
}
defer resp.Body.Close()

var peers []Peer
json.NewDecoder(resp.Body).Decode(&peers)
```

#### Basic Cryptography Concepts
- Public/private key pairs
- Digital signatures
- What "encrypted" means

You don't need to understand the math, just the concepts.

### Don't Worry About (You'll Learn)
- Peer-to-peer networking details
- VPN internals
- Advanced cryptography
- libp2p specifics
- Cross-platform development

## Development Environment Setup

### 1. Install Go

#### Windows
1. Go to https://golang.org/dl/
2. Download Windows installer (.msi)
3. Run installer with default settings
4. Open Command Prompt and test:
   ```cmd
   go version
   ```

#### macOS
```bash
# Option 1: Download from website
# Go to https://golang.org/dl/ and download .pkg

# Option 2: Use Homebrew
brew install go

# Test installation
go version
```

#### Linux (Ubuntu/Debian)
```bash
# Option 1: Package manager (may be older version)
sudo apt update
sudo apt install golang-go

# Option 2: Official installer (recommended)
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Add to PATH in ~/.bashrc or ~/.profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Test
go version
```

### 2. Install Git

#### Windows
1. Download from https://git-scm.com/
2. Install with default settings
3. Use "Git Bash" for terminal

#### macOS
```bash
# Usually pre-installed, or:
xcode-select --install
# Or with Homebrew:
brew install git
```

#### Linux
```bash
sudo apt update
sudo apt install git
```

Test Git installation:
```bash
git --version
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

### 3. Choose an Editor/IDE

#### VS Code (Recommended for beginners)
1. Download from https://code.visualstudio.com/
2. Install Go extension:
   - Open VS Code
   - Press Ctrl+Shift+X (Cmd+Shift+X on Mac)
   - Search "Go" and install the official Go extension
3. Install Go tools when prompted

#### Other Good Options
- **GoLand** (JetBrains) - Professional IDE
- **Vim/Neovim** - For experienced terminal users
- **Sublime Text** - Lightweight editor

### 4. Clone the AWL Repository

```bash
# Create a development directory
mkdir ~/development
cd ~/development

# Clone the repository
git clone https://github.com/anywherelan/awl.git
cd awl

# Download dependencies
go mod tidy

# Test basic build (might fail due to missing frontend)
go build ./cmd/awl
```

### 5. Additional Tools (Optional but Helpful)

#### Network Analysis Tools
```bash
# Linux/macOS
sudo apt install wireshark tcpdump  # Linux
brew install wireshark              # macOS

# Windows
# Download Wireshark from https://www.wireshark.org/
```

#### Development Tools
```bash
# Go development tools
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/cmd/goimports@latest

# Check installation
staticcheck --version
```

## Learning Path and Timeline

### Week 1: Foundation Review
**Goal**: Ensure Go and networking basics are solid

#### Day 1-2: Go Review
- Complete [Tour of Go](https://tour.golang.org/) if new to Go
- Practice with basic programs:

```go
// Exercise: Simple TCP echo server
// Exercise: JSON configuration reader
// Exercise: Basic HTTP client
```

#### Day 3-4: Networking Concepts
- Learn about IP addresses and subnets
- Understand TCP vs UDP
- Try basic socket programming

#### Day 5-7: Environment Setup
- Set up development environment
- Clone AWL repository
- Run tutorial Phase 1 examples
- Read AWL README and architecture docs

### Week 2: Cryptography Basics
**Goal**: Understand security concepts used in AWL

#### Day 1-3: Cryptography Concepts
- Public/private key cryptography
- Digital signatures
- TLS/SSL basics

#### Day 4-7: Hands-on Practice
- Run tutorial Phase 2 examples
- Experiment with Ed25519 keys
- Understand authentication protocols

### Week 3: Virtual Networking
**Goal**: Learn about VPN and virtual interfaces

#### Day 1-3: TUN/TAP Concepts
- Virtual network interfaces
- Packet capture and injection
- IP routing basics

#### Day 4-7: Practical Implementation
- Run tutorial Phase 3 examples (needs admin privileges)
- Set up virtual interfaces
- Practice packet analysis

### Week 4: P2P Networking
**Goal**: Understand peer-to-peer systems

#### Day 1-3: P2P Concepts
- Peer discovery mechanisms
- NAT traversal
- Distributed hash tables

#### Day 4-7: libp2p Implementation
- Run tutorial Phase 4 examples
- Study libp2p documentation
- Build simple P2P applications

## Validation Checklist

Before proceeding with AWL development, ensure you can:

### Basic Go Skills
- [ ] Write a simple HTTP server
- [ ] Use goroutines and channels
- [ ] Handle JSON data
- [ ] Implement interfaces
- [ ] Handle errors properly
- [ ] Use third-party packages with `go mod`

### Networking Understanding
- [ ] Explain the difference between TCP and UDP
- [ ] Understand what an IP address and port are
- [ ] Know what a subnet mask is (like /24)
- [ ] Can use tools like `ping`, `netstat`, or `ss`

### Development Environment
- [ ] Go compiler works (`go version`)
- [ ] Can build and run Go programs
- [ ] Git is configured and working
- [ ] Editor has Go support and syntax highlighting
- [ ] Can install Go packages (`go get`)

### System Administration (for VPN examples)
- [ ] Can run commands with elevated privileges
- [ ] Understand basic network configuration
- [ ] Can create and configure network interfaces

## Troubleshooting Common Issues

### Go Installation Issues
```bash
# If "go: command not found"
echo $PATH                    # Check if Go is in PATH
which go                      # Find Go location
export PATH=$PATH:/usr/local/go/bin  # Add to PATH

# If modules don't work
go env GOPATH                 # Check Go workspace
go env GOMOD                  # Check module mode
```

### Permission Issues
```bash
# For TUN interface examples on Linux/macOS
sudo -v                       # Test sudo access

# On Windows
# Right-click Command Prompt -> "Run as Administrator"
```

### Network Issues
```bash
# Check if ports are available
netstat -tuln | grep :8080   # Linux/macOS
netstat -an | findstr :8080  # Windows

# Check firewall settings
sudo ufw status               # Linux
# Windows: Check Windows Defender Firewall
```

### Build Issues
```bash
# Clear module cache
go clean -modcache

# Verify dependencies
go mod verify
go mod tidy

# Check Go environment
go env
```

## Getting Help

### Documentation Resources
- [Go Documentation](https://golang.org/doc/)
- [Go Blog](https://blog.golang.org/)
- [AWL Repository](https://github.com/anywherelan/awl)

### Community Resources
- [Go Forum](https://forum.golangbridge.org/)
- [Go Slack](https://gophers.slack.com/)
- [Stack Overflow](https://stackoverflow.com/questions/tagged/go) (tag: go)

### Learning Communities
- [Go Wiki](https://github.com/golang/go/wiki)
- [Awesome Go](https://github.com/avelino/awesome-go) - Curated list of Go packages
- [Go Time Podcast](https://changelog.com/gotime)

## Ready to Start?

Once you've completed this setup and feel comfortable with the prerequisites, you're ready to dive into the [LEARNING_GUIDE.md](./LEARNING_GUIDE.md) and start with the hands-on tutorials!

The journey from novice to understanding mesh VPN implementation is challenging but rewarding. Take your time with each phase, and don't hesitate to revisit concepts as needed.

Good luck! 🚀