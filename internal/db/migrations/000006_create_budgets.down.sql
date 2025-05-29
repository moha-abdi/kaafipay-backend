DROP TRIGGER IF EXISTS update_budget_categories_updated_at ON budget_categories;
DROP TRIGGER IF EXISTS update_budgets_updated_at ON budgets;

DROP INDEX IF EXISTS idx_budget_categories_category;
DROP INDEX IF EXISTS idx_budget_categories_budget;
DROP INDEX IF EXISTS idx_budgets_user;

DROP TABLE IF EXISTS budget_categories;
DROP TABLE IF EXISTS budgets; 