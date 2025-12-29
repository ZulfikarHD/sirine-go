# 🔌 API Testing Guide

Panduan lengkap untuk testing API endpoints dengan cURL, Postman, dan integration testing.

---

## 📋 Daftar Isi

1. [Testing dengan cURL](#testing-dengan-curl)
2. [Testing dengan Postman](#testing-dengan-postman)
3. [Integration Testing](#integration-testing)

---

## 💻 Testing dengan cURL

cURL adalah command-line tool untuk quick API testing.

### **Basic Requests**

#### **Health Check**

```bash
curl http://localhost:8080/health

# Expected response:
# {"status":"ok","message":"Server berjalan dengan baik"}
```

#### **GET Request**

```bash
# Get all users
curl http://localhost:8080/api/users

# Get user by ID
curl http://localhost:8080/api/users/1

# Get with query parameters
curl "http://localhost:8080/api/users?page=1&limit=10&role=admin"
```

#### **POST Request (Create)**

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "password123",
    "role": "user",
    "department": "IT",
    "is_active": true
  }'
```

#### **PUT Request (Update)**

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Updated Name",
    "department": "HR"
  }'
```

#### **DELETE Request**

```bash
curl -X DELETE http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### **Authentication Testing**

#### **Login**

```bash
# Login to get token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'

# Save token dari response
# Response: {"data":{"token":"eyJhbGc...","user":{...}}}
```

#### **Using Token untuk Authenticated Requests**

```bash
# Set token as variable
TOKEN="eyJhbGciOiJIUzI1NiIs..."

# Use token in requests
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

#### **Refresh Token**

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Authorization: Bearer $TOKEN"
```

#### **Logout**

```bash
curl -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

### **Advanced cURL Options**

#### **Pretty Print JSON Response**

```bash
curl http://localhost:8080/api/users | jq

# jq filters
curl http://localhost:8080/api/users | jq '.data[] | {name, email}'
```

#### **Save Response to File**

```bash
curl http://localhost:8080/api/users -o response.json
```

#### **Show Response Headers**

```bash
curl -i http://localhost:8080/api/users
# or
curl -v http://localhost:8080/api/users  # Verbose (includes request headers too)
```

#### **Test Response Time**

```bash
curl -w "\nTime: %{time_total}s\n" http://localhost:8080/api/users
```

---

## 📮 Testing dengan Postman

Postman provides GUI untuk comprehensive API testing.

### **Setup Postman Collection**

#### **1. Create Collection**

1. Open Postman
2. Click "New Collection"
3. Name: "Sirine Go API"
4. Add description: "Complete API tests for Sirine Go App"

#### **2. Setup Environment**

Create environment dengan variables:

**Environment Name:** Sirine Go - Local

| Variable | Initial Value | Current Value |
|----------|---------------|---------------|
| `base_url` | `http://localhost:8080` | `http://localhost:8080` |
| `token` | `` | (Will be set after login) |
| `user_id` | `1` | `1` |

#### **3. Setup Authentication**

Create login request yang saves token:

**Request:** Login
- Method: `POST`
- URL: `{{base_url}}/api/auth/login`
- Body (JSON):
```json
{
  "email": "admin@example.com",
  "password": "admin123"
}
```

**Tests Script** (save token automatically):
```javascript
// Save token to environment
if (pm.response.code === 200) {
    const response = pm.response.json()
    pm.environment.set("token", response.data.token)
    pm.environment.set("user_id", response.data.user.id)
    console.log("Token saved:", response.data.token)
}

// Test assertions
pm.test("Status code is 200", function() {
    pm.response.to.have.status(200)
})

pm.test("Response has token", function() {
    const response = pm.response.json()
    pm.expect(response.data.token).to.exist
})
```

### **Example Requests**

#### **GET All Users**

- Method: `GET`
- URL: `{{base_url}}/api/users`
- Headers:
  - `Authorization`: `Bearer {{token}}`

**Tests:**
```javascript
pm.test("Status code is 200", function() {
    pm.response.to.have.status(200)
})

pm.test("Response has data array", function() {
    const response = pm.response.json()
    pm.expect(response.data).to.be.an('array')
})

pm.test("Response has meta pagination", function() {
    const response = pm.response.json()
    pm.expect(response.meta).to.exist
    pm.expect(response.meta.total).to.be.a('number')
})
```

#### **GET User by ID**

- Method: `GET`
- URL: `{{base_url}}/api/users/{{user_id}}`
- Headers:
  - `Authorization`: `Bearer {{token}}`

**Tests:**
```javascript
pm.test("Status code is 200", function() {
    pm.response.to.have.status(200)
})

pm.test("User has required fields", function() {
    const user = pm.response.json().data
    pm.expect(user).to.have.property('id')
    pm.expect(user).to.have.property('name')
    pm.expect(user).to.have.property('email')
    pm.expect(user).to.have.property('role')
})
```

#### **POST Create User**

- Method: `POST`
- URL: `{{base_url}}/api/users`
- Headers:
  - `Authorization`: `Bearer {{token}}`
  - `Content-Type`: `application/json`
- Body (JSON):
```json
{
  "name": "Test User {{$randomInt}}",
  "email": "test{{$timestamp}}@example.com",
  "password": "password123",
  "role": "user",
  "department": "IT",
  "is_active": true
}
```

**Tests:**
```javascript
pm.test("Status code is 201", function() {
    pm.response.to.have.status(201)
})

pm.test("User created successfully", function() {
    const response = pm.response.json()
    pm.expect(response.data).to.have.property('id')
    // Save new user ID for subsequent tests
    pm.environment.set("new_user_id", response.data.id)
})
```

#### **PUT Update User**

- Method: `PUT`
- URL: `{{base_url}}/api/users/{{new_user_id}}`
- Headers:
  - `Authorization`: `Bearer {{token}}`
  - `Content-Type`: `application/json`
- Body (JSON):
```json
{
  "name": "Updated Name",
  "department": "HR"
}
```

#### **DELETE User**

- Method: `DELETE`
- URL: `{{base_url}}/api/users/{{new_user_id}}`
- Headers:
  - `Authorization`: `Bearer {{token}}`

**Tests:**
```javascript
pm.test("Status code is 200", function() {
    pm.response.to.have.status(200)
})

pm.test("Delete message correct", function() {
    const response = pm.response.json()
    pm.expect(response.message).to.include('berhasil dihapus')
})
```

### **Automated Test Flow**

Create test runner yang executes multiple requests in sequence:

1. **Login** → Save token
2. **Create User** → Save user ID
3. **Get User** → Verify created
4. **Update User** → Verify updated
5. **Delete User** → Verify deleted
6. **Logout** → Clean up

**Run Collection:**
1. Collection → ... → Run
2. Select all requests
3. Click "Run Sirine Go API"

---

## 🔗 Integration Testing

Integration tests verify full flow dari frontend → API → database.

### **Scenario 1: User CRUD Flow**

Test complete Create → Read → Update → Delete flow:

```bash
#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🧪 Starting Integration Test: User CRUD Flow"

# 1. Login to get token
echo "\n1️⃣ Login..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token')

if [ "$TOKEN" != "null" ]; then
    echo "${GREEN}✅ Login successful${NC}"
else
    echo "${RED}❌ Login failed${NC}"
    exit 1
fi

# 2. Create user
echo "\n2️⃣ Creating user..."
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Integration Test User",
    "email": "integration@test.com",
    "password": "password123",
    "role": "user",
    "department": "IT",
    "is_active": true
  }')

USER_ID=$(echo $CREATE_RESPONSE | jq -r '.data.id')

if [ "$USER_ID" != "null" ]; then
    echo "${GREEN}✅ User created with ID: $USER_ID${NC}"
else
    echo "${RED}❌ User creation failed${NC}"
    exit 1
fi

# 3. Get user to verify
echo "\n3️⃣ Fetching user..."
GET_RESPONSE=$(curl -s http://localhost:8080/api/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN")

USER_NAME=$(echo $GET_RESPONSE | jq -r '.data.name')

if [ "$USER_NAME" == "Integration Test User" ]; then
    echo "${GREEN}✅ User fetched successfully${NC}"
else
    echo "${RED}❌ User fetch failed${NC}"
    exit 1
fi

# 4. Update user
echo "\n4️⃣ Updating user..."
UPDATE_RESPONSE=$(curl -s -X PUT http://localhost:8080/api/users/$USER_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Updated Integration User"}')

if echo $UPDATE_RESPONSE | jq -e '.data.name == "Updated Integration User"' > /dev/null; then
    echo "${GREEN}✅ User updated successfully${NC}"
else
    echo "${RED}❌ User update failed${NC}"
    exit 1
fi

# 5. Verify in database
echo "\n5️⃣ Verifying in database..."
DB_CHECK=$(mysql -u root -p'your_password' sirine_go -e \
  "SELECT name FROM users WHERE id=$USER_ID" -s -N)

if [ "$DB_CHECK" == "Updated Integration User" ]; then
    echo "${GREEN}✅ Database verification successful${NC}"
else
    echo "${RED}❌ Database verification failed${NC}"
    exit 1
fi

# 6. Delete user
echo "\n6️⃣ Deleting user..."
DELETE_RESPONSE=$(curl -s -X DELETE http://localhost:8080/api/users/$USER_ID \
  -H "Authorization: Bearer $TOKEN")

if echo $DELETE_RESPONSE | jq -e '.message' | grep -q "berhasil"; then
    echo "${GREEN}✅ User deleted successfully${NC}"
else
    echo "${RED}❌ User deletion failed${NC}"
    exit 1
fi

# 7. Verify soft delete in database
echo "\n7️⃣ Verifying soft delete..."
DB_DELETED=$(mysql -u root -p'your_password' sirine_go -e \
  "SELECT deleted_at FROM users WHERE id=$USER_ID" -s -N)

if [ "$DB_DELETED" != "NULL" ]; then
    echo "${GREEN}✅ Soft delete verified${NC}"
else
    echo "${RED}❌ Soft delete verification failed${NC}"
    exit 1
fi

echo "\n${GREEN}🎉 All integration tests passed!${NC}"
```

**Run script:**
```bash
chmod +x integration-test.sh
./integration-test.sh
```

### **Scenario 2: Authentication Flow**

Test login → access protected route → refresh token → logout:

```bash
#!/bin/bash

echo "🔐 Testing Authentication Flow"

# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.data.token')

echo "Token obtained: ${TOKEN:0:20}..."

# 2. Access protected route
curl -s http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN" | jq

# 3. Refresh token
NEW_TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/refresh \
  -H "Authorization: Bearer $TOKEN" \
  | jq -r '.data.token')

echo "New token obtained: ${NEW_TOKEN:0:20}..."

# 4. Logout
curl -s -X POST http://localhost:8080/api/auth/logout \
  -H "Authorization: Bearer $NEW_TOKEN" | jq

echo "✅ Authentication flow completed"
```

---

## ✅ Best Practices

### **1. Use Environment Variables**

Avoid hardcoding URLs dan credentials:

```bash
# Set environment variables
export API_URL="http://localhost:8080"
export API_TOKEN="your_token_here"

# Use in requests
curl $API_URL/api/users \
  -H "Authorization: Bearer $API_TOKEN"
```

### **2. Validate Response Structure**

Always check response format:

```javascript
// Postman test
pm.test("Response structure valid", function() {
    const response = pm.response.json()
    pm.expect(response).to.have.property('data')
    pm.expect(response).to.have.property('message')
})
```

### **3. Test Error Cases**

Don't just test happy path:

```bash
# Test invalid input
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Missing required fields"}'
# Expected: 400 Bad Request

# Test unauthorized access
curl http://localhost:8080/api/users
# Expected: 401 Unauthorized

# Test not found
curl http://localhost:8080/api/users/99999
# Expected: 404 Not Found
```

### **4. Clean Up Test Data**

Always cleanup test data setelah integration tests:

```bash
# Delete test users
curl -X DELETE http://localhost:8080/api/users/$TEST_USER_ID

# Or truncate test database
mysql -u root -p sirine_go_test -e "TRUNCATE TABLE users;"
```

---

## 📚 Related Documentation

- [overview.md](./overview.md) - Testing strategy
- [backend-testing.md](./backend-testing.md) - Backend unit tests
- [frontend-testing.md](./frontend-testing.md) - Frontend tests
- [../03-development/api-documentation.md](../03-development/api-documentation.md) - Complete API reference

---

## 📞 Support

Jika ada pertanyaan tentang API testing:
- Developer: Zulfikar Hidayatullah
- Phone: +62 857-1583-8733

---

**Last Updated:** 28 Desember 2025  
**Version:** 2.0.0 (Phase 2 Restructure)
