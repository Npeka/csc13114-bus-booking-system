-- Add transaction type and refund fields to transactions table
ALTER TABLE transactions 
ADD COLUMN IF NOT EXISTS transaction_type VARCHAR(10) NOT NULL DEFAULT 'IN',
ADD COLUMN IF NOT EXISTS refund_status VARCHAR(20),
ADD COLUMN IF NOT EXISTS refund_amount INTEGER;

-- Create index for transaction_type
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_type ON transactions(transaction_type);

-- Create index for refund_status
CREATE INDEX IF NOT EXISTS idx_transactions_refund_status ON transactions(refund_status) WHERE refund_status IS NOT NULL;

-- Add check constraints
ALTER TABLE transactions
ADD CONSTRAINT chk_transaction_type CHECK (transaction_type IN ('IN', 'OUT'));

ALTER TABLE transactions
ADD CONSTRAINT chk_refund_status CHECK (
    refund_status IS NULL OR 
    refund_status IN ('PENDING', 'PROCESSING', 'COMPLETED', 'FAILED', 'CANCELLED')
);
