# 🔧 Backend Development

Dokumentasi lengkap untuk backend development dalam Sirine Go App yang dibangun dengan Go, Gin framework, dan MySQL.

## 📚 Guides

### Getting Started
**[Getting Started Guide](./getting-started.md)**  
Panduan awal untuk setup development environment backend, yaitu:
- Prerequisites dan instalasi
- Setup database dan migrations
- Development commands
- Project structure overview

### Service Pattern Architecture
**[Service Pattern Guide](./service-pattern.md)**  
Panduan implementasi Service Pattern yang merupakan arsitektur utama aplikasi, mencakup:
- Konsep dan prinsip Service Pattern
- Layering: Handler → Service → Model
- Best practices untuk struktur code
- Contoh implementasi lengkap

### Middleware Development
**[Middleware Guide](./middleware.md)**  
Panduan untuk membuat dan menggunakan middleware, antara lain:
- Authentication middleware dengan JWT
- Role-based authorization
- Request logging dan CORS
- Custom middleware patterns

## 🏗️ Tech Stack

Backend Sirine Go App menggunakan:
- **Go 1.21+** - Programming language
- **Gin** - Web framework untuk routing dan middleware
- **GORM** - ORM untuk database operations
- **MySQL 8.0+** - Relational database
- **JWT** - Token-based authentication
- **Bcrypt** - Password hashing

## 📋 Quick Links

- [API Reference](../../04-api-reference/README.md)
- [Database Models](../../05-guides/database/models.md)
- [Configuration Guide](../../05-guides/configuration/README.md)
- [Testing Guide](../../06-testing/backend-testing.md)

## 🎯 Development Philosophy

Backend development dalam Sirine Go mengikuti prinsip:
1. **Service Pattern** untuk separation of concerns
2. **RESTful API** design yang consistent
3. **Security-first** approach dengan JWT dan validation
4. **Activity Logging** untuk audit trail
5. **Error Handling** yang comprehensive

---

**Last Updated:** 28 Desember 2025
