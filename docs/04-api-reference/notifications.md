# 🔔 Notifications API Reference

Complete API reference untuk notification endpoints dalam Sirine Go App.

**Base URL:** `http://localhost:8080/api/notifications`

---

## 📋 Overview

Notification system memungkinkan users untuk receive dan manage notifications dalam aplikasi, dengan fitur:
- List notifications dengan filter unread/all
- Mark notifications as read (individual atau bulk)
- Delete notifications
- Get unread count untuk badge display

**Authentication:** Semua endpoints memerlukan valid JWT token.

**Notification Types:**
- `INFO` - Informational notifications
- `SUCCESS` - Success notifications (achievements, etc)
- `WARNING` - Warning notifications
- `ERROR` - Error notifications

---

## 🔑 Endpoints

### 1. List Notifications

Get semua notifications untuk current user dengan optional filter.

```http
GET /api/notifications?unread_only=false
Authorization: Bearer {token}
```

**Query Parameters:**
- `unread_only` (optional) - Filter unread only (default: false)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "title": "Selamat Datang",
      "message": "Selamat datang di Sirine Go App",
      "type": "INFO",
      "is_read": false,
      "read_at": null,
      "created_at": "2025-12-28T10:00:00+07:00"
    },
    {
      "id": 2,
      "user_id": 1,
      "title": "Achievement Unlocked",
      "message": "Anda mendapatkan achievement First Login!",
      "type": "SUCCESS",
      "is_read": true,
      "read_at": "2025-12-28T10:05:00+07:00",
      "created_at": "2025-12-28T10:00:00+07:00"
    }
  ]
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid

---

### 2. Get Unread Count

Get jumlah unread notifications untuk display badge.

```http
GET /api/notifications/unread-count
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "count": 5
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid

---

### 3. Get Recent Notifications

Get recent notifications dengan limit untuk display di notification center.

```http
GET /api/notifications/recent?limit=5
Authorization: Bearer {token}
```

**Query Parameters:**
- `limit` (optional) - Number of notifications (default: 5, max: 20)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 5,
      "title": "Achievement Unlocked",
      "message": "Anda mendapatkan achievement First Login!",
      "type": "SUCCESS",
      "is_read": false,
      "created_at": "2025-12-28T10:00:00+07:00"
    }
  ]
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid

---

### 4. Mark Notification as Read

Mark individual notification as read.

```http
PUT /api/notifications/:id/read
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "message": "Notifikasi ditandai sudah dibaca"
}
```

**HTTP Status Codes:**
- `200 OK` - Mark berhasil
- `401 Unauthorized` - Token invalid
- `404 Not Found` - Notification not found

---

### 5. Mark All as Read

Mark semua notifications as read untuk current user.

```http
PUT /api/notifications/read-all
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "message": "Semua notifikasi ditandai sudah dibaca"
}
```

**HTTP Status Codes:**
- `200 OK` - Mark berhasil
- `401 Unauthorized` - Token invalid

---

### 6. Delete Notification

Delete individual notification.

```http
DELETE /api/notifications/:id
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "message": "Notifikasi berhasil dihapus"
}
```

**HTTP Status Codes:**
- `200 OK` - Delete berhasil
- `401 Unauthorized` - Token invalid
- `404 Not Found` - Notification not found

---

## 📊 Notification Types & Icons

| Type | Icon | Usage | Example |
|------|------|-------|---------|
| `INFO` | ℹ️ | General information | "Selamat datang di aplikasi" |
| `SUCCESS` | ✅ | Success messages, achievements | "Achievement unlocked!" |
| `WARNING` | ⚠️ | Warnings, reminders | "Password akan expired dalam 7 hari" |
| `ERROR` | ❌ | Error notifications | "Gagal mengupload file" |

---

## 🔔 Notification Triggers

Notifications otomatis dibuat saat:

### Authentication Events
- First login (achievement notification)
- Password changed successfully
- Account locked due to failed attempts

### Profile Events
- Profile updated successfully
- Profile photo uploaded

### Achievement Events
- New achievement unlocked
- Level up (Bronze → Silver → Gold → Platinum)

### Admin Actions (for target user)
- Account created by admin
- Password reset by admin
- Role/status changed by admin

---

## 🧪 Testing Examples

### cURL Examples

**Get All Notifications:**
```bash
curl http://localhost:8080/api/notifications \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Get Unread Only:**
```bash
curl http://localhost:8080/api/notifications?unread_only=true \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Get Unread Count:**
```bash
curl http://localhost:8080/api/notifications/unread-count \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Mark as Read:**
```bash
curl -X PUT http://localhost:8080/api/notifications/1/read \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Mark All as Read:**
```bash
curl -X PUT http://localhost:8080/api/notifications/read-all \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Delete Notification:**
```bash
curl -X DELETE http://localhost:8080/api/notifications/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### JavaScript/Axios Examples

**Get Notifications:**
```javascript
const response = await axios.get('/api/notifications')
const notifications = response.data.data
```

**Get Unread Count:**
```javascript
const response = await axios.get('/api/notifications/unread-count')
const unreadCount = response.data.data.count
```

**Mark as Read:**
```javascript
await axios.put(`/api/notifications/${notificationId}/read`)
```

**Real-time Updates (with Polling):**
```javascript
// Poll unread count every 30 seconds
setInterval(async () => {
  const response = await axios.get('/api/notifications/unread-count')
  updateBadge(response.data.data.count)
}, 30000)
```

---

## 🎨 Frontend Integration Tips

### Notification Badge Component
```vue
<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const unreadCount = ref(0)

const fetchUnreadCount = async () => {
  const response = await axios.get('/api/notifications/unread-count')
  unreadCount.value = response.data.data.count
}

onMounted(() => {
  fetchUnreadCount()
  // Poll every 30 seconds
  setInterval(fetchUnreadCount, 30000)
})
</script>

<template>
  <div class="relative">
    <IconBell />
    <span v-if="unreadCount > 0" class="badge">
      {{ unreadCount > 99 ? '99+' : unreadCount }}
    </span>
  </div>
</template>
```

### Notification List Component
```vue
<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const notifications = ref([])

const fetchNotifications = async (unreadOnly = false) => {
  const response = await axios.get('/api/notifications', {
    params: { unread_only: unreadOnly }
  })
  notifications.value = response.data.data
}

const markAsRead = async (id) => {
  await axios.put(`/api/notifications/${id}/read`)
  await fetchNotifications()
}

const markAllAsRead = async () => {
  await axios.put('/api/notifications/read-all')
  await fetchNotifications()
}

onMounted(() => {
  fetchNotifications()
})
</script>
```

---

## 📚 Related Documentation

- [Authentication API](./authentication.md)
- [Profile API](./profile.md)
- [Achievements API](./achievements.md)
- [Error Handling Guide](../05-guides/error-handling.md)

---

**Last Updated:** 28 Desember 2025  
**Sprint:** Sprint 4  
**Status:** ✅ Production Ready
