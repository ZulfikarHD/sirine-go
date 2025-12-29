# Gin Validation Guide - Server-Side Validation

## Overview

Gin Validation Guide merupakan panduan komprehensif untuk server-side validation di aplikasi Sirine Go yang bertujuan untuk memastikan data integrity dan user experience yang baik, yaitu: validation menggunakan go-playground/validator library melalui struct tags, automatic error translation ke Bahasa Indonesia, dan consistent error response format yang frontend-friendly.

Gin menggunakan **go-playground/validator** library untuk server-side validation, yang sangat mirip dengan Laravel Request Validation. Validation dilakukan melalui struct tags dengan syntax `binding:"rules"`.

## Comparison: Laravel vs Gin

### Laravel Request Validation
```php
$request->validate([
    'email' => 'required|email|max:255',
    'password' => 'required|min:8|confirmed',
    'age' => 'required|numeric|min:18|max:100',
    'role' => 'required|in:admin,user,manager',
]);
```

### Gin Validation (Equivalent)
```go
type UserRequest struct {
    Email           string `json:"email" binding:"required,email,max=255"`
    Password        string `json:"password" binding:"required,min=8"`
    PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
    Age             int    `json:"age" binding:"required,numeric,min=18,max=100"`
    Role            string `json:"role" binding:"required,oneof=admin user manager"`
}
```

---

## Common Validation Rules

### 1. Required & Optional
```go
type Example struct {
    Name     string `json:"name" binding:"required"`           // Required field
    Nickname string `json:"nickname"`                          // Optional field
    Bio      string `json:"bio" binding:"omitempty,min=10"`    // Optional, but if provided min 10 chars
}
```

**Laravel Equivalent:**
```php
'name' => 'required',
'nickname' => 'nullable',
'bio' => 'nullable|min:10',
```

---

### 2. String Validation
```go
type StringValidation struct {
    Email       string `json:"email" binding:"required,email"`
    URL         string `json:"url" binding:"required,url"`
    Alpha       string `json:"alpha" binding:"required,alpha"`           // Only letters
    Alphanumeric string `json:"alphanumeric" binding:"required,alphanum"` // Letters + numbers
    MinLength   string `json:"min_length" binding:"required,min=5"`
    MaxLength   string `json:"max_length" binding:"required,max=100"`
    Length      string `json:"length" binding:"required,len=10"`         // Exact length
    Contains    string `json:"contains" binding:"required,contains=@"`   // Must contain @
}
```

**Laravel Equivalent:**
```php
'email' => 'required|email',
'url' => 'required|url',
'alpha' => 'required|alpha',
'alphanumeric' => 'required|alpha_num',
'min_length' => 'required|min:5',
'max_length' => 'required|max:100',
'length' => 'required|size:10',
```

---

### 3. Numeric Validation
```go
type NumericValidation struct {
    Age      int     `json:"age" binding:"required,min=18,max=100"`
    Price    float64 `json:"price" binding:"required,gt=0"`              // Greater than 0
    Discount float64 `json:"discount" binding:"required,gte=0,lte=100"`  // 0-100
    Quantity int     `json:"quantity" binding:"required,numeric"`
}
```

**Laravel Equivalent:**
```php
'age' => 'required|min:18|max:100',
'price' => 'required|gt:0',
'discount' => 'required|gte:0|lte:100',
'quantity' => 'required|numeric',
```

---

