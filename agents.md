# 🤖 AI Agents Guidance - Root Directory

## 📋 **Overview**

This document provides guidance for AI agents working with the AI TUI project at the root level. This is a multi-component repository containing various projects and tools.

## 🎯 **Primary Responsibilities**

### **AI Agents Working at Root Level Should:**

#### **1. Project Coordination**

- **Monitor Multiple Projects**: Track status of AI TUI, Qwen Code, EGUI Template, and other components
- **Dependency Management**: Ensure all projects have compatible dependencies
- **Build Coordination**: Manage cross-project build processes
- **Release Management**: Coordinate releases across multiple components

#### **2. Architecture Oversight**

- **Maintain Consistency**: Ensure architectural patterns are consistent across projects
- **Code Quality**: Enforce coding standards and best practices
- **Documentation**: Keep project documentation up-to-date and accurate
- **Testing Strategy**: Oversee testing approaches across all components

#### **3. Integration Management**

- **Cross-Project Integration**: Ensure different projects work together when needed
- **API Compatibility**: Maintain compatibility between different project versions
- **Shared Resources**: Manage shared resources like documentation, scripts, and tools

## 🏗️ **Project Structure Understanding**

### **Key Components**

```
qwen-code/
├── 🌌 AI TUI Project           # Main terminal application
│   ├── main_ai_tui.go         # Core application
│   ├── *_test.go              # Comprehensive test suite
│   └── ai-tui                # Built executable
├── 📦 Qwen Code               # VS Code extension
│   ├── packages/              # Extension packages
│   ├── docs/                  # Documentation
│   └── integration-tests/     # Test suites
├── 🎨 EGUI Template           # Rust GUI template
│   ├── demo/                  # Example implementations
│   └── widget_gallery.rs      # UI components
├── 🔧 Development Tools        # Build and deployment tools
│   ├── Makefile               # Build automation
│   ├── Dockerfile             # Container builds
│   └── scripts/               # Utility scripts
└── 📚 Documentation         # Project documentation
    ├── README.md              # Main project docs
    ├── CONTRIBUTING.md         # Contribution guidelines
    └── SECURITY.md            # Security policies
```

## 🎮 **Agent Interaction Guidelines**

### **When Working with Multiple Projects:**

#### **1. Context Switching**

- **Project Identification**: Always identify which project you're working on
- **Dependency Awareness**: Understand how projects depend on each other
- **Build Order**: Know the correct build order for multi-project changes

#### **2. Change Management**

- **Impact Assessment**: Evaluate changes across all affected projects
- **Version Coordination**: Ensure version numbers are properly synchronized
- **Testing Strategy**: Test changes in isolation and then in integration

#### **3. Communication**

- **Clear Documentation**: Document changes that affect multiple projects
- **Change Logs**: Maintain comprehensive change logs
- **Stakeholder Notification**: Inform relevant stakeholders of cross-project changes

## 🔧 **Technical Guidelines**

### **Build Processes**

```bash
# AI TUI
go build -o ai-tui main_ai_tui.go

# Qwen Code (if applicable)
npm install && npm run build

# EGUI Template
cargo build --release
```

### **Testing Commands**

```bash
# AI TUI Tests
go test -v ./...

# Integration Tests
npm run test:integration

# Cross-Project Testing
make test-all
```

### **Quality Assurance**

- **Code Reviews**: All changes must be reviewed
- **Automated Testing**: CI/CD pipelines must pass
- **Documentation Updates**: Update relevant documentation
- **Performance Impact**: Measure performance impact of changes

## 🚨 **Common Pitfalls to Avoid**

### **Multi-Project Issues**

1. **Dependency Conflicts**: Don't create incompatible dependency versions
2. **Breaking Changes**: Avoid breaking changes without proper coordination
3. **Documentation Drift**: Keep documentation synchronized with code
4. **Test Coverage**: Don't reduce test coverage across projects

### **Best Practices**

1. **Incremental Changes**: Make small, incremental changes
2. **Rollback Plans**: Always have rollback plans for significant changes
3. **Communication**: Over-communicate changes and their impacts
4. **Testing**: Test thoroughly before deploying changes

## 📊 **Monitoring & Metrics**

### **Key Metrics to Track**

- **Build Success Rate**: Across all projects
- **Test Coverage**: Maintain 85%+ coverage
- **Performance**: Monitor performance regressions
- **Dependency Health**: Track dependency vulnerabilities

### **Alerting**

- **Build Failures**: Immediate notification for build failures
- **Test Failures**: Alert on test failures
- **Security Issues**: Immediate notification for security vulnerabilities
- **Performance Regressions**: Alert on performance degradation

## 🤝 **Collaboration Guidelines**

### **Working with Other Agents**

- **Clear Communication**: Use clear, unambiguous language
- **Context Sharing**: Share relevant context when switching projects
- **Conflict Resolution**: Have clear processes for resolving conflicts
- **Knowledge Sharing**: Document learnings and best practices

### **Handoff Procedures**

- **Status Updates**: Provide clear status updates during handoffs
- **Context Transfer**: Transfer all relevant context
- **Open Issues**: Clearly identify any open issues or blockers
- **Next Steps**: Document clear next steps

## 📞 **Support & Escalation**

### **When to Escalate**

- **Build Failures**: If builds fail and you can't resolve
- **Security Issues**: Any security-related concerns
- **Performance Issues**: Significant performance regressions
- **Architecture Decisions**: Major architectural changes

### **Escalation Process**

1. **Document Issue**: Clearly document the problem
2. **Attempt Resolution**: Show attempted solutions
3. **Provide Context**: Include all relevant context
4. **Request Help**: Clearly specify what help is needed

---

## 🎯 **Success Criteria**

AI agents working at the root level are successful when:

1. ✅ **All Projects Build**: All components build successfully
2. ✅ **Tests Pass**: All test suites pass
3. ✅ **Documentation Current**: Documentation is up-to-date
4. ✅ **Dependencies Healthy**: All dependencies are secure and compatible
5. ✅ **Performance Maintained**: No performance regressions
6. ✅ **Integration Working**: Cross-project integrations work correctly

## 📈 **Continuous Improvement**

### **Regular Reviews**

- **Weekly Reviews**: Review project status and progress
- **Monthly Assessments**: Assess overall project health
- **Quarterly Planning**: Plan improvements and new features
- **Annual Retrospectives**: Review annual progress and plan next year

### **Optimization Opportunities**

- **Build Performance**: Optimize build times and processes
- **Test Efficiency**: Improve test speed and coverage
- **Documentation**: Enhance documentation quality and accessibility
- **Developer Experience**: Improve overall developer experience

---

**Last Updated**: 2025-11-15  
**Maintained By**: AI Development Team  
**Version**: 1.0.0
