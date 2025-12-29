# 🧪 Testing Documentation

Complete testing guides untuk Sirine Go App mencakup backend, frontend, API, dan performance testing.

---

## 📚 Testing Guides

### **📖 [overview.md](./overview.md)**
Testing strategy, manual testing checklist, dan overall testing approach.

**Kapan dibaca:**
- Pertama kali setup testing
- Ingin memahami testing strategy
- Butuh manual testing checklist
- Persiapan deployment

**Isi:**
- Testing pyramid & strategy
- Manual testing procedures
- PWA & Service Worker testing
- Browser compatibility checklist
- Pre-deployment testing checklist

---

### **🔧 [backend-testing.md](./backend-testing.md)**
Panduan unit testing untuk Go backend menggunakan Go testing framework dan testify.

**Kapan dibaca:**
- Membuat backend tests
- Testing services & handlers
- Setup test database
- Troubleshooting backend tests

**Isi:**
- Setup testing environment
- Service layer unit tests
- Handler/HTTP tests dengan Gin
- Running tests & coverage reports
- Best practices & patterns

---

### **🎨 [frontend-testing.md](./frontend-testing.md)**
Panduan testing frontend Vue 3 dengan Vitest, component testing, dan PWA testing.

**Kapan dibaca:**
- Membuat frontend tests
- Testing Vue components
- Testing composables
- PWA Service Worker testing

**Isi:**
- Setup Vitest configuration
- Component testing dengan @vue/test-utils
- Composable testing patterns
- PWA & offline testing
- Running tests & coverage

---

### **🔌 [api-testing.md](./api-testing.md)**
Panduan API testing dengan cURL, Postman, dan integration testing.

**Kapan dibaca:**
- Testing API endpoints
- Setup Postman collection
- Integration testing flow
- API documentation verification

**Isi:**
- cURL testing commands
- Postman setup & automation
- Integration test scripts
- End-to-end testing scenarios

---

### **🚀 [performance-testing.md](./performance-testing.md)**
Panduan performance testing, load testing, dan browser compatibility.

**Kapan dibaca:**
- Optimization requirements
- Pre-production deployment
- Performance troubleshooting
- Cross-browser testing

**Isi:**
- Backend load testing (Apache Bench)
- Frontend performance (Lighthouse)
- Bundle size analysis
- Browser compatibility matrix
- Performance benchmarks

---

## 🗂️ Test Scenarios

Feature-specific test scenarios tersedia di folder `test-scenarios/`:

- **[user-management-testing.md](./user-management-testing.md)** - Complete test scenarios untuk User Management & Profile features

---

## 🚀 Quick Start

### **1. Backend Testing**

```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./services
go test ./handlers
```

### **2. Frontend Testing**

```bash
cd frontend

# Run tests
yarn test

# Run with coverage
yarn test:coverage

# Run in UI mode
yarn test:ui
```

### **3. API Testing**

```bash
# Quick health check
curl http://localhost:8080/health

# Test API endpoint
curl http://localhost:8080/api/users

# Or use Postman collection (see api-testing.md)
```

### **4. Performance Testing**

```bash
# Backend load test
ab -n 1000 -c 10 http://localhost:8080/health

# Frontend audit
# Build frontend → Open in Chrome → Lighthouse tab → Generate report
```

---

## 📊 Coverage Goals

### **Backend**
- **Services:** 80%+
- **Handlers:** 70%+
- **Overall:** 75%+

```bash
cd backend && go test -cover ./...
```

### **Frontend**
- **Components:** 70%+
- **Composables:** 80%+
- **Overall:** 70%+

```bash
cd frontend && yarn test:coverage
```

---

## ✅ Pre-Deployment Checklist

Sebelum deploy ke production, pastikan:

### **Backend Testing**
- [ ] All unit tests pass (`go test ./...`)
- [ ] Coverage ≥ 75%
- [ ] Load test passed (1000+ req/s)
- [ ] No memory leaks

### **Frontend Testing**
- [ ] All component tests pass (`yarn test`)
- [ ] Coverage ≥ 70%
- [ ] Lighthouse scores ≥ 90
- [ ] Bundle size < 200KB

### **Integration Testing**
- [ ] Full CRUD flow works
- [ ] Authentication flow works
- [ ] API integration verified
- [ ] Database persistence tested

### **Performance**
- [ ] API response time < 100ms
- [ ] Page load time < 3s
- [ ] 60fps animations
- [ ] No console errors

### **Cross-Browser**
- [ ] Chrome (latest)
- [ ] Firefox (latest)
- [ ] Safari (latest)
- [ ] Chrome Android
- [ ] Safari iOS

---

## 🔄 Testing Workflow

```
1. Write Code
    ↓
2. Write Tests (TDD preferred)
    ↓
3. Run Unit Tests
    ↓
4. Run Integration Tests
    ↓
5. Manual Testing
    ↓
6. Performance Testing
    ↓
7. Cross-Browser Testing
    ↓
8. Deploy
```

---

## 📚 Related Documentation

### **Development**
- [../03-development/api-documentation.md](../03-development/api-documentation.md) - API reference untuk testing
- [../03-development/customization-guide.md](../03-development/customization-guide.md) - Adding new features

### **Deployment**
- [../08-deployment/production-deployment.md](../08-deployment/production-deployment.md) - Production deployment guide

### **Troubleshooting**
- [../09-troubleshooting/faq.md](../09-troubleshooting/faq.md) - Common issues & solutions

---

## 💡 Testing Tips

### **1. Test Early, Test Often**
Write tests alongside code development, bukan setelahnya.

### **2. Follow Testing Pyramid**
Banyak unit tests (fast), beberapa integration tests (medium), sedikit E2E tests (slow).

### **3. Test Behavior, Not Implementation**
Focus pada what user sees/does, bukan internal code details.

### **4. Keep Tests Simple**
Setiap test harus focus pada satu thing dan easy to understand.

### **5. Mock External Dependencies**
Mock API calls, databases, third-party services untuk isolated tests.

### **6. Use Descriptive Names**
Test names should explain what & why:
- ✅ `TestUserService_Create_DuplicateEmail_ReturnsError`
- ❌ `TestCreate`

---

## 📞 Support

Jika ada pertanyaan tentang testing:
- **Developer:** Zulfikar Hidayatullah
- **Phone:** +62 857-1583-8733
- **Timezone:** Asia/Jakarta (WIB)

---

## 📝 Testing Documentation Structure

```
06-testing/
├── README.md                        # This file - Testing hub
├── overview.md                      # Testing strategy & manual testing
├── backend-testing.md               # Go unit testing guide
├── frontend-testing.md              # Vue/Vitest testing guide
├── api-testing.md                   # API & integration testing
├── performance-testing.md           # Performance & compatibility
└── test-scenarios/                  # Feature-specific tests
    └── user-management-testing.md   # User Management test cases
```

---

**Happy Testing! 🧪**

Comprehensive testing ensures high-quality, reliable aplikasi yang ready untuk production.

---

**Last Updated:** 28 Desember 2025  
**Version:** 2.0.0 (Phase 2 Restructure)
