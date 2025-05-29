-- Drop indexes
DROP INDEX IF EXISTS idx_whatsapp_sessions_status;
DROP INDEX IF EXISTS idx_mfa_tokens_expires_at;
DROP INDEX IF EXISTS idx_mfa_tokens_token;
DROP INDEX IF EXISTS idx_mfa_codes_expires_at;
DROP INDEX IF EXISTS idx_mfa_codes_phone;

-- Drop trigger
DROP TRIGGER IF EXISTS update_whatsapp_sessions_updated_at ON whatsapp_sessions;

-- Drop tables
DROP TABLE IF EXISTS whatsapp_sessions;
DROP TABLE IF EXISTS mfa_tokens;
DROP TABLE IF EXISTS mfa_codes;
