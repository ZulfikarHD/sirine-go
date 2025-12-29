# 📊 Activity Logs API Reference

Complete API reference untuk activity logs endpoints dalam Sirine Go App.

**Base URL:** `http://localhost:8080/api/admin/activity-logs`

---

## 📋 Overview

Activity Logs system menyediakan audit trail untuk semua activities dalam aplikasi, dengan fitur:
- Comprehensive logging untuk semua CRUD operations
- Advanced filtering dan searching
- User activity tracking
- Activity statistics dan reports

**Authentication:** Semua endpoints memerlukan valid JWT token.  
**Authorization:** ADMIN atau MANAGER role only.

**Logged Actions:**
- `CREATE` - Create new resource
- `UPDATE` - Update existing resource
- `DELETE` - Delete resource
- `LOGIN` - User login
- `LOGOUT` - User logout
- `PASSWORD_CHANGE` - Password changed
- `PHOTO_UPLOAD` - Profile photo uploaded

---

## 🔑 Endpoints

### 1. List Activity Logs

Get all activity logs dengan advanced filtering dan pagination.

```http
GET /api/admin/activity-logs?page=1&page_size=20&action=UPDATE&entity_type=users&search=admin&start_date=2025-12-01&end_date=2025-12-31
Authorization: Bearer {token}
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Items per page (default: 20, max: 100)
- `action` (optional) - Filter by action (CREATE, UPDATE, DELETE, etc)
- `entity_type` (optional) - Filter by entity type (users, profile, achievements, etc)
- `search` (optional) - Search entity_id atau entity_type
- `start_date` (optional) - Filter from date (YYYY-MM-DD)
- `end_date` (optional) - Filter until date (YYYY-MM-DD)

**Response:**
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 1,
        "user_id": 1,
        "user": {
          "full_name": "Super Admin",
          "nip": "99999"
        },
        "action": "UPDATE",
        "entity_type": "users",
        "entity_id": "10",
        "changes": {
          "before": {
            "role": "STAFF_KHAZWAL",
            "status": "ACTIVE"
          },
          "after": {
            "role": "MANAGER_KHAZWAL",
            "status": "ACTIVE"
          }
        },
        "ip_address": "127.0.0.1",
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
        "created_at": "2025-12-28T10:00:00+07:00"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid
- `403 Forbidden` - Insufficient permissions

---

### 2. Get Activity Log Detail

Get detail dari single activity log.

```http
GET /api/admin/activity-logs/:id
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user": {
      "id": 1,
      "full_name": "Super Admin",
      "nip": "99999",
      "role": "ADMIN"
    },
    "action": "UPDATE",
    "entity_type": "users",
    "entity_id": "10",
    "changes": {
      "before": {
        "role": "STAFF_KHAZWAL",
        "status": "ACTIVE",
        "department": "KHAZWAL"
      },
      "after": {
        "role": "MANAGER_KHAZWAL",
        "status": "ACTIVE",
        "department": "KHAZWAL"
      }
    },
    "ip_address": "127.0.0.1",
    "user_agent": "Mozilla/5.0...",
    "created_at": "2025-12-28T10:00:00+07:00"
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Log not found

---

### 3. Get User Activity Logs

Get all activity logs untuk specific user.

