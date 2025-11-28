-- Remove PayOS related columns from transactions table
DROP INDEX IF EXISTS idx_transactions_status;
DROP INDEX IF EXISTS idx_transactions_order_code;

ALTER TABLE transactions 
DROP COLUMN IF EXISTS transaction_time,
DROP COLUMN IF EXISTS reference,
DROP COLUMN IF EXISTS qr_code,
DROP COLUMN IF EXISTS checkout_url,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS payment_link_id,
DROP COLUMN IF EXISTS order_code;
