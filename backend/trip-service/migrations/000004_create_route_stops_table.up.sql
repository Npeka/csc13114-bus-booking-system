-- Create route_stops table
CREATE TABLE IF NOT EXISTS route_stops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    route_id UUID NOT NULL,
    stop_order INTEGER NOT NULL CHECK (stop_order >= 1),
    stop_type VARCHAR(20) NOT NULL,
    location VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    offset_minutes INTEGER NOT NULL DEFAULT 0 CHECK (offset_minutes >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    -- Foreign key with CASCADE
    CONSTRAINT fk_route_stops_route 
        FOREIGN KEY (route_id) 
        REFERENCES routes(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);

-- Create partial unique index for route_id + stop_order (excludes soft-deleted records)
CREATE UNIQUE INDEX idx_route_stops_route_order_active 
ON route_stops(route_id, stop_order) 
WHERE deleted_at IS NULL;

-- Create indexes for queries
CREATE INDEX idx_route_stops_route ON route_stops(route_id);
CREATE INDEX idx_route_stops_location ON route_stops(location);
CREATE INDEX idx_route_stops_deleted_at ON route_stops(deleted_at);
CREATE INDEX idx_route_stops_is_active ON route_stops(is_active) WHERE is_active = true;