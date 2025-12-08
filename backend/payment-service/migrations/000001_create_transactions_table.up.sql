-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Core transaction fields
    booking_id UUID NOT NULL,
    user_id UUID NOT NULL,
    amount INTEGER NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'VND',
    payment_method VARCHAR(50) NOT NULL DEFAULT 'PAYOS',
    
    -- PayOS integration fields
    order_code BIGINT UNIQUE,
    payment_link_id VARCHAR(255),
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    checkout_url TEXT,
    qr_code TEXT,
    reference VARCHAR(255),
    transaction_time BIGINT
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_transactions_booking_id ON transactions(booking_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_order_code ON transactions(order_code);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions(deleted_at);
