-- Create buses table
CREATE TABLE IF NOT EXISTS buses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plate_number VARCHAR(20) NOT NULL UNIQUE,
    model VARCHAR(255) NOT NULL,
    bus_type VARCHAR(20) NOT NULL DEFAULT 'standard',
    seat_capacity INTEGER NOT NULL,
    amenities TEXT[] DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT buses_bus_type_check CHECK (bus_type IN ('standard', 'vip', 'sleeper', 'double_decker')),
    CONSTRAINT buses_seat_capacity_check CHECK (seat_capacity >= 1 AND seat_capacity <= 100)
);

-- Indexes for buses
CREATE INDEX idx_buses_plate_number ON buses(plate_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_buses_bus_type ON buses(bus_type) WHERE is_active = true;
CREATE INDEX idx_buses_is_active ON buses(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_buses_deleted_at ON buses(deleted_at);

-- Comments
COMMENT ON TABLE buses IS 'Physical buses in the fleet';
COMMENT ON COLUMN buses.bus_type IS 'Type of bus: standard, vip, sleeper, double_decker';
COMMENT ON COLUMN buses.seat_capacity IS 'Total number of seats in the bus';
COMMENT ON COLUMN buses.amenities IS 'Array of amenities like wifi, ac, tv, etc';
