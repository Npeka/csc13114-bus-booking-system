-- Create seats table
CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bus_id UUID NOT NULL REFERENCES buses(id) ON DELETE CASCADE,
    seat_number VARCHAR(10) NOT NULL,
    seat_type VARCHAR(20) NOT NULL DEFAULT 'standard',
    floor INTEGER NOT NULL DEFAULT 1,
    "row" INTEGER NOT NULL,
    "column" INTEGER NOT NULL,
    price_multiplier DECIMAL(3,2) NOT NULL DEFAULT 1.0,
    is_available BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT seats_unique_bus_seat UNIQUE(bus_id, seat_number),
    CONSTRAINT seats_seat_type_check CHECK (seat_type IN ('standard', 'vip', 'sleeper')),
    CONSTRAINT seats_floor_check CHECK (floor IN (1, 2)),
    CONSTRAINT seats_row_check CHECK ("row" > 0 AND "row" <= 20),
    CONSTRAINT seats_column_check CHECK ("column" > 0 AND "column" <= 5),
    CONSTRAINT seats_price_multiplier_check CHECK (price_multiplier >= 0.5 AND price_multiplier <= 5.0)
);

-- Indexes for seats
CREATE INDEX idx_seats_bus_id ON seats(bus_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_seats_seat_number ON seats(seat_number);
CREATE INDEX idx_seats_is_available ON seats(is_available) WHERE deleted_at IS NULL;
CREATE INDEX idx_seats_bus_available ON seats(bus_id, is_available) WHERE deleted_at IS NULL;
CREATE INDEX idx_seats_deleted_at ON seats(deleted_at);

-- Comments
COMMENT ON TABLE seats IS 'Individual seats on buses';
COMMENT ON COLUMN seats.seat_type IS 'Type of seat: standard, vip, sleeper';
COMMENT ON COLUMN seats.price_multiplier IS 'Price multiplier for this seat (0.5 to 5.0)';
COMMENT ON COLUMN seats.floor IS 'Floor number for double-decker buses (1 or 2)';
