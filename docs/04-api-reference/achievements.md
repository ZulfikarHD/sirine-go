# 🎮 Achievements API Reference

Complete API reference untuk achievements & gamification endpoints dalam Sirine Go App.

**Base URL:** `http://localhost:8080/api`

---

## 📋 Overview

Gamification system menggunakan achievements dan points untuk engage users dengan fitur:
- Achievement tracking (locked/unlocked)
- Points accumulation system
- Level progression (Bronze → Silver → Gold → Platinum)
- User stats dan progress tracking

**Authentication:** Semua endpoints memerlukan valid JWT token.

**Levels & Requirements:**
- **Bronze:** 0-99 points
- **Silver:** 100-499 points
- **Gold:** 500-999 points
- **Platinum:** 1000+ points

---

## 🏆 Achievement Categories

| Category | Description | Example Achievements |
|----------|-------------|---------------------|
| `LOGIN` | Login-related | First Login, Daily Login Streak |
| `PROFILE` | Profile completion | Profile Complete, Photo Upload |
| `SOCIAL` | Social interactions | (Future: Comments, Likes) |
| `MILESTONE` | Milestones | Account Anniversary |

---

## 🔑 User Endpoints

### 1. List All Achievements

Get list semua available achievements dalam system.

```http
GET /api/achievements
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "code": "FIRST_LOGIN",
      "name": "First Login",
      "description": "Login untuk pertama kalinya",
      "icon": "🎉",
      "points": 10,
      "category": "LOGIN",
      "is_active": true
    },
    {
      "id": 2,
      "code": "PROFILE_COMPLETE",
      "name": "Profile Complete",
      "description": "Lengkapi profil dengan foto dan informasi",
      "icon": "✨",
      "points": 20,
      "category": "PROFILE",
      "is_active": true
    },
    {
      "id": 3,
      "code": "PHOTO_UPLOADED",
      "name": "Photo Master",
      "description": "Upload foto profil pertama kali",
      "icon": "📸",
      "points": 15,
      "category": "PROFILE",
      "is_active": true
    }
  ]
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid

---

### 2. Get User Achievements

Get user achievements dengan unlock status untuk current user.

```http
GET /api/profile/achievements
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "achievement": {
        "id": 1,
        "code": "FIRST_LOGIN",
        "name": "First Login",
        "description": "Login untuk pertama kalinya",
        "icon": "🎉",
        "points": 10,
        "category": "LOGIN"
      },
      "unlocked": true,
      "unlocked_at": "2025-12-28T10:00:00+07:00"
    },
    {
      "achievement": {
        "id": 2,
        "code": "PROFILE_COMPLETE",
        "name": "Profile Complete",
        "description": "Lengkapi profil dengan foto dan informasi",
        "icon": "✨",
        "points": 20,
        "category": "PROFILE"
      },
      "unlocked": false,
      "unlocked_at": null
    },
    {
      "achievement": {
        "id": 3,
        "code": "PHOTO_UPLOADED",
        "name": "Photo Master",
        "description": "Upload foto profil pertama kali",
        "icon": "📸",
        "points": 15,
        "category": "PROFILE"
      },
      "unlocked": true,
      "unlocked_at": "2025-12-28T10:15:00+07:00"
    }
  ]
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid

---

### 3. Get User Gamification Stats

Get comprehensive gamification statistics untuk current user.

```http
GET /api/profile/stats
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_points": 150,
    "current_level": "Silver",
    "next_level": "Gold",
    "points_to_next_level": 350,
    "achievements_unlocked": 5,
    "achievements_total": 6,
    "completion_percentage": 83.33,
    "recent_achievements": [
      {
        "name": "Photo Master",
        "points": 15,
        "unlocked_at": "2025-12-28T10:15:00+07:00"
      }
    ]
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid

---

## 👨‍💼 Admin Endpoints

### 4. Award Achievement (Admin Only)

Manually award achievement ke specific user.

```http
POST /api/admin/achievements/award
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "user_id": 10,
  "achievement_code": "FIRST_LOGIN"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "achievement": "First Login",
    "points_awarded": 10,
    "user_total_points": 10,
    "user_level": "Bronze",
    "notification_sent": true
  },
  "message": "Achievement berhasil diberikan"
}
```

**HTTP Status Codes:**
- `200 OK` - Award berhasil
- `400 Bad Request` - Invalid achievement code atau already unlocked
- `401 Unauthorized` - Token invalid
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - User or achievement not found

---

### 5. Get User Achievements (Admin)

Admin view untuk melihat achievements dari specific user.

```http
GET /api/admin/users/:id/achievements
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 10,
      "full_name": "John Doe",
      "nip": "12345",
      "total_points": 25,
      "level": "Bronze"
    },
    "achievements": [
      {
        "achievement": {
          "code": "FIRST_LOGIN",
          "name": "First Login",
          "icon": "🎉",
          "points": 10
        },
        "unlocked": true,
        "unlocked_at": "2025-12-28T10:00:00+07:00"
      },
      {
        "achievement": {
          "code": "PHOTO_UPLOADED",
          "name": "Photo Master",
          "icon": "📸",
          "points": 15
        },
        "unlocked": true,
        "unlocked_at": "2025-12-28T10:15:00+07:00"
      }
    ],
    "stats": {
      "unlocked_count": 2,
      "total_count": 6,
      "completion_percentage": 33.33
    }
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Request berhasil
- `401 Unauthorized` - Token invalid
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - User not found

