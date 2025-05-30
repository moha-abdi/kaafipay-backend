package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ValidationError represents a validation error for model fields
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// BudgetRule represents a single rule for categorizing transactions
type BudgetRule struct {
	Type     string `json:"type"`     // 'description', 'merchant', 'amount'
	Operator string `json:"operator"` // 'contains', 'equals', 'greater', 'less'
	Value    string `json:"value"`    // The value to match against
}

// BudgetCategory represents a budget category with rules for auto-categorization
type BudgetCategory struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID       `json:"userId" gorm:"type:uuid;not null"`
	Name      string          `json:"name" gorm:"type:varchar(50);not null"`
	Icon      string          `json:"icon" gorm:"type:varchar(50);not null"`
	Budget    float64         `json:"budget" gorm:"type:decimal(10,2);not null"`
	Rules     []BudgetRule    `json:"rules" gorm:"-"`                                         // Handled via RulesJSON
	RulesJSON json.RawMessage `json:"-" gorm:"column:rules;type:jsonb;not null;default:'[]'"` // Actual DB column
	CreatedAt time.Time       `json:"createdAt" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time       `json:"updatedAt" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt  `json:"-" gorm:"index"`
}

// BeforeSave converts Rules to RulesJSON before saving
func (bc *BudgetCategory) BeforeSave(tx *gorm.DB) error {
	rulesJSON, err := json.Marshal(bc.Rules)
	if err != nil {
		return err
	}
	bc.RulesJSON = rulesJSON
	return nil
}

// AfterFind converts RulesJSON to Rules after fetching
func (bc *BudgetCategory) AfterFind(tx *gorm.DB) error {
	return json.Unmarshal(bc.RulesJSON, &bc.Rules)
}

// Validate checks if the budget category data is valid
func (bc *BudgetCategory) Validate() error {
	if bc.Name == "" || len(bc.Name) > 50 {
		return ValidationError{Field: "name", Message: "Name must be between 1 and 50 characters"}
	}
	if bc.Icon == "" || len(bc.Icon) > 50 {
		return ValidationError{Field: "icon", Message: "Icon must be between 1 and 50 characters"}
	}
	if bc.Budget <= 0 {
		return ValidationError{Field: "budget", Message: "Budget must be greater than 0"}
	}
	if len(bc.Rules) == 0 {
		return ValidationError{Field: "rules", Message: "At least one rule is required"}
	}
	if len(bc.Rules) > 10 {
		return ValidationError{Field: "rules", Message: "Maximum 10 rules allowed"}
	}

	validTypes := map[string]bool{"description": true, "merchant": true, "amount": true}
	validOperators := map[string]bool{"contains": true, "equals": true, "greater": true, "less": true}

	for i, rule := range bc.Rules {
		if !validTypes[rule.Type] {
			return ValidationError{Field: "rules", Message: "Invalid rule type at index " + string(rune(i))}
		}
		if !validOperators[rule.Operator] {
			return ValidationError{Field: "rules", Message: "Invalid rule operator at index " + string(rune(i))}
		}
		if rule.Value == "" || len(rule.Value) > 100 {
			return ValidationError{Field: "rules", Message: "Rule value must be between 1 and 100 characters at index " + string(rune(i))}
		}
	}

	return nil
}

// TableName specifies the table name for the model
func (BudgetCategory) TableName() string {
	return "budget_categories"
}
