# 🤖 AI Agents Guidance - AI TUI Core

## 📋 **Overview**

This document provides guidance for AI agents working with the core AI TUI application files. This includes the main application, test suites, and build artifacts.

## 🎯 **Primary Responsibilities**

### **AI Agents Working with AI TUI Core Should:**

#### **1. Code Maintenance**

- **Bug Fixes**: Address issues in main application and tests
- **Feature Development**: Implement new features and enhancements
- **Performance Optimization**: Improve animation and UI performance
- **Code Quality**: Maintain high code quality and standards

#### **2. Testing Management**

- **Test Execution**: Run comprehensive test suites
- **Test Development**: Create new tests for new functionality
- **Coverage Analysis**: Maintain 85%+ test coverage
- **Bug Reproduction**: Create tests for reported bugs

#### **3. Build and Deployment**

- **Build Management**: Ensure application builds correctly
- **Cross-Platform**: Maintain compatibility across platforms
- **Release Preparation**: Prepare builds for release
- **Dependency Management**: Keep dependencies updated and secure

## 🏗️ **Core File Structure**

### **Application Files**

```
AI TUI Core/
├── 📄 main_ai_tui.go           # Main application entry point
│   ├── 🌊 Animation System       # Underwater world implementation
│   ├── 🎨 UI Framework          # Bubble Tea TUI implementation
│   ├── 🤖 Logging System        # Conversation and event logging
│   └── 🔧 Core Logic           # Main application logic
├── 🧪 Test Suite               # Comprehensive test files
│   ├── logging_test.go          # Logging system tests
│   ├── animation_test.go        # Animation engine tests
│   ├── ui_test.go              # UI component tests
│   ├── integration_test.go       # Integration workflow tests
│   ├── performance_test.go      # Performance and stress tests
│   └── edge_case_test.go       # Edge case and error tests
├── 📊 Data Directory          # Runtime data storage
│   └── .ai-tui-data/         # Conversation and event logs
└── 🔧 Build Artifacts          # Generated files
    ├── ai-tui                 # Compiled executable
    └── test results            # Test output files
```

## 🔧 **Technical Guidelines**

### **Code Standards**

- **Go Conventions**: Follow standard Go coding conventions
- **Error Handling**: Use proper error handling patterns
- **Documentation**: Include comprehensive comments
- **Testing**: Write tests for all new functionality

### **Performance Requirements**

- **Animation FPS**: Maintain 60 FPS target
- **Memory Usage**: Keep under 10MB for normal usage
- **Startup Time**: Initialize in under 100ms
- **Response Time**: UI responses under 5ms

### **Quality Standards**

- **Test Coverage**: Maintain 85%+ coverage
- **Linting**: Pass all Go linting rules
- **Security**: No security vulnerabilities
- **Compatibility**: Support Go 1.19+

## 🧪 **Testing Guidelines**

### **Test Categories**

```
Testing Strategy:
├── 🧪 Unit Tests (35+ functions)
│   ├── Logging System (8 tests)
│   ├── Animation Engine (12 tests)
│   └── UI Components (15+ tests)
├── 🔗 Integration Tests (10 functions)
│   ├── Full Workflows
│   ├── Data Persistence
│   └── Concurrent Access
├── ⚡ Performance Tests (12 functions)
│   ├── Stress Testing
│   ├── Memory Usage
│   └── Benchmarking
└── 🚨 Edge Case Tests (15+ functions)
    ├── Error Handling
    ├── Boundary Conditions
    └── Corruption Scenarios
```

### **Test Execution**

```bash
# Run all tests
go test -v ./...

# Run specific test category
go test -v logging_test.go main_ai_tui.go
go test -v animation_test.go main_ai_tui.go
go test -v ui_test.go main_ai_tui.go

# Run with coverage
go test -cover -v ./...

# Run benchmarks
go test -bench=. -v ./...
```

## 🚀 **Build and Release**

### **Build Commands**

```bash
# Development build
go build -o ai-tui main_ai_tui.go

# Production build
go build -ldflags="-s -w" -o ai-tui main_ai_tui.go

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o ai-tui-linux main_ai_tui.go
GOOS=darwin GOARCH=amd64 go build -o ai-tui-macos main_ai_tui.go
GOOS=windows GOARCH=amd64 go build -o ai-tui.exe main_ai_tui.go
```

### **Release Checklist**

- [ ] All tests pass
- [ ] Code coverage ≥85%
- [ ] No security vulnerabilities
- [ ] Documentation updated
- [ ] Version number updated
- [ ] Cross-platform builds successful
- [ ] Performance benchmarks met

## 🔄 **Development Workflow**

### **Feature Development**

1. **Requirements Analysis**: Understand feature requirements
2. **Test Planning**: Plan tests before implementation
3. **Implementation**: Write code following standards
4. **Testing**: Implement comprehensive tests
5. **Documentation**: Update relevant documentation
6. **Review**: Code review and quality checks
7. **Integration**: Test with existing functionality

