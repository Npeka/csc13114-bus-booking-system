-- Migration: Create route_stops table
-- Description: Add support for multiple pickup/dropoff points per route

CREATE TABLE IF NOT EXISTS route_stops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    stop_order INTEGER NOT NULL,
    stop_type VARCHAR(20) NOT NULL CHECK (stop_type IN ('pickup', 'dropoff', 'both')),
    location VARCHAR(255) NOT NULL,
    address TEXT,
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    offset_minutes INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT unique_route_stop_order UNIQUE (route_id, stop_order)
);

-- Indexes for performance
CREATE INDEX idx_route_stops_route ON route_stops(route_id, stop_order);
CREATE INDEX idx_route_stops_location ON route_stops(location);
CREATE INDEX idx_route_stops_deleted ON route_stops(deleted_at);

-- Comments
COMMENT ON TABLE route_stops IS 'Pickup and dropoff points for routes';
COMMENT ON COLUMN route_stops.stop_order IS 'Order of stop in the route (1, 2, 3...)';
COMMENT ON COLUMN route_stops.stop_type IS 'Type of stop: pickup, dropoff, or both';
COMMENT ON COLUMN route_stops.offset_minutes IS 'Minutes from route departure time';