---

## 🎯 Achievement Triggers

Achievements otomatis ter-unlock saat user melakukan actions tertentu:

### LOGIN Category

| Achievement | Code | Points | Trigger |
|-------------|------|--------|---------|
| **First Login** | `FIRST_LOGIN` | 10 | Login pertama kali |
| **Daily Login Streak** | `LOGIN_STREAK_7` | 50 | Login 7 hari berturut-turut (future) |

### PROFILE Category

| Achievement | Code | Points | Trigger |
|-------------|------|--------|---------|
| **Profile Complete** | `PROFILE_COMPLETE` | 20 | Lengkapi email, phone, dan foto |
| **Photo Master** | `PHOTO_UPLOADED` | 15 | Upload foto profil pertama kali |

### MILESTONE Category

| Achievement | Code | Points | Trigger |
|-------------|------|--------|---------|
| **Account Anniversary** | `ANNIVERSARY_1YEAR` | 100 | Akun berusia 1 tahun (future) |

---

## 🏅 Level System

### Level Requirements

```javascript
const getLevelByPoints = (points) => {
  if (points >= 1000) return 'Platinum'
  if (points >= 500) return 'Gold'
  if (points >= 100) return 'Silver'
  return 'Bronze'
}
```

### Level Benefits

| Level | Points Required | Badge Color | Benefits |
|-------|----------------|-------------|----------|
| **Bronze** | 0-99 | 🥉 Bronze | Basic access |
| **Silver** | 100-499 | 🥈 Silver | Priority support |
| **Gold** | 500-999 | 🥇 Gold | Premium features |
| **Platinum** | 1000+ | 💎 Platinum | All features + VIP |

---

## 🔔 Notifications

Saat achievement unlocked, system otomatis:
1. **Update user points** dan recalculate level
2. **Create notification** dengan details:
   - Title: "Achievement Unlocked"
   - Message: "{Achievement Name}! +{Points} points"
   - Type: SUCCESS
3. **Trigger haptic feedback** di mobile (jika supported)

---

## 🧪 Testing Examples

### cURL Examples

**Get All Achievements:**
```bash
curl http://localhost:8080/api/achievements \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Get User Achievements:**
```bash
curl http://localhost:8080/api/profile/achievements \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Get User Stats:**
```bash
curl http://localhost:8080/api/profile/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Award Achievement (Admin):**
```bash
curl -X POST http://localhost:8080/api/admin/achievements/award \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 10,
    "achievement_code": "FIRST_LOGIN"
  }'
```

### JavaScript/Axios Examples

**Get User Achievements:**
```javascript
const response = await axios.get('/api/profile/achievements')
const achievements = response.data.data

// Filter unlocked
const unlockedAchievements = achievements.filter(a => a.unlocked)

// Filter locked
const lockedAchievements = achievements.filter(a => !a.unlocked)
```

**Get User Stats:**
```javascript
const response = await axios.get('/api/profile/stats')
const stats = response.data.data

console.log(`Level: ${stats.current_level}`)
console.log(`Points: ${stats.total_points}`)
console.log(`Completion: ${stats.completion_percentage}%`)
```

**Award Achievement (Admin):**
```javascript
const response = await axios.post('/api/admin/achievements/award', {
  user_id: userId,
  achievement_code: 'FIRST_LOGIN'
})

console.log(response.data.message)
```

---

## 🎨 Frontend Integration Tips

### Achievement Grid Component
```vue
<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'

const achievements = ref([])

const unlockedAchievements = computed(() => 
  achievements.value.filter(a => a.unlocked)
)

const lockedAchievements = computed(() => 
  achievements.value.filter(a => !a.unlocked)
)

const fetchAchievements = async () => {
  const response = await axios.get('/api/profile/achievements')
  achievements.value = response.data.data
}

onMounted(() => {
  fetchAchievements()
})
</script>

<template>
  <div class="grid grid-cols-3 gap-4">
    <div 
      v-for="item in achievements" 
      :key="item.achievement.id"
      :class="{ 'opacity-50 grayscale': !item.unlocked }"
    >
      <div class="achievement-card">
        <div class="text-4xl">{{ item.achievement.icon }}</div>
        <h3>{{ item.achievement.name }}</h3>
        <p>+{{ item.achievement.points }} points</p>
        <span v-if="item.unlocked" class="badge">Unlocked</span>
      </div>
    </div>
  </div>
</template>
```

### Progress Card Component
```vue
<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const stats = ref(null)

const fetchStats = async () => {
  const response = await axios.get('/api/profile/stats')
  stats.value = response.data.data
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <div v-if="stats" class="progress-card">
    <h2>{{ stats.current_level }}</h2>
    <div class="points">{{ stats.total_points }} points</div>
    
    <div class="progress-bar">
      <div 
        class="progress-fill" 
        :style="{ width: `${stats.completion_percentage}%` }"
      />
    </div>
    
    <p class="next-level">
      {{ stats.points_to_next_level }} points to {{ stats.next_level }}
    </p>
  </div>
</template>
```

---

## 📚 Related Documentation

- [Profile API](./profile.md)
- [Notifications API](./notifications.md)
- [User Management API](./user-management.md)
- [Error Handling Guide](../05-guides/error-handling.md)

---

**Last Updated:** 28 Desember 2025  
**Sprint:** Sprint 5  
**Status:** ✅ Production Ready
