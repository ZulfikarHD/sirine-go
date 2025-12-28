# Sprint 4 Testing Guide - Notifications & Audit Logs

**Sprint**: Advanced Features - Notifications & Audit (Week 4)
**Date**: December 28, 2025
**Status**: Implementation Complete

---

## Overview

Sprint 4 menambahkan in-app notifications dan audit logs system dengan fitur:
- Real-time notification updates (polling setiap 30 detik)
- Notification bell dengan badge count
- Notification center untuk kelola notifikasi
- Activity logs viewer untuk admin (audit trail)
- JSON diff viewer untuk changes detail

---

## Prerequisites untuk Testing

### 1. Backend Setup
```bash
# Pastikan backend server running
cd backend
go run cmd/main.go
```

### 2. Frontend Setup
```bash
# Pastikan frontend dev server running
cd frontend
yarn dev
```

### 3. Test Users
- **Admin**: NIP: 99999, Password: Admin@123
- **Staff/User**: (buat user baru via admin panel)

---

## Test Scenarios

### A. Notification System Testing

#### Test 1: Notification Bell Display
**Objective**: Verify notification bell menampilkan unread count dengan benar

**Steps**:
1. Login sebagai user (admin atau staff)
2. Perhatikan notification bell di navbar (kanan atas)
3. Verify badge count muncul jika ada unread notifications
4. Klik notification bell
5. Verify dropdown muncul dengan smooth animation (Motion-V spring)

**Expected Results**:
- ✓ Bell icon visible di navbar
- ✓ Badge count muncul dengan gradient indigo-fuchsia
- ✓ Badge animate dengan spring bounce saat count berubah
- ✓ Dropdown muncul dengan smooth scale + opacity animation
- ✓ Recent 5 notifications ditampilkan

**Test Data**: Initial state mungkin 0 notifications (normal untuk fresh install)

---

#### Test 2: Create Test Notification (Manual API Call)
**Objective**: Buat test notification untuk verify notification flow

**Steps (via Postman/Thunder Client)**:
1. GET `/api/auth/me` untuk dapatkan user_id
2. POST `/api/notifications` (Admin only endpoint - untuk testing)
   ```json
   {
     "user_id": 1,
     "title": "Test Notification",
     "message": "Ini adalah notifikasi test untuk Sprint 4",
     "type": "INFO"
   }
   ```
3. Verify response success
4. Refresh halaman frontend
5. Check notification bell badge count update

**Expected Results**:
- ✓ API return 201 Created dengan notification data
- ✓ Badge count di frontend update otomatis (max 30 detik delay karena polling)
- ✓ Notification muncul di dropdown bell

**Alternative**: Buat notification otomatis via:
- Create new user (admin action) → welcome notification
- Update user profile → update notification
- Change password → security notification

---

#### Test 3: Mark Notification as Read
**Objective**: Verify mark as read functionality dengan optimistic update

**Steps**:
1. Klik notification bell
2. Pilih salah satu unread notification (yang ada blue dot indicator)
3. Klik "Tandai Dibaca" button
4. Verify perubahan UI instant (optimistic update)
5. Verify badge count berkurang

**Expected Results**:
- ✓ Notification langsung berubah status (blue dot hilang)
- ✓ Badge count berkurang -1
- ✓ UI update instant tanpa loading
- ✓ Haptic feedback (vibrate 10ms) pada mobile device
- ✓ API call `/api/notifications/:id/read` success (verify via Network tab)

---

#### Test 4: Mark All Notifications as Read
**Objective**: Verify bulk mark all as read

**Steps**:
1. Pastikan ada minimal 2+ unread notifications
2. Klik notification bell
3. Klik "Tandai Semua Dibaca" button (header dropdown)
4. Verify semua notification berubah status
5. Verify badge count jadi 0

**Expected Results**:
- ✓ Semua notifications instantly marked as read
- ✓ Badge count langsung 0
- ✓ Blue dots hilang dari semua items
- ✓ Haptic feedback pattern ([10, 50, 10] ms)
- ✓ API call `/api/notifications/read-all` success

---

#### Test 5: Notification Center Page
**Objective**: Verify notification center dengan tabs dan filters

