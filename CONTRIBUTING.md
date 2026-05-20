# Contributing to open-source-vpn

Thank you for your interest in contributing to open-source-vpn! This document provides guidelines and instructions for contributing.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code:

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on what is best for the community
- Show empathy towards other community members

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Environment details** (OS, Go version, etc.)
- **Logs and error messages**
- **Configuration** (sanitized, no private keys!)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear title and description**
- **Use case** and motivation
- **Proposed solution** or implementation approach
- **Alternatives considered**
- **Additional context** (mockups, examples, etc.)

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the coding standards
3. **Add tests** for new functionality
4. **Update documentation** as needed
5. **Ensure tests pass** (`make test`)
6. **Format your code** (`make fmt`)
7. **Run linter** (`make lint`)
8. **Commit with clear messages**
9. **Push to your fork** and submit a pull request

#### Pull Request Guidelines

- Keep changes focused and atomic
- Write clear commit messages
- Reference related issues
- Update CHANGELOG.md
- Ensure CI passes
- Request review from maintainers

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make
- golangci-lint (for linting)

### Setup

```bash
# Clone your fork
git clone https://github.com/lahcenassmira/open-source-vpn.git
cd open-source-vpn

# Add upstream remote
git remote add upstream https://github.com/open-source-vpn/vpn.git

# Install dependencies
go mod download

# Build
make build

# Run tests
make test
```

### Project Structure

```
open-source-vpn/
├── cmd/                 # Command-line applications
│   ├── server/         # Server implementation
│   └── client/         # Client implementation
├── internal/           # Private application code
│   ├── crypto/        # Cryptography
│   ├── tunnel/        # TUN device management
│   ├── network/       # Networking and routing
│   ├── protocol/      # VPN protocol
│   └── config/        # Configuration
├── pkg/               # Public libraries
│   ├── logger/       # Logging
│   └── metrics/      # Metrics
├── configs/          # Example configurations
├── docker/           # Docker files
├── docs/             # Documentation
└── tests/            # Integration tests
```

## Coding Standards

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines and:

- Use `gofmt` for formatting
- Follow Go naming conventions
- Write clear, self-documenting code
- Add comments for exported functions
- Keep functions small and focused
- Handle errors explicitly

### Code Examples

**Good:**
```go
// EncryptPacket encrypts a packet for transmission
func (c *Connection) EncryptPacket(packetType network.PacketType, payload []byte) ([]byte, error) {
    if c.closed.Load() {
        return nil, fmt.Errorf("connection closed")
    }
    
    // Generate nonce
    nonce := make([]byte, crypto.NonceSize)
    if _, err := rand.Read(nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }
    
    // Encrypt payload
    encrypted, err := c.sendCipher.Encrypt(nonce, payload, nil)
    if err != nil {
        return nil, fmt.Errorf("encryption failed: %w", err)
    }
    
    return encrypted, nil
}
```

**Bad:**
```go
// bad naming, no error handling, unclear logic
func (c *Connection) enc(t network.PacketType, p []byte) []byte {
    n := make([]byte, crypto.NonceSize)
    rand.Read(n)
    e, _ := c.sendCipher.Encrypt(n, p, nil)
    return e
}
```

### Testing

- Write unit tests for new functionality
- Aim for >80% code coverage
- Use table-driven tests where appropriate
- Mock external dependencies
- Test error cases

**Example:**
```go
func TestEncryptPacket(t *testing.T) {
    tests := []struct {
        name    string
        payload []byte
        wantErr bool
    }{
        {
            name:    "valid payload",
            payload: []byte("test data"),
            wantErr: false,
        },
        {
            name:    "empty payload",
            payload: []byte{},
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Documentation

- Add godoc comments for exported types and functions
- Update README.md for user-facing changes
- Update docs/ for architectural changes
- Include examples in documentation

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(crypto): add X25519 key exchange

Implement X25519 elliptic curve key exchange for the Noise protocol
handshake. This provides better performance than traditional DH.

Closes #123
```

```
fix(server): handle connection timeout correctly

Previously, expired connections were not cleaned up properly,
leading to memory leaks. This fix adds a periodic cleanup routine.

Fixes #456
```

## Security

### Reporting Security Issues

**DO NOT** open public issues for security vulnerabilities. Instead:

1. Email security@open-source-vpn.example.com
2. Include detailed description
3. Provide steps to reproduce
4. Allow time for fix before disclosure

### Security Guidelines

- Never commit private keys or secrets
- Use secure coding practices
- Validate all inputs
- Handle errors securely
- Follow crypto best practices
- Keep dependencies updated

## Review Process

1. **Automated checks** run on all PRs (tests, linting, etc.)
2. **Code review** by at least one maintainer
3. **Testing** on multiple platforms if needed
4. **Documentation review** for user-facing changes
5. **Approval** and merge by maintainer

### Review Checklist

- [ ] Code follows style guidelines
- [ ] Tests pass and coverage is adequate
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] No security issues introduced
- [ ] Performance impact is acceptable
- [ ] Backward compatibility maintained

## Release Process

1. Update version in code
2. Update CHANGELOG.md
3. Create release branch
4. Run full test suite
5. Build binaries for all platforms
6. Create GitHub release
7. Update documentation
8. Announce release

## Getting Help

- **Documentation**: Check docs/ directory
- **Issues**: Search existing issues
- **Discussions**: Use GitHub Discussions
- **Chat**: Join our community chat (link TBD)

## Recognition

Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project README

Thank you for contributing to open-source-vpn! 🎉
