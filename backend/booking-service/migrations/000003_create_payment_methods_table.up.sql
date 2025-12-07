-- Create payment_methods table
CREATE TABLE IF NOT EXISTS payment_methods (
    -- Standard fields
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Business fields
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Create indexes
CREATE INDEX idx_payment_methods_code ON payment_methods(code);
CREATE INDEX idx_payment_methods_is_active ON payment_methods(is_active);
CREATE INDEX idx_payment_methods_deleted_at ON payment_methods(deleted_at);

-- Add updated_at trigger
CREATE TRIGGER update_payment_methods_updated_at BEFORE UPDATE ON payment_methods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default payment methods
INSERT INTO payment_methods (name, code, description, is_active) VALUES
    ('Tiền mặt', 'cash', 'Thanh toán bằng tiền mặt khi lên xe', TRUE),
    ('Chuyển khoản ngân hàng', 'bank_transfer', 'Chuyển khoản qua ngân hàng', TRUE),
    ('Ví điện tử MoMo', 'momo', 'Thanh toán qua ví MoMo', TRUE),
    ('Ví điện tử ZaloPay', 'zalopay', 'Thanh toán qua ví ZaloPay', TRUE),
    ('Thẻ tín dụng/ghi nợ', 'card', 'Thanh toán bằng thẻ Visa/Mastercard', TRUE);
