CREATE TABLE provider_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    linked_account_id UUID REFERENCES linked_accounts(id),
    provider_transaction_id VARCHAR(255),
    transaction_type VARCHAR(50) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    description TEXT,
    merchant_name VARCHAR(255),
    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
    balance_after DECIMAL(12,2),
    category_id UUID REFERENCES transaction_categories(id),
    provider_metadata JSONB,
    sync_status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(linked_account_id, provider_transaction_id)
);

-- Add indexes for common queries
CREATE INDEX idx_provider_transactions_linked_account ON provider_transactions(linked_account_id);
CREATE INDEX idx_provider_transactions_category ON provider_transactions(category_id);
CREATE INDEX idx_provider_transactions_date ON provider_transactions(transaction_date);
CREATE INDEX idx_provider_transactions_status ON provider_transactions(sync_status);

-- Add trigger for updated_at
CREATE TRIGGER update_provider_transactions_updated_at
    BEFORE UPDATE ON provider_transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 