package validx

import "testing"

func TestAnyPassesWhenOneRuleMatches(t *testing.T) {
	rule := Any[string](
		MinLen(10),
		RequiredString(),
	)

	ok, err := rule("x")
	if !ok {
		t.Fatalf("expected Any to pass when one rule matches, got err=%v", err)
	}
}

func TestAnyReturnsLastErrorWhenAllFail(t *testing.T) {
	rule := Any[string](
		MinLen(3),
		OneOf("x"),
	)

	ok, err := rule("ab")
	if ok {
		t.Fatal("expected Any to fail when all rules fail")
	}
	if err == nil {
		t.Fatal("expected Any to return a non-nil error")
	}
	if err.Code != "one_of" {
		t.Fatalf("expected last error code one_of, got %q", err.Code)
	}
}

func TestAnyWithNoRulesPasses(t *testing.T) {
	rule := Any[string]()

	ok, err := rule("anything")
	if !ok || err != nil {
		t.Fatalf("expected Any with no rules to pass, got ok=%v err=%v", ok, err)
	}
}

func TestPatternNilRegexFailsGracefully(t *testing.T) {
	rule := Pattern(nil)

	ok, err := rule("abc")
	if ok {
		t.Fatal("expected Pattern(nil) to fail")
	}
	if err == nil {
		t.Fatal("expected non-nil RuleError for Pattern(nil)")
	}
	if err.Code != "pattern" {
		t.Fatalf("expected error code pattern, got %q", err.Code)
	}
}
