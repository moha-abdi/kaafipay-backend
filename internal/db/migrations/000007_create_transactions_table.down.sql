DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP INDEX IF EXISTS idx_transactions_created_at;
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_receiver;
DROP INDEX IF EXISTS idx_transactions_sender;
DROP TABLE IF EXISTS transactions; 