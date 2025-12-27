import { vi } from 'vitest'

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
global.localStorage = localStorageMock

// Mock navigator.vibrate untuk haptic feedback tests
global.navigator.vibrate = vi.fn()

// Mock window.scrollTo
global.scrollTo = vi.fn()
