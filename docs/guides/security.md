# 🛡️ Security Guide

## Overview

Security Guide ini menjelaskan standar keamanan dan best practices yang diterapkan dalam Sirine Go App untuk melindungi data sensitif dan mencegah serangan umum.

## 🔐 Core Security Concepts

### 1. Authentication
- **JWT (JSON Web Tokens)**: Digunakan untuk stateless authentication. Token memiliki masa berlaku pendek (15 menit).
- **Refresh Tokens**: Token jangka panjang (30 hari) yang disimpan di database (hashed) untuk memperbarui access token tanpa login ulang.
- **Bcrypt Hashing**: Password user di-hash menggunakan bcrypt dengan cost 12 sebelum disimpan di database.

### 2. Authorization
- **Role-Based Access Control (RBAC)**: Akses ke endpoint dibatasi berdasarkan role user (ADMIN, MANAGER, STAFF, dll).
- **Middleware**: `AuthMiddleware` memvalidasi token, dan `RoleMiddleware` memverifikasi hak akses.

### 3. Session Management
- **Concurrent Logins**: Setiap login membuat sesi baru. Logout menghapus sesi spesifik.
- **Revocation**: Sesi dapat dicabut (revoked) dari sisi server jika terdeteksi aktivitas mencurigakan.

## 🛡️ Vulnerability Protection

### Cross-Site Scripting (XSS)
- **Frontend**: Vue 3 secara otomatis meng-escape konten dalam data binding `{{ }}`.
- **Backend**: Input sanitization dilakukan sebelum data diproses.
- **Prevention**: Jangan gunakan `v-html` untuk user-generated content.

### SQL Injection
- **GORM**: Menggunakan parameterized queries secara default, yang mencegah SQL injection.
- **Prevention**: Hindari penggunaan raw SQL query dengan string concatenation. Gunakan prepared statements.

### Cross-Site Request Forgery (CSRF)
- **Token Based**: Karena menggunakan JWT (Bearer Token) di header `Authorization`, aplikasi secara inheren lebih resisten terhadap CSRF dibandingkan cookie-based auth.
- **Prevention**: Pastikan token tidak disimpan di cookies yang tidak secure.

### Rate Limiting & Brute Force
- **Account Lockout**: Akun dikunci selama 15 menit setelah 5 kali percobaan login gagal.
- **Rate Limiting**: (Future) Implementasi rate limiting per IP address.

## 📝 Secure Coding Practices

### Data Handling
- **Sensitive Data**: Jangan pernah log password, token, atau PII (Personally Identifiable Information) dalam plain text.
- **Error Messages**: Berikan pesan error yang generik ke user ("Invalid credentials") dan log detail error di server.

### Configuration
- **Environment Variables**: Simpan credentials (DB password, JWT secret) di file `.env`. Jangan commit file `.env` ke repository.
- **HTTPS**: Gunakan HTTPS di production untuk mengenkripsi traffic antara client dan server.

## 🚨 Incident Response

### Reporting Vulnerabilities
Jika Anda menemukan celah keamanan, segera laporkan ke team lead atau security officer.

### Log Monitoring
Pantau `activity_logs` untuk mendeteksi pola aneh seperti:
- Multiple login failures dari IP yang sama.
- Akses ke resource sensitif di luar jam kerja.

## 🔗 Related Documentation
- [Authentication Implementation](./authentication/implementation.md)
- [Authentication Testing](./authentication/testing.md)
