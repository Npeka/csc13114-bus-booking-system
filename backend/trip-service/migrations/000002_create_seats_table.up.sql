-- Create seats table
CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    bus_id UUID NOT NULL,
    seat_number VARCHAR(10) NOT NULL,
    row INTEGER NOT NULL CHECK (row >= 1),
    "column" INTEGER NOT NULL CHECK ("column" >= 1),
    seat_type VARCHAR(20) NOT NULL,
    price_multiplier DECIMAL(3,2) NOT NULL DEFAULT 1.0 CHECK (price_multiplier >= 0.5 AND price_multiplier <= 5.0),
    is_available BOOLEAN NOT NULL DEFAULT true,
    floor INTEGER NOT NULL DEFAULT 1 CHECK (floor >= 1 AND floor <= 2),
    
    -- Foreign key with CASCADE
    CONSTRAINT fk_seats_bus 
        FOREIGN KEY (bus_id) 
        REFERENCES buses(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);

-- Create partial unique indexes (excludes soft-deleted records)
CREATE UNIQUE INDEX idx_seats_bus_seat_number_active 
ON seats(bus_id, seat_number) 
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX idx_seats_bus_position_active 
ON seats(bus_id, row, "column", floor) 
WHERE deleted_at IS NULL;

-- Create indexes for queries
CREATE INDEX idx_seats_bus_id ON seats(bus_id);
CREATE INDEX idx_seats_deleted_at ON seats(deleted_at);
CREATE INDEX idx_seats_is_available ON seats(is_available) WHERE is_available = true;