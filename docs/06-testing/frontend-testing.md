# 🎨 Frontend Testing Guide

Panduan lengkap untuk testing frontend Vue 3 dengan Vitest dan PWA testing.

---

## 📋 Daftar Isi

1. [Setup Vitest](#setup-vitest)
2. [Component Testing](#component-testing)
3. [Composable Testing](#composable-testing)
4. [PWA Testing](#pwa-testing)
5. [Running Tests](#running-tests)

---

## ⚙️ Setup Vitest

Frontend testing menggunakan **Vitest** (built-in dengan Vite) dan **@vue/test-utils** untuk component testing.

### **Install Dependencies**

```bash
cd frontend

# Install testing dependencies
yarn add -D vitest @vue/test-utils jsdom @vitest/ui
```

### **Configure Vitest**

**Update `vite.config.js`:**

```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/__tests__/setup.js', // Optional
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html', 'json'],
      exclude: [
        'node_modules/',
        'src/__tests__/',
      ]
    }
  }
})
```

### **Update package.json**

Add test scripts:

```json
{
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "test": "vitest",
    "test:ui": "vitest --ui",
    "test:coverage": "vitest --coverage",
    "test:run": "vitest run"
  }
}
```

### **Create Test Setup File** (Optional)

```javascript
// frontend/src/__tests__/setup.js
import { config } from '@vue/test-utils'

// Global test configuration
config.global.mocks = {
  $t: (key) => key, // Mock i18n if used
}

// Mock window.matchMedia (for responsive tests)
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => {},
  }),
})
```

---

## 🧩 Component Testing

### **Component Test Structure**

```
frontend/src/components/
├── UserCard.vue
├── __tests__/
│   └── UserCard.test.js       # Test file
├── ModalConfirm.vue
└── __tests__/
    └── ModalConfirm.test.js   # Test file
```

### **Example: Component Test**

```javascript
// frontend/src/components/__tests__/UserCard.test.js
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import UserCard from '../UserCard.vue'

describe('UserCard', () => {
  const mockUser = {
    id: 1,
    name: 'John Doe',
    email: 'john@example.com',
    role: 'admin',
    department: 'IT',
    is_active: true
  }

  it('renders user data correctly', () => {
    const wrapper = mount(UserCard, {
      props: { user: mockUser }
    })

    expect(wrapper.text()).toContain('John Doe')
    expect(wrapper.text()).toContain('john@example.com')
    expect(wrapper.text()).toContain('admin')
  })

  it('shows active badge when user is active', () => {
    const wrapper = mount(UserCard, {
      props: { user: mockUser }
    })

    const badge = wrapper.find('[data-test="status-badge"]')
    expect(badge.exists()).toBe(true)
    expect(badge.text()).toContain('Aktif')
    expect(badge.classes()).toContain('bg-green-100')
  })

  it('emits edit event when edit button clicked', async () => {
    const wrapper = mount(UserCard, {
      props: { user: mockUser }
    })

    const editButton = wrapper.find('[data-test="edit-button"]')
    await editButton.trigger('click')

    expect(wrapper.emitted('edit')).toBeTruthy()
    expect(wrapper.emitted('edit')[0]).toEqual([mockUser])
  })

  it('emits delete event when delete button clicked', async () => {
    const wrapper = mount(UserCard, {
      props: { user: mockUser }
    })

    const deleteButton = wrapper.find('[data-test="delete-button"]')
    await deleteButton.trigger('click')

    expect(wrapper.emitted('delete')).toBeTruthy()
    expect(wrapper.emitted('delete')[0]).toEqual([mockUser.id])
  })

  it('applies correct role badge color', () => {
    const wrapper = mount(UserCard, {
      props: { user: mockUser }
    })

    const roleBadge = wrapper.find('[data-test="role-badge"]')
    // Admin role should have purple badge
    expect(roleBadge.classes()).toContain('bg-purple-100')
  })
})
```

### **Example: Modal Component Test**

```javascript
// frontend/src/components/__tests__/ModalConfirm.test.js
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import ModalConfirm from '../ModalConfirm.vue'

describe('ModalConfirm', () => {
  it('renders with correct title and message', () => {
    const wrapper = mount(ModalConfirm, {
      props: {
        show: true,
        title: 'Konfirmasi Hapus',
        message: 'Apakah Anda yakin ingin menghapus user ini?'
      }
    })

    expect(wrapper.text()).toContain('Konfirmasi Hapus')
    expect(wrapper.text()).toContain('Apakah Anda yakin')
  })

  it('emits confirm event when confirm button clicked', async () => {
    const wrapper = mount(ModalConfirm, {
      props: {
        show: true,
        title: 'Confirm',
        message: 'Are you sure?'
      }
    })

    const confirmButton = wrapper.find('[data-test="confirm-button"]')
    await confirmButton.trigger('click')

    expect(wrapper.emitted('confirm')).toBeTruthy()
  })

  it('emits cancel event when cancel button clicked', async () => {
    const wrapper = mount(ModalConfirm, {
      props: {
        show: true,
        title: 'Confirm',
        message: 'Are you sure?'
      }
    })

    const cancelButton = wrapper.find('[data-test="cancel-button"]')
    await cancelButton.trigger('click')

    expect(wrapper.emitted('cancel')).toBeTruthy()
  })

  it('does not render when show prop is false', () => {
    const wrapper = mount(ModalConfirm, {
      props: {
        show: false,
        title: 'Test',
        message: 'Test message'
      }
    })

    const modal = wrapper.find('[data-test="modal"]')
    expect(modal.exists()).toBe(false)
  })
})
```

---

## 🪝 Composable Testing

Test reusable logic dalam composables untuk memastikan correctness.

### **Composable Test Structure**

```
frontend/src/composables/
├── useUsers.js
├── __tests__/
│   └── useUsers.test.js       # Test file
├── useAuth.js
└── __tests__/
    └── useAuth.test.js        # Test file
```

### **Example: Composable Test**

```javascript
// frontend/src/composables/__tests__/useUsers.test.js
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useUsers } from '../useUsers'
import { ref } from 'vue'

// Mock API module
vi.mock('../useApi', () => ({
  useApi: () => ({
    get: vi.fn(() => Promise.resolve({ 
      data: { 
        data: [
          { id: 1, name: 'User 1' },
          { id: 2, name: 'User 2' }
        ] 
      } 
    })),
    post: vi.fn(() => Promise.resolve({ 
      data: { 
        data: { id: 3, name: 'New User' } 
      } 
    })),
    put: vi.fn(() => Promise.resolve({ 
      data: { 
        data: { id: 1, name: 'Updated User' } 
      } 
    })),
    delete: vi.fn(() => Promise.resolve({ 
      data: { message: 'User deleted' } 
    }))
  })
}))

describe('useUsers', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('fetches users successfully', async () => {
    const { users, fetchUsers, loading } = useUsers()

    expect(loading.value).toBe(false)
    
    await fetchUsers()

    expect(loading.value).toBe(false)
    expect(users.value).toHaveLength(2)
    expect(users.value[0].name).toBe('User 1')
  })

  it('creates user successfully', async () => {
    const { createUser, loading, error } = useUsers()

    const newUser = {
      name: 'New User',
      email: 'new@example.com',
      role: 'user'
    }

    const result = await createUser(newUser)

    expect(result).toBeDefined()
    expect(result.id).toBe(3)
    expect(error.value).toBeNull()
  })

  it('updates user successfully', async () => {
    const { updateUser, error } = useUsers()

    const updates = {
      id: 1,
      name: 'Updated User'
    }

    const result = await updateUser(updates)

    expect(result).toBeDefined()
    expect(result.name).toBe('Updated User')
    expect(error.value).toBeNull()
  })

  it('deletes user successfully', async () => {
    const { deleteUser, error } = useUsers()

    await deleteUser(1)

    expect(error.value).toBeNull()
  })

  it('handles fetch error correctly', async () => {
    // Mock API error
    vi.mock('../useApi', () => ({
      useApi: () => ({
        get: vi.fn(() => Promise.reject(new Error('Network error')))
      })
    }))

    const { fetchUsers, error } = useUsers()

    await fetchUsers()

    expect(error.value).toBeDefined()
  })
})
```

### **Example: Auth Composable Test**

```javascript
// frontend/src/composables/__tests__/useAuth.test.js
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useAuth } from '../useAuth'

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
global.localStorage = localStorageMock

describe('useAuth', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.getItem.mockReturnValue(null)
  })

  it('initializes with no user when no token in localStorage', () => {
    const { user, isAuthenticated } = useAuth()

    expect(user.value).toBeNull()
    expect(isAuthenticated.value).toBe(false)
  })

  it('login sets user and token', async () => {
    const { login, user, isAuthenticated } = useAuth()

    const credentials = {
      email: 'test@example.com',
      password: 'password123'
    }

    await login(credentials)

    expect(user.value).toBeDefined()
    expect(isAuthenticated.value).toBe(true)
    expect(localStorageMock.setItem).toHaveBeenCalledWith('token', expect.any(String))
  })

  it('logout clears user and token', () => {
    const { logout, user, isAuthenticated } = useAuth()

    logout()

    expect(user.value).toBeNull()
    expect(isAuthenticated.value).toBe(false)
    expect(localStorageMock.removeItem).toHaveBeenCalledWith('token')
  })

  it('checks user role correctly', () => {
    const { user, isAdmin, isManager } = useAuth()

    // Set user as admin
    user.value = { role: 'admin' }

    expect(isAdmin.value).toBe(true)
    expect(isManager.value).toBe(false)
  })
})
```

---

## 🌐 PWA Testing

### **Service Worker Testing**

#### **Check Service Worker Registration**

```javascript
// frontend/src/__tests__/service-worker.test.js
import { describe, it, expect, vi } from 'vitest'

describe('Service Worker', () => {
  it('registers service worker in production', async () => {
    // Mock navigator.serviceWorker
    const mockRegister = vi.fn(() => Promise.resolve({
      active: { state: 'activated' }
    }))

    global.navigator.serviceWorker = {
      register: mockRegister
    }

    // Import register function
    const { registerSW } = await import('../registerSW')

    await registerSW()

    expect(mockRegister).toHaveBeenCalledWith('/sw.js')
  })

  it('handles registration failure gracefully', async () => {
    const mockRegister = vi.fn(() => Promise.reject(new Error('Registration failed')))

    global.navigator.serviceWorker = {
      register: mockRegister
    }

    const { registerSW } = await import('../registerSW')

    // Should not throw
    await expect(registerSW()).resolves.not.toThrow()
  })
})
```

### **Offline Detection Testing**

```javascript
// frontend/src/composables/__tests__/useOnline.test.js
import { describe, it, expect, vi } from 'vitest'
import { useOnline } from '../useOnline'

describe('useOnline', () => {
  it('detects online status correctly', () => {
    // Mock navigator.onLine
    Object.defineProperty(navigator, 'onLine', {
      writable: true,
      value: true
    })

    const { isOnline } = useOnline()

    expect(isOnline.value).toBe(true)
  })

  it('detects offline status correctly', () => {
    Object.defineProperty(navigator, 'onLine', {
      writable: true,
      value: false
    })

    const { isOnline } = useOnline()

    expect(isOnline.value).toBe(false)
  })

  it('updates status on online event', () => {
    const { isOnline } = useOnline()

    // Trigger online event
    const event = new Event('online')
    window.dispatchEvent(event)

    expect(isOnline.value).toBe(true)
  })

  it('updates status on offline event', () => {
    const { isOnline } = useOnline()

    // Trigger offline event
    const event = new Event('offline')
    window.dispatchEvent(event)

    expect(isOnline.value).toBe(false)
  })
})
```

---

## ▶️ Running Tests

### **Run All Tests**

```bash
cd frontend

# Run all tests
yarn test

# Run tests in watch mode (re-run on file changes)
yarn test --watch

# Run tests once (CI mode)
yarn test:run
```

### **Run Specific Tests**

```bash
# Run tests for specific file
yarn test UserCard.test.js

# Run tests matching pattern
yarn test useAuth

# Run tests in specific folder
yarn test src/composables/__tests__
```

### **Test Coverage**

```bash
# Generate coverage report
yarn test:coverage

# Coverage report will be in coverage/ folder
# Open coverage/index.html in browser
```

**Expected Coverage:**

```
File                  | % Stmts | % Branch | % Funcs | % Lines
----------------------|---------|----------|---------|--------
All files             |   75.32 |    68.45 |   78.12 |   76.89
components/           |   72.18 |    65.32 |   75.45 |   73.21
composables/          |   82.45 |    75.18 |   85.32 |   83.67
views/                |   68.92 |    61.45 |   70.18 |   69.78
```

### **UI Mode** (Interactive)

Vitest UI provides interactive test runner dengan visual feedback:

```bash
yarn test:ui

# Opens browser at http://localhost:51204/__vitest__/
```

**Features:**
- ✅ Visual test tree
- ✅ Filter tests by name/file
- ✅ View test execution time
- ✅ View console logs per test
- ✅ Re-run specific tests

---

## ✅ Best Practices

### **1. Use data-test Attributes**

Add test-specific attributes untuk stable selectors:

```vue
<template>
  <button 
    data-test="edit-button"
    class="btn-primary"
    @click="handleEdit"
  >
    Edit
  </button>
</template>
```

```javascript
// In test
const button = wrapper.find('[data-test="edit-button"]')
await button.trigger('click')
```

### **2. Test User Behavior, Not Implementation**

Focus pada what user sees/does, bukan internal implementation:

```javascript
// ❌ Bad - Testing implementation
it('calls handleSubmit method', async () => {
  const handleSubmit = vi.fn()
  wrapper.vm.handleSubmit = handleSubmit
  // ...
})

// ✅ Good - Testing behavior
it('submits form when button clicked', async () => {
  await wrapper.find('form').trigger('submit')
  expect(wrapper.emitted('submit')).toBeTruthy()
})
```

### **3. Mock External Dependencies**

Mock API calls, localStorage, dan external services:

```javascript
vi.mock('../useApi', () => ({
  useApi: () => ({
    get: vi.fn(() => Promise.resolve({ data: [] }))
  })
}))
```

### **4. Clean Up After Tests**

Use beforeEach/afterEach untuk cleanup:

```javascript
describe('UserList', () => {
  let wrapper

  beforeEach(() => {
    wrapper = mount(UserList)
  })

  afterEach(() => {
    wrapper.unmount()
  })

  // Tests...
})
```

### **5. Test Accessibility**

Ensure components accessible:

```javascript
it('has proper ARIA labels', () => {
  const button = wrapper.find('button')
  expect(button.attributes('aria-label')).toBe('Close modal')
})
```

---

## 📚 Related Documentation

- [overview.md](./overview.md) - Testing strategy
- [backend-testing.md](./backend-testing.md) - Backend testing
- [api-testing.md](./api-testing.md) - API testing

---

## 📞 Support

Jika ada pertanyaan tentang frontend testing:
- Developer: Zulfikar Hidayatullah
- Phone: +62 857-1583-8733

---

**Last Updated:** 28 Desember 2025  
**Version:** 2.0.0 (Phase 2 Restructure)
