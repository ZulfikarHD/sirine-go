# Sprint 2: Quick Reference Card

## 🚀 Quick Commands

```bash
# Start Backend
cd backend && go run cmd/server/main.go

# Start Frontend
cd frontend && yarn dev

# Build Backend
cd backend && go build -o bin/server cmd/server/main.go

# Run Tests (future)
cd backend && go test ./...
cd frontend && yarn test
```

---

## 📁 New Files Created

### Backend
```
backend/
├── services/
│   └── user_service.go          # User CRUD operations
├── handlers/
│   ├── user_handler.go          # Admin user management endpoints
│   └── profile_handler.go       # Self-service profile endpoints
└── middleware/
    └── activity_logger.go       # Auto-log middleware
```

### Frontend
```
frontend/src/
├── stores/
│   └── user.js                  # User management state
├── views/
│   ├── admin/users/
│   │   └── UserList.vue         # User management page
│   └── profile/
│       └── EditProfile.vue      # Profile editor
├── components/
│   └── admin/
│       ├── RoleBadge.vue        # Role badge component
│       └── UserFormModal.vue    # Create/Edit modal
└── router/
    └── index.js                 # Updated routes
```

---

## 🔌 API Endpoints

### User Management (Admin/Manager)
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/users` | List users | Admin/Manager |
| GET | `/api/users/search?q=` | Search users | Admin/Manager |
| GET | `/api/users/:id` | Get user detail | Admin/Manager |
| POST | `/api/users` | Create user | Admin |
| PUT | `/api/users/:id` | Update user | Admin |
| DELETE | `/api/users/:id` | Delete user | Admin |
| POST | `/api/users/bulk-delete` | Bulk delete | Admin |
| POST | `/api/users/bulk-update-status` | Bulk update | Admin |

### Profile (Self-Service)
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/profile` | Get own profile | All users |
| PUT | `/api/profile` | Update own profile | All users |

---

## 🎨 Frontend Routes

| Path | Component | Access |
|------|-----------|--------|
| `/admin/users` | UserList.vue | Admin/Manager |
| `/profile` | Profile.vue | All users |
| `/profile/edit` | EditProfile.vue | All users |

---

## 🔐 Role-Based Access

### Admin/Manager
- ✅ View all users
- ✅ Create users
- ✅ Edit users (all fields)
- ✅ Delete users
- ✅ Search & filter users
- ✅ Bulk operations
- ✅ View own profile
- ✅ Edit own profile

### All Other Roles
- ❌ Cannot access user management
- ✅ View own profile
- ✅ Edit own profile (name, email, phone only)

---

## 📊 Database Tables Used

### Users Table
```sql
users (
  id, nip, full_name, email, phone,
  password_hash, role, department, shift,
  profile_photo_url, status, must_change_password,
  last_login_at, failed_login_attempts, locked_until,
  created_at, updated_at, deleted_at
)
```

### Activity Logs Table
```sql
activity_logs (
  id, user_id, action, entity_type, entity_id,
  changes (JSON), ip_address, user_agent,
  created_at
)
```

---

## 🎯 Key Features

### User Management
- ✅ CRUD operations
- ✅ Search by NIP/name/email
- ✅ Filter by role/department/status
- ✅ Pagination (20 per page)
- ✅ Auto-generated passwords (12 char)
- ✅ Soft delete
- ✅ Bulk operations

### Profile Management
- ✅ View own profile
- ✅ Edit name/email/phone
- ✅ Read-only NIP/role/department
- ✅ Validation

### Activity Logging
- ✅ Auto-log CREATE/UPDATE/DELETE
- ✅ Before/after changes (JSON)
- ✅ IP address & user agent
- ✅ Background processing

---

## 🔧 Common Code Patterns

### Backend: Create Service
```go
userService := services.NewUserService(db, passwordService)
```

### Backend: Create Handler
```go
userHandler := handlers.NewUserHandler(userService)
```

### Backend: Apply Middleware
```go
users.Use(middleware.AuthMiddleware(db, cfg))
users.Use(middleware.RequireRole("ADMIN", "MANAGER"))
users.Use(middleware.ActivityLogger(db))
```

