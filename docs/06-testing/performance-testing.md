# 🚀 Performance Testing Guide

Panduan untuk testing performance aplikasi, load testing, dan browser compatibility.

---

## 📋 Daftar Isi

1. [Backend Performance](#backend-performance)
2. [Frontend Performance](#frontend-performance)
3. [Browser Compatibility](#browser-compatibility)

---

## ⚡ Backend Performance

### **Load Testing dengan Apache Bench**

Apache Bench (ab) adalah command-line tool untuk load testing HTTP servers.

#### **Install Apache Bench**

```bash
# Ubuntu/Debian
sudo apt install apache2-utils

# macOS (included with Apache)
# Already installed

# Verify installation
ab -V
```

#### **Basic Load Test**

```bash
# Test health endpoint
ab -n 1000 -c 10 http://localhost:8080/health

# Parameters:
# -n 1000 = Total 1000 requests
# -c 10   = 10 concurrent requests
```

**Output Explanation:**

```
Concurrency Level:      10
Time taken for tests:   0.523 seconds
Complete requests:      1000
Failed requests:        0
Total transferred:      156000 bytes
HTML transferred:       43000 bytes
Requests per second:    1912.05 [#/sec] (mean)
Time per request:       5.230 [ms] (mean)
Time per request:       0.523 [ms] (mean, across all concurrent requests)
Transfer rate:          291.39 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       1
Processing:     2    5   2.3      4      18
Waiting:        2    5   2.3      4      18
Total:          2    5   2.3      5      18

Percentage of the requests served within a certain time (ms)
  50%      5
  66%      6
  75%      6
  80%      7
  90%      8
  95%     10
  98%     13
  99%     15
 100%     18 (longest request)
```

#### **Test API Endpoints**

```bash
# Test GET endpoint
ab -n 1000 -c 10 http://localhost:8080/api/users

# Test with authentication header
ab -n 1000 -c 10 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/profile

# Test POST endpoint dengan data
ab -n 100 -c 10 \
  -p data.json \
  -T "application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/users
```

**Create data.json:**
```json
{
  "name": "Load Test User",
  "email": "loadtest@example.com",
  "password": "password123",
  "role": "user",
  "department": "IT"
}
```

#### **Expected Backend Performance**

Target metrics untuk Gin backend:

| Metric | Target | Acceptable |
|--------|--------|------------|
| **Requests/second** | 1000+ | 500+ |
| **Mean response time** | <10ms | <50ms |
| **95th percentile** | <20ms | <100ms |
| **Failed requests** | 0 | <1% |
| **CPU usage** | <50% | <80% |
| **Memory usage** | <500MB | <1GB |

#### **Stress Testing**

Test aplikasi under extreme load untuk find breaking point:

```bash
# Gradually increase load
ab -n 10000 -c 100 http://localhost:8080/health  # Heavy load
ab -n 50000 -c 500 http://localhost:8080/health  # Extreme load

# Monitor server resources:
# CPU: htop
# Memory: free -h
# Connections: netstat -an | grep :8080 | wc -l
```

---

## 🎨 Frontend Performance

### **Lighthouse Audit**

Lighthouse provides comprehensive performance analysis untuk web apps.

#### **Run Lighthouse Audit**

1. **Build production version:**
   ```bash
   cd frontend
   yarn build
   ```

2. **Serve production build:**
   ```bash
   # Backend serves frontend dist/ folder
   cd backend
   go run cmd/server/main.go
   ```

3. **Open app:** `http://localhost:8080`

4. **Run Lighthouse:**
   - Open Chrome DevTools (F12)
   - Go to "Lighthouse" tab
   - Select categories:
     - ✅ Performance
     - ✅ Progressive Web App
     - ✅ Best Practices
     - ✅ Accessibility
     - ✅ SEO
   - Select "Desktop" or "Mobile"
   - Click "Generate report"

#### **Target Lighthouse Scores**

| Category | Target | Minimum |
|----------|--------|---------|
| **Performance** | 90+ | 80+ |
| **PWA** | 90+ | 85+ |
| **Best Practices** | 95+ | 90+ |
| **Accessibility** | 90+ | 85+ |
| **SEO** | 95+ | 90+ |

#### **Key Performance Metrics**

**Core Web Vitals:**

1. **LCP (Largest Contentful Paint)**
   - Target: <2.5s
   - Acceptable: <4.0s
   - Measures: Loading performance

2. **FID (First Input Delay)**
   - Target: <100ms
   - Acceptable: <300ms
   - Measures: Interactivity

3. **CLS (Cumulative Layout Shift)**
   - Target: <0.1
   - Acceptable: <0.25
   - Measures: Visual stability

**Other Metrics:**

| Metric | Target | Description |
|--------|--------|-------------|
| **FCP (First Contentful Paint)** | <1.8s | Time to first content |
| **SI (Speed Index)** | <3.4s | How quickly content visible |
| **TTI (Time to Interactive)** | <3.8s | Time until fully interactive |
| **TBT (Total Blocking Time)** | <200ms | Time main thread blocked |

#### **Lighthouse Optimization Tips**

**Performance Optimizations:**

1. **Minimize JavaScript:**
   ```bash
   # Check bundle size
   cd frontend
   yarn build
   ls -lh dist/assets/
   ```

   - Target: index.*.js < 200KB (gzipped < 70KB)
   - Target: index.*.css < 50KB (gzipped < 10KB)

2. **Enable Compression:**
   ```go
   // backend/main.go
   router.Use(gzip.Gzip(gzip.DefaultCompression))
   ```

3. **Lazy Load Images:**
   ```vue
   <img loading="lazy" src="..." alt="..." />
   ```

4. **Preload Critical Assets:**
   ```html
   <link rel="preload" href="/fonts/inter.woff2" as="font" type="font/woff2" crossorigin>
   ```

### **Bundle Size Analysis**

Analyze frontend bundle untuk identify optimization opportunities:

```bash
cd frontend

# Build with analysis
yarn build

# Check dist folder sizes
du -sh dist/*
ls -lh dist/assets/

# Install bundle analyzer (optional)
yarn add -D rollup-plugin-visualizer

# Add to vite.config.js:
# import { visualizer } from 'rollup-plugin-visualizer'
# plugins: [vue(), visualizer()]

# Build and open stats.html
yarn build
```

**Expected Bundle Sizes:**

```
dist/
├── assets/
│   ├── index-a1b2c3d4.js      ~150KB (gzipped: ~50KB)
│   ├── index-e5f6g7h8.css     ~30KB  (gzipped: ~8KB)
│   ├── vendor-i9j0k1l2.js     ~80KB  (gzipped: ~25KB)
│   └── ... (fonts, images)
└── index.html                  ~2KB
```

**Optimization Strategies:**

1. **Code Splitting:**
   ```javascript
   // Lazy load routes
   const UserManagement = () => import('./views/UserManagement.vue')
   ```

2. **Tree Shaking:**
   - Import only what you need
   - Avoid `import *` statements

3. **Remove Unused Dependencies:**
   ```bash
   yarn remove unused-package
   ```

### **Runtime Performance Monitoring**

#### **Chrome DevTools Performance**

1. Open DevTools (F12) → Performance tab
2. Click record (⚫)
3. Interact with app (navigate, click, scroll)
4. Stop recording
5. Analyze:
   - FPS (target: 60fps)
   - CPU usage
   - Memory allocations
   - Long tasks (>50ms)

#### **Vue DevTools Performance**

1. Install Vue DevTools extension
2. Open DevTools → Vue tab
3. Performance tab
4. Record component render times
5. Identify slow components

---

## 🌍 Browser Compatibility

### **Desktop Browser Testing**

Test aplikasi di berbagai desktop browsers:

#### **Chrome (Primary)**

- **Version:** 90+
- **Market Share:** ~65%
- **Testing Priority:** HIGH

**Test Checklist:**
- [ ] App loads correctly
- [ ] All features functional
- [ ] Animations smooth (60fps)
- [ ] PWA installable
- [ ] Service Worker works
- [ ] DevTools console clean (no errors)
- [ ] Memory usage acceptable

#### **Firefox**

- **Version:** 88+
- **Market Share:** ~10%
- **Testing Priority:** MEDIUM

**Known Issues to Test:**
- [ ] Backdrop-filter support (limited)
- [ ] CSS Grid compatibility
- [ ] Flexbox gaps
- [ ] Date input handling

#### **Safari (macOS)**

- **Version:** 14+
- **Market Share:** ~20%
- **Testing Priority:** HIGH (iOS ecosystem)

**Known Issues to Test:**
- [ ] Webkit-specific prefixes
- [ ] PWA installation (Add to Dock)
- [ ] Date input appearance
- [ ] Flexbox bugs
- [ ] Service Worker behavior

#### **Edge**

- **Version:** 90+ (Chromium-based)
- **Market Share:** ~5%
- **Testing Priority:** LOW

**Test Checklist:**
- [ ] Chromium parity (should match Chrome)
- [ ] Windows-specific issues

### **Mobile Browser Testing**

#### **Chrome Android**

- **Version:** 90+
- **Device:** Physical or Emulator

**Test Checklist:**
- [ ] Responsive layout (320px - 768px)
- [ ] Touch targets ≥44x44px
- [ ] Tap delays minimal
- [ ] Pinch zoom disabled (where appropriate)
- [ ] Virtual keyboard handling
- [ ] PWA "Add to Home Screen" works
- [ ] Offline mode functional
- [ ] Haptic feedback works
- [ ] Performance smooth (no jank)

#### **Safari iOS**

- **Version:** 14+
- **Device:** Physical iPhone or Simulator

**Test Checklist:**
- [ ] iOS-specific gestures work
- [ ] Safe area handling (notch)
- [ ] PWA installation via Share → Add to Home Screen
- [ ] Standalone mode works
- [ ] Status bar styling correct
- [ ] Input focus behavior
- [ ] Scroll bounce (pull-to-refresh conflict)

### **Automated Cross-Browser Testing**

Use tools untuk automated cross-browser testing:

#### **BrowserStack** (Paid)

```bash
# Example: Test on multiple browsers automatically
# browserstack.config.js
module.exports = {
  browsers: [
    'chrome:latest',
    'firefox:latest',
    'safari:latest',
    'edge:latest'
  ]
}
```

#### **Playwright** (Free, Open Source)

```bash
# Install
yarn add -D @playwright/test

# Run tests on multiple browsers
npx playwright test --project=chromium
npx playwright test --project=firefox
npx playwright test --project=webkit
```

### **Responsive Design Testing**

Test di berbagai screen sizes:

| Device Type | Viewport Width | Priority |
|-------------|----------------|----------|
| Mobile Portrait | 320px - 480px | HIGH |
| Mobile Landscape | 568px - 768px | MEDIUM |
| Tablet Portrait | 768px - 1024px | MEDIUM |
| Tablet Landscape | 1024px - 1366px | LOW |
| Desktop | 1366px - 1920px | HIGH |
| Large Desktop | 1920px+ | LOW |

**Chrome DevTools Device Emulation:**
1. Open DevTools (F12)
2. Toggle device toolbar (Ctrl+Shift+M)
3. Select device presets or custom dimensions
4. Test all features at each size

---

## 📊 Performance Benchmarks

### **Backend Benchmarks**

Run regular benchmarks untuk track performance over time:

```bash
# Save benchmark results
ab -n 1000 -c 10 http://localhost:8080/health > benchmarks/$(date +%Y%m%d)_health.txt
ab -n 1000 -c 10 http://localhost:8080/api/users > benchmarks/$(date +%Y%m%d)_users.txt

# Compare with previous benchmarks
diff benchmarks/20251201_health.txt benchmarks/20251228_health.txt
```

### **Frontend Benchmarks**

```bash
# Run Lighthouse and save results
lighthouse http://localhost:8080 \
  --output=json \
  --output-path=./benchmarks/lighthouse_$(date +%Y%m%d).json

# Compare scores
jq '.categories[].score' benchmarks/lighthouse_20251228.json
```

---

## ✅ Performance Checklist

Before production deployment:

### **Backend Performance**
- [ ] API responses < 100ms (average)
- [ ] Can handle 1000+ req/s
- [ ] No memory leaks (stable memory over 24h)
- [ ] Database queries optimized (indexes, no N+1)
- [ ] Compression enabled (Gzip)
- [ ] Caching implemented (where appropriate)

### **Frontend Performance**
- [ ] Lighthouse Performance ≥ 90
- [ ] Bundle size < 200KB (gzipped)
- [ ] First Contentful Paint < 1.8s
- [ ] Time to Interactive < 3.8s
- [ ] No layout shift (CLS < 0.1)
- [ ] 60fps animations
- [ ] Lazy loading images

### **Cross-Browser**
- [ ] Chrome (latest) - all features work
- [ ] Firefox (latest) - all features work
- [ ] Safari (latest) - all features work
- [ ] Chrome Android - responsive & smooth
- [ ] Safari iOS - responsive & smooth

---

## 📚 Related Documentation

- [overview.md](./overview.md) - Testing strategy
- [backend-testing.md](./backend-testing.md) - Backend tests
- [frontend-testing.md](./frontend-testing.md) - Frontend tests
- [api-testing.md](./api-testing.md) - API testing

---

## 📞 Support

Jika ada pertanyaan tentang performance testing:
- Developer: Zulfikar Hidayatullah
- Phone: +62 857-1583-8733

---

**Last Updated:** 28 Desember 2025  
**Version:** 2.0.0 (Phase 2 Restructure)
