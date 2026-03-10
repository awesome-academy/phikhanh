package utils

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Đăng ký custom validators
	validate.RegisterValidation("citizen_id", validateCitizenID)
	validate.RegisterValidation("vn_phone", validateVNPhone)
	validate.RegisterValidation("strong_password", validateStrongPassword)
	validate.RegisterValidation("past_date", validatePastDate)
}

// RegisterCustomValidators - Đăng ký custom validators cho Gin
func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("citizen_id", validateCitizenID)
		v.RegisterValidation("vn_phone", validateVNPhone)
		v.RegisterValidation("strong_password", validateStrongPassword)
		v.RegisterValidation("past_date", validatePastDate)
	}
}

// Validate CitizenID - phải đủ 12 chữ số
func validateCitizenID(fl validator.FieldLevel) bool {
	citizenID := fl.Field().String()
	matched, _ := regexp.MatchString(`^\d{12}$`, citizenID)
	return matched
}

// Validate số điện thoại Việt Nam
func validateVNPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^(0(3|5|7|8|9)[0-9]{8}|84(3|5|7|8|9)[0-9]{8})$`, phone)
	return matched
}

// Validate password mạnh - tối thiểu 8 ký tự, có chữ hoa và ký tự đặc biệt
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	hasUpper, _ := regexp.MatchString(`[A-Z]`, password)
	if !hasUpper {
		return false
	}

	hasSpecial, _ := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>]`, password)
	if !hasSpecial {
		return false
	}

	return true
}

// Validate ngày trong quá khứ - hỗ trợ string, *string, time.Time, *time.Time
func validatePastDate(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Handle pointer types
	if field.Kind() == reflect.Ptr {
		// Nếu là nil và có omitempty thì valid
		if field.IsNil() {
			return true
		}
		// Lấy giá trị thực từ pointer
		field = field.Elem()
	}

	var date time.Time
	var err error

	// Handle different types
	switch field.Kind() {
	case reflect.String:
		dateStr := field.String()
		if dateStr == "" {
			return true // Empty string is valid with omitempty
		}
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return false
		}
	case reflect.Struct:
		// Handle time.Time
		if field.Type() == reflect.TypeOf(time.Time{}) {
			date = field.Interface().(time.Time)
		} else {
			return false
		}
	default:
		return false
	}

	// Check if date is in the past
	return date.Before(time.Now())
}

// GetValidator - Trả về validator instance
func GetValidator() *validator.Validate {
	return validate
}

// ValidateStruct - Validate struct và trả về message lỗi
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationErrorsMap - Trả về map[field]error_message với form field names
func FormatValidationErrorsMap(err error) map[string]string {
	errorMap := make(map[string]string)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		errorMap["general"] = err.Error()
		return errorMap
	}

	// Map struct field name -> form field name (matching HTML input name attributes)
	fieldNameMap := map[string]string{
		"CitizenID":    "citizen_id",
		"Name":         "name",
		"Email":        "email",
		"Password":     "password",
		"Role":         "role",
		"Phone":        "phone",
		"DateOfBirth":  "date_of_birth",
		"Gender":       "gender",
		"DepartmentID": "department_id",
		"Address":      "address",
		"NewStatus":    "new_status",
		"Notes":        "notes",
		"Code":         "code",
	}

	// Human-readable label cho error messages
	fieldLabelMap := map[string]string{
		"CitizenID":    "Citizen ID",
		"Name":         "Name",
		"Email":        "Email",
		"Password":     "Password",
		"Role":         "Role",
		"Phone":        "Phone",
		"DateOfBirth":  "Date of Birth",
		"Gender":       "Gender",
		"DepartmentID": "Department",
		"Address":      "Address",
		"NewStatus":    "Status",
		"Notes":        "Notes",
		"Code":         "Code",
	}

	for _, e := range validationErrors {
		structField := e.Field()

		formField := toSnakeCase(structField)
		if mapped, ok := fieldNameMap[structField]; ok {
			formField = mapped
		}

		label := structField
		if mapped, ok := fieldLabelMap[structField]; ok {
			label = mapped
		}

		errorMap[formField] = formatFieldError(label, e)
	}

	return errorMap
}

// formatFieldError - Format error message cho một field
func formatFieldError(label string, e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return label + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return label + " must be at least " + e.Param() + " characters"
	case "max":
		return label + " must not exceed " + e.Param() + " characters"
	case "oneof":
		return label + " must be one of: " + strings.ReplaceAll(e.Param(), " ", ", ")
	case "citizen_id":
		return "Must be exactly 12 digits"
	case "strong_password":
		return "Must be at least 8 characters with uppercase letter and special character"
	case "vn_phone":
		return "Invalid Vietnamese phone number (e.g. 0901234567)"
	case "past_date":
		return "Must be a past date (YYYY-MM-DD)"
	default:
		return label + " is invalid"
	}
}

// getJSONFieldName - Lấy JSON tag name từ struct field
func getJSONFieldName(e validator.FieldError) string {
	field := e.Field()

	// Namespace có dạng "RegisterRequest.CitizenID"
	namespace := e.Namespace()
	parts := strings.Split(namespace, ".")
	if len(parts) > 1 {
		// Lấy tên struct và field
		structName := parts[0]
		fieldName := parts[1]

		// Reflect để lấy JSON tag
		if e.StructNamespace() != "" {
			// Tìm struct type từ namespace
			if structField, ok := getStructFieldByName(structName, fieldName); ok {
				if jsonTag := structField.Tag.Get("json"); jsonTag != "" {
					// Parse JSON tag (format: "field_name,omitempty")
					tagParts := strings.Split(jsonTag, ",")
					if len(tagParts) > 0 && tagParts[0] != "" {
						return tagParts[0]
					}
				}
			}
		}
	}

	// Fallback: convert to snake_case
	return toSnakeCase(field)
}

// getStructFieldByName - Helper để lấy struct field (simplified)
func getStructFieldByName(structName, fieldName string) (reflect.StructField, bool) {
	// This is a simplified version - in production, you might need type registry
	// For now, return empty to fallback to snake_case
	return reflect.StructField{}, false
}

// toSnakeCase - Convert CamelCase to snake_case
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
