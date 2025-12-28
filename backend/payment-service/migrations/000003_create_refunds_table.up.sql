-- Create refunds table (separate from transactions)
CREATE TABLE IF NOT EXISTS refunds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL,
    transaction_id UUID NOT NULL,
    user_id UUID NOT NULL,
    refund_amount INT NOT NULL CHECK (refund_amount > 0),
    refund_status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (refund_status IN ('PENDING', 'PROCESSING', 'COMPLETED', 'REJECTED')),
    refund_reason TEXT NOT NULL,
    rejected_reason TEXT,
    processed_by UUID,
    processed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Foreign key constraint
    CONSTRAINT fk_refunds_transaction FOREIGN KEY (transaction_id) 
        REFERENCES transactions(id) ON DELETE CASCADE
);

-- ===== INDEXES FOR OPTIMIZATION =====

-- Primary lookup indexes
CREATE INDEX idx_refunds_booking_id ON refunds(booking_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_refunds_transaction_id ON refunds(transaction_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_refunds_user_id ON refunds(user_id) WHERE deleted_at IS NULL;

-- Status filter (most common query)
CREATE INDEX idx_refunds_status ON refunds(refund_status) WHERE deleted_at IS NULL;

-- Pending refunds (admin dashboard)
CREATE INDEX idx_refunds_pending ON refunds(refund_status, created_at DESC) 
    WHERE refund_status = 'PENDING' AND deleted_at IS NULL;

-- Date range queries (for reports/exports)
CREATE INDEX idx_refunds_created_at ON refunds(created_at DESC) WHERE deleted_at IS NULL;

-- Admin processed refunds lookup
CREATE INDEX idx_refunds_processed_by ON refunds(processed_by, processed_at DESC) 
    WHERE processed_by IS NOT NULL AND deleted_at IS NULL;

-- Composite index for common admin filters (status + date range)
CREATE INDEX idx_refunds_status_created ON refunds(refund_status, created_at DESC) 
    WHERE deleted_at IS NULL;

-- Unique constraint: one refund per booking
CREATE UNIQUE INDEX idx_refunds_unique_booking ON refunds(booking_id) 
    WHERE deleted_at IS NULL;

-- ===== COMMENTS =====
COMMENT ON TABLE refunds IS 'Refund requests for cancelled bookings';
COMMENT ON COLUMN refunds.refund_status IS 'PENDING | PROCESSING | COMPLETED | REJECTED';
COMMENT ON COLUMN refunds.refund_amount IS 'Amount to refund in VND (integer)';
COMMENT ON COLUMN refunds.refund_reason IS 'User provided reason for refund request';
COMMENT ON COLUMN refunds.rejected_reason IS 'Admin reason if rejected';
COMMENT ON COLUMN refunds.processed_by IS 'Admin user ID who approved/rejected';
COMMENT ON COLUMN refunds.processed_at IS 'Timestamp when admin processed the refund';