**Steps**:
1. Dari notification bell dropdown, klik "Lihat Semua Notifikasi"
2. Verify redirect ke `/notifications`
3. Test tab switching: "Semua" vs "Belum Dibaca"
4. Verify staggered animations saat notifications load
5. Test actions pada each notification:
   - Mark as read
   - Delete notification

**Expected Results**:
- ✓ Page load dengan entrance animation (fadeUp)
- ✓ Notifications display dalam cards dengan proper styling
- ✓ Tabs berfungsi dengan smooth transition
- ✓ Staggered animation (delay 0.05s per item)
- ✓ Icon badges color-coded by type (blue=INFO, green=SUCCESS, red=ERROR, etc)
- ✓ Time ago format dalam Bahasa Indonesia (e.g. "2 menit yang lalu")
- ✓ Empty state muncul jika no notifications

**Test Different Notification Types**:
- INFO: Blue icon & background
- SUCCESS: Green icon & background
- WARNING: Yellow icon & background
- ERROR: Red icon & background
- ACHIEVEMENT: Purple icon & background

---

#### Test 6: Real-time Polling
**Objective**: Verify notification count auto-update via polling

**Steps**:
1. Login dan biarkan halaman terbuka
2. Buat notification baru via API (separate browser/Postman)
3. Wait max 30 seconds
4. Verify badge count update otomatis tanpa refresh page

**Expected Results**:
- ✓ Polling dimulai saat component mounted
- ✓ Count update setiap 30 detik
- ✓ Polling stop saat component unmounted (logout/navigate away)
- ✓ No console errors

**Dev Check**: Inspect Network tab, verify periodic requests ke `/api/notifications/unread-count`

---

#### Test 7: Delete Notification
**Objective**: Verify delete notification functionality

**Steps**:
1. Buka notification center
2. Klik "Hapus" pada salah satu notification
3. Verify notification hilang dari list
4. Verify badge count update jika notification was unread
5. Refresh page, verify notification tetap terhapus

**Expected Results**:
- ✓ Notification immediately removed dari UI
- ✓ Badge count update accordingly
- ✓ Haptic feedback (30ms vibrate)
- ✓ Success toast "Notifikasi berhasil dihapus"
- ✓ API call `/api/notifications/:id` DELETE success

---

### B. Activity Logs / Audit Trail Testing

#### Test 8: Access Activity Logs Page (Admin Only)
**Objective**: Verify role-based access untuk audit logs

**Steps**:
1. Login sebagai **Admin** atau **Manager**
2. Navigate ke sidebar → Manajemen → "Audit Logs"
3. Verify redirect ke `/admin/audit`
4. Verify page load dengan activity logs table

**Expected Results**:
- ✓ Menu "Audit Logs" visible untuk Admin/Manager
- ✓ Page accessible tanpa error
- ✓ Table display dengan proper columns: Waktu, User, Aksi, Entity, Detail
- ✓ Entrance animation smooth (fadeUp dengan stagger)

**Negative Test**:
1. Login sebagai **Staff** (non-admin)
2. Try access `/admin/audit` directly via URL
3. Expected: Redirect ke dashboard (authorization failed)

---

#### Test 9: View Activity Logs with Filters
**Objective**: Verify filters dan pagination functionality

**Steps**:
1. Di Activity Logs page, test filters:
   - Jenis Aksi: Select "CREATE"
   - Entity Type: Select "users"
   - Search: Type "1" (entity ID)
2. Verify table update dengan filtered results
3. Check pagination controls
4. Navigate to page 2 (jika available)
5. Change page size (10 → 20)

**Expected Results**:
- ✓ Filters applied instantly
- ✓ Table update dengan smooth animation
- ✓ Pagination info accurate ("Halaman 1 dari X")
- ✓ Page navigation berfungsi (prev/next buttons)
- ✓ Page size change reset ke page 1
- ✓ "Reset Filter" button visible saat ada active filters
- ✓ Haptic feedback pada filter change

---

#### Test 10: View Activity Log Changes Detail
**Objective**: Verify JSON diff viewer untuk before/after changes

**Steps**:
1. Cari activity log dengan action "UPDATE" (dari user edit)
2. Klik "Lihat Changes" button
3. Verify row expand dengan animation
4. Verify before/after JSON display
5. Verify field-by-field comparison jika data is object
6. Klik "Sembunyikan" untuk collapse

