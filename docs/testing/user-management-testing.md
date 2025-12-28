# User Management & Profile Testing Guide

**Feature**: User Management & Profile System  
**Version**: 1.1.0  
**Last Updated**: 28 Desember 2025

---

## Overview

Testing guide ini mencakup comprehensive test scenarios untuk User Management & Profile features, yaitu: backend API testing, frontend integration testing, end-to-end user flows, dan performance testing.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Backend API Testing](#backend-api-testing)
3. [Frontend Component Testing](#frontend-component-testing)
4. [Integration Testing](#integration-testing)
5. [End-to-End Testing](#end-to-end-testing)
6. [Performance Testing](#performance-testing)
7. [Security Testing](#security-testing)
8. [Accessibility Testing](#accessibility-testing)
9. [Test Automation](#test-automation)

---

## Prerequisites

### Environment Setup

**Backend**:
```bash
cd backend
go test ./... -v
```

**Frontend**:
```bash
cd frontend
yarn test
yarn test:e2e
```

### Test Data

**Admin Account**:
- NIP: `99999`
- Password: `admin123`
- Role: ADMIN

**Test User Account**:
- NIP: `12345`
- Password: (generated during test)
- Role: STAFF_KHAZWAL

### Required Tools

- **Backend**: Go testing framework, testify
- **Frontend**: Vitest, Testing Library, Playwright
- **API Testing**: cURL, Postman, or HTTPie
- **Performance**: k6, Apache Bench
- **Browser**: Chrome DevTools, Lighthouse

---

## Backend API Testing

### 1. User Management Endpoints

#### 1.1 GET /api/users - List Users

**Test Case 1: List all users (default pagination)**

```bash
# Setup
TOKEN=$(curl -s http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip":"99999","password":"admin123"}' \
  | jq -r '.data.access_token')

# Test
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
```json
{
  "success": true,
  "data": {
    "users": [...],
    "total": 1,
    "page": 1,
    "per_page": 20,
    "total_pages": 1
  }
}
```

**Assertions**:
- ✅ Status code: 200
- ✅ Response contains `users` array
- ✅ Pagination metadata present
- ✅ Users sorted by created_at DESC
- ✅ No password fields in response

**Test Case 2: Filter by role**

```bash
curl "http://localhost:8080/api/users?role=ADMIN" \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
- ✅ Only ADMIN users returned
- ✅ Total count matches filtered results

**Test Case 3: Search by name**

```bash
curl "http://localhost:8080/api/users?search=Admin" \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
- ✅ Results contain "Admin" in name or NIP
- ✅ Case-insensitive search works

**Test Case 4: Pagination**

```bash
curl "http://localhost:8080/api/users?page=2&per_page=10" \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
- ✅ Page 2 data returned
- ✅ Max 10 users per page
- ✅ Correct total_pages calculation

**Test Case 5: Unauthorized access**

```bash
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer invalid_token" \
  | jq
```

**Expected Result**:
- ✅ Status code: 401
- ✅ Error message: "Unauthorized"

---

#### 1.2 POST /api/users - Create User

**Test Case 1: Create user with valid data**

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "12345",
    "full_name": "Test User",
    "email": "test@sirine.local",
    "phone": "081234567890",
    "role": "STAFF_KHAZWAL",
    "department": "KHAZWAL",
    "shift": "PAGI"
  }' | jq
```

**Expected Result**:
```json
{
  "success": true,
  "message": "User berhasil dibuat",
  "data": {
    "user": {
      "id": 2,
      "nip": "12345",
      "full_name": "Test User",
      "status": "ACTIVE"
    },
    "generated_password": "A7b@xY3k9Mz2"
  }
}
```

**Assertions**:
- ✅ Status code: 201
- ✅ User created in database
- ✅ Generated password returned (12 chars, mixed case, numbers, symbols)
- ✅ `must_change_password` set to true
- ✅ Activity log created (CREATE action)
- ✅ Password hashed with bcrypt

**Test Case 2: Duplicate NIP**

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "99999",
    "full_name": "Duplicate",
    "email": "dup@sirine.local",
    "phone": "081111111111",
    "role": "STAFF_KHAZWAL",
    "department": "KHAZWAL"
  }' | jq
```

**Expected Result**:
- ✅ Status code: 400
- ✅ Error: "NIP sudah terdaftar dalam sistem"

**Test Case 3: Duplicate email**

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "11111",
    "full_name": "Test",
    "email": "admin@sirine.local",
    "phone": "081111111111",
    "role": "STAFF_KHAZWAL",
    "department": "KHAZWAL"
  }' | jq
```

**Expected Result**:
- ✅ Status code: 400
- ✅ Error: "email sudah terdaftar dalam sistem"

**Test Case 4: Invalid phone format**

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nip": "11111",
    "full_name": "Test",
    "email": "test2@sirine.local",
    "phone": "123456",
    "role": "STAFF_KHAZWAL",
    "department": "KHAZWAL"
  }' | jq
```

**Expected Result**:
- ✅ Status code: 400
- ✅ Error: "phone must start with 08 and be 10-15 digits"

**Test Case 5: Non-admin access**

```bash
# Login as staff
STAFF_TOKEN=$(curl -s http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"nip":"12345","password":"generated_password"}' \
  | jq -r '.data.access_token')

# Try to create user
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{...}' | jq
```

**Expected Result**:
- ✅ Status code: 403
- ✅ Error: "Anda tidak memiliki akses ke resource ini"

---

#### 1.3 PUT /api/users/:id - Update User

**Test Case 1: Update user data**

```bash
curl -X PUT http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Updated Name",
    "role": "OPERATOR_CETAK",
    "department": "CETAK"
  }' | jq
```

**Expected Result**:
```json
{
  "success": true,
  "message": "User berhasil diupdate",
  "data": {
    "id": 2,
    "full_name": "Updated Name",
    "role": "OPERATOR_CETAK",
    "department": "CETAK",
    "updated_at": "2025-12-28T15:00:00+07:00"
  }
}
```

**Assertions**:
- ✅ Status code: 200
- ✅ User data updated in database
- ✅ Activity log created dengan before/after values
- ✅ Only specified fields updated
- ✅ `updated_at` timestamp changed

**Test Case 2: Update to duplicate email**

```bash
curl -X PUT http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@sirine.local"}' | jq
```

**Expected Result**:
- ✅ Status code: 400
- ✅ Error: "email sudah digunakan oleh user lain"

**Test Case 3: Update non-existent user**

```bash
curl -X PUT http://localhost:8080/api/users/999 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"full_name": "Test"}' | jq
```

**Expected Result**:
- ✅ Status code: 404
- ✅ Error: "user dengan ID 999 tidak ditemukan"

---

#### 1.4 DELETE /api/users/:id - Delete User

**Test Case 1: Soft delete user**

```bash
curl -X DELETE http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
```json
{
  "success": true,
  "message": "User berhasil dihapus"
}
```

**Assertions**:
- ✅ Status code: 200
- ✅ User `deleted_at` timestamp set
- ✅ User NOT permanently removed (soft delete)
- ✅ Activity log created (DELETE action)
- ✅ User tidak muncul di GET /api/users

**Test Case 2: Try to delete self**

```bash
# Get current user ID
USER_ID=$(curl -s http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN" \
  | jq -r '.data.id')

# Try to delete self
curl -X DELETE http://localhost:8080/api/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
- ✅ Status code: 400
- ✅ Error: "Anda tidak dapat menghapus akun sendiri"

**Test Case 3: Delete non-existent user**

```bash
curl -X DELETE http://localhost:8080/api/users/999 \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
- ✅ Status code: 404
- ✅ Error: "user dengan ID 999 tidak ditemukan"

---

### 2. Profile Management Endpoints

#### 2.1 GET /api/profile - Get Own Profile

**Test Case 1: Get profile (authenticated)**

```bash
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN" \
  | jq
```

**Expected Result**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "nip": "99999",
    "full_name": "Administrator",
    "email": "admin@sirine.local",
    "phone": "081234567890",
    "role": "ADMIN",
    "department": "KHAZWAL",
    "shift": "PAGI",
    "status": "ACTIVE",
    "last_login_at": "2025-12-28T10:30:00+07:00"
  }
}
```

**Assertions**:
- ✅ Status code: 200
- ✅ Returns current user data
- ✅ No password field in response
- ✅ All profile fields present

**Test Case 2: Unauthenticated access**

```bash
curl http://localhost:8080/api/profile | jq
```

**Expected Result**:
- ✅ Status code: 401
- ✅ Error: "Unauthorized"

---

#### 2.2 PUT /api/profile - Update Own Profile

**Test Case 1: Update allowed fields**

```bash
curl -X PUT http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Updated Admin",
    "email": "newemail@sirine.local",
    "phone": "089876543210"
  }' | jq
```

**Expected Result**:
```json
{
  "success": true,
  "message": "Profile berhasil diupdate",
  "data": {
    "id": 1,
    "full_name": "Updated Admin",
    "email": "newemail@sirine.local",
    "phone": "089876543210",
    "updated_at": "2025-12-28T15:30:00+07:00"
  }
}
```

**Assertions**:
- ✅ Status code: 200
- ✅ Only allowed fields updated (name, email, phone)
- ✅ Activity log created
- ✅ Role, department, shift unchanged

**Test Case 2: Try to update restricted fields**

```bash
curl -X PUT http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test",
    "role": "ADMIN",
    "department": "CETAK"
  }' | jq
```

**Expected Result**:
- ✅ Status code: 200
- ✅ Only `full_name` updated
- ✅ `role` and `department` ignored (tidak berubah)

**Test Case 3: Duplicate email**

```bash
curl -X PUT http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email": "existing@sirine.local"}' | jq
```

**Expected Result**:
- ✅ Status code: 400
- ✅ Error: "email sudah digunakan oleh user lain"

---

### 3. Bulk Operations

#### 3.1 POST /api/users/bulk-delete

**Test Case 1: Bulk delete multiple users**

```bash
curl -X POST http://localhost:8080/api/users/bulk-delete \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_ids": [2, 3, 5]}' | jq
```

**Expected Result**:
```json
{
  "success": true,
  "message": "3 users berhasil dihapus"
}
```

**Assertions**:
- ✅ Status code: 200
- ✅ All specified users soft deleted
- ✅ Activity log created untuk each user
- ✅ Cannot include current user ID

**Test Case 2: Try to include self in bulk delete**

```bash
USER_ID=$(curl -s http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN" \
  | jq -r '.data.id')

curl -X POST http://localhost:8080/api/users/bulk-delete \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"user_ids\": [2, $USER_ID]}" | jq
```

**Expected Result**:
- ✅ Status code: 400
- ✅ Error: "Anda tidak dapat menghapus akun sendiri"

---

#### 3.2 POST /api/users/bulk-update-status

**Test Case 1: Bulk update status**

```bash
curl -X POST http://localhost:8080/api/users/bulk-update-status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_ids": [2, 3, 5],
    "status": "INACTIVE"
  }' | jq
```

**Expected Result**:
```json
{
  "success": true,
  "message": "Status 3 users berhasil diupdate"
}
```

**Assertions**:
- ✅ Status code: 200
- ✅ All specified users status updated
- ✅ Activity log created untuk each user
- ✅ Valid status values only (ACTIVE/INACTIVE/SUSPENDED)

---

## Frontend Component Testing

### 1. UserList.vue Component

**Test File**: `frontend/src/views/admin/users/__tests__/UserList.spec.js`

```javascript
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import UserList from '../UserList.vue'

describe('UserList.vue', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders user table', async () => {
    const wrapper = mount(UserList)
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('fetches users on mount', async () => {
    const mockFetch = vi.fn()
    // Test implementation
  })

  it('applies filters correctly', async () => {
    // Test filter functionality
  })

  it('handles pagination', async () => {
    // Test pagination
  })

  it('opens create modal', async () => {
    const wrapper = mount(UserList)
    await wrapper.find('[data-testid="create-user-btn"]').trigger('click')
    expect(wrapper.vm.showModal).toBe(true)
  })
})
```

**Test Cases**:
- ✅ Component renders correctly
- ✅ Fetches users on mount
- ✅ Displays user data in table
- ✅ Filter by role works
- ✅ Filter by department works
- ✅ Search functionality works
- ✅ Pagination controls work
- ✅ Create button opens modal
- ✅ Edit button opens modal dengan data
- ✅ Delete button shows confirmation
- ✅ Bulk select works
- ✅ Loading states shown
- ✅ Error states handled

---

### 2. UserFormModal.vue Component

**Test Cases**:
- ✅ Modal opens/closes correctly
- ✅ Form fields render
- ✅ Validation works (NIP, email, phone)
- ✅ Submit creates user (create mode)
- ✅ Submit updates user (edit mode)
- ✅ Generated password displayed (create)
- ✅ Copy password button works
- ✅ Error messages shown
- ✅ Loading state during submit
- ✅ Modal closes after success

---

### 3. Profile.vue Component

**Test Cases**:
- ✅ Displays user profile data
- ✅ Shows avatar or default
- ✅ Role badge displayed correctly
- ✅ Edit button navigates to edit page
- ✅ All info sections visible
- ✅ Responsive layout works

---

### 4. EditProfile.vue Component

**Test Cases**:
- ✅ Form pre-filled dengan current data
- ✅ Editable fields enabled
- ✅ Read-only fields disabled
- ✅ Validation works
- ✅ Submit updates profile
- ✅ Success redirects to profile
- ✅ Error messages shown
- ✅ Cancel button works

---

## Integration Testing

### 1. User Management Flow

**Scenario**: Admin creates, edits, and deletes a user

```javascript
describe('User Management Integration', () => {
  it('complete CRUD flow', async () => {
    // 1. Login as admin
    await login('99999', 'admin123')
    
    // 2. Navigate to user management
    await navigateTo('/admin/users')
    
    // 3. Create user
    await clickButton('Tambah User')
    await fillForm({
      nip: '12345',
      full_name: 'Test User',
      email: 'test@sirine.local',
      phone: '081234567890',
      role: 'STAFF_KHAZWAL',
      department: 'KHAZWAL'
    })
    await clickButton('Simpan')
    
    // 4. Verify user created
    expect(await findInTable('12345')).toBeTruthy()
    
    // 5. Edit user
    await clickEditButton('12345')
    await updateField('full_name', 'Updated Name')
    await clickButton('Update')
    
    // 6. Verify update
    expect(await findInTable('Updated Name')).toBeTruthy()
    
    // 7. Delete user
    await clickDeleteButton('12345')
    await confirmDialog()
    
    // 8. Verify deletion
    expect(await findInTable('12345')).toBeFalsy()
  })
})
```

---

### 2. Profile Management Flow

**Scenario**: User updates own profile

```javascript
describe('Profile Management Integration', () => {
  it('user updates profile', async () => {
    // 1. Login as user
    await login('12345', 'password')
    
    // 2. Navigate to profile
    await navigateTo('/profile')
    
    // 3. Click edit
    await clickButton('Edit Profile')
    
    // 4. Update fields
    await updateField('full_name', 'New Name')
    await updateField('email', 'newemail@sirine.local')
    await updateField('phone', '089876543210')
    
    // 5. Submit
    await clickButton('Simpan')
    
    // 6. Verify redirect to profile
    expect(currentPath()).toBe('/profile')
    
    // 7. Verify updated data
    expect(await findText('New Name')).toBeTruthy()
    expect(await findText('newemail@sirine.local')).toBeTruthy()
  })
})
```

---

## End-to-End Testing

### Using Playwright

**Test File**: `frontend/e2e/user-management.spec.js`

```javascript
import { test, expect } from '@playwright/test'

test.describe('User Management E2E', () => {
  test.beforeEach(async ({ page }) => {
    // Login as admin
    await page.goto('http://localhost:5173/login')
    await page.fill('[name="nip"]', '99999')
    await page.fill('[name="password"]', 'admin123')
    await page.click('button[type="submit"]')
    await page.waitForURL('**/dashboard')
  })

  test('create user flow', async ({ page }) => {
    // Navigate to user management
    await page.click('text=Manajemen User')
    await expect(page).toHaveURL('**/admin/users')

    // Click create button
    await page.click('text=Tambah User')
    
    // Fill form
    await page.fill('[name="nip"]', '12345')
    await page.fill('[name="full_name"]', 'Test User')
    await page.fill('[name="email"]', 'test@sirine.local')
    await page.fill('[name="phone"]', '081234567890')
    await page.selectOption('[name="role"]', 'STAFF_KHAZWAL')
    await page.selectOption('[name="department"]', 'KHAZWAL')
    
    // Submit
    await page.click('text=Simpan')
    
    // Wait for success
    await expect(page.locator('.toast-success')).toBeVisible()
    
    // Verify password modal
    await expect(page.locator('text=Password:')).toBeVisible()
    
    // Copy password
    const password = await page.locator('[data-testid="generated-password"]').textContent()
    expect(password).toHaveLength(12)
    
    // Close modal
    await page.click('text=Tutup')
    
    // Verify user in table
    await expect(page.locator('text=12345')).toBeVisible()
    await expect(page.locator('text=Test User')).toBeVisible()
  })

  test('edit user flow', async ({ page }) => {
    await page.goto('http://localhost:5173/admin/users')
    
    // Click edit button
    await page.click('[data-testid="edit-user-12345"]')
    
    // Update name
    await page.fill('[name="full_name"]', 'Updated Name')
    
    // Submit
    await page.click('text=Update')
    
    // Verify success
    await expect(page.locator('.toast-success')).toBeVisible()
    await expect(page.locator('text=Updated Name')).toBeVisible()
  })

  test('delete user flow', async ({ page }) => {
    await page.goto('http://localhost:5173/admin/users')
    
    // Click delete button
    await page.click('[data-testid="delete-user-12345"]')
    
    // Confirm deletion
    await page.click('text=Hapus')
    
    // Verify success
    await expect(page.locator('.toast-success')).toBeVisible()
    
    // Verify user removed
    await expect(page.locator('text=12345')).not.toBeVisible()
  })

  test('profile update flow', async ({ page }) => {
    await page.goto('http://localhost:5173/profile')
    
    // Click edit
    await page.click('text=Edit Profile')
    
    // Update fields
    await page.fill('[name="full_name"]', 'New Admin Name')
    await page.fill('[name="email"]', 'newadmin@sirine.local')
    
    // Submit
    await page.click('text=Simpan')
    
    // Verify redirect
    await expect(page).toHaveURL('**/profile')
    
    // Verify updated data
    await expect(page.locator('text=New Admin Name')).toBeVisible()
  })
})
```

**Run E2E Tests**:
```bash
cd frontend
yarn test:e2e
```

---

## Performance Testing

### 1. API Performance

**Using k6**:

```javascript
// k6-user-management.js
import http from 'k6/http'
import { check, sleep } from 'k6'

export let options = {
  stages: [
    { duration: '30s', target: 20 },  // Ramp up
    { duration: '1m', target: 50 },   // Stay at 50 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% requests < 500ms
    http_req_failed: ['rate<0.01'],   // Error rate < 1%
  },
}

export default function () {
  const token = 'YOUR_TOKEN'
  const headers = {
    'Authorization': `Bearer ${token}`,
  }

  // Test GET /api/users
  let res = http.get('http://localhost:8080/api/users', { headers })
  check(res, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  })

  sleep(1)
}
```

**Run Test**:
```bash
k6 run k6-user-management.js
```

**Performance Targets**:
- ✅ GET /api/users: < 200ms (p95)
- ✅ POST /api/users: < 500ms (p95)
- ✅ PUT /api/users/:id: < 300ms (p95)
- ✅ DELETE /api/users/:id: < 200ms (p95)
- ✅ Throughput: > 100 req/sec
- ✅ Error rate: < 1%

---

### 2. Frontend Performance

**Using Lighthouse**:

```bash
# Install Lighthouse CLI
npm install -g lighthouse

# Run audit
lighthouse http://localhost:5173/admin/users \
  --output html \
  --output-path ./lighthouse-report.html
```

**Performance Targets**:
- ✅ Performance Score: > 90
- ✅ First Contentful Paint: < 1.5s
- ✅ Time to Interactive: < 3s
- ✅ Largest Contentful Paint: < 2.5s
- ✅ Cumulative Layout Shift: < 0.1

---

## Security Testing

### 1. Authentication & Authorization

**Test Cases**:

```bash
# Test 1: Access without token
curl http://localhost:8080/api/users
# Expected: 401 Unauthorized

# Test 2: Access with invalid token
curl http://localhost:8080/api/users \
  -H "Authorization: Bearer invalid_token"
# Expected: 401 Unauthorized

# Test 3: Staff tries to create user
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $STAFF_TOKEN" \
  -d '{...}'
# Expected: 403 Forbidden

# Test 4: Manager tries to delete user
curl -X DELETE http://localhost:8080/api/users/2 \
  -H "Authorization: Bearer $MANAGER_TOKEN"
# Expected: 403 Forbidden
```

**Assertions**:
- ✅ All endpoints require authentication
- ✅ Role-based access control enforced
- ✅ Cannot access other users' profiles
- ✅ Cannot delete self
- ✅ Cannot escalate own privileges

---

### 2. Input Validation

**Test Cases**:

```bash
# SQL Injection attempt
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nip":"99999 OR 1=1","full_name":"Test",...}'
# Expected: 400 Bad Request

# XSS attempt
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"full_name":"<script>alert(1)</script>",...}'
# Expected: 400 Bad Request or sanitized

# Email validation
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"email":"invalid-email",...}'
# Expected: 400 Bad Request
```

**Assertions**:
- ✅ SQL injection prevented
- ✅ XSS prevented
- ✅ Input sanitized
- ✅ Validation errors clear

---

### 3. Password Security

**Test Cases**:
- ✅ Passwords hashed dengan bcrypt
- ✅ Generated passwords strong (12 chars, mixed)
- ✅ Passwords never returned di API responses
- ✅ `must_change_password` flag set for new users

---

## Accessibility Testing

### 1. Keyboard Navigation

**Test Cases**:
- ✅ Tab through all form fields
- ✅ Enter submits forms
- ✅ Esc closes modals
- ✅ Arrow keys navigate lists
- ✅ Space toggles checkboxes

**Manual Test**:
1. Navigate to `/admin/users`
2. Press Tab repeatedly
3. Verify focus indicators visible
4. Verify logical tab order
5. Press Enter on "Tambah User"
6. Verify modal opens
7. Press Esc
8. Verify modal closes

---

### 2. Screen Reader

**Test Cases**:
- ✅ Form labels announced
- ✅ Error messages announced
- ✅ Success messages announced
- ✅ Table headers announced
- ✅ Button purposes clear

**Manual Test** (using NVDA/JAWS):
1. Navigate to user list
2. Verify table structure announced
3. Navigate to form
4. Verify labels read correctly
5. Submit dengan errors
6. Verify errors announced

---

### 3. Color Contrast

**Test Cases**:
- ✅ Text contrast ratio > 4.5:1
- ✅ Interactive elements contrast > 3:1
- ✅ Focus indicators visible
- ✅ Error states clear without color alone

**Tool**: Chrome DevTools Accessibility Panel

---

## Test Automation

### CI/CD Pipeline

**GitHub Actions** (`.github/workflows/test.yml`):

```yaml
name: Test User Management

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run backend tests
        run: |
          cd backend
          go test ./... -v -cover

  frontend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install dependencies
        run: |
          cd frontend
          yarn install
      - name: Run unit tests
        run: |
          cd frontend
          yarn test
      - name: Run E2E tests
        run: |
          cd frontend
          yarn test:e2e

  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Start backend
        run: |
          cd backend
          go run cmd/server/main.go &
      - name: Start frontend
        run: |
          cd frontend
          yarn dev &
      - name: Wait for services
        run: sleep 10
      - name: Run integration tests
        run: |
          cd frontend
          yarn test:integration
```

---

## Test Coverage

### Backend Coverage

```bash
cd backend
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Target Coverage**:
- ✅ Overall: > 80%
- ✅ Services: > 90%
- ✅ Handlers: > 85%
- ✅ Models: > 75%

---

### Frontend Coverage

```bash
cd frontend
yarn test:coverage
```

**Target Coverage**:
- ✅ Overall: > 80%
- ✅ Components: > 85%
- ✅ Stores: > 90%
- ✅ Composables: > 85%

---

## Troubleshooting

### Common Issues

**Issue 1: CORS errors during testing**
```
Solution: Ensure backend CORS middleware configured correctly
Check: backend/middleware/cors.go
```

**Issue 2: 404 errors for API calls**
```
Solution: Verify baseURL in frontend/src/composables/useApi.js
Should be: http://localhost:8080/api
```

**Issue 3: Animations not working**
```
Solution: Check motion-v installed
Run: cd frontend && yarn add motion-v
```

**Issue 4: Tests failing due to timing**
```
Solution: Add proper waits in E2E tests
Use: await page.waitForSelector()
```

---

## Related Documentation

- **API Reference**: [user-management.md](../api/user-management.md)
- **Admin Journey**: [admin-user-management.md](../user-journeys/user-management/admin-user-management.md)
- **Profile Journey**: [user-profile-management.md](../user-journeys/user-management/user-profile-management.md)
- **Sprint 2 Summary**: [SPRINT2_SUMMARY.md](../../SPRINT2_SUMMARY.md)

---

**Last Updated**: 28 Desember 2025  
**Version**: 1.1.0 - Sprint 2  
**Status**: ✅ Complete
