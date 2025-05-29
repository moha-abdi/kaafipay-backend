DROP INDEX IF EXISTS idx_account_verification_sessions_provider;
DROP INDEX IF EXISTS idx_account_verification_sessions_user;
DROP INDEX IF EXISTS idx_account_auth_credentials_account;
DROP INDEX IF EXISTS idx_linked_accounts_provider;
DROP INDEX IF EXISTS idx_linked_accounts_user;

DROP TRIGGER IF EXISTS update_account_verification_sessions_updated_at ON account_verification_sessions;
DROP TRIGGER IF EXISTS update_account_auth_credentials_updated_at ON account_auth_credentials;
DROP TRIGGER IF EXISTS update_linked_accounts_updated_at ON linked_accounts;
DROP TRIGGER IF EXISTS update_service_providers_updated_at ON service_providers;

DROP TABLE IF EXISTS account_verification_sessions;
DROP TABLE IF EXISTS account_auth_credentials;
DROP TABLE IF EXISTS linked_accounts;
DROP TABLE IF EXISTS service_providers; 