### Frontend: Use Store
```javascript
import { useUserStore } from '@/stores/user'
const userStore = useUserStore()
await userStore.fetchUsers(1)
```

### Frontend: Navigate
```javascript
import { useRouter } from 'vue-router'
const router = useRouter()
router.push('/profile/edit')
```

---

## 🐛 Troubleshooting

### Backend Issues

**Error: "NIP sudah terdaftar"**
- Check database: `SELECT * FROM users WHERE nip = 'xxxxx'`
- May include soft-deleted users

**Error: "email sudah terdaftar"**
- Check database: `SELECT * FROM users WHERE email = 'xxx@xxx.com'`
- May include soft-deleted users

**Activity logs not created**
- Check middleware applied to route
- Check context values set in handler
- Check database connection

### Frontend Issues

**Users not loading**
- Check API endpoint: `http://localhost:8080/api/users`
- Check auth token in localStorage
- Check browser console for errors

**Modal not showing**
- Check `showFormModal` state
- Check Teleport target exists
- Check z-index conflicts

**Filters not working**
- Check `userStore.setFilters()` called
- Check API query params sent
- Check backend filter logic

---

## 📝 Validation Rules

### User Creation
- **NIP**: Required, max 5 char, unique
- **Full Name**: Required, min 3 char, max 100 char
- **Email**: Required, valid email format, unique
- **Phone**: Required, min 10 char, max 15 char, format 08xxx
- **Role**: Required, enum value
- **Department**: Required, enum value
- **Shift**: Optional, enum value, default PAGI

### Profile Update
- **Full Name**: Required, min 3 char, max 100 char
- **Email**: Required, valid email format, unique (exclude self)
- **Phone**: Required, min 10 char, max 15 char, format 08xxx

---

## 🎨 UI Components

### RoleBadge
```vue
<RoleBadge :role="user.role" />
```

**Colors:**
- Admin: Indigo
- Manager: Purple
- Staff Khazwal: Blue
- Operator Cetak: Fuchsia
- QC Inspector: Emerald
- Verifikator: Amber
- Staff Khazkhir: Cyan

### UserFormModal
```vue
<UserFormModal
  v-if="showFormModal"
  :user="selectedUser"
  @close="closeFormModal"
  @success="handleFormSuccess"
/>
```

**Props:**
- `user`: User object untuk edit mode (null untuk create)

**Events:**
- `@close`: Modal closed
- `@success`: User created/updated successfully

---

## 🔍 Useful Queries

### Find User by NIP
```sql
SELECT * FROM users WHERE nip = '12345';
```

### View Recent Activity Logs
```sql
SELECT 
  al.*,
  u.full_name as user_name
FROM activity_logs al
JOIN users u ON al.user_id = u.id
ORDER BY al.created_at DESC
LIMIT 10;
```

### View Soft Deleted Users
```sql
SELECT * FROM users WHERE deleted_at IS NOT NULL;
```

### Count Users by Role
```sql
SELECT role, COUNT(*) as count
FROM users
WHERE deleted_at IS NULL
GROUP BY role;
```

### View User Changes
```sql
SELECT 
  id,
  action,
  entity_type,
  entity_id,
  JSON_PRETTY(changes) as changes
FROM activity_logs
WHERE action = 'UPDATE' AND entity_type = 'users'
ORDER BY created_at DESC
LIMIT 5;
```

---

## 📚 Related Documentation

- [Sprint 2 Summary](SPRINT2_SUMMARY.md) - Complete implementation details
- [Sprint 2 Testing Guide](SPRINT2_TESTING_GUIDE.md) - Testing scenarios
- [Sprint 1 Summary](SPRINT1_SUMMARY.md) - Authentication foundation
- [Main README](README.md) - Project overview
- [Sprint Plan](/.cursor/plans/sprint_plan_-_authentication_fa6ccc79.plan.md) - Full 6-sprint plan

---

## 🎯 Next Sprint Preview

**Sprint 3: Password Management & Security**
- Change password (self-service)
- Forgot password flow
- Reset password via email
- First-time login force change
- Password strength indicator
- Token refresh mechanism

---

**Quick Reference Version**: 1.0
**Last Updated**: December 28, 2025
**Sprint**: 2 of 6
