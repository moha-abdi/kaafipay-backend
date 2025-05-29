-- Drop existing tables if they exist
DROP TABLE IF EXISTS linked_account_syncs CASCADE;
DROP TABLE IF EXISTS linked_accounts CASCADE;
DROP TABLE IF EXISTS user_accounts CASCADE;
DROP TABLE IF EXISTS account_providers CASCADE;

-- Drop existing types if they exist
DROP TYPE IF EXISTS account_provider CASCADE;

-- Create account provider enum type
CREATE TYPE account_provider AS ENUM (
    'ZAAD',
    'EDAHAB',
    'SAHAL',
    'EVCPLUS',
    'SOMNET',
    'SOLTELCO'
);

-- Table for storing linked accounts
CREATE TABLE linked_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    provider account_provider NOT NULL,
    account_id VARCHAR(255) NOT NULL,
    account_number VARCHAR(255) NOT NULL,
    account_title VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    currency_code VARCHAR(10) NOT NULL,
    currency_name VARCHAR(50) NOT NULL,
    currency_symbol VARCHAR(10) NOT NULL,
    is_default_account BOOLEAN DEFAULT false,
    -- Provider specific authentication details (encrypted)
    provider_username VARCHAR(255) NOT NULL,
    provider_password VARCHAR(255) NOT NULL,
    device_id VARCHAR(255) NOT NULL,
    -- Additional provider-specific details
    customer_id VARCHAR(255),
    subscription_id VARCHAR(255),
    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_sync_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    
    -- Constraints
    UNIQUE(user_id, provider, account_number)
);

-- Table for storing account sync history
CREATE TABLE linked_account_syncs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    linked_account_id UUID NOT NULL REFERENCES linked_accounts(id),
    sync_status VARCHAR(50) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_linked_accounts_user_id ON linked_accounts(user_id);
CREATE INDEX idx_linked_accounts_provider ON linked_accounts(provider);
CREATE INDEX idx_linked_account_syncs_account_id ON linked_account_syncs(linked_account_id);

-- Trigger to update updated_at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON linked_accounts
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