```http
GET /api/admin/activity-logs/user/:id?page=1&page_size=20
Authorization: Bearer {token}
```

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Items per page (default: 20, max: 100)

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 10,
      "full_name": "John Doe",
      "nip": "12345",
      "role": "STAFF_KHAZWAL"
    },
    "logs": [
      {
        "id": 5,
        "action": "UPDATE",
        "entity_type": "profile",
        "changes": {...},
        "created_at": "2025-12-28T10:00:00+07:00"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 20
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - User not found

---

### 4. Get Activity Statistics

Get statistics summary dari activity logs.

```http
GET /api/admin/activity-logs/stats?start_date=2025-12-01&end_date=2025-12-31
Authorization: Bearer {token}
```

**Query Parameters:**
- `start_date` (optional) - Filter from date (YYYY-MM-DD)
- `end_date` (optional) - Filter until date (YYYY-MM-DD)

**Response:**
```json
{
  "success": true,
  "data": {
    "total_logs": 1500,
    "by_action": {
      "CREATE": 300,
      "UPDATE": 800,
      "DELETE": 400,
      "LOGIN": 500,
      "LOGOUT": 450,
      "PASSWORD_CHANGE": 50
    },
    "by_entity_type": {
      "users": 600,
      "profile": 400,
      "achievements": 500
    },
    "most_active_users": [
      {
        "user_id": 1,
        "user_name": "Super Admin",
        "activity_count": 150
      },
      {
        "user_id": 2,
        "user_name": "Manager KHAZWAL",
        "activity_count": 120
      }
    ],
    "activity_by_date": [
      {
        "date": "2025-12-28",
        "count": 50
      },
      {
        "date": "2025-12-27",
        "count": 45
      }
    ]
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid
- `403 Forbidden` - Insufficient permissions

---

## 📊 Entity Types

Activity logs track operations pada various entity types:

| Entity Type | Description | Actions |
|-------------|-------------|---------|
| `users` | User CRUD operations | CREATE, UPDATE, DELETE |
| `profile` | Profile updates | UPDATE |
| `password` | Password changes | PASSWORD_CHANGE |
| `photo` | Photo uploads | PHOTO_UPLOAD |
| `achievements` | Achievement awards | CREATE |
| `auth` | Authentication | LOGIN, LOGOUT |

---

## 🔍 Changes Tracking

Activity logs track before/after values untuk UPDATE operations:

**Example - User Role Change:**
```json
{
  "action": "UPDATE",
  "entity_type": "users",
  "changes": {
    "before": {
      "role": "STAFF_KHAZWAL",
      "department": "KHAZWAL"
    },
    "after": {
      "role": "MANAGER_KHAZWAL",
      "department": "KHAZWAL"
    }
  }
}
```

**Example - Profile Update:**
```json
{
  "action": "UPDATE",
  "entity_type": "profile",
  "changes": {
    "before": {
      "email": "old@email.com",
      "phone": "081111111111"
    },
    "after": {
      "email": "new@email.com",
      "phone": "082222222222"
    }
  }
}
```

---

## 🧪 Testing Examples

### cURL Examples

**Get All Logs:**
```bash
curl http://localhost:8080/api/admin/activity-logs \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Filter by Action:**
```bash
curl "http://localhost:8080/api/admin/activity-logs?action=UPDATE" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Filter by Date Range:**
```bash
curl "http://localhost:8080/api/admin/activity-logs?start_date=2025-12-01&end_date=2025-12-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Get User Logs:**
```bash
curl http://localhost:8080/api/admin/activity-logs/user/10 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Get Statistics:**
```bash
curl http://localhost:8080/api/admin/activity-logs/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### JavaScript/Axios Examples

**Get Logs with Filters:**
```javascript
const response = await axios.get('/api/admin/activity-logs', {
  params: {
    page: 1,
    page_size: 20,
    action: 'UPDATE',
    entity_type: 'users',
    start_date: '2025-12-01',
    end_date: '2025-12-31'
  }
})

const logs = response.data.data.logs
```

**Get User Activity:**
```javascript
const response = await axios.get(`/api/admin/activity-logs/user/${userId}`)
const userLogs = response.data.data.logs
```

**Get Statistics:**
```javascript
const response = await axios.get('/api/admin/activity-logs/stats')
const stats = response.data.data
```

---

## 🎨 Frontend Integration Tips

### Activity Logs Table Component
```vue
<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const logs = ref([])
const filters = ref({
  action: '',
  entity_type: '',
  start_date: '',
  end_date: ''
})

const fetchLogs = async () => {
  const response = await axios.get('/api/admin/activity-logs', {
    params: filters.value
  })
  logs.value = response.data.data.logs
}

const formatChanges = (changes) => {
  // Display before/after comparison
  return Object.keys(changes.after).map(key => ({
    field: key,
    before: changes.before[key],
    after: changes.after[key]
  }))
}

onMounted(() => {
  fetchLogs()
})
</script>
```

### Activity Statistics Dashboard
```vue
<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const stats = ref(null)

const fetchStats = async () => {
  const response = await axios.get('/api/admin/activity-logs/stats')
  stats.value = response.data.data
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <div v-if="stats">
    <StatCard title="Total Logs" :value="stats.total_logs" />
    <PieChart :data="stats.by_action" />
    <BarChart :data="stats.activity_by_date" />
  </div>
</template>
```

---

## 📚 Related Documentation

- [User Management API](./user-management.md)
- [Authentication API](./authentication.md)
- [Error Handling Guide](../05-guides/error-handling.md)
- [Security Best Practices](../05-guides/security.md)

---

**Last Updated:** 28 Desember 2025  
**Sprint:** Sprint 4  
**Status:** ✅ Production Ready
