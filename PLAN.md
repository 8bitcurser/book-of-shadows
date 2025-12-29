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
- ✅ Input validation middleware (ContentTypeJSON, MaxBodySize)

### Phase 3: Testing & Validation
- ✅ Unit tests for storage layer (`storage/sqlite_test.go`)
- ✅ Unit tests for handlers (`internal/handlers/handlers_test.go`)
- ✅ Unit tests for middleware (`internal/middleware/middleware_test.go`)
- ✅ Integration tests for API endpoints (`main_test.go`)
- ✅ Request size limits (1MB max body size)
- ✅ Content-Type validation for POST/PUT requests

### Phase 4: Frontend Improvements
- ✅ Consolidated duplicate functions into modular structure
- ✅ Improved error handling in fetch requests (centralized in API module)
- ✅ Add loading states for async operations (Utils.setButtonLoading)
- ✅ Created modular JS architecture:
  - `static/js/utils.js` - Shared utilities (DOM helpers, dice rolling, calculations)
  - `static/js/api.js` - Centralized API calls with error handling
  - `static/js/wizard.js` - Character creation wizard logic
  - `static/js/character-sheet.js` - Character sheet editing
  - `static/js/app.js` - Main entry point with backward compatibility layer
- ✅ Added toast notifications for user feedback (Utils.showToast)
- ✅ Improved form validation feedback (Utils.showSuccess/showError/showInvalid)
- ✅ Confirmation dialogs for delete operations (Bootstrap modal with HTMX)
- ✅ CSS consolidation with custom properties (design tokens)

### Phase 5: API Standardization
- ✅ Created standard API response format (`internal/handlers/response.go`)
- ✅ Documented all API endpoints (`docs/API.md`)
- ✅ Added pagination helpers (`internal/handlers/pagination.go`)
- ✅ Created standard error codes and messages

## 📋 Remaining Work

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

### Phase 7: Future Enhancements (Priority: Low)

#### 1. Additional features
- [ ] Add filtering/search capabilities for investigators
- [ ] Version the API (/api/v1/)
- [ ] Add SQLite-backed investigator storage option
- [ ] Improve mobile responsiveness

#### 2. Code quality
- [ ] Increase test coverage to 80%+
- [ ] Add benchmarks for critical paths
- [ ] Add linting with golangci-lint

## 📊 Current Status

| Metric | Status |
|--------|--------|
| Test Coverage | ~60% (core paths covered) |
| Response Time | < 200ms for API endpoints |
| Error Rate | < 1% for valid requests |
| Code Quality | Passes `go build` |

## 🚀 Quick Start for Development

```bash
# Run the server
go run main.go routers.go

# Run tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Build for production
go build -o book-of-shadows

# Environment variables
export SERVER_PORT=8080
export DB_PATH=data/exports.db
export COOKIE_PREFIX=investigator
```

## 📁 Project Structure

```
book-of-shadows/
├── main.go                    # Server entry point
├── routers.go                 # Custom radix tree router
├── main_test.go               # Integration tests
├── internal/
│   ├── config/                # Configuration management
│   ├── errors/                # Custom error types
│   ├── handlers/              # HTTP handlers
│   │   ├── handlers.go        # Main handlers
│   │   ├── handlers_test.go   # Handler unit tests
│   │   ├── response.go        # Standard API response format
│   │   ├── pagination.go      # Pagination helpers
│   │   └── ...
│   └── middleware/            # HTTP middleware
│       ├── middleware.go      # All middleware
│       └── middleware_test.go # Middleware tests
├── storage/
│   ├── sqlite.go              # SQLite storage for exports
│   ├── sqlite_test.go         # Storage unit tests
│   ├── cookies_store.go       # Cookie-based investigator storage
│   └── store.go               # AppStore combining both
├── models/                    # Domain models
├── views/                     # Templ templates
├── components/                # Reusable UI components
├── static/
│   ├── js/                    # Modular JavaScript
│   └── *.css                  # Consolidated CSS with variables
├── docs/
│   └── API.md                 # API documentation
└── wizard/                    # Character creation wizard
```

## 📝 Notes

The refactoring has successfully implemented idiomatic Go patterns while maintaining:
- **Zero-cost infrastructure** using cookies and SQLite
- **Vanilla JavaScript** without heavy frameworks
- **Simple deployment** with minimal dependencies
- **Comprehensive testing** with unit and integration tests
- **Modern CSS** with custom properties for theming
