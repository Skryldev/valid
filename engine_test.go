package validx

import "testing"

type testUser struct {
	Email string
	Age   int
	Roles []string
}

func hasFieldCode(errs []Error, field, code string) bool {
	for _, e := range errs {
		if e.Field == field && e.Code == code {
			return true
		}
	}
	return false
}

func TestValidateAppliesRulesAcrossTypes(t *testing.T) {
	s := New[testUser]().Schema()

	Field(s, "email", func(u *testUser) string { return u.Email }).
		Use(RequiredString()).
		Use(Email())
	Field(s, "age", func(u *testUser) int { return u.Age }).
		Use(Min(18))
	Field(s, "roles", func(u *testUser) []string { return u.Roles }).
		Use(Unique[string]())

	err := s.Validate(&testUser{
		Email: "",
		Age:   15,
		Roles: []string{"admin", "admin"},
	})
	if err == nil {
		t.Fatal("expected validation errors, got nil")
	}

	all := err.All()
	if !hasFieldCode(all, "email", "required") {
		t.Fatalf("expected required error for email, got %+v", all)
	}
	if !hasFieldCode(all, "age", "min") {
		t.Fatalf("expected min error for age, got %+v", all)
	}
	if !hasFieldCode(all, "roles", "unique") {
		t.Fatalf("expected unique error for roles, got %+v", all)
	}
}

func TestValidateReturnsNilForValidInput(t *testing.T) {
	s := New[testUser]().Schema()

	Field(s, "email", func(u *testUser) string { return u.Email }).
		Use(RequiredString()).
		Use(Email())
	Field(s, "age", func(u *testUser) int { return u.Age }).
		Use(Between(18, 65))
	Field(s, "roles", func(u *testUser) []string { return u.Roles }).
		Use(MinItems[string](1)).
		Use(Unique[string]())

	err := s.Validate(&testUser{
		Email: "user@company.com",
		Age:   30,
		Roles: []string{"admin"},
	})
	if err != nil {
		t.Fatalf("expected nil error, got %+v", err.All())
	}
}

func TestValidateHandlesNilRule(t *testing.T) {
	s := New[testUser]().Schema()

	var nilRule Rule[string]
	Field(s, "email", func(u *testUser) string { return u.Email }).
		Use(nilRule)

	err := s.Validate(&testUser{Email: "user@company.com"})
	if err == nil {
		t.Fatal("expected invalid_rule error, got nil")
	}

	if !hasFieldCode(err.All(), "email", "invalid_rule") {
		t.Fatalf("expected invalid_rule code, got %+v", err.All())
	}
}
