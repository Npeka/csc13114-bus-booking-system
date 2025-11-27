-- Migration: Update seats table with enhanced seat map support
-- Description: Add row, column, floor, and price multiplier for seat configuration

-- Drop old seats table if exists (basic version)
DROP TABLE IF EXISTS seats CASCADE;

-- Create new seats table with full seat map support
CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bus_id UUID NOT NULL REFERENCES buses(id) ON DELETE CASCADE ON UPDATE CASCADE,
    seat_number VARCHAR(10) NOT NULL,
    row INTEGER NOT NULL CHECK (row >= 1),
    "column" INTEGER NOT NULL CHECK ("column" >= 1),
    seat_type VARCHAR(20) NOT NULL CHECK (seat_type IN ('standard', 'vip', 'sleeper')),
    price_multiplier DECIMAL(3,2) NOT NULL DEFAULT 1.0 CHECK (price_multiplier >= 0.5 AND price_multiplier <= 5.0),
    is_available BOOLEAN NOT NULL DEFAULT true,
    floor INTEGER NOT NULL DEFAULT 1 CHECK (floor IN (1, 2)),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT unique_bus_seat_number UNIQUE (bus_id, seat_number),
    CONSTRAINT unique_bus_position UNIQUE (bus_id, row, "column", floor)
);

-- Indexes for performance
CREATE INDEX idx_seats_bus ON seats(bus_id);
CREATE INDEX idx_seats_deleted ON seats(deleted_at);

-- Comments
COMMENT ON TABLE seats IS 'Seat configuration for buses with visual layout support';
COMMENT ON COLUMN seats.seat_number IS 'Seat identifier (e.g., A1, B2)';
COMMENT ON COLUMN seats.row IS 'Row number in seat layout';
COMMENT ON COLUMN seats."column" IS 'Column number in seat layout';
COMMENT ON COLUMN seats.seat_type IS 'Type of seat: standard, vip, or sleeper';
COMMENT ON COLUMN seats.price_multiplier IS 'Price multiplier for this seat type (1.0 = base price)';
COMMENT ON COLUMN seats.floor IS 'Floor number for double-decker buses (1 or 2)';
