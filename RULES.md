# Qwen Code - Rules and Guidelines

## 🎯 **Project Rules and Standards**

This document defines the rules, guidelines, and standards for working with the Qwen Code project. All contributors, AI agents, and automated systems must follow these rules.

---

## 📋 **Core Principles**

### 1. **Code Quality First**

- All code must pass linting and type checking
- Test coverage minimum: 80%
- No hardcoded secrets or credentials
- Follow established coding conventions

### 2. **Documentation Required**

- All public APIs must have documentation
- Complex logic requires inline comments
- New features need user documentation
- Maintain README and guides

### 3. **Security by Design**

- Validate all user inputs
- Use secure coding practices
- Regular security audits required
- Report vulnerabilities immediately

### 4. **Performance Matters**

- Monitor resource usage
- Optimize for speed and memory
- Profile before optimizing
- Maintain performance benchmarks

---

## 🔧 **Development Rules**

### **Code Standards**

```typescript
// Use TypeScript strict mode
// Follow ESLint rules
// Use Prettier for formatting
// Prefer explicit types over 'any'

// Example function structure
/**
 * Brief description of function purpose
 * @param param1 - Description of parameter
 * @returns Description of return value
 * @throws {Error} When and why error occurs
 */
export function functionName(param1: string): ReturnType {
  // Implementation
}
```

### **File Organization**

```
src/
├── components/     # React/UI components
├── services/       # Business logic
├── utils/          # Helper functions
├── types/          # TypeScript definitions
├── tests/          # Test files
└── docs/           # Documentation
```

### **Naming Conventions**

- **Files**: kebab-case (e.g., `my-component.ts`)
- **Functions**: camelCase (e.g., `myFunction`)
- **Constants**: UPPER_SNAKE_CASE (e.g., `API_URL`)
- **Classes**: PascalCase (e.g., `MyClass`)
- **Interfaces**: PascalCase with 'I' prefix (e.g., `IMyInterface`)

---

## 🧪 **Testing Rules**

### **Test Requirements**

- Unit tests for all functions
- Integration tests for workflows
- E2E tests for user journeys
- Performance tests for critical paths

### **Test Structure**

```typescript
describe('Component/Function', () => {
  beforeEach(() => {
    // Setup
  });

  it('should handle expected case', () => {
    // Test
  });

  it('should handle edge case', () => {
    // Edge case test
  });

  it('should handle error case', () => {
    // Error handling test
  });
});
```

### **Coverage Requirements**

- **Statements**: 80% minimum
- **Branches**: 75% minimum
- **Functions**: 85% minimum
- **Lines**: 80% minimum

---

## 📝 **Documentation Rules**

### **Required Documentation**

1. **README.md**: Project overview and setup
2. **API Documentation**: All public APIs
3. **User Guides**: How to use features
4. **Development Guide**: Contributing guidelines
5. **Architecture Documentation**: System design

### **Documentation Standards**

- Use clear, concise language
- Include code examples
- Maintain consistent formatting
- Update with code changes

---

## 🔒 **Security Rules**

### **Input Validation**

```typescript
// Always validate inputs
function processInput(input: string): void {
  if (!input || input.trim().length === 0) {
    throw new Error('Input cannot be empty');
  }

  if (input.length > MAX_LENGTH) {
    throw new Error('Input too long');
  }

  // Process validated input
}
```

### **Security Checklist**

- [ ] No hardcoded secrets
- [ ] Input validation implemented
- [ ] SQL injection protection
- [ ] XSS prevention
- [ ] CSRF protection
- [ ] Authentication/authorization
- [ ] HTTPS enforcement
- [ ] Security headers set

---

## 🚀 **Deployment Rules**

### **Pre-deployment Checklist**

- [ ] All tests passing
- [ ] Code coverage meets requirements
- [ ] Security scan passed
- [ ] Performance benchmarks met
- [ ] Documentation updated
- [ ] Version number updated
- [ ] Change log updated
- [ ] Backup created

### **Environment Rules**

- **Development**: Use test data and mock services
- **Staging**: Mirror production configuration
- **Production**: Use real data and services
- **No direct access** to production databases

---

## 🤖 **AI Agent Rules**

### **Agent Behavior Guidelines**

1. **Be Helpful**: Provide accurate, useful information
2. **Be Safe**: Don't execute harmful commands
3. **Be Clear**: Use simple, understandable language
4. **Be Honest**: Admit limitations and uncertainties
5. **Be Efficient**: Minimize resource usage

### **Agent Capabilities**

- ✅ **Allowed**: Code analysis, documentation, testing
- ✅ **Allowed**: File operations within project bounds
- ✅ **Allowed**: Running tests and builds
- ❌ **Forbidden**: System-level commands
- ❌ **Forbidden**: Network access to external sites
- ❌ **Forbidden**: Accessing user files outside project

