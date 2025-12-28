-- Create bank_accounts table
CREATE TABLE IF NOT EXISTS bank_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Core fields
    user_id UUID NOT NULL,
    bank_code VARCHAR(20) NOT NULL,
    account_number VARCHAR(50) NOT NULL,
    account_holder VARCHAR(100) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_bank_accounts_user_id ON bank_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_primary ON bank_accounts(is_primary);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_deleted_at ON bank_accounts(deleted_at);

-- Create unique constraint to ensure only one primary account per user
CREATE UNIQUE INDEX IF NOT EXISTS idx_bank_accounts_user_primary 
    ON bank_accounts(user_id) 
    WHERE is_primary = TRUE AND deleted_at IS NULL;

-- Add comment
COMMENT ON TABLE bank_accounts IS 'User bank accounts for refund purposes';
