-- Drop indexes first
DROP INDEX IF EXISTS idx_linked_accounts_deleted_at;
DROP INDEX IF EXISTS idx_linked_account_syncs_deleted_at;

-- Drop columns
ALTER TABLE linked_accounts DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE linked_account_syncs DROP COLUMN IF EXISTS deleted_at;
