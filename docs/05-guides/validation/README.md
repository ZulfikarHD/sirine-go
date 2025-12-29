# ✅ Validation System

## Overview

Validation di Sirine Go mengikuti pola **Laravel-style** menggunakan Gin framework, memudahkan validasi input request dengan syntax yang familiar dan ekspresif.

## Guides

### 1. [Validation Guide](./guide.md)
Konsep dasar dan cara penggunaan validator.
- Basic Usage
- Available Rules
- Custom Validators
- Error Handling

### 2. [Examples & Test Cases](./examples.md)
Contoh praktis penggunaan validasi untuk berbagai skenario.
- User Registration
- Product Input
- Complex Logic

---

## Quick Example

```go
req := struct {
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=18"`
}{}

if err := c.ShouldBindJSON(&req); err != nil {
    // Handle validation error
}
```
