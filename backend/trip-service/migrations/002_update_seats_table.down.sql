-- Migration: Rollback seats table update

DROP INDEX IF EXISTS idx_seats_deleted;
DROP INDEX IF EXISTS idx_seats_bus;

DROP TABLE IF EXISTS seats CASCADE;

-- Recreate basic seats table (if needed for rollback)
CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bus_id UUID NOT NULL REFERENCES buses(id) ON DELETE CASCADE,
    seat_code VARCHAR(10) NOT NULL,
    seat_type VARCHAR(50) NOT NULL DEFAULT 'standard',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
