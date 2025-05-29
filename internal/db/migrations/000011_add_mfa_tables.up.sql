-- Create MFA codes table
CREATE TABLE mfa_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(6) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create MFA tokens table
CREATE TABLE mfa_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token VARCHAR(64) NOT NULL UNIQUE,
    phone VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create WhatsApp sessions table
CREATE TABLE whatsapp_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id VARCHAR(100) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add indexes
CREATE INDEX idx_mfa_codes_phone ON mfa_codes(phone);
CREATE INDEX idx_mfa_codes_expires_at ON mfa_codes(expires_at);
CREATE INDEX idx_mfa_tokens_token ON mfa_tokens(token);
CREATE INDEX idx_mfa_tokens_expires_at ON mfa_tokens(expires_at);
CREATE INDEX idx_whatsapp_sessions_status ON whatsapp_sessions(status);

-- Add trigger for updated_at
CREATE TRIGGER update_whatsapp_sessions_updated_at
    BEFORE UPDATE ON whatsapp_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
