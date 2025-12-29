# Call of Cthulhu Character Sheet Project - Remaining Work Plan

## ✅ Completed

### Phase 1: Critical Fixes
- ✅ Database connection pooling with singleton pattern
- ✅ Custom error types and proper error handling
- ✅ Fixed resource leaks with proper defer statements
- ✅ Dependency injection for handlers
- ✅ Environment-based configuration management
- ✅ Graceful server shutdown

### Phase 2: Handler Migration
- ✅ Move PDF export handler to new handler structure
- ✅ Move archetype occupations handler to new structure
- ✅ Removed old handlers.go and export.go files
- ✅ Created comprehensive Makefile for build/deploy
- ✅ Update wizard package to use dependency injection
- ✅ Remove old storage/cookies.go and sql.go
- ✅ Update bug reporting to use new error handling
- ✅ Removed old bugReporting package

### Phase 2.1: Middleware Layer
- ✅ Request logging middleware
- ✅ Error recovery middleware (panic recovery)
- ✅ Security headers middleware
- ✅ Request ID middleware

## 📋 Remaining Work

### Phase 3: Testing & Validation (Priority: High)

#### 1. Add comprehensive tests
- [ ] Unit tests for storage layer
- [ ] Unit tests for handlers
- [ ] Integration tests for API endpoints
- [ ] Test database operations and cleanup

#### 2. Input validation
- [ ] Validate character attribute ranges (1-100)
- [ ] Validate skill values
- [ ] Sanitize string inputs
- [ ] Add request size limits

### Phase 4: Frontend Improvements (Priority: Medium)

#### 1. JavaScript refactoring
- [x] Consolidated duplicate functions into modular structure
- [x] Improved error handling in fetch requests (centralized in API module)
- [x] Add loading states for async operations (Utils.setButtonLoading)
- [x] Created modular JS architecture:
  - `static/js/utils.js` - Shared utilities (DOM helpers, dice rolling, calculations)
  - `static/js/api.js` - Centralized API calls with error handling
  - `static/js/wizard.js` - Character creation wizard logic
  - `static/js/character-sheet.js` - Character sheet editing
  - `static/js/app.js` - Main entry point with backward compatibility layer

#### 2. UI/UX enhancements
- [x] Added toast notifications for user feedback (Utils.showToast)
- [x] Improved form validation feedback (Utils.showSuccess/showError/showInvalid)
- [ ] Add confirmation dialogs for delete operations

### Phase 5: API Standardization (Priority: Medium)

#### 1. Consistent API responses
- [ ] Standardize JSON field naming (use camelCase throughout)
- [ ] Create standard error response format
- [ ] Document all API endpoints

#### 2. API improvements
- [ ] Add pagination for investigator lists
- [ ] Add filtering/search capabilities
- [ ] Version the API (/api/v1/)

### Phase 6: Security & Performance (Priority: Low)

#### 1. Security enhancements
- [ ] Add rate limiting middleware
- [ ] Implement CSRF protection
- [ ] Add Content Security Policy headers
- [ ] Validate cookie sizes before setting

#### 2. Performance optimization
- [ ] Cache static data (archetypes, occupations)
- [ ] Optimize PDF generation
- [ ] Add compression middleware
- [ ] Implement connection pooling for high traffic

## 🎯 Next Sprint (Immediate Focus)

1. **Write unit tests for storage layer**
2. **Write unit tests for handlers**
3. **Add input validation middleware**
4. **Add confirmation dialogs for delete operations**
5. **Standardize JSON field naming (camelCase)**

## 📊 Success Metrics

- **Test Coverage**: Minimum 70% for critical paths
- **Response Time**: < 200ms for API endpoints
- **Error Rate**: < 1% for valid requests
- **Code Quality**: Pass `go vet` and `golint`

## 🚀 Quick Start for Development

```bash
# Run the server
go run main.go

# Run tests (once written)
go test ./...

# Build for production
go build -o book-of-shadows

# Environment variables
export SERVER_PORT=8080
export DB_PATH=data/exports.db
export COOKIE_PREFIX=investigator
```

## 📝 Notes

The refactoring has successfully implemented idiomatic Go patterns while maintaining:
- **Zero-cost infrastructure** using cookies and SQLite
- **Vanilla JavaScript** without heavy frameworks
- **Simple deployment** with minimal dependencies