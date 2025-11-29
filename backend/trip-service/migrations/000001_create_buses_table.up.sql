-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "citext";

-- Create buses table
CREATE TABLE IF NOT EXISTS buses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    plate_number VARCHAR(20) NOT NULL,
    model VARCHAR(255) NOT NULL,
    seat_capacity INTEGER NOT NULL CHECK (seat_capacity >= 1 AND seat_capacity <= 100),
    amenities TEXT[],
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Create partial unique index for plate_number (excludes soft-deleted records)
CREATE UNIQUE INDEX idx_buses_plate_number_active 
ON buses(plate_number) 
WHERE deleted_at IS NULL;

-- Create index on deleted_at for soft delete queries
CREATE INDEX idx_buses_deleted_at ON buses(deleted_at);

-- Create index on is_active for filtering
CREATE INDEX idx_buses_is_active ON buses(is_active) WHERE is_active = true;