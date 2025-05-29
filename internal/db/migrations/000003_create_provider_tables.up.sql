CREATE TABLE service_providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    provider_type VARCHAR(50) NOT NULL,
    icon_url VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    auth_flow_type VARCHAR(50) NOT NULL,
    country_code VARCHAR(2),
    currency_code VARCHAR(3),
    api_config JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE linked_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    provider_id UUID REFERENCES service_providers(id),
    provider_account_id VARCHAR(255),
    account_name VARCHAR(255),
    phone_number VARCHAR(50),
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    auth_status VARCHAR(50) NOT NULL,
    last_sync_at TIMESTAMP WITH TIME ZONE,
    current_balance DECIMAL(12,2),
    provider_metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, provider_id, provider_account_id)
);

CREATE TABLE account_auth_credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    linked_account_id UUID REFERENCES linked_accounts(id),
    auth_type VARCHAR(50) NOT NULL,
    credential_value TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE account_verification_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    provider_id UUID REFERENCES service_providers(id),
    session_token VARCHAR(255) UNIQUE NOT NULL,
    verification_step VARCHAR(50) NOT NULL,
    phone_number VARCHAR(50),
    verification_data JSONB,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add triggers for updated_at
CREATE TRIGGER update_service_providers_updated_at
    BEFORE UPDATE ON service_providers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_linked_accounts_updated_at
    BEFORE UPDATE ON linked_accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_account_auth_credentials_updated_at
    BEFORE UPDATE ON account_auth_credentials
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_account_verification_sessions_updated_at
    BEFORE UPDATE ON account_verification_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add indexes
CREATE INDEX idx_linked_accounts_user ON linked_accounts(user_id);
CREATE INDEX idx_linked_accounts_provider ON linked_accounts(provider_id);
CREATE INDEX idx_account_auth_credentials_account ON account_auth_credentials(linked_account_id);
CREATE INDEX idx_account_verification_sessions_user ON account_verification_sessions(user_id);
CREATE INDEX idx_account_verification_sessions_provider ON account_verification_sessions(provider_id); 