### 4. Enum Validation (oneof)
```go
type EnumValidation struct {
    Role       string `json:"role" binding:"required,oneof=ADMIN MANAGER STAFF"`
    Status     string `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
    Department string `json:"department" binding:"required,oneof=KHAZWAL CETAK VERIFIKASI KHAZKHIR"`
}
```

**Laravel Equivalent:**
```php
'role' => 'required|in:ADMIN,MANAGER,STAFF',
'status' => 'required|in:ACTIVE,INACTIVE',
'department' => 'required|in:KHAZWAL,CETAK,VERIFIKASI,KHAZKHIR',
```

---

### 5. Comparison Validation
```go
type ComparisonValidation struct {
    Password        string `json:"password" binding:"required,min=8"`
    PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
    StartDate       string `json:"start_date" binding:"required"`
    EndDate         string `json:"end_date" binding:"required,gtfield=StartDate"`  // Greater than StartDate
}
```

**Laravel Equivalent:**
```php
'password' => 'required|min:8',
'password_confirm' => 'required|same:password',
'start_date' => 'required',
'end_date' => 'required|after:start_date',
```

---

### 6. Array/Slice Validation
```go
type ArrayValidation struct {
    Tags       []string `json:"tags" binding:"required,min=1,max=10,dive,min=2,max=20"`
    // min=1,max=10 -> array must have 1-10 items
    // dive -> validate each item
    // min=2,max=20 -> each item must be 2-20 chars
    
    UserIDs    []uint64 `json:"user_ids" binding:"required,dive,required,gt=0"`
    // Each user_id must be > 0
}
```

**Laravel Equivalent:**
```php
'tags' => 'required|array|min:1|max:10',
'tags.*' => 'required|min:2|max:20',
'user_ids' => 'required|array',
'user_ids.*' => 'required|gt:0',
```

---

### 7. Conditional Validation
```go
type ConditionalValidation struct {
    Type    string `json:"type" binding:"required,oneof=email phone"`
    Email   string `json:"email" binding:"required_if=Type email,omitempty,email"`
    Phone   string `json:"phone" binding:"required_if=Type phone,omitempty,min=10"`
}
```

**Laravel Equivalent:**
```php
'type' => 'required|in:email,phone',
'email' => 'required_if:type,email|email',
'phone' => 'required_if:type,phone|min:10',
```

---

## Real-World Examples dari Sirine Go

### 1. Login Request (Current Implementation)
```go
type LoginRequest struct {
    NIP        string `json:"nip" binding:"required"`
    Password   string `json:"password" binding:"required"`
    RememberMe bool   `json:"remember_me"`
}
```

### 2. Create User Request (Sprint 2)
```go
type CreateUserRequest struct {
    NIP        string `json:"nip" binding:"required,max=5"`
    FullName   string `json:"full_name" binding:"required,min=3,max=100"`
    Email      string `json:"email" binding:"required,email,max=255"`
    Phone      string `json:"phone" binding:"required,min=10,max=15"`
    Role       string `json:"role" binding:"required,oneof=ADMIN MANAGER STAFF_KHAZWAL OPERATOR_CETAK QC_INSPECTOR VERIFIKATOR STAFF_KHAZKHIR"`
    Department string `json:"department" binding:"required,oneof=KHAZWAL CETAK VERIFIKASI KHAZKHIR"`
    Shift      string `json:"shift" binding:"omitempty,oneof=PAGI SIANG MALAM"`
}
```

### 3. Update Profile Request (Sprint 2)
```go
type UpdateProfileRequest struct {
    FullName string `json:"full_name" binding:"required,min=3,max=100"`
    Email    string `json:"email" binding:"required,email,max=255"`
    Phone    string `json:"phone" binding:"required,min=10,max=15"`
}
```

### 4. Change Password Request (Sprint 3)
```go
type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" binding:"required"`
    NewPassword     string `json:"new_password" binding:"required,min=8,containsany=!@#$%^&*(),containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=0123456789"`
    ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}
```

### 5. Forgot Password Request (Sprint 3)
```go
type ForgotPasswordRequest struct {
    NIP   string `json:"nip" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}
```

---

## Custom Validation Messages

### Default Error Handling (Current)
```go
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "success": false,
        "message": "Data yang dikirim tidak valid",
        "error":   err.Error(),
    })
    return
}
```

### Enhanced Error Handling dengan Custom Messages
```go
// Helper function untuk translate validation errors
func TranslateValidationError(err error) map[string]string {
    errors := make(map[string]string)
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            field := e.Field()
            
            switch e.Tag() {
            case "required":
                errors[field] = fmt.Sprintf("%s harus diisi", field)
            case "email":
                errors[field] = fmt.Sprintf("%s harus berupa email yang valid", field)
            case "min":
                errors[field] = fmt.Sprintf("%s minimal %s karakter", field, e.Param())
            case "max":
                errors[field] = fmt.Sprintf("%s maksimal %s karakter", field, e.Param())
            case "oneof":
                errors[field] = fmt.Sprintf("%s harus salah satu dari: %s", field, e.Param())
            case "eqfield":
                errors[field] = fmt.Sprintf("%s harus sama dengan %s", field, e.Param())
            default:
                errors[field] = fmt.Sprintf("%s tidak valid", field)
            }
        }
    }
    
    return errors
}

// Usage in handler
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
        "success": false,
        "message": "Validasi gagal",
        "errors":  TranslateValidationError(err),
    })
    return
}
```

**Response Example:**
```json
{
    "success": false,
    "message": "Validasi gagal",
    "errors": {
        "email": "email harus berupa email yang valid",
        "password": "password minimal 8 karakter",
        "role": "role harus salah satu dari: ADMIN MANAGER STAFF"
    }
}
```

---

## Custom Validators

### Register Custom Validator
```go
// Di main.go atau config
import "github.com/go-playground/validator/v10"

func RegisterCustomValidators() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        // Custom validator untuk NIP format
        v.RegisterValidation("nip_format", func(fl validator.FieldLevel) bool {
            nip := fl.Field().String()
            // NIP harus 5 karakter alphanumeric
            matched, _ := regexp.MatchString(`^[A-Za-z0-9]{5}$`, nip)
            return matched
        })
        
        // Custom validator untuk phone Indonesia
        v.RegisterValidation("phone_id", func(fl validator.FieldLevel) bool {
            phone := fl.Field().String()
            // Phone harus dimulai dengan 08 atau +62
            matched, _ := regexp.MatchString(`^(\+62|08)[0-9]{8,12}$`, phone)
            return matched
        })
    }
}

// Usage
type UserRequest struct {
    NIP   string `json:"nip" binding:"required,nip_format"`
    Phone string `json:"phone" binding:"required,phone_id"`
}
```

---

## Complete Validation Tags Reference

| Tag | Description | Laravel Equivalent |
|-----|-------------|-------------------|
| `required` | Field must be present | `required` |
| `omitempty` | Skip validation if empty | `nullable` |
| `email` | Must be valid email | `email` |
| `url` | Must be valid URL | `url` |
| `alpha` | Only letters | `alpha` |
| `alphanum` | Letters + numbers | `alpha_num` |
| `numeric` | Must be numeric | `numeric` |
| `min=X` | Minimum length/value | `min:X` |
| `max=X` | Maximum length/value | `max:X` |
| `len=X` | Exact length | `size:X` |
| `eq=X` | Equal to value | `same:X` |
| `ne=X` | Not equal to value | `different:X` |
| `gt=X` | Greater than | `gt:X` |
| `gte=X` | Greater than or equal | `gte:X` |
| `lt=X` | Less than | `lt:X` |
| `lte=X` | Less than or equal | `lte:X` |
| `oneof=A B C` | Must be one of values | `in:A,B,C` |
| `eqfield=Field` | Equal to another field | `same:field` |
| `nefield=Field` | Not equal to another field | `different:field` |
| `gtfield=Field` | Greater than another field | `after:field` |
| `ltfield=Field` | Less than another field | `before:field` |
| `contains=X` | Must contain substring | `contains:X` |
| `excludes=X` | Must not contain | - |
| `startswith=X` | Must start with | `starts_with:X` |
| `endswith=X` | Must end with | `ends_with:X` |
| `dive` | Validate array elements | `array`, `*` |
| `required_if=Field Value` | Required if field equals value | `required_if:field,value` |
| `required_unless=Field Value` | Required unless field equals value | `required_unless:field,value` |
| `required_with=Field` | Required if field is present | `required_with:field` |
| `required_without=Field` | Required if field is not present | `required_without:field` |

---

## Best Practices

### 1. Separate Request Structs
```go
// ✅ Good: Separate structs untuk different operations
type CreateUserRequest struct { ... }
type UpdateUserRequest struct { ... }
type UpdateProfileRequest struct { ... }

// ❌ Bad: Single struct untuk semua operations
type UserRequest struct { ... }
```

### 2. Use Descriptive Field Names
```go
// ✅ Good
type LoginRequest struct {
    NIP      string `json:"nip" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// ❌ Bad
type LoginRequest struct {
    N string `json:"n" binding:"required"`
    P string `json:"p" binding:"required"`
}
```

### 3. Validate at Handler Level
```go
// ✅ Good: Validate immediately di handler
func (h *Handler) Create(c *gin.Context) {
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Handle validation error
        return
    }
    // Process request
}

