-- Create seat_locks table for temporary seat reservations during booking process
CREATE TABLE IF NOT EXISTS seat_locks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    session_id VARCHAR(255) NOT NULL,
    locked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Index for checking if a specific seat is locked
CREATE INDEX IF NOT EXISTS idx_seat_locks_trip_seat ON seat_locks(trip_id, seat_id) WHERE deleted_at IS NULL;

-- Index for efficient expiry queries
CREATE INDEX IF NOT EXISTS idx_seat_locks_expires_at ON seat_locks(expires_at) WHERE deleted_at IS NULL;

-- Index for session-based unlocking
CREATE INDEX IF NOT EXISTS idx_seat_locks_session_id ON seat_locks(session_id) WHERE deleted_at IS NULL;

-- Composite index for active locked seats check
CREATE INDEX IF NOT EXISTS idx_seat_locks_active ON seat_locks(trip_id, seat_id, expires_at) WHERE deleted_at IS NULL;
