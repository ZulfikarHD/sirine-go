# 🗺️ Authentication User Journeys: Overview

**Feature**: Authentication System  
**Sprint**: 1  
**Last Updated**: 27 Desember 2025

---

## 📋 Overview

User Journeys ini mendokumentasikan visual flow untuk authentication system, mencakup admin login, staff login, error handling, session management, dan logout.

## 🎯 Journey Map Index

| Journey ID | Journey Name | User Type | Status | File |
|------------|--------------|-----------|--------|------|
| J1 | Happy Path - Admin Login | Admin | ✅ Complete | [login-flow.md](./login-flow.md) |
| J2 | Happy Path - Staff Login | Staff | ✅ Complete | [login-flow.md](./login-flow.md) |
| J3 | Error - Invalid Credentials | Any | ✅ Complete | [error-scenarios.md](./error-scenarios.md) |
| J4 | Error - Account Locked | Any | ✅ Complete | [error-scenarios.md](./error-scenarios.md) |
| J5 | Session Persistence | Any | ✅ Complete | [session-management.md](./session-management.md) |
| J6 | Token Expiry & Refresh | Any | ✅ Complete | [session-management.md](./session-management.md) |
| J7 | Logout Flow | Any | ✅ Complete | [logout-flow.md](./logout-flow.md) |
| J8 | Protected Route Access | Any | ✅ Complete | [login-flow.md](./login-flow.md) |

---

## 📊 User Experience Metrics

| Metric | Target | Notes |
|--------|--------|-------|
| Login Time (P50) | < 3s | From form submit to dashboard |
| Login Time (P95) | < 5s | Including slow networks |
| Session Restore | < 100ms | Page refresh experience |
| Token Refresh | < 500ms | Transparent to user |
| Logout Time | < 2s | Complete cleanup |
| Error Recovery | < 5s | User can retry immediately |

---

## 🔗 Related Documentation

- **Implementation**: [Authentication Implementation](../../guides/authentication/implementation.md)
- **API Reference**: [Authentication API](../../guides/authentication/api-reference.md)
- **Testing**: [Authentication Testing](../../guides/authentication/testing.md)
