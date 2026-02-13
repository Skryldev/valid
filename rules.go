package validx

import (
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

//
// 🔹 RuleError حرفه‌ای برای reporting
//

type RuleError struct {
	Code    string // machine-friendly code
	Message string // human-friendly message
}

func (e RuleError) Error() string {
	return e.Message
}

//
// 🔹 Rule Signature
//

type Rule[T any] func(value T) (bool, *RuleError)

func ok[T any]() (bool, *RuleError) { return true, nil }
func fail(code, msg string) (bool, *RuleError) {
	return false, &RuleError{Code: code, Message: msg}
}

//
// 🔹 STRING RULES (backend-relevant)
//

// RequiredString: مقدار غیرخالی و غیرspace-only
func RequiredString() Rule[string] {
	return func(v string) (bool, *RuleError) {
		if strings.TrimSpace(v) == "" {
			return fail("required", "value is required")
		}
		return ok[string]()
	}
}

// MinLen: طول حداقل string (unicode safe)
func MinLen(min int) Rule[string] {
	return func(v string) (bool, *RuleError) {
		if utf8.RuneCountInString(v) < min {
			return fail("min_length", "minimum length not satisfied")
		}
		return ok[string]()
	}
}

// MaxLen: طول حداکثر string
func MaxLen(max int) Rule[string] {
	return func(v string) (bool, *RuleError) {
		if utf8.RuneCountInString(v) > max {
			return fail("max_length", "maximum length exceeded")
		}
		return ok[string]()
	}
}

// Email: RFC-compliant
func Email() Rule[string] {
	return func(v string) (bool, *RuleError) {
		if v == "" {
			return ok[string]()
		}
		_, err := mail.ParseAddress(v)
		if err != nil {
			return fail("email", "invalid email format")
		}
		return ok[string]()
	}
}

// URL: absolute URL
func URL() Rule[string] {
	return func(v string) (bool, *RuleError) {
		if v == "" {
			return ok[string]()
		}
		u, err := url.ParseRequestURI(v)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fail("url", "invalid url")
		}
		return ok[string]()
	}
}

// Pattern: regex validation
func Pattern(rx *regexp.Regexp) Rule[string] {
	return func(v string) (bool, *RuleError) {
		if rx == nil {
			return fail("pattern", "pattern is not configured")
		}
		if !rx.MatchString(v) {
			return fail("pattern", "pattern mismatch")
		}
		return ok[string]()
	}
}

// OneOf: enum validation
func OneOf(options ...string) Rule[string] {
	set := make(map[string]struct{}, len(options))
	for _, o := range options {
		set[o] = struct{}{}
	}
	return func(v string) (bool, *RuleError) {
		if _, ok := set[v]; !ok {
			return fail("one_of", "value not allowed")
		}
		return ok[string]()
	}
}

//
// 🔹 NUMBER RULES
//

type Number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Min[T Number](min T) Rule[T] {
	return func(v T) (bool, *RuleError) {
		if v < min {
			return fail("min", "value below minimum")
		}
		return ok[T]()
	}
}

func Max[T Number](max T) Rule[T] {
	return func(v T) (bool, *RuleError) {
		if v > max {
			return fail("max", "value above maximum")
		}
		return ok[T]()
	}
}

func Between[T Number](min, max T) Rule[T] {
	return func(v T) (bool, *RuleError) {
		if v < min || v > max {
			return fail("between", "value out of range")
		}
		return ok[T]()
	}
}

//
// 🔹 SLICE / COLLECTION RULES
//

func MinItems[T any](min int) Rule[[]T] {
	return func(v []T) (bool, *RuleError) {
		if len(v) < min {
			return fail("min_items", "not enough items")
		}
		return ok[[]T]()
	}
}

func MaxItems[T any](max int) Rule[[]T] {
	return func(v []T) (bool, *RuleError) {
		if len(v) > max {
			return fail("max_items", "too many items")
		}
		return ok[[]T]()
	}
}

func Unique[T comparable]() Rule[[]T] {
	return func(v []T) (bool, *RuleError) {
		seen := make(map[T]struct{}, len(v))
		for _, item := range v {
			if _, ok := seen[item]; ok {
				return fail("unique", "duplicate values found")
			}
			seen[item] = struct{}{}
		}
		return ok[[]T]()
	}
}

//
// 🔹 COMPOSABLE RULES
//

// All: همه ruleها باید pass کنند
func All[T any](rules ...Rule[T]) Rule[T] {
	return func(v T) (bool, *RuleError) {
		for _, r := range rules {
			if ok, err := r(v); !ok {
				return false, err
			}
		}
		return ok[T]()
	}
}

// Any: حداقل یکی pass کند
func Any[T any](rules ...Rule[T]) Rule[T] {
	return func(v T) (bool, *RuleError) {
		if len(rules) == 0 {
			return ok[T]()
		}
		var lastErr *RuleError
		for _, r := range rules {
			if passed, err := r(v); passed {
				return ok[T]()
			} else if err != nil {
				lastErr = err
			}
		}
		if lastErr != nil {
			return false, lastErr
		}
		return fail("any", "none of the rules matched")
	}
}

// Optional: اگر مقدار zero بود validation skip شود
func Optional[T comparable](rule Rule[T]) Rule[T] {
	var zero T
	return func(v T) (bool, *RuleError) {
		if v == zero {
			return ok[T]()
		}
		return rule(v)
	}
}
