# 🧪 Testing Guide - Sirine Go App

Panduan lengkap untuk testing aplikasi Sirine Go App.

---

## 📋 Daftar Isi

1. [Manual Testing](#manual-testing)
2. [Backend Testing](#backend-testing)
3. [Frontend Testing](#frontend-testing)
4. [API Testing](#api-testing)
5. [Integration Testing](#integration-testing)
6. [PWA Testing](#pwa-testing)
7. [Performance Testing](#performance-testing)
8. [Browser Compatibility](#browser-compatibility)

---

## ✅ Manual Testing

### Quick Health Check

Setelah setup, lakukan quick test ini:

```bash
# 1. Backend health
curl http://localhost:8080/health

# Expected: {"status":"ok","message":"Server berjalan dengan baik"}

# 2. API response
curl http://localhost:8080/api/examples

# Expected: {"data":[]} atau array of examples

# 3. Frontend accessible
curl -I http://localhost:5173

# Expected: HTTP/1.1 200 OK
```

---

### UI Testing Checklist

Buka browser: `http://localhost:5173`

#### **Layout & Design**
- [ ] Header muncul dengan benar
- [ ] Status indicator (Online/Offline) visible
- [ ] Button "Tambah Data Baru" terlihat jelas
- [ ] Responsive di berbagai ukuran layar
- [ ] Tailwind CSS styles applied
- [ ] No layout broken atau overlap

#### **Functionality - Create**
- [ ] Klik "Tambah Data Baru" → Form muncul
- [ ] Form fields: Judul, Konten, Aktif checkbox
- [ ] Input validation bekerja (required fields)
- [ ] Submit button → Loading indicator muncul
- [ ] Success → Data muncul di list
- [ ] Error → Error message muncul
- [ ] Form close setelah success

#### **Functionality - Read**
- [ ] List data muncul dengan benar
- [ ] Card menampilkan: Title, Content, Status, Timestamps
- [ ] Empty state muncul jika tidak ada data
- [ ] Animasi smooth saat card muncul (Motion-v)

#### **Functionality - Update**
- [ ] Klik "Edit" button → Form muncul dengan data
- [ ] Data pre-filled di form
- [ ] Update data → Success message
- [ ] Data terupdate di list

#### **Functionality - Delete**
- [ ] Klik "Hapus" button → Confirmation dialog
- [ ] Confirm → Data terhapus
- [ ] Cancel → Data tidak berubah
- [ ] Success message muncul

---

## 🔧 Backend Testing

### Unit Testing dengan Go

Create test files untuk setiap service/handler.

#### **1. Service Test Example**

```go
// backend/services/example_service_test.go
package services

import (
    "testing"
    "sirine-go/backend/database"
    "sirine-go/backend/models"
    "github.com/stretchr/testify/assert"
)

func TestExampleService_Create(t *testing.T) {
    // Setup test database
    database.ConnectTest()
    defer database.CloseTest()
    
    // Create service
    service := NewExampleService()
    
    // Test data
    example := &models.Example{
        Title:    "Test Example",
        Content:  "Test Content",
        IsActive: true,
    }
    
    // Test create
    err := service.Create(example)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotZero(t, example.ID)
    assert.Equal(t, "Test Example", example.Title)
}

func TestExampleService_GetAll(t *testing.T) {
    database.ConnectTest()
    defer database.CloseTest()
    
    service := NewExampleService()
    
    // Test get all
    examples, err := service.GetAll()
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, examples)
}
```

#### **2. Handler Test Example**

```go
// backend/handlers/example_handler_test.go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "bytes"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestExampleHandler_GetAll(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    
    service := services.NewExampleService()
    handler := NewExampleHandler(service)
    
    r.GET("/api/examples", handler.GetAll)
    
    // Test request
    req, _ := http.NewRequest("GET", "/api/examples", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "data")
}

func TestExampleHandler_Create(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    
    service := services.NewExampleService()
    handler := NewExampleHandler(service)
    
    r.POST("/api/examples", handler.Create)
    
    // Test data
    data := map[string]interface{}{
        "title":     "Test",
        "content":   "Test Content",
        "is_active": true,
    }
    jsonData, _ := json.Marshal(data)
    
    // Test request
    req, _ := http.NewRequest("POST", "/api/examples", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusCreated, w.Code)
    assert.Contains(t, w.Body.String(), "berhasil dibuat")
}
```

#### **3. Run Backend Tests**

```bash
# Run all tests
cd backend
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./services
go test ./handlers

# Verbose output
go test -v ./...

# Run specific test
go test -run TestExampleService_Create ./services
```

---

## 🎨 Frontend Testing

### Component Testing dengan Vitest

Frontend testing menggunakan Vitest (built-in dengan Vite).

#### **1. Setup Vitest**

```bash
cd frontend
yarn add -D vitest @vue/test-utils jsdom
```

**Update `vite.config.js`:**
```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    globals: true,
    environment: 'jsdom'
  }
})
```

**Update `package.json`:**
```json
{
  "scripts": {
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage"
  }
}
```

#### **2. Component Test Example**

```javascript
// frontend/src/components/__tests__/ExampleCard.test.js
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ExampleCard from '../ExampleCard.vue'

describe('ExampleCard', () => {
  it('renders example data correctly', () => {
    const example = {
      id: 1,
      title: 'Test Title',
      content: 'Test Content',
      is_active: true
    }
    
    const wrapper = mount(ExampleCard, {
      props: { example }
    })
    
    expect(wrapper.text()).toContain('Test Title')
    expect(wrapper.text()).toContain('Test Content')
  })
  
  it('emits edit event when edit button clicked', async () => {
    const example = { id: 1, title: 'Test' }
    const wrapper = mount(ExampleCard, {
      props: { example }
    })
    
    await wrapper.find('[data-test="edit-button"]').trigger('click')
    
    expect(wrapper.emitted('edit')).toBeTruthy()
    expect(wrapper.emitted('edit')[0]).toEqual([example])
  })
  
  it('shows active badge when is_active is true', () => {
    const example = { id: 1, title: 'Test', is_active: true }
    const wrapper = mount(ExampleCard, {
      props: { example }
    })
    
    expect(wrapper.text()).toContain('Aktif')
  })
})
```

#### **3. Composable Test Example**

```javascript
// frontend/src/composables/__tests__/useExamples.test.js
import { describe, it, expect, vi } from 'vitest'
import { useExamples } from '../useExamples'

// Mock API
vi.mock('../useApi', () => ({
  useApi: () => ({
    get: vi.fn(() => Promise.resolve({ data: { data: [] } })),
    post: vi.fn(() => Promise.resolve({ data: { data: {} } }))
  })
}))

describe('useExamples', () => {
  it('fetches examples successfully', async () => {
    const { examples, fetchExamples } = useExamples()
    
    await fetchExamples()
    
    expect(examples.value).toBeInstanceOf(Array)
  })
  
  it('creates example successfully', async () => {
    const { createExample } = useExamples()
    
    const newExample = {
      title: 'New Test',
      content: 'Content',
      is_active: true
    }
    
    const result = await createExample(newExample)
    
    expect(result).toBeDefined()
  })
})
```

#### **4. Run Frontend Tests**

```bash
cd frontend

# Run tests
yarn test

# Watch mode
yarn test --watch

# Coverage
yarn test:coverage

# UI mode (interactive)
yarn test:ui
```

---

## 🔌 API Testing

### Testing dengan cURL

#### **Health Check**
```bash
curl http://localhost:8080/health
```

#### **GET All Examples**
```bash
curl http://localhost:8080/api/examples
```

#### **GET by ID**
```bash
curl http://localhost:8080/api/examples/1
```

#### **POST Create**
```bash
curl -X POST http://localhost:8080/api/examples \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Example",
    "content": "Test content here",
    "is_active": true
  }'
```

#### **PUT Update**
```bash
curl -X PUT http://localhost:8080/api/examples/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content",
    "is_active": false
  }'
```

#### **DELETE**
```bash
curl -X DELETE http://localhost:8080/api/examples/1
```

---

### Testing dengan Postman

#### **1. Create Collection**

1. Open Postman
2. Create new Collection: "Sirine Go API"
3. Add requests untuk setiap endpoint

#### **2. Environment Variables**

Create environment dengan variables:
- `base_url`: `http://localhost:8080`
- `example_id`: `1`

#### **3. Example Requests**

**Health Check:**
```
GET {{base_url}}/health
```

**Get All Examples:**
```
GET {{base_url}}/api/examples
```

**Create Example:**
```
POST {{base_url}}/api/examples
Content-Type: application/json

{
  "title": "Test from Postman",
  "content": "Content here",
  "is_active": true
}

# Save response ID
Tests > Set variable:
pm.environment.set("example_id", pm.response.json().data.id)
```

**Update Example:**
```
PUT {{base_url}}/api/examples/{{example_id}}
Content-Type: application/json

{
  "title": "Updated from Postman"
}
```

**Delete Example:**
```
DELETE {{base_url}}/api/examples/{{example_id}}
```

#### **4. Automated Tests in Postman**

Add tests di Postman:

```javascript
// Health Check tests
pm.test("Status code is 200", function() {
    pm.response.to.have.status(200);
});

pm.test("Response has status ok", function() {
    var json = pm.response.json();
    pm.expect(json.status).to.eql("ok");
});

// Create tests
pm.test("Status code is 201", function() {
    pm.response.to.have.status(201);
});

pm.test("Response has data", function() {
    var json = pm.response.json();
    pm.expect(json.data).to.exist;
    pm.expect(json.data.id).to.exist;
});
```

---

## 🔗 Integration Testing

Test full flow dari frontend ke backend ke database.

### Scenario 1: Create and Verify

```bash
# 1. Create via API
response=$(curl -X POST http://localhost:8080/api/examples \
  -H "Content-Type: application/json" \
  -d '{"title":"Integration Test","content":"Test","is_active":true}')

# 2. Extract ID
id=$(echo $response | jq -r '.data.id')

# 3. Verify via GET
curl http://localhost:8080/api/examples/$id

# 4. Verify in database
mysql -u root -p sirine_go -e "SELECT * FROM examples WHERE id=$id;"
```

### Scenario 2: Update and Verify

```bash
# 1. Update
curl -X PUT http://localhost:8080/api/examples/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated via Test"}'

# 2. Verify
curl http://localhost:8080/api/examples/1

# 3. Check database
mysql -u root -p sirine_go -e "SELECT title FROM examples WHERE id=1;"
```

### Scenario 3: Delete and Verify

```bash
# 1. Delete
curl -X DELETE http://localhost:8080/api/examples/1

# 2. Try to get (should 404)
curl http://localhost:8080/api/examples/1

# 3. Check soft delete in database
mysql -u root -p sirine_go -e "SELECT id, title, deleted_at FROM examples WHERE id=1;"
```

---

## 🌐 PWA Testing

### Service Worker Testing

#### **1. Register Service Worker**

1. Build frontend: `cd frontend && yarn build`
2. Start backend (will serve dist): `make dev-backend`
3. Open browser: `http://localhost:8080`
4. Open DevTools (F12) → Application → Service Workers
5. Verify: Service Worker status = "activated and running"

#### **2. Cache Testing**

**Check cache storage:**
1. DevTools → Application → Cache Storage
2. Verify caches exist:
   - `workbox-precache-v2-...` (static assets)
   - `workbox-runtime-...` (API responses)

**Test caching:**
```bash
# 1. Load page online
# 2. Check Network tab - resources loaded from network
# 3. Reload page
# 4. Check Network tab - resources loaded from Service Worker
```

#### **3. Offline Testing**

1. Open app: `http://localhost:8080`
2. Load completely (all data fetched)
3. DevTools (F12) → Network tab
4. Set "Offline" (dropdown at top)
5. Reload page
6. ✅ Page should load from cache
7. ✅ Status indicator shows "Offline"
8. ✅ Previously loaded data accessible

#### **4. PWA Installation Testing**

**Desktop:**
1. Chrome address bar → Install icon (⊕)
2. Click "Install"
3. App opens in standalone window
4. Verify app in Applications folder

**Mobile:**
1. Open in Chrome mobile
2. Menu (⋮) → "Add to Home Screen"
3. Tap "Add"
4. Icon appears on home screen
5. Tap icon → Opens as app

---

## 🚀 Performance Testing

### Backend Performance

#### **Load Testing dengan Apache Bench**

```bash
# Install
sudo apt install apache2-utils

# Test health endpoint
ab -n 1000 -c 10 http://localhost:8080/health

# Test API endpoint
ab -n 1000 -c 10 http://localhost:8080/api/examples

# Results analysis:
# - Requests per second
# - Time per request
# - Failed requests
```

#### **Expected Results**

```
Requests per second:    1000+ req/s
Time per request:       < 10ms (average)
Failed requests:        0
```

---

### Frontend Performance

#### **Lighthouse Audit**

1. Open app: `http://localhost:8080` (production build)
2. DevTools (F12) → Lighthouse tab
3. Select:
   - ✅ Performance
   - ✅ Progressive Web App
   - ✅ Best Practices
   - ✅ Accessibility
   - ✅ SEO
4. Click "Generate report"

**Expected Scores:**
- Performance: 90+
- PWA: 90+
- Best Practices: 90+
- Accessibility: 85+
- SEO: 90+

#### **Bundle Size Analysis**

```bash
cd frontend
yarn build

# Check bundle sizes
ls -lh dist/assets/

# Expected:
# - index.*.js: < 200KB (gzipped < 70KB)
# - index.*.css: < 50KB (gzipped < 10KB)
```

---

## 🌍 Browser Compatibility

### Desktop Testing

Test di browsers:
- ✅ Chrome 90+ (primary)
- ✅ Firefox 88+ (secondary)
- ✅ Safari 14+ (macOS)
- ✅ Edge 90+ (Windows)

**Test checklist per browser:**
- [ ] App loads correctly
- [ ] All features work
- [ ] Animations smooth
- [ ] PWA installable
- [ ] Service Worker works
- [ ] No console errors

### Mobile Testing

Test di devices:
- ✅ Chrome Android 90+
- ✅ Safari iOS 14+

**Test checklist mobile:**
- [ ] Responsive layout
- [ ] Touch interactions work
- [ ] Gestures (swipe, pinch) work
- [ ] PWA installable
- [ ] Offline mode works
- [ ] Performance good (no lag)

---

## 📊 Test Coverage Goals

### Backend Coverage Target
- Services: 80%+
- Handlers: 70%+
- Overall: 75%+

```bash
cd backend
go test -cover ./...
```

### Frontend Coverage Target
- Components: 70%+
- Composables: 80%+
- Overall: 70%+

```bash
cd frontend
yarn test:coverage
```

---

## ✅ Testing Checklist

Sebelum deploy to production, pastikan:

### Backend
- [ ] All unit tests pass
- [ ] All endpoints tested
- [ ] Error handling tested
- [ ] Database connections work
- [ ] Performance acceptable

### Frontend
- [ ] All component tests pass
- [ ] UI responsive di semua screen sizes
- [ ] All user flows work
- [ ] No console errors
- [ ] Animations smooth

### Integration
- [ ] Full CRUD flow works
- [ ] API integration works
- [ ] Database persistence works
- [ ] Error messages correct

### PWA
- [ ] Service Worker registered
- [ ] Offline mode works
- [ ] App installable
- [ ] Cache working properly

### Performance
- [ ] Lighthouse scores 90+
- [ ] Load time < 3s
- [ ] API response < 100ms
- [ ] No memory leaks

---

## 📞 Support

Jika ada pertanyaan tentang testing:
- Developer: Zulfikar Hidayatullah
- Phone: +62 857-1583-8733

**Related Documentation:**
- [FAQ.md](./FAQ.md) - Common issues
- [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) - API reference
- [SETUP_GUIDE.md](./SETUP_GUIDE.md) - Setup guide

---

**Last Updated:** 27 Desember 2025  
**Version:** 1.0.0
