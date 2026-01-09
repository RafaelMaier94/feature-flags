package domain

import (
	"testing"
)

func TestPercentageRule_Validate(t *testing.T) {
	// Test valid percentage
	t.Run("valid percentage", func(t *testing.T) {
		rule := &PercentageRule{Percentage: 50}
		err := rule.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test percentage too low
	t.Run("percentage below 0", func(t *testing.T) {
		rule := &PercentageRule{Percentage: -1}
		err := rule.Validate()
		if err == nil {
			t.Error("expected error for negative percentage, got nil")
		}
	})

	// Test percentage too high
	t.Run("percentage above 100", func(t *testing.T) {
		rule := &PercentageRule{Percentage: 101}
		err := rule.Validate()
		if err == nil {
			t.Error("expected error for percentage > 100, got nil")
		}
	})

	// Test edge cases
	t.Run("percentage at 0", func(t *testing.T) {
		rule := &PercentageRule{Percentage: 0}
		err := rule.Validate()
		if err != nil {
			t.Errorf("0 should be valid, got error: %v", err)
		}
	})

	t.Run("percentage at 100", func(t *testing.T) {
		rule := &PercentageRule{Percentage: 100}
		err := rule.Validate()
		if err != nil {
			t.Errorf("100 should be valid, got error: %v", err)
		}
	})
}

func TestUserIDRule_Validate(t *testing.T) {
	// Test valid user IDs
	t.Run("valid user IDs", func(t *testing.T) {
		rule := &UserIDRule{UserIDs: []string{"user1", "user2"}}
		err := rule.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test empty user IDs
	t.Run("empty user IDs", func(t *testing.T) {
		rule := &UserIDRule{UserIDs: []string{}}
		err := rule.Validate()
		if err == nil {
			t.Error("expected error for empty user IDs, got nil")
		}
	})

	// Test nil user IDs
	t.Run("nil user IDs", func(t *testing.T) {
		rule := &UserIDRule{UserIDs: nil}
		err := rule.Validate()
		if err == nil {
			t.Error("expected error for nil user IDs, got nil")
		}
	})
}

func TestRule_Validate(t *testing.T) {
	// Test valid rule with percentage
	t.Run("valid percentage rule", func(t *testing.T) {
		rule := Rule{
			Evaluator: &PercentageRule{Percentage: 50},
		}
		err := rule.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test valid rule with user IDs
	t.Run("valid user ID rule", func(t *testing.T) {
		rule := Rule{
			Evaluator: &UserIDRule{UserIDs: []string{"user1"}},
		}
		err := rule.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test nil evaluator
	t.Run("nil evaluator", func(t *testing.T) {
		rule := Rule{Evaluator: nil}
		err := rule.Validate()
		if err == nil {
			t.Error("expected error for nil evaluator, got nil")
		}
	})

	// Test invalid percentage rule
	t.Run("invalid percentage rule", func(t *testing.T) {
		rule := Rule{
			Evaluator: &PercentageRule{Percentage: 150},
		}
		err := rule.Validate()
		if err == nil {
			t.Error("expected error for invalid percentage, got nil")
		}
	})
}

func TestFeatureFlag_Validate(t *testing.T) {
	// Test valid feature flag
	t.Run("valid feature flag", func(t *testing.T) {
		flag := &FeatureFlag{
			Key:     "test-feature",
			Enabled: true,
			Rules: []Rule{
				{Evaluator: &PercentageRule{Percentage: 50}},
			},
			Version: 1,
		}
		err := flag.Validate()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test empty key
	t.Run("empty key", func(t *testing.T) {
		flag := &FeatureFlag{
			Key:     "",
			Enabled: true,
			Version: 1,
		}
		err := flag.Validate()
		if err == nil {
			t.Error("expected error for empty key, got nil")
		}
	})

	// Test invalid rule
	t.Run("invalid rule", func(t *testing.T) {
		flag := &FeatureFlag{
			Key:     "test",
			Enabled: true,
			Rules: []Rule{
				{Evaluator: &PercentageRule{Percentage: 150}},
			},
			Version: 1,
		}
		err := flag.Validate()
		if err == nil {
			t.Error("expected error for invalid rule, got nil")
		}
	})

	// Test no rules (should be valid)
	t.Run("no rules", func(t *testing.T) {
		flag := &FeatureFlag{
			Key:     "test",
			Enabled: true,
			Rules:   []Rule{},
			Version: 1,
		}
		err := flag.Validate()
		if err != nil {
			t.Errorf("no rules should be valid, got error: %v", err)
		}
	})
}

func TestRuleType(t *testing.T) {
	t.Run("percentage rule type", func(t *testing.T) {
		rule := &PercentageRule{Percentage: 50}
		if rule.Type() != RuleTypePercentage {
			t.Errorf("expected %s, got %s", RuleTypePercentage, rule.Type())
		}
	})

	t.Run("user ID rule type", func(t *testing.T) {
		rule := &UserIDRule{UserIDs: []string{"user1"}}
		if rule.Type() != RuleTypeUserID {
			t.Errorf("expected %s, got %s", RuleTypeUserID, rule.Type())
		}
	})
}
