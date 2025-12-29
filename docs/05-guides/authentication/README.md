# 🔐 Authentication System

## Overview

Sistem autentikasi Sirine Go menggunakan JWT (JSON Web Tokens) dengan mekanisme refresh token untuk keamanan dan user experience yang optimal.

## Documentation Map

Dokumentasi autentikasi dibagi menjadi beberapa bagian:

### 1. 🏗️ [Implementation Guide](./implementation.md)
Penjelasan teknis tentang architecture, security flows, dan implementation details.
- Architecture Overview
- Security Features (Rate limiting, lockout)
- Role-based Access Control (RBAC)

### 2. 🔌 [API Reference](./api-reference.md)
Dokumentasi lengkap endpoints autentikasi.
- POST `/api/auth/login`
- POST `/api/auth/refresh`
- POST `/api/auth/logout`
- GET `/api/auth/me`

### 3. 🗺️ [User Journeys](../../user-journeys/authentication/overview.md)
Visualisasi flow pengguna (User Experience).
- Login Flows (Admin vs Staff)
- Error Scenarios
- Session Management

### 4. 🧪 [Testing Guide](./testing.md)
Panduan testing untuk modul autentikasi.
- Unit Tests
- Integration Tests
- Manual Testing Checklist

---

## Quick Links

- [User Journeys Overview](../../user-journeys/authentication/overview.md)
- [Database Models](../database/models.md) (User, UserSession)
