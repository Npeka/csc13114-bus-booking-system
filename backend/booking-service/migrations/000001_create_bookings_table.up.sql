-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    -- Standard fields
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Business fields
    booking_reference VARCHAR(20) UNIQUE NOT NULL,
    trip_id UUID NOT NULL,
    user_id UUID NOT NULL,
    
    -- Pricing
    total_amount DECIMAL(10,2) NOT NULL,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    
    -- Payment info - PayOS integration handled by Payment Service
    payment_order_id VARCHAR(255),
    
    -- Timestamps
    expires_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    
    -- Optional
    cancellation_reason TEXT,
    notes TEXT,
    
    -- Constraints
    CONSTRAINT chk_booking_status CHECK (status IN ('pending', 'confirmed', 'cancelled', 'expired')),
    CONSTRAINT chk_payment_status CHECK (payment_status IN ('pending', 'paid', 'refunded', 'failed'))
);

-- Create indexes
CREATE INDEX idx_bookings_trip_id ON bookings(trip_id);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_payment_status ON bookings(payment_status);
CREATE INDEX idx_bookings_booking_reference ON bookings(booking_reference);
CREATE INDEX idx_bookings_expires_at ON bookings(expires_at);
CREATE INDEX idx_bookings_deleted_at ON bookings(deleted_at);

-- Add updated_at trigger
CREATE TRIGGER update_bookings_updated_at BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
