-- Create route_stops table
CREATE TABLE IF NOT EXISTS route_stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    stop_order INT NOT NULL,
    stop_type VARCHAR(50) NOT NULL CHECK (stop_type IN ('pickup', 'dropoff', 'both')),
    location VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    offset_minutes INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT route_stops_unique_route_order UNIQUE(route_id, stop_order),
    CONSTRAINT route_stops_order_check CHECK (stop_order > 0),
    CONSTRAINT route_stops_offset_check CHECK (offset_minutes >= 0)
);

-- Indexes for route_stops
CREATE INDEX idx_route_stops_route_id ON route_stops(route_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_route_stops_route_order ON route_stops(route_id, stop_order) WHERE deleted_at IS NULL;
CREATE INDEX idx_route_stops_location ON route_stops(location);
CREATE INDEX idx_route_stops_deleted_at ON route_stops(deleted_at);

-- Comments
COMMENT ON TABLE route_stops IS 'Stops along a route (pickup/dropoff points)';
COMMENT ON COLUMN route_stops.stop_type IS 'Type of stop: pickup, dropoff, or both';
COMMENT ON COLUMN route_stops.stop_order IS 'Order of stop in route sequence (1-indexed)';
COMMENT ON COLUMN route_stops.location IS 'Name of the stop location';
COMMENT ON COLUMN route_stops.offset_minutes IS 'Minutes from route origin to reach this stop';