**Expected Results**:
- ✓ Row expand dengan smooth height animation
- ✓ Before data (red border) vs After data (green border) side-by-side
- ✓ Field changes highlighted dalam detail section
- ✓ Password fields masked (show •••••••• instead of hash)
- ✓ JSON formatted dengan proper indentation
- ✓ Collapse animation smooth saat close

**Test Data Setup**:
```bash
# Create activity log dengan changes
1. Go to Admin → Manajemen User
2. Edit salah satu user (change name atau email)
3. Save changes
4. Navigate ke Audit Logs
5. Find UPDATE action untuk users entity
```

---

#### Test 11: Activity Log Action Types & Color Coding
**Objective**: Verify action badges color-coded correctly

**Steps**:
1. Generate different activity log types:
   - CREATE: Buat user baru
   - UPDATE: Edit user profile
   - DELETE: Delete user (soft delete)
   - LOGIN: Login action
   - LOGOUT: Logout action
   - PASSWORD_CHANGE: Change password
2. View audit logs page
3. Verify each action type memiliki correct badge color

**Expected Results**:
- ✓ CREATE: Green badge dengan text "Buat"
- ✓ UPDATE: Blue badge dengan text "Update"
- ✓ DELETE: Red badge dengan text "Hapus"
- ✓ LOGIN: Purple badge dengan text "Login"
- ✓ LOGOUT: Gray badge dengan text "Logout"
- ✓ PASSWORD_CHANGE: Yellow badge dengan text "Ganti Password"

---

#### Test 12: Activity Log Metadata
**Objective**: Verify additional metadata display (IP, User Agent)

**Steps**:
1. Expand activity log detail yang ada changes
2. Scroll ke bagian bawah expanded section
3. Verify "Additional Info" section muncul

**Expected Results**:
- ✓ IP Address displayed (e.g., "127.0.0.1" atau actual IP)
- ✓ User Agent displayed (browser info)
- ✓ Text formatted dengan proper labels

**Note**: IP Address dan User Agent diset saat activity logged via middleware

---

#### Test 13: Search Activity Logs
**Objective**: Verify search functionality

**Steps**:
1. Di Activity Logs page, focus search input
2. Type entity ID atau entity type (e.g., "users")
3. Wait 300ms (debounce delay)
4. Verify table update dengan search results
5. Clear search input
6. Verify table reset to all logs

**Expected Results**:
- ✓ Search debounced (tidak immediate request per keystroke)
- ✓ Results filtered correctly based on search term
- ✓ Search icon visible di input field
- ✓ Reset pagination ke page 1 saat search

---

#### Test 14: Pagination Edge Cases
**Objective**: Test pagination dengan different scenarios

**Test Cases**:

**Case 1: Single Page**
- Setup: Have < 20 logs (default page size)
- Expected: No pagination controls shown

**Case 2: Multiple Pages**
- Setup: Have 50+ logs
- Expected: Pagination controls visible dengan correct total_pages

**Case 3: Last Page**
- Navigate ke last page
- Expected: "Next" button disabled

**Case 4: First Page**
- Expected: "Previous" button disabled

**Case 5: Page Size Change**
- Change dari 20 → 50
- Expected: Fewer pages, reset ke page 1

---

### C. Integration Tests

#### Test 15: Create User → Check Notifications & Audit
**Objective**: End-to-end test notification + audit integration

**Steps**:
1. Login sebagai Admin
2. Create new user via Admin → Manajemen User
3. Check notifications (bell) untuk welcome notification ke new user
4. Navigate ke Audit Logs
5. Verify CREATE action logged untuk users entity
6. Expand detail, verify before=null, after=user data

**Expected Results**:
- ✓ User created successfully
- ✓ (Optional) Welcome notification sent ke new user
- ✓ Activity log recorded dengan:
  - Action: CREATE
  - Entity Type: users
  - Entity ID: new user's ID
  - User: Admin yang buat
  - Changes: before=null, after=full user object
- ✓ Timestamp accurate

---

#### Test 16: Update Profile → Check Audit
**Objective**: Verify profile update generates audit log

