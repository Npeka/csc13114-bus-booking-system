-- Add PayOS related columns to transactions table
ALTER TABLE transactions 
ADD COLUMN IF NOT EXISTS order_code BIGINT UNIQUE,
ADD COLUMN IF NOT EXISTS payment_link_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'PENDING',
ADD COLUMN IF NOT EXISTS checkout_url TEXT,
ADD COLUMN IF NOT EXISTS qr_code TEXT,
ADD COLUMN IF NOT EXISTS reference VARCHAR(255),
ADD COLUMN IF NOT EXISTS transaction_time BIGINT;

-- Create index on order_code for faster lookups
CREATE INDEX IF NOT EXISTS idx_transactions_order_code ON transactions(order_code);

-- Create index on status for filtering
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