### **Bug Fix Process**

1. **Issue Analysis**: Understand bug and root cause
2. **Test Creation**: Create test that reproduces bug
3. **Fix Implementation**: Implement minimal fix
4. **Verification**: Ensure fix works and test passes
5. **Regression Testing**: Test for unintended side effects
6. **Documentation**: Update documentation if needed

## 📊 **Monitoring and Metrics**

### **Key Performance Indicators**

- **Test Pass Rate**: Target 100% for all test suites
- **Code Coverage**: Maintain 85%+ coverage
- **Build Success Rate**: 100% successful builds
- **Performance Benchmarks**: Meet or exceed performance targets
- **Bug Fix Time**: Resolve bugs within target timeframe

### **Quality Metrics**

- **Code Quality**: Maintain high code quality scores
- **Security**: Zero security vulnerabilities
- **Documentation**: Keep documentation current and accurate
- **User Satisfaction**: Monitor and improve user experience

## 🚨 **Common Issues and Solutions**

### **Frequently Encountered Problems**

#### **1. Animation Performance**

- **Issue**: Frame rate drops below 60 FPS
- **Solution**: Optimize particle rendering, reduce object count
- **Prevention**: Profile regularly, optimize bottlenecks

#### **2. Memory Leaks**

- **Issue**: Memory usage increases over time
- **Solution**: Check for unreleased resources, fix object lifecycle
- **Prevention**: Use memory profiling tools regularly

#### **3. Test Failures**

- **Issue**: Tests failing due to environment issues
- **Solution**: Use isolated test environments, mock external dependencies
- **Prevention**: Make tests deterministic and environment-independent

#### **4. Build Issues**

- **Issue**: Cross-platform build failures
- **Solution**: Use CI/CD for multi-platform builds
- **Prevention**: Test builds on all target platforms

## 🛡️ **Security Considerations**

### **Security Best Practices**

- **Input Validation**: Validate all user inputs
- **File Access**: Secure file system operations
- **Dependency Management**: Keep dependencies updated
- **Error Handling**: Don't expose sensitive information in errors

### **Security Monitoring**

- **Vulnerability Scanning**: Regular security scans
- **Dependency Checks**: Monitor for vulnerable dependencies
- **Code Review**: Security-focused code reviews
- **Incident Response**: Plan for security incidents

## 🤝 **Collaboration Guidelines**

### **Working with Other Agents**

- **Clear Communication**: Use clear, descriptive commit messages
- **Code Reviews**: Review all changes thoroughly
- **Knowledge Sharing**: Document learnings and solutions
- **Conflict Resolution**: Handle conflicts constructively

### **Coordination with Documentation Team**

- **Change Notification**: Inform documentation team of changes
- **API Changes**: Coordinate API changes with documentation
- **Feature Updates**: Provide documentation updates for new features
- **Bug Fixes**: Share bug fix information

## 📈 **Continuous Improvement**

### **Regular Activities**

- **Code Reviews**: Participate in regular code reviews
- **Learning**: Stay updated with Go and TUI best practices
- **Tooling**: Improve development tools and processes
- **Performance**: Regular performance optimization

### **Innovation Opportunities**

- **New Features**: Suggest and implement new features
- **User Experience**: Improve user interface and experience
- **Performance**: Optimize for better performance
- **Integration**: Explore integration opportunities

## 📞 **Support and Escalation**

### **When to Escalate**

- **Critical Bugs**: Security issues or crashes
- **Performance Issues**: Significant performance regressions
- **Build Failures**: Unable to resolve build issues
- **Test Failures**: Persistent test failures

### **Escalation Process**

1. **Document Issue**: Clearly document the problem
2. **Attempt Resolution**: Show attempted solutions
3. **Provide Context**: Include all relevant context
4. **Request Help**: Clearly specify what help is needed

---

## 🎯 **Success Criteria**

AI TUI core agents are successful when:

1. ✅ **Application Builds**: All builds complete successfully
2. ✅ **Tests Pass**: All test suites pass with 85%+ coverage
3. ✅ **Performance Meets Targets**: 60 FPS, <10MB memory, <100ms startup
4. ✅ **Quality Standards**: Code quality, security, and documentation standards met
5. ✅ **User Experience**: Application provides excellent user experience
6. ✅ **Reliability**: Application is stable and reliable

## 📊 **Current Status**

- **Application Status**: 🟢 Production Ready
- **Test Coverage**: 🟢 85%
- **Performance**: 🟢 All targets met
- **Quality**: 🟢 High quality standards maintained
- **Documentation**: 🟢 Current and comprehensive

---

**Last Updated**: 2025-11-15  
**Maintained By**: AI TUI Development Team  
**Version**: 1.0.0  
**Status**: 🟢 Production Ready
