package main

import (
	"fmt"
	"regexp"
	"strings"

	validx "github.com/askari/gpm/valid"
)

// RegisterInput یک DTO نمونه برای ثبت‌نام کاربر است.
type RegisterInput struct {
	Email        string
	Password     string
	Age          int
	Roles        []string
	Website      string
	ReferralCode string
}

var (
	corporateEmailRx = regexp.MustCompile(`@company\.com$`)
	hasDigitRx       = regexp.MustCompile(`[0-9]`)
	hasUpperRx       = regexp.MustCompile(`[A-Z]`)
	referralCodeRx   = regexp.MustCompile(`^[A-Z]{3}-[A-Z0-9]{5}$`)
)

// CorporateEmail: ایمیل باید دامنه سازمانی داشته باشد.
func CorporateEmail() validx.Rule[string] {
	return func(v string) (bool, *validx.RuleError) {
		if !corporateEmailRx.MatchString(v) {
			return false, &validx.RuleError{
				Code:    "corporate_email",
				Message: "email must end with @company.com",
			}
		}
		return true, nil
	}
}

// StrongPassword: حداقل 8 کاراکتر + یک عدد + یک حرف بزرگ.
func StrongPassword() validx.Rule[string] {
	return validx.All(
		validx.MinLen(8),
		func(v string) (bool, *validx.RuleError) {
			if !hasDigitRx.MatchString(v) {
				return false, &validx.RuleError{
					Code:    "password_digit",
					Message: "password must contain at least one digit",
				}
			}
			return true, nil
		},
		func(v string) (bool, *validx.RuleError) {
			if !hasUpperRx.MatchString(v) {
				return false, &validx.RuleError{
					Code:    "password_upper",
					Message: "password must contain at least one uppercase letter",
				}
			}
			return true, nil
		},
	)
}

// AllowedRoles: همه roleها باید در لیست مجاز باشند.
func AllowedRoles() validx.Rule[[]string] {
	allowed := map[string]struct{}{
		"admin":   {},
		"editor":  {},
		"support": {},
		"user":    {},
	}

	return func(v []string) (bool, *validx.RuleError) {
		for _, role := range v {
			if _, ok := allowed[role]; !ok {
				return false, &validx.RuleError{
					Code:    "role_invalid",
					Message: "roles contain unsupported value: " + role,
				}
			}
		}
		return true, nil
	}
}

// RegisterSchema نمونه کامل تعریف schema با Ruleهای آماده و سفارشی.
func RegisterSchema() *validx.Schema[RegisterInput] {
	s := validx.New[RegisterInput]().Schema()

	validx.Field(s, "email", func(i *RegisterInput) string { return i.Email }).
		Use(validx.RequiredString()).
		Use(validx.Email()).
		Use(CorporateEmail())

	validx.Field(s, "password", func(i *RegisterInput) string { return i.Password }).
		Use(validx.RequiredString()).
		Use(StrongPassword())

	validx.Field(s, "age", func(i *RegisterInput) int { return i.Age }).
		Use(validx.Between(18, 70))

	validx.Field(s, "roles", func(i *RegisterInput) []string { return i.Roles }).
		Use(validx.MinItems[string](1)).
		Use(validx.Unique[string]()).
		Use(AllowedRoles())

	// مقدار خالی قابل قبول است، ولی اگر ست شود باید URL معتبر باشد.
	validx.Field(s, "website", func(i *RegisterInput) string { return i.Website }).
		Use(validx.Optional(validx.URL()))

	// مقدار خالی قابل قبول است، ولی اگر ست شود باید الگوی کد معرف را رعایت کند.
	validx.Field(s, "referral_code", func(i *RegisterInput) string { return i.ReferralCode }).
		Use(validx.Optional(validx.Pattern(referralCodeRx)))

	return s
}

func printValidationResult(title string, schema *validx.Schema[RegisterInput], input RegisterInput) {
	fmt.Println(strings.Repeat("=", 72))
	fmt.Println(title)

	err := schema.Validate(&input)
	if err == nil {
		fmt.Println("Validation: PASSED")
		return
	}

	fmt.Println("Validation: FAILED")
	for _, e := range err.All() {
		fmt.Printf("- field=%s code=%s message=%s\n", e.Field, e.Code, e.Message)
	}
}

func main() {
	schema := RegisterSchema()

	invalidInput := RegisterInput{
		Email:        "ali@gmail.com",
		Password:     "weakpass",
		Age:          16,
		Roles:        []string{"admin", "admin", "owner"},
		Website:      "not-a-url",
		ReferralCode: "ab-12",
	}

	validInput := RegisterInput{
		Email:        "ali@company.com",
		Password:     "StrongPass1",
		Age:          30,
		Roles:        []string{"admin", "support"},
		Website:      "https://company.com/docs",
		ReferralCode: "ABC-9X2Q1",
	}

	printValidationResult("Scenario 1: Invalid Input", schema, invalidInput)
	printValidationResult("Scenario 2: Valid Input", schema, validInput)
}
