# Contributing to Remiges SMTP

Thank you for your interest in contributing to Remiges SMTP! This document provides guidelines and instructions for contributing to the project.

## Table of Contents
1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Code Style](#code-style)
5. [Testing](#testing)
6. [Documentation](#documentation)
7. [Submitting Changes](#submitting-changes)
8. [Review Process](#review-process)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and considerate of others.

## Getting Started

1. **Fork the Repository**
   ```bash
   git clone https://github.com/stlimtat/remiges-smtp.git
   cd remiges-smtp
   ```

2. **Set Up Development Environment**
   ```bash
   # Install dependencies
   go mod download

   # Install development tools
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. **Create a Development Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

1. **Make Changes**
   - Follow the code style guidelines
   - Write tests for new features
   - Update documentation as needed

2. **Run Tests**
   ```bash
   # Run all tests
   bazel test //...

   # Run specific test
   bazel test //pkg/smtp:go_default_test
   ```

3. **Check Code Quality**
   ```bash
   # Run linter
   golangci-lint run
   ```

## Code Style

### Go Code
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Maximum line length: 100 characters
- Use meaningful variable and function names

### Documentation
- Document all exported functions and types
- Use clear, concise comments
- Update README and other documentation files when adding features

### Commit Messages
- Use present tense ("Add feature" not "Added feature")
- First line should be 50 characters or less
- Include a blank line between the subject and body
- Reference issues and pull requests in the body

Example:
```
Add support for TLS connections

- Implement TLS configuration options
- Add tests for TLS handshake
- Update documentation

Fixes #123
```

## Testing

1. **Write Tests**
   - Unit tests for all new code
   - Integration tests for complex features
   - Test coverage should not decrease

2. **Run Tests**
   ```bash
   # Run all tests
   bazel test //...

   # Run tests with coverage
   bazel coverage //...
   ```

## Documentation

1. **Update Documentation**
   - Keep README.md up to date
   - Document new features in appropriate .md files
   - Add examples for new features

2. **Documentation Structure**
   - Use clear headings and subheadings
   - Include code examples where appropriate
   - Keep documentation concise and focused

## Submitting Changes

1. **Push Changes**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create Pull Request**
   - Use the PR template
   - Describe changes clearly
   - Reference related issues
   - Request reviews from maintainers

3. **PR Checklist**
   - [ ] Tests added/updated
   - [ ] Documentation updated
   - [ ] Code style followed
   - [ ] All tests pass
   - [ ] No merge conflicts

## Review Process

1. **Code Review**
   - Address reviewer comments
   - Make requested changes
   - Keep PR up to date with main branch

2. **Merge Process**
   - Squash commits if requested
   - Wait for all checks to pass
   - Get approval from maintainers

## Need Help?

- Open an [issue](https://github.com/stlimtat/remiges-smtp/issues)
- Join our [community chat](https://github.com/stlimtat/remiges-smtp/discussions)
- Check the [FAQ](./FAQ.md)

Thank you for contributing to Remiges SMTP!