// ❌ Bad: Pass unvalidated data ke service
func (h *Handler) Create(c *gin.Context) {
    var req map[string]interface{}
    c.BindJSON(&req)
    h.service.Create(req) // No validation!
}
```

### 4. Return Detailed Errors
```go
// ✅ Good: Return field-specific errors
{
    "success": false,
    "message": "Validasi gagal",
    "errors": {
        "email": "email harus valid",
        "password": "password minimal 8 karakter"
    }
}

// ❌ Bad: Generic error
{
    "success": false,
    "message": "Invalid input"
}
```

---

## Recommendation untuk Sprint 2

Untuk Sprint 2 (User Management), saya recommend membuat validation helper:

```go
// backend/utils/validation.go
package utils

import (
    "fmt"
    "github.com/go-playground/validator/v10"
)

func TranslateValidationErrors(err error) map[string]string {
    errors := make(map[string]string)
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            errors[e.Field()] = GetErrorMessage(e)
        }
    }
    
    return errors
}

func GetErrorMessage(e validator.FieldError) string {
    field := e.Field()
    
    switch e.Tag() {
    case "required":
        return fmt.Sprintf("%s harus diisi", field)
    case "email":
        return "Email harus valid"
    case "min":
        return fmt.Sprintf("%s minimal %s karakter", field, e.Param())
    case "max":
        return fmt.Sprintf("%s maksimal %s karakter", field, e.Param())
    case "oneof":
        return fmt.Sprintf("%s tidak valid", field)
    case "eqfield":
        return fmt.Sprintf("%s harus sama dengan %s", field, e.Param())
    default:
        return fmt.Sprintf("%s tidak valid", field)
    }
}
```

---

## Summary

**Yes, Gin has powerful server-side validation** yang sangat mirip dengan Laravel Request Validation:

✅ **Declarative validation** via struct tags  
✅ **Rich validation rules** (50+ built-in validators)  
✅ **Custom validators** support  
✅ **Nested validation** untuk complex objects  
✅ **Array validation** dengan `dive`  
✅ **Conditional validation** dengan `required_if`, etc  
✅ **Field comparison** dengan `eqfield`, `gtfield`, etc  

**Main Difference:**
- Laravel: Validation rules di controller/FormRequest
- Gin: Validation rules di struct tags (lebih type-safe)

**Advantage Gin:**
- Type-safe at compile time
- Auto-completion di IDE
- Better performance (no reflection at runtime untuk type checking)

**Advantage Laravel:**
- More readable untuk non-Go developers
- Easier to customize messages
- More flexible conditional logic
