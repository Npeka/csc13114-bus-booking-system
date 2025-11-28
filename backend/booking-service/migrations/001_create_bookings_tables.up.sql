-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_reference VARCHAR(10) UNIQUE NOT NULL,
    trip_id UUID NOT NULL,
    user_id UUID,
    guest_email VARCHAR(255),
    guest_phone VARCHAR(20),
    guest_name VARCHAR(255),
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(50),
    payment_id VARCHAR(255),
    expires_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    cancellation_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Create indexes for bookings
CREATE INDEX IF NOT EXISTS idx_bookings_reference ON bookings(booking_reference);
CREATE INDEX IF NOT EXISTS idx_bookings_trip ON bookings(trip_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_guest_email ON bookings(guest_email);
CREATE INDEX IF NOT EXISTS idx_bookings_expires ON bookings(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_bookings_deleted ON bookings(deleted_at);

-- Create passengers table
CREATE TABLE IF NOT EXISTS passengers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    seat_id UUID NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    id_number VARCHAR(50),
    phone_number VARCHAR(20),
    seat_number VARCHAR(10) NOT NULL,
    seat_type VARCHAR(20) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for passengers
CREATE INDEX IF NOT EXISTS idx_passengers_booking ON passengers(booking_id);
CREATE INDEX IF NOT EXISTS idx_passengers_seat ON passengers(seat_id);

-- Create seat_locks table
CREATE TABLE IF NOT EXISTS seat_locks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    locked_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    UNIQUE(trip_id, seat_id)
);

-- Create indexes for seat_locks
CREATE INDEX IF NOT EXISTS idx_seat_locks_trip_seat ON seat_locks(trip_id, seat_id);
CREATE INDEX IF NOT EXISTS idx_seat_locks_session ON seat_locks(session_id);
CREATE INDEX IF NOT EXISTS idx_seat_locks_expires ON seat_locks(expires_at);