**Steps**:
1. Login as any user
2. Go to Profile → Edit
3. Change full_name or email
4. Save changes
5. Admin: Check Audit Logs
6. Find UPDATE action untuk profile entity
7. Expand changes, verify before/after values

**Expected Results**:
- ✓ Profile updated successfully
- ✓ Activity log created dengan UPDATE action
- ✓ Changes JSON shows:
  - before: {full_name: "Old Name", ...}
  - after: {full_name: "New Name", ...}
- ✓ Only changed fields highlighted

---

#### Test 17: Password Change → Security Notification
**Objective**: Verify password change triggers notification + audit

**Steps**:
1. Login as user
2. Go to Profile → Ganti Password
3. Change password successfully
4. Check notifications
5. Admin: Check audit logs

**Expected Results**:
- ✓ Password changed successfully
- ✓ Security notification sent: "Password Anda telah diubah"
- ✓ Notification type: WARNING (yellow)
- ✓ Activity log ACTION: PASSWORD_CHANGE
- ✓ Changes: password hash masked in display

---

#### Test 18: Logout → Check Audit
**Objective**: Verify logout action logged

**Steps**:
1. Login as any user
2. Click user dropdown → Keluar
3. After redirect to login
4. Re-login as Admin
5. Check audit logs
6. Find LOGOUT action

**Expected Results**:
- ✓ Logout successful
- ✓ Redirect ke /login
- ✓ Activity log created dengan LOGOUT action
- ✓ No changes data (expected untuk logout)

---

### D. Mobile Responsiveness Testing

#### Test 19: Notification Bell on Mobile
**Objective**: Verify notification UI pada mobile devices

**Steps**:
1. Open DevTools → Toggle device toolbar
2. Select mobile device (e.g., iPhone 12)
3. Test notification bell:
   - Click bell
   - Verify dropdown width appropriate
   - Test scroll dalam notification list
   - Test mark as read action

**Expected Results**:
- ✓ Bell visible dan clickable
- ✓ Dropdown width max-w-80 (320px) pas untuk mobile
- ✓ Dropdown positioned correctly (tidak keluar viewport)
- ✓ Scroll works dalam notification list
- ✓ Touch targets adequate (min 44x44px)
- ✓ Haptic feedback berfungsi (jika device support)

---

#### Test 20: Activity Logs Table on Mobile
**Objective**: Verify table responsiveness

**Steps**:
1. Mobile device mode (DevTools)
2. Navigate ke Activity Logs
3. Test filters collapse properly
4. Test table horizontal scroll
5. Test expand/collapse rows

**Expected Results**:
- ✓ Filters stack vertically pada mobile
- ✓ Table scrollable horizontal (overflow-x-auto)
- ✓ Critical columns visible (Waktu, User, Aksi)
- ✓ Expand button accessible
- ✓ Pagination controls stacked untuk mobile

---

### E. Performance Testing

#### Test 21: Notification Polling Performance
**Objective**: Verify polling tidak cause performance issues

**Steps**:
1. Open DevTools → Performance tab
2. Start recording
3. Wait 2 minutes (4 polling cycles at 30s interval)
4. Stop recording
5. Analyze network activity dan JS execution

**Expected Results**:
- ✓ Polling requests lightweight (< 1KB)
- ✓ No memory leaks
- ✓ No excessive re-renders
- ✓ CPU usage minimal during idle polling
- ✓ Network waterfall shows regular 30s intervals

---

#### Test 22: Large Activity Logs Dataset
**Objective**: Test pagination dengan large dataset

**Steps**:
1. Generate 200+ activity logs (via using system extensively)
2. Navigate ke Activity Logs page
3. Test pagination navigation (go to page 5, 10, etc)
4. Test page size 100
5. Monitor load times

**Expected Results**:
- ✓ Initial load < 2s
- ✓ Page navigation smooth (< 1s)
- ✓ No UI freezing
- ✓ Pagination controls update correctly
- ✓ Staggered animation performant (tidak lag)

---

### F. Error Handling & Edge Cases

#### Test 23: Network Failure Handling
**Objective**: Verify graceful error handling

