-- Remove transaction type and refund fields from transactions table
ALTER TABLE transactions
DROP CONSTRAINT IF EXISTS chk_refund_status,
DROP CONSTRAINT IF EXISTS chk_transaction_type;

DROP INDEX IF EXISTS idx_transactions_refund_status;
DROP INDEX IF EXISTS idx_transactions_transaction_type;

ALTER TABLE transactions
DROP COLUMN IF EXISTS refund_amount,
DROP COLUMN IF EXISTS refund_status,
DROP COLUMN IF EXISTS transaction_type;
