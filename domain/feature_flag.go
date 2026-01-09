package domain

import (
	"errors"
)

// RuleType represents the type of rule
type RuleType string

const (
	RuleTypePercentage RuleType = "percentage"
	RuleTypeUserID     RuleType = "user_id"
)

type RuleEvaluator interface {
	Type() RuleType;
	Validate() error;
}

type PercentageRule struct {
	Percentage int32;
}
func (r *PercentageRule) Type() RuleType {
	return RuleTypePercentage;
}
func (r *PercentageRule) Validate() error{
	if r.Percentage < 0 || r.Percentage > 100 {
		return errors.New("must be 0-100")
	}
	return nil;
}

type UserIDRule struct{
	UserIDs []string;
}
func (r *UserIDRule) Type() RuleType{
	return RuleTypeUserID
}
func (r *UserIDRule) Validate() error{
	if len(r.UserIDs) == 0 {
		return errors.New("must have at least one user ID")
	}
	return nil;
}

type Rule struct {
	Evaluator RuleEvaluator;
}

func (r *Rule) Validate() error {
	if r.Evaluator == nil {
		return errors.New("rule evaluator cannot be nil");
	}
	return r.Evaluator.Validate();
}

type FeatureFlag struct {
	Key     string;
	Enabled bool;
	Rules   []Rule;
	Version int64;
}

func (f *FeatureFlag) Validate() error {
	if f.Key == "" {
		return errors.New("key cannot be empty");
	}

	for _, rule := range f.Rules {
		if err := rule.Validate(); err != nil {
			return err;
		}
	}

	return nil;
}