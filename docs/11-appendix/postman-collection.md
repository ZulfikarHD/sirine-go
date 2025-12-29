# 📮 Postman Collection Guide - Sirine Go

Complete guide untuk setup dan menggunakan Postman collection untuk test Sirine Go API endpoints.

---

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Collection Structure](#collection-structure)
3. [Environment Setup](#environment-setup)
4. [Authentication Flow](#authentication-flow)
5. [Request Examples](#request-examples)
6. [Testing Workflows](#testing-workflows)
7. [Best Practices](#best-practices)

---

## ✅ Prerequisites

### Required Software
- **Postman Desktop App** atau **Postman Web**
  - Download: https://www.postman.com/downloads/
  - Version: Latest (recommended)

### Application Setup
- Backend server running di `http://localhost:8080`
- Database initialized dengan test data
- At least 1 admin user untuk testing

---

## 📁 Collection Structure

### Recommended Folder Organization

```
Sirine Go API
├── 🔐 Authentication
│   ├── Login
│   ├── Refresh Token
│   ├── Get Current User
│   └── Logout
│
├── 👤 User Management (Admin)
│   ├── List Users
│   ├── Search Users
│   ├── Get User Detail
│   ├── Create User
│   ├── Update User
│   ├── Delete User
│   ├── Bulk Delete Users
│   └── Bulk Update Status
│
├── 👨 Profile Management
│   ├── Get Profile
│   ├── Update Profile
│   ├── Upload Photo
│   ├── Delete Photo
│   └── Get Activity Logs
│
├── 🔑 Password Management
│   ├── Change Password
│   ├── Forgot Password
│   ├── Reset Password
│   └── Admin Force Reset
│
├── 🔔 Notifications
│   ├── List Notifications
│   ├── Get Unread Count
│   ├── Get Recent Notifications
│   ├── Mark as Read
│   ├── Mark All as Read
│   ├── Delete Notification
│   └── Create Notification (Admin)
│
├── 📋 Activity Logs (Admin)
│   ├── List Activity Logs
│   ├── Get Log Detail
│   ├── Get User Activity
│   └── Get Statistics
│
├── 🎮 Achievements
│   ├── List Achievements
│   ├── Get User Achievements
│   ├── Get User Stats
│   ├── Award Achievement (Admin)
│   └── Get User Achievements (Admin)
│
├── 📦 Bulk Operations (Admin)
│   ├── Import Users (CSV)
│   └── Export Users (CSV)
│
└── ❤️ Health Check
    └── Health Check
```

---

## ⚙️ Environment Setup

### Create Environment

1. Click **Environments** di Postman sidebar
2. Click **Create Environment**
3. Name: `Sirine Go - Development`

### Environment Variables

```json
{
  "base_url": "http://localhost:8080",
  "access_token": "",
  "refresh_token": "",
  "user_id": "",
  "test_user_id": "",
  "test_nip": "NIP001",
  "test_password": "Password123!",
  "admin_nip": "ADMIN001",
  "admin_password": "AdminPass123!"
}
```

### Variable Descriptions

| Variable | Description | Example |
|----------|-------------|---------|
| `base_url` | Backend server URL | `http://localhost:8080` |
| `access_token` | JWT access token (auto-set after login) | `eyJhbGciOiJIUzI1Ni...` |
| `refresh_token` | JWT refresh token (auto-set after login) | `eyJhbGciOiJIUzI1Ni...` |
| `user_id` | Logged-in user ID (auto-set) | `1` |
| `test_user_id` | Test user ID untuk CRUD operations | `2` |
| `test_nip` | Test user NIP | `NIP001` |
| `test_password` | Test user password | `Password123!` |
| `admin_nip` | Admin NIP | `ADMIN001` |
| `admin_password` | Admin password | `AdminPass123!` |

---

## 🔐 Authentication Flow

### 1. Login Request

**Endpoint:**
```
POST {{base_url}}/api/auth/login
```

**Headers:**
```json
{
  "Content-Type": "application/json"
}
```

**Body:**
```json
{
  "nip": "{{admin_nip}}",
  "password": "{{admin_password}}"
}
```

**Tests Script (Auto-save tokens):**
```javascript
// Parse response
const response = pm.response.json();

// Save tokens to environment
if (response.access_token) {
    pm.environment.set("access_token", response.access_token);
}

if (response.refresh_token) {
    pm.environment.set("refresh_token", response.refresh_token);
}

// Save user ID
if (response.user && response.user.id) {
    pm.environment.set("user_id", response.user.id);
}

// Test assertions
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has access_token", function () {
    pm.expect(response.access_token).to.exist;
});

pm.test("Response has user data", function () {
    pm.expect(response.user).to.exist;
});
```

### 2. Authenticated Requests

**Authorization Header (Add to all authenticated requests):**
```
Authorization: Bearer {{access_token}}
```

**Pre-request Script (Optional - Auto-refresh):**
```javascript
// Check if token expired (optional)
const currentTime = Date.now() / 1000;
const tokenExpiry = pm.environment.get("token_expiry");

if (tokenExpiry && currentTime > tokenExpiry) {
    // Trigger refresh token request
    pm.sendRequest({
        url: pm.environment.get("base_url") + "/api/auth/refresh",
        method: 'POST',
        header: {
            'Authorization': 'Bearer ' + pm.environment.get("refresh_token")
        }
    }, function (err, res) {
        if (!err && res.code === 200) {
            const data = res.json();
            pm.environment.set("access_token", data.access_token);
        }
    });
}
```

### 3. Refresh Token

**Endpoint:**
```
POST {{base_url}}/api/auth/refresh
```

**Headers:**
```
Authorization: Bearer {{refresh_token}}
```

**Tests Script:**
```javascript
const response = pm.response.json();

if (response.access_token) {
    pm.environment.set("access_token", response.access_token);
}

pm.test("New access token received", function () {
    pm.expect(response.access_token).to.exist;
});
```

### 4. Logout

**Endpoint:**
```
POST {{base_url}}/api/auth/logout
```

**Headers:**
```
Authorization: Bearer {{access_token}}
```

**Tests Script:**
```javascript
// Clear tokens
pm.environment.set("access_token", "");
pm.environment.set("refresh_token", "");
pm.environment.set("user_id", "");

pm.test("Logout successful", function () {
    pm.response.to.have.status(200);
});
```

---

## 📝 Request Examples

### User Management

#### List Users (Paginated)

```
GET {{base_url}}/api/users?page=1&limit=20&role=ADMIN&department=KHAZWAL&status=active
Authorization: Bearer {{access_token}}
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `role` (optional): Filter by role
- `department` (optional): Filter by department
- `status` (optional): Filter by status (active/inactive)
- `search` (optional): Search by NIP or name

#### Create User

```
POST {{base_url}}/api/users
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "nip": "NIP999",
  "full_name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "08123456789",
  "role": "STAFF_KHAZWAL",
  "department": "KHAZWAL",
  "shift": "PAGI"
}
```

**Tests Script:**
```javascript
const response = pm.response.json();

// Save created user ID for later tests
pm.environment.set("test_user_id", response.data.id);

pm.test("User created successfully", function () {
    pm.response.to.have.status(201);
});

pm.test("Generated password exists", function () {
    pm.expect(response.generated_password).to.exist;
});

// Log generated password for testing
console.log("Generated Password:", response.generated_password);
```

#### Update User

```
PUT {{base_url}}/api/users/{{test_user_id}}
Authorization: Bearer {{access_token}}
Content-Type: application/json

{
  "full_name": "John Updated",
  "email": "john.updated@example.com",
  "phone": "08129999999",
  "role": "MANAGER_KHAZWAL"
}
```

#### Delete User

```
DELETE {{base_url}}/api/users/{{test_user_id}}
Authorization: Bearer {{access_token}}
```

### Profile Management

#### Upload Photo

```
POST {{base_url}}/api/profile/photo
Authorization: Bearer {{access_token}}
Content-Type: multipart/form-data

Body:
- photo: [Select File] (JPG, PNG, WebP, max 5MB)
```

### Notifications

#### Get Unread Count

```
GET {{base_url}}/api/notifications/unread-count
Authorization: Bearer {{access_token}}
```

**Tests Script:**
```javascript
const response = pm.response.json();

pm.test("Unread count is number", function () {
    pm.expect(response.unread_count).to.be.a('number');
});
```

---

## 🔄 Testing Workflows

### Workflow 1: Complete User Lifecycle

**Order:**
1. Login as Admin
2. Create User
3. Get User Detail
4. Update User
5. Get Updated User
6. Delete User
7. Verify User Deleted

**Runner Configuration:**
- Select folder: "User Management (Admin)"
- Iterations: 1
- Delay: 100ms between requests

### Workflow 2: Authentication Flow

**Order:**
1. Login
2. Get Current User
3. Refresh Token
4. Get Current User (with new token)
5. Logout

### Workflow 3: Profile Management

**Order:**
1. Login
2. Get Profile
3. Update Profile
4. Upload Photo
5. Get Profile (verify changes)
6. Get Activity Logs

---

## ✅ Best Practices

### 1. Use Environment Variables
```
✅ GOOD: {{base_url}}/api/users
❌ BAD:  http://localhost:8080/api/users
```

### 2. Add Tests to Requests
```javascript
// Always check status code
pm.test("Status code is 200", () => {
    pm.response.to.have.status(200);
});

// Validate response structure
pm.test("Response has data", () => {
    const response = pm.response.json();
    pm.expect(response.data).to.exist;
});

// Check response time
pm.test("Response time < 500ms", () => {
    pm.expect(pm.response.responseTime).to.be.below(500);
});
```

### 3. Use Pre-request Scripts
```javascript
// Set timestamp for testing
pm.environment.set("timestamp", Date.now());

// Generate random data
pm.environment.set("random_nip", "NIP" + Math.floor(Math.random() * 10000));
```

### 4. Organize Collections
- Use folders untuk group related endpoints
- Use meaningful names
- Add descriptions to requests
- Document expected responses

### 5. Handle Errors Gracefully
```javascript
pm.test("Error handling", () => {
    if (pm.response.code !== 200) {
        const error = pm.response.json();
        console.log("Error:", error.error);
    }
});
```

---

## 📊 Test Automation

### Collection Runner

1. Click **Collections** → Select collection
2. Click **Run collection**
3. Configure:
   - Environment: Sirine Go - Development
   - Iterations: 1
   - Delay: 100ms
   - Data file: (optional CSV)

### Newman (CLI)

```bash
# Install Newman
npm install -g newman

# Export collection & environment dari Postman

# Run collection
newman run Sirine_Go_API.postman_collection.json \
  -e Sirine_Go_Dev.postman_environment.json \
  --reporters cli,html \
  --reporter-html-export report.html
```

---

## 🔍 Debugging Tips

### View Console Logs
```javascript
console.log("Response:", pm.response.json());
console.log("Status:", pm.response.code);
console.log("Headers:", pm.response.headers);
```

### Inspect Variables
```javascript
console.log("Access Token:", pm.environment.get("access_token"));
console.log("User ID:", pm.environment.get("user_id"));
```

### Network Debugging
- Open Postman Console (View → Show Postman Console)
- View raw request/response
- Check headers
- View timing information

---

## 📥 Import/Export

### Export Collection
1. Click **...** di collection
2. Select **Export**
3. Choose **Collection v2.1**
4. Save file

### Import Collection
1. Click **Import** button
2. Drag & drop file atau select file
3. Review and import

---

## 📞 Support

Untuk questions about API testing:

**Developer:** Zulfikar Hidayatullah  
**Phone:** +62 857-1583-8733

---

## 🔗 Related Documentation

- [API Documentation](../03-development/api-documentation.md)
- [Authentication Guide](../04-api-reference/authentication.md)
- [User Management API](../04-api-reference/user-management.md)
- [Testing Guide](../06-testing/api-testing.md)

---

## 📦 Download Collection

**Coming Soon:** Pre-configured Postman collection akan tersedia untuk download.

**Current Method:** Create collection manually menggunakan guide ini sebagai reference.

---

**Last Updated:** 29 Desember 2025  
**Version:** 1.0.0
