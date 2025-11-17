# Qwen Code - Usage Guide & Procedures

## 📚 **Table of Contents**

1. [Quick Start](#quick-start)
2. [Installation](#installation)
3. [Basic Usage](#basic-usage)
4. [Advanced Features](#advanced-features)
5. [AI Agent Integration](#ai-agent-integration)
6. [Development Procedures](#development-procedures)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)

## 🚀 **Quick Start**

### **Prerequisites**

- Node.js >= 20.0.0
- Git
- Docker (optional, for sandbox features)

### **Installation Methods**

#### **Method 1: NPM (Recommended)**

```bash
npm install -g @qwen-code/qwen-code
```

#### **Method 2: GitHub Releases**

```bash
# Download the latest release
curl -L https://github.com/QwenLM/qwen-code/releases/latest/download/qwen-code-linux-amd64.tar.gz | tar xz
sudo mv qwen /usr/local/bin/
```

#### **Method 3: Source Build**

```bash
git clone https://github.com/QwenLM/qwen-code.git
cd qwen-code
npm install
npm run build
npm run start
```

### **First Run**

```bash
# Initialize configuration
qwen init

# Start interactive session
qwen chat

# Or use with specific AI provider
qwen chat --provider openai --model gpt-4
```

## 🔧 **Basic Usage**

### **Command Line Interface**

#### **Core Commands**

```bash
# Start interactive chat
qwen chat

# Process single file
qwen process path/to/file.js

# Batch process directory
qwen process --recursive src/

# Generate code
qwen generate "Create a REST API endpoint" --language python

# Review code
qwen review path/to/file.py --style google

# Run tests
qwen test --coverage
```

#### **Configuration**

```bash
# Set default AI provider
qwen config set provider openrouter

# Set API key
qwen config set openrouter_api_key your_key_here

# Set default model
qwen config set model meta-llama/llama-3.1-405b-instruct

# List configuration
qwen config list

# Reset configuration
qwen config reset
```

### **Interactive Mode**

#### **Chat Interface**

```bash
# Start with specific context
qwen chat --context "I'm working on a React project"

# Load project context
qwen chat --project ./my-project

# Include files in context
qwen chat --include src/*.js --include README.md

# Use specific agent
qwen chat --agent code-reviewer
```

#### **Multi-Agent Collaboration**

```bash
# Start multi-agent session
qwen multi-agent

# Define agent roles
qwen multi-agent --roles "developer,reviewer,tester"

# Load agent configuration
qwen multi-agent --config agents.yaml
```

## 🎯 **Advanced Features**

### **Code Generation**

#### **Template-Based Generation**

```bash
# Generate from template
qwen generate --template react-component --name UserProfile

# Custom template
qwen generate --template ./templates/api-endpoint.js --endpoint users
```

#### **Context-Aware Generation**

```bash
# Generate with project context
qwen generate --project ./my-app "Add user authentication"

# Generate with specific files
qwen generate --include models/user.js "Create user controller"
```

### **Code Analysis**

#### **Static Analysis**

```bash
# Analyze code quality
qwen analyze --quality src/

# Security analysis
qwen analyze --security src/

# Performance analysis
qwen analyze --performance src/

# Dependency analysis
qwen analyze --dependencies
```

#### **Code Review**

```bash
# Automated code review
qwen review --branch feature/new-feature

# Review specific commit
qwen review --commit abc123

# Custom review rules
qwen review --rules .review-rules.json
```

### **Testing Integration**

#### **Test Generation**

```bash
# Generate unit tests
qwen test-generate --type unit src/utils.js

# Generate integration tests
qwen test-generate --type integration src/api/

# Generate E2E tests
qwen test-generate --type e2e --framework cypress
```

#### **Test Execution**

```bash
# Run tests with AI analysis
qwen test run --analyze

# Generate coverage reports
qwen test run --coverage --report html

# Performance testing
qwen test run --performance --threshold 100ms
```

## 🤖 **AI Agent Integration**

### **Agent Configuration**

#### **Creating Custom Agents**

```yaml
# agents.yaml
agents:
  code-reviewer:
    provider: openrouter
    model: meta-llama/llama-3.1-405b-instruct
    system_prompt: 'You are a senior code reviewer. Focus on code quality, security, and best practices.'
    temperature: 0.3

  documentation-writer:
    provider: openai
    model: gpt-4
    system_prompt: 'You are a technical writer. Create clear, comprehensive documentation.'
    temperature: 0.5
```

#### **Agent Workflows**

```yaml
# workflows.yaml
workflows:
  code-review:
    agents:
      - name: code-reviewer
        input: 'Review this code for quality and security'
      - name: documentation-writer
        input: 'Document the reviewed code changes'

  feature-development:
    agents:
      - name: architect
        input: 'Design the architecture for this feature'
      - name: developer
        input: 'Implement the feature based on the architecture'
      - name: tester
        input: 'Create comprehensive tests for the implementation'
```

### **Multi-Agent Sessions**

#### **Starting Multi-Agent Mode**

```bash
# Interactive multi-agent session
qwen multi-agent --workflow code-review

# Batch processing with agents
qwen multi-agent --workflow feature-development --input "Add user authentication"

# Custom agent orchestration
qwen multi-agent --agents "architect,developer,tester" --input "Build a REST API"
```

## 👥 **Development Procedures**

### **Setting Up Development Environment**

#### **Local Development**

```bash
# Clone repository
git clone https://github.com/QwenLM/qwen-code.git
cd qwen-code

# Install dependencies
npm install

# Setup development environment
npm run dev:setup

# Start development server
npm run dev
```

#### **Docker Development**

```bash
# Build development container
docker build -t qwen-code:dev .

# Run development container
docker run -it --rm \
  -v $(pwd):/app \
  -p 3000:3000 \
  qwen-code:dev
```

### **Code Contribution Workflow**

#### **Pre-Commit Procedures**

```bash
# Run pre-commit checks
npm run pre-commit

# Manual pre-commit steps
npm run format
npm run lint
npm run test
npm run typecheck
```

#### **Branch Management**

```bash
# Create feature branch
git checkout -b feature/new-feature

# Commit changes
git add .
git commit -m "feat: add new feature"

# Push and create PR
git push origin feature/new-feature
```

### **Testing Procedures**

#### **Running Tests**

```bash
# Run all tests
npm test

# Run specific test suites
npm run test:unit
npm run test:integration
npm run test:e2e

# Run tests with coverage
npm run test:coverage
```

#### **Test Development**

```bash
# Generate test files
qwen test-generate --type unit src/utils.js

# Run tests in watch mode
npm run test:watch

# Debug tests
npm run test:debug
```

## 🔍 **Troubleshooting**

### **Common Issues**

#### **Installation Problems**

```bash
# Clear npm cache
npm cache clean --force

# Remove node_modules
rm -rf node_modules package-lock.json

# Reinstall
npm install
```

#### **API Key Issues**

```bash
# Check API key configuration
qwen config list

# Test API connection
qwen test --connection

# Reset API keys
qwen config reset
```

#### **Performance Issues**

```bash
# Check system resources
qwen doctor

# Optimize configuration
qwen optimize --performance

# Clear cache
qwen cache clear
```

### **Debug Mode**

#### **Enable Debugging**

```bash
# Enable debug logging
export DEBUG=qwen:*
qwen chat --debug

# Verbose output
qwen --verbose process file.js

# Trace mode
qwen --trace generate "Create a function"
```

#### **Log Analysis**

```bash
# View logs
qwen logs show

# Export logs
qwen logs export --format json

# Analyze logs
qwen logs analyze --errors
```

## 📋 **Best Practices**

### **Code Quality**

#### **Writing Prompts**

```bash
# Good: Specific and contextual
qwen generate "Create a React component for user profile with props: name, email, avatar"

# Bad: Vague and generic
qwen generate "Make a component"
```

#### **Context Management**

```bash
# Include relevant files
qwen chat --include src/*.js --include README.md

# Use project context
qwen chat --project ./my-project

# Set specific context
qwen chat --context "Working on e-commerce checkout flow"
```

### **Security**

#### **API Key Management**

```bash
# Use environment variables
export QWEN_API_KEY=your_key_here

# Use secure storage
qwen config set api_key --secure

# Rotate keys regularly
qwen config rotate api_key
```

#### **Code Security**

```bash
# Security analysis
qwen analyze --security src/

# Vulnerability scanning
qwen scan --vulnerabilities

# Secure code generation
qwen generate --secure "Create authentication system"
```

### **Performance**

#### **Optimization**

```bash
# Use caching
qwen cache enable

# Optimize for speed
qwen optimize --speed

# Batch operations
qwen process --batch src/*.js
```

#### **Resource Management**

```bash
# Monitor usage
qwen monitor --resources

# Set limits
qwen config set max_tokens 1000

# Use streaming
qwen generate --stream "Long content generation"
```

---

## 📞 **Support & Resources**

### **Documentation**

- [Main Documentation](./AI_TUI_DOCUMENTATION.md)
- [API Reference](./docs/api.md)
- [Examples](./examples/)
- [FAQ](./docs/faq.md)

### **Community**

- [GitHub Issues](https://github.com/QwenLM/qwen-code/issues)
- [Discussions](https://github.com/QwenLM/qwen-code/discussions)
- [Discord Server](https://discord.gg/qwen-code)

### **Getting Help**

```bash
# Get help
qwen --help

# Get command help
qwen chat --help

# Show version
qwen --version

# System information
qwen doctor
```

---

**Last Updated**: 2025-11-17  
**Version**: 0.2.2  
**Maintained By**: Qwen Code Team
