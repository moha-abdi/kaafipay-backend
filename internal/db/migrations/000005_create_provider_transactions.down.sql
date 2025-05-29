DROP TRIGGER IF EXISTS update_provider_transactions_updated_at ON provider_transactions;
DROP INDEX IF EXISTS idx_provider_transactions_status;
DROP INDEX IF EXISTS idx_provider_transactions_date;
DROP INDEX IF EXISTS idx_provider_transactions_category;
DROP INDEX IF EXISTS idx_provider_transactions_linked_account;
DROP TABLE IF EXISTS provider_transactions; 