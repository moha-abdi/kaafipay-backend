DROP TRIGGER IF EXISTS update_payment_methods_updated_at ON payment_methods;
DROP INDEX IF EXISTS idx_payment_methods_status;
DROP INDEX IF EXISTS idx_payment_methods_type;
DROP INDEX IF EXISTS idx_payment_methods_user;
DROP TABLE IF EXISTS payment_methods; 