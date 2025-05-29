-- Drop tables in reverse order
DROP TABLE IF EXISTS linked_account_syncs CASCADE;
DROP TABLE IF EXISTS linked_accounts CASCADE;

-- Drop the trigger function
DROP FUNCTION IF EXISTS trigger_set_timestamp CASCADE;

-- Drop the enum type
DROP TYPE IF EXISTS account_provider CASCADE;
