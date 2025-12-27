package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

// TranslateValidationErrors mengkonversi validation errors ke format yang user-friendly
// dengan pesan dalam Bahasa Indonesia
func TranslateValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = GetErrorMessage(e)
		}
	}
	
	return errors
}

// GetErrorMessage menghasilkan error message dalam Bahasa Indonesia
// berdasarkan validation tag yang gagal
func GetErrorMessage(e validator.FieldError) string {
	field := e.Field()
	param := e.Param()
	
	// Mapping field names ke label yang lebih user-friendly
	fieldLabels := map[string]string{
		"NIP":        "NIP",
		"Email":      "Email",
		"Password":   "Password",
		"FullName":   "Nama Lengkap",
		"Phone":      "Nomor Telepon",
		"Role":       "Role",
		"Department": "Departemen",
		"Shift":      "Shift",
	}
	
	label := fieldLabels[field]
	if label == "" {
		label = field
	}
	
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s harus diisi", label)
	
	case "email":
		return fmt.Sprintf("%s harus berupa alamat email yang valid", label)
	
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", label, param)
	
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", label, param)
	
	case "len":
		return fmt.Sprintf("%s harus %s karakter", label, param)
	
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka", label)
	
	case "alpha":
		return fmt.Sprintf("%s hanya boleh berisi huruf", label)
	
	case "alphanum":
		return fmt.Sprintf("%s hanya boleh berisi huruf dan angka", label)
	
	case "oneof":
		return fmt.Sprintf("%s tidak valid", label)
	
	case "eqfield":
		return fmt.Sprintf("%s harus sama dengan %s", label, param)
	
	case "nefield":
		return fmt.Sprintf("%s tidak boleh sama dengan %s", label, param)
	
	case "gt":
		return fmt.Sprintf("%s harus lebih besar dari %s", label, param)
	
	case "gte":
		return fmt.Sprintf("%s harus lebih besar atau sama dengan %s", label, param)
	
	case "lt":
		return fmt.Sprintf("%s harus lebih kecil dari %s", label, param)
	
	case "lte":
		return fmt.Sprintf("%s harus lebih kecil atau sama dengan %s", label, param)
	
	case "url":
		return fmt.Sprintf("%s harus berupa URL yang valid", label)
	
	case "contains":
		return fmt.Sprintf("%s harus mengandung '%s'", label, param)
	
	case "containsany":
		return fmt.Sprintf("%s harus mengandung salah satu karakter: %s", label, param)
	
	case "startswith":
		return fmt.Sprintf("%s harus dimulai dengan '%s'", label, param)
	
	case "endswith":
		return fmt.Sprintf("%s harus diakhiri dengan '%s'", label, param)
	
	case "required_if":
		return fmt.Sprintf("%s harus diisi", label)
	
	case "required_unless":
		return fmt.Sprintf("%s harus diisi", label)
	
	case "required_with":
		return fmt.Sprintf("%s harus diisi jika %s ada", label, param)
	
	case "required_without":
		return fmt.Sprintf("%s harus diisi jika %s tidak ada", label, param)
	
	default:
		return fmt.Sprintf("%s tidak valid", label)
	}
}

// ValidationErrorResponse merupakan struktur response untuk validation errors
type ValidationErrorResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

// NewValidationErrorResponse membuat response untuk validation errors
func NewValidationErrorResponse(err error) ValidationErrorResponse {
	return ValidationErrorResponse{
		Success: false,
		Message: "Data yang dikirim tidak valid",
		Errors:  TranslateValidationErrors(err),
	}
}
