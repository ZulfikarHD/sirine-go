# 🧪 Testing Overview - Sirine Go App

Panduan testing strategy dan manual testing checklist untuk Sirine Go App.

---

## 📋 Daftar Isi

1. [Testing Strategy](#testing-strategy)
2. [Manual Testing](#manual-testing)
3. [Testing Checklist](#testing-checklist)

---

## 🎯 Testing Strategy

Aplikasi Sirine Go menggunakan multi-layered testing approach untuk memastikan quality dan reliability, yaitu:

### **Testing Pyramid**

```
        🔺 E2E Tests (Manual)
       🔺🔺 Integration Tests  
      🔺🔺🔺 API Tests
     🔺🔺🔺🔺 Unit Tests (Base)
```

### **Testing Types**

1. **Unit Testing**
   - Backend: Go tests untuk services & handlers
   - Frontend: Vitest untuk components & composables
   - **Coverage Target:** 75%+ (Backend), 70%+ (Frontend)
   - **Documentation:** [backend-testing.md](./backend-testing.md), [frontend-testing.md](./frontend-testing.md)

2. **API Testing**
   - cURL commands untuk quick testing
   - Postman collection untuk comprehensive testing
   - **Documentation:** [api-testing.md](./api-testing.md)

3. **Integration Testing**
   - Full flow testing (Frontend → API → Database)
   - User journey testing
   - **Documentation:** [api-testing.md](./api-testing.md#integration-testing)

4. **Performance Testing**
   - Load testing (Apache Bench)
   - Lighthouse audits
   - Bundle size analysis
   - **Documentation:** [performance-testing.md](./performance-testing.md)

5. **Manual Testing**
   - UI/UX testing
   - Browser compatibility
   - PWA functionality
   - **Documentation:** This file

---

## ✅ Manual Testing

### Quick Health Check

Setelah setup, lakukan quick test ini untuk memastikan aplikasi berjalan dengan baik:

```bash
# 1. Backend health check
curl http://localhost:8080/health

# Expected: {"status":"ok","message":"Server berjalan dengan baik"}

# 2. API response test
curl http://localhost:8080/api/users

# Expected: {"data":[],"meta":{...}} atau array of users

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
- [ ] Navigation menu terlihat jelas
- [ ] Responsive di berbagai ukuran layar (Desktop, Tablet, Mobile)
- [ ] Tailwind CSS styles applied dengan benar
- [ ] Glass effect cards rendered dengan backdrop blur
- [ ] Gradient colors (Indigo & Fuchsia) applied correctly
- [ ] No layout broken atau overlap
- [ ] Typography readable (font sizes appropriate)

#### **Animations & Motion**
- [ ] Page entrance animations smooth (Motion-v)
- [ ] Card stagger animations work pada list items
- [ ] Button press feedback (scale to 0.97) terasa responsive
- [ ] Modal animations (fade + scale) smooth
- [ ] Loading skeleton animations work
- [ ] Haptic feedback work di mobile (jika supported)
- [ ] No janky animations atau visual glitches

#### **Functionality - Authentication**
- [ ] Login form muncul dengan benar
- [ ] Input validation bekerja (email format, required fields)
- [ ] Submit button → Loading indicator muncul
- [ ] Success → Redirect ke dashboard
- [ ] Error → Error message muncul dengan jelas (Indonesian)
- [ ] Logout button works dan redirect ke login
- [ ] Session persistence works (refresh page tetap logged in)

#### **Functionality - User Management (Admin)**
- [ ] User list muncul dengan pagination
- [ ] Search functionality works (real-time dengan debounce)
- [ ] Filter by role works
- [ ] Filter by department works
- [ ] Filter by status works
- [ ] Create user button → Form modal muncul
- [ ] Edit user → Form pre-filled dengan data
- [ ] Delete user → Confirmation dialog muncul
- [ ] Bulk operations work (bulk delete, bulk update status)
- [ ] CSV import works
- [ ] CSV export downloads file correctly

#### **Functionality - Profile Management**
- [ ] Profile page accessible
- [ ] Profile data displayed correctly
- [ ] Edit profile → Update berhasil
- [ ] Change password form works
- [ ] Photo upload works (drag & drop + file picker)
- [ ] Photo preview muncul sebelum upload
- [ ] Photo delete works
- [ ] Activity logs visible

#### **Functionality - Notifications**
- [ ] Notification bell shows unread count
- [ ] Notification center opens on click
- [ ] Mark as read works
- [ ] Mark all as read works
- [ ] Delete notification works
- [ ] Filter notifications by type works
- [ ] Real-time updates work (30s polling)

#### **Functionality - Gamification**
- [ ] Achievements visible di profile
- [ ] Points counter animated
- [ ] Level badge displayed correctly (Bronze/Silver/Gold/Platinum)
- [ ] Achievement progress bars work
- [ ] Achievement unlock animation triggered

---

### Browser Compatibility Testing

Test aplikasi di berbagai browsers untuk memastikan compatibility:

#### **Desktop Testing**

Test di browsers berikut:
- ✅ **Chrome 90+** (primary browser - full support)
- ✅ **Firefox 88+** (secondary - test animations)
- ✅ **Safari 14+** (macOS - test webkit issues)
- ✅ **Edge 90+** (Windows - test compatibility)

**Test checklist per browser:**
- [ ] App loads correctly tanpa errors
- [ ] All features work properly
- [ ] Animations smooth (no janky motion)
- [ ] PWA installable
- [ ] Service Worker registered dan activated
- [ ] No console errors di DevTools
- [ ] Fonts rendered correctly
- [ ] Colors displayed correctly (gradient support)

#### **Mobile Testing**

Test di devices berikut:
- ✅ **Chrome Android 90+**
- ✅ **Safari iOS 14+**

**Test checklist mobile:**
- [ ] Responsive layout works (mobile-first design)
- [ ] Touch interactions responsive
- [ ] Tap targets sufficient size (min 44x44px)
- [ ] Gestures work (swipe, pinch zoom disabled appropriately)
- [ ] PWA installable ("Add to Home Screen")
- [ ] Offline mode works
- [ ] Performance smooth (60fps animations)
- [ ] No horizontal scroll
- [ ] Bottom navigation accessible
- [ ] Virtual keyboard doesn't break layout
- [ ] Haptic feedback works (vibrate API)

---

### PWA Testing

#### **Service Worker Testing**

1. **Build for production:**
   ```bash
   cd frontend && yarn build
   ```

2. **Serve production build:**
   ```bash
   make dev-backend
   # Backend serves frontend dist/ folder
   ```

3. **Verify Service Worker:**
   - Open: `http://localhost:8080`
   - Open DevTools (F12) → Application tab → Service Workers
   - ✅ Status: "activated and running"
   - ✅ Source: `/sw.js`

4. **Check Cache Storage:**
   - DevTools → Application → Cache Storage
   - ✅ Verify caches exist:
     - `workbox-precache-v2-...` (static assets)
     - `workbox-runtime-...` (API responses)

#### **Offline Testing**

1. Load aplikasi online (semua data fetched)
2. Open DevTools (F12) → Network tab
3. Set dropdown ke "Offline"
4. Reload page (Cmd+R / Ctrl+R)
5. ✅ Page loads from cache
6. ✅ Status indicator shows "Offline" (red badge)
7. ✅ Previously loaded data accessible
8. ✅ CRUD operations queued (if implemented)
9. Set back to "Online"
10. ✅ Status indicator updates to "Online" (green badge)
11. ✅ Queued operations synced

#### **PWA Installation Testing**

**Desktop (Chrome):**
1. Open app di Chrome
2. Address bar → Install icon (⊕) muncul
3. Click "Install"
4. ✅ App opens dalam standalone window
5. ✅ App appears in Applications/Start Menu
6. ✅ Launch works dari desktop icon

**Mobile (Android):**
1. Open app di Chrome mobile
2. Menu (⋮) → "Add to Home Screen"
3. Tap "Add"
4. ✅ Icon appears on home screen
5. ✅ Tap icon → Opens as standalone app (no browser UI)
6. ✅ Splash screen shows (if configured)

**Mobile (iOS):**
1. Open app di Safari iOS
2. Share button (↑) → "Add to Home Screen"
3. Tap "Add"
4. ✅ Icon appears on home screen
5. ✅ Tap icon → Opens as standalone app

---

## ✅ Testing Checklist

Sebelum deploy ke production, pastikan semua checklist ini completed:

### **Backend Testing**
- [ ] All unit tests pass (`go test ./...`)
- [ ] Test coverage ≥ 75%
- [ ] All API endpoints tested
- [ ] Error handling tested untuk semua edge cases
- [ ] Database connections work (production credentials)
- [ ] Migrations run successfully
- [ ] Seeders work (if needed)
- [ ] Performance acceptable (load testing passed)

### **Frontend Testing**
- [ ] All component tests pass (`yarn test`)
- [ ] Test coverage ≥ 70%
- [ ] UI responsive di all screen sizes (320px - 2560px)
- [ ] All user flows tested manually
- [ ] No console errors atau warnings
- [ ] Animations smooth (60fps)
- [ ] Build succeeds (`yarn build`)
- [ ] Production bundle size acceptable (<500KB total)

### **API Integration Testing**
- [ ] Full CRUD flow works end-to-end
- [ ] Authentication flow works (login, logout, refresh)
- [ ] Authorization works (role-based access control)
- [ ] API responses match documentation
- [ ] Error responses formatted correctly (Indonesian)
- [ ] Database persistence verified
- [ ] File uploads work (profile photos)
- [ ] CSV import/export work

### **PWA Testing**
- [ ] Service Worker registered successfully
- [ ] Offline mode works (app loads without network)
- [ ] App installable pada desktop & mobile
- [ ] Cache strategy working properly (NetworkFirst for API)
- [ ] Manifest file valid (`manifest.json`)
- [ ] Icons provided untuk all sizes

### **Performance Testing**
- [ ] Lighthouse Performance score ≥ 90
- [ ] Lighthouse PWA score ≥ 90
- [ ] Lighthouse Best Practices ≥ 90
- [ ] Lighthouse Accessibility ≥ 85
- [ ] Page load time < 3s (on 3G)
- [ ] API response time < 100ms (average)
- [ ] No memory leaks detected
- [ ] Backend handles 1000+ req/s

### **Security Testing**
- [ ] JWT authentication works correctly
- [ ] Password hashing secure (bcrypt cost 12)
- [ ] CORS configured correctly
- [ ] Rate limiting works (429 responses)
- [ ] SQL injection prevented (GORM parameterized queries)
- [ ] XSS prevented (Vue escaping)
- [ ] CSRF protection (if needed)
- [ ] Secure headers configured (production)

### **Cross-Browser Testing**
- [ ] Chrome (latest) - all features work
- [ ] Firefox (latest) - all features work
- [ ] Safari (latest) - all features work
- [ ] Edge (latest) - all features work
- [ ] Chrome Android - responsive & functional
- [ ] Safari iOS - responsive & functional

---

## 📊 Test Coverage Goals

### **Backend Coverage Target**
- **Services:** 80%+
- **Handlers:** 70%+
- **Overall:** 75%+

```bash
cd backend
go test -cover ./...
```

### **Frontend Coverage Target**
- **Components:** 70%+
- **Composables:** 80%+
- **Overall:** 70%+

```bash
cd frontend
yarn test:coverage
```

---

## 📚 Related Documentation

- [backend-testing.md](./backend-testing.md) - Backend unit testing guide
- [frontend-testing.md](./frontend-testing.md) - Frontend testing guide
- [api-testing.md](./api-testing.md) - API & integration testing
- [performance-testing.md](./performance-testing.md) - Performance testing
- [Test Scenarios](./test-scenarios/) - Feature-specific test cases

---

## 📞 Support

Jika ada pertanyaan tentang testing:
- Developer: Zulfikar Hidayatullah
- Phone: +62 857-1583-8733
- Timezone: Asia/Jakarta (WIB)

---

**Last Updated:** 28 Desember 2025  
**Version:** 2.0.0 (Phase 2 Restructure)