**Test Case 1: Offline Mode**
1. Open DevTools → Network → Offline
2. Try mark notification as read
3. Expected: Optimistic update rollback, error message

**Test Case 2: 401 Unauthorized**
1. Manually clear auth token (localStorage)
2. Try any notification action
3. Expected: Redirect ke login

**Test Case 3: 500 Server Error**
1. (Simulate via backend or mock)
2. Try fetch notifications
3. Expected: Error toast "Gagal memuat notifikasi"

---

#### Test 24: Empty States
**Objective**: Verify empty states display properly

**Test Case 1: No Notifications**
- Expected: Bell with no badge, empty state in dropdown & center

**Test Case 2: No Unread Notifications**
- Expected: Badge hidden, "Tandai Semua Dibaca" button hidden

**Test Case 3: No Activity Logs (Filtered)**
- Apply strict filter that returns 0 results
- Expected: Empty state message in table

---

#### Test 25: Concurrent Updates
**Objective**: Test optimistic updates dengan race conditions

**Steps**:
1. Open 2 browser windows, same user logged in
2. Window 1: Mark notification as read
3. Window 2: Immediately try delete same notification
4. Verify behavior

**Expected**: Both windows should sync via polling, one action succeeds first

---

## API Endpoints Reference

### Notification Endpoints
- `GET /api/notifications` - Get all user notifications
  - Query: `?unread_only=true`
- `GET /api/notifications/unread-count` - Get unread count
- `GET /api/notifications/recent?limit=5` - Get recent notifications
- `PUT /api/notifications/:id/read` - Mark as read
- `PUT /api/notifications/read-all` - Mark all as read
- `DELETE /api/notifications/:id` - Delete notification

### Activity Log Endpoints (Admin Only)
- `GET /api/admin/activity-logs` - Get all activity logs
  - Query params: `page`, `page_size`, `action`, `entity_type`, `search`
- `GET /api/admin/activity-logs/:id` - Get log detail
- `GET /api/admin/activity-logs/user/:id` - Get logs by user
- `GET /api/profile/activity` - Get own activity logs

---

## Known Issues & Limitations

### Current Limitations:
1. **Polling Interval**: Fixed 30 seconds (not configurable from UI)
2. **Notification Types**: Limited to 5 types (INFO, SUCCESS, WARNING, ERROR, ACHIEVEMENT)
3. **Real-time**: Uses polling instead of WebSocket (acceptable for MVP)
4. **Search**: Basic search pada entity_id dan entity_type only

### Future Enhancements (Post-Sprint 4):
- WebSocket untuk true real-time notifications
- Notification sound/desktop push notifications
- Export activity logs to CSV/PDF
- Advanced filters (date range picker dengan calendar)
- User activity heatmap visualization

---

## Bug Reporting

Jika menemukan bug saat testing, report dengan format:

**Bug Title**: [Component] Short description

**Steps to Reproduce**:
1. Step 1
2. Step 2
3. ...

**Expected Behavior**: What should happen

**Actual Behavior**: What actually happened

**Environment**:
- Browser: Chrome 120 / Firefox 121 / Safari 17
- OS: Windows 11 / macOS 14 / Ubuntu 22.04
- Screen Size: Desktop / Tablet / Mobile

**Screenshots**: (attach jika visual bug)

**Console Errors**: (attach console logs jika ada)

---

## Success Criteria

Sprint 4 dianggap success jika:

- ✓ Semua test cases di atas PASS
- ✓ Notification bell berfungsi dengan badge count accurate
- ✓ Real-time polling berjalan tanpa errors
- ✓ Activity logs accessible untuk Admin/Manager only
- ✓ JSON diff viewer menampilkan changes dengan jelas
- ✓ No critical bugs
- ✓ Mobile responsive berfungsi dengan baik
- ✓ Performance acceptable (load < 2s, actions < 500ms)

---

## Next Steps

After testing selesai:
1. Document bugs yang ditemukan
2. Fix critical bugs immediately
3. Schedule non-critical fixes untuk Sprint 5
4. Update sprint plan dengan learnings
5. Prepare for Sprint 5: Enhancements & Gamification

---

**Testing Completed By**: [Your Name]
**Date**: December 28, 2025
**Status**: ✅ Ready for manual testing
