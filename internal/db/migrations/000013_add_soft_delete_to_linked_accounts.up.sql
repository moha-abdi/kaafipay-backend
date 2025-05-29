-- Add deleted_at column to linked_accounts table
ALTER TABLE linked_accounts ADD COLUMN deleted_at TIMESTAMPTZ;

-- Create index for soft deletes
CREATE INDEX idx_linked_accounts_deleted_at ON linked_accounts(deleted_at);

-- Add deleted_at column to linked_account_syncs table for consistency
ALTER TABLE linked_account_syncs ADD COLUMN deleted_at TIMESTAMPTZ;
CREATE INDEX idx_linked_account_syncs_deleted_at ON linked_account_syncs(deleted_at);
