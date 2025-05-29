CREATE TABLE transaction_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(7),
    parent_id UUID REFERENCES transaction_categories(id),
    is_system BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE merchant_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    merchant_name VARCHAR(255) NOT NULL,
    category_id UUID REFERENCES transaction_categories(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(merchant_name)
); 