### **Agent Communication**

- Use structured responses
- Provide context for actions
- Explain reasoning when helpful
- Ask for clarification when needed

---

## 📊 **Performance Rules**

### **Performance Targets**

- **Startup Time**: < 3 seconds
- **API Response**: < 200ms
- **Memory Usage**: < 512MB
- **CPU Usage**: < 50% during normal operation
- **Build Time**: < 5 minutes

### **Monitoring Requirements**

- Track key metrics
- Set up alerts for anomalies
- Log performance data
- Regular performance reviews

---

## 🔧 **Tool Usage Rules**

### **Development Tools**

- **IDE**: VS Code with recommended extensions
- **Version Control**: Git with feature branches
- **Package Manager**: npm for Node.js projects
- **Build Tools**: Use configured build scripts
- **Linting**: ESLint and Prettier required

### **Code Review Process**

1. Create pull request
2. Automated checks must pass
3. At least one human review required
4. Address all feedback
5. Merge after approval

---

## 📋 **File Management Rules**

### **File Naming**

- Use descriptive names
- Follow established patterns
- Avoid spaces in filenames
- Use appropriate extensions

### **File Organization**

- Group related files together
- Keep directory structure shallow
- Use index files for exports
- Separate concerns properly

### **Git Rules**

- **Commit Messages**: Use conventional format
- **Branch Names**: Use feature/fix prefixes
- **Pull Requests**: Provide clear descriptions
- **Merge Strategy**: Use squash and merge

---

## 🎨 **UI/UX Rules**

### **Design Principles**

- **Consistency**: Follow established patterns
- **Accessibility**: WCAG 2.1 AA compliance
- **Responsive**: Mobile-first design
- **Performance**: Optimize for speed

### **Component Rules**

- Reusable components preferred
- Proper error states
- Loading indicators
- Clear user feedback

---

## 📚 **Documentation Usage Instructions**

### **For Developers**

1. **Start with README.md**: Project overview and setup
2. **Read USAGE_AND_PROCEDURES.md**: How to use the system
3. **Consult features.md**: Understand available features
4. **Follow development.md**: Development guidelines
5. **Reference SRS.md**: Requirements and specifications

### **For Users**

1. **README.md**: Quick start guide
2. **USAGE_AND_PROCEDURES.md**: Detailed usage instructions
3. **features.md**: Feature descriptions and examples

### **For Contributors**

1. **development.md**: Development setup and guidelines
2. **tasks.md**: Task management procedures
3. **journal.md**: Session logging and tracking
4. **project_summary.md**: Project overview and status

### **For AI Agents**

1. **agents.md**: AI agent guidelines and protocols
2. **SRS.md**: System requirements and constraints
3. **rules.md**: This file - all rules and standards
4. **features.md**: Feature specifications and capabilities

---

## 🚨 **Emergency Procedures**

### **Security Incident**

1. Immediately stop affected services
2. Notify security team
3. Document incident details
4. Implement fixes
5. Review and improve procedures

### **Production Outage**

1. Assess impact and scope
2. Communicate with stakeholders
3. Implement fixes
4. Monitor recovery
5. Conduct post-mortem

---

## 📈 **Continuous Improvement**

### **Regular Reviews**

- **Weekly**: Code quality metrics
- **Monthly**: Security assessments
- **Quarterly**: Performance reviews
- **Annually**: Architecture evaluation

### **Feedback Process**

1. Collect feedback from users
2. Analyze and categorize
3. Prioritize improvements
4. Implement changes
5. Measure impact

---

## ⚖️ **Compliance and Legal**

### **License Compliance**

- Respect all open source licenses
- Track third-party dependencies
- Include license notices
- Review legal requirements

### **Data Protection**

- Follow GDPR/CCPA guidelines
- Minimize data collection
- Implement proper security
- Respect user privacy

---

## 🎯 **Success Metrics**

### **Quality Metrics**

- Code coverage > 80%
- Zero critical security issues
- Performance targets met
- User satisfaction > 4.5/5

### **Development Metrics**

- Feature delivery time
- Bug fix response time
- Code review turnaround
- Documentation completeness

---

## 📞 **Support and Escalation**

### **Getting Help**

1. Check documentation first
2. Search existing issues
3. Create detailed issue report
4. Follow escalation procedures

### **Issue Reporting**

- Provide clear description
- Include reproduction steps
- Add environment details
- Attach relevant logs

---

**Last Updated**: 2025-11-17  
**Version**: 1.0.0  
**Maintainers**: Qwen Code Development Team

---

_These rules ensure consistent, high-quality development of the Qwen Code project. All contributors and automated systems must follow these guidelines._
