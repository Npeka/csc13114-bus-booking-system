-- Create routes table
CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    distance_km DECIMAL(10,2) NOT NULL,
    estimated_minutes INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT routes_distance_check CHECK (distance_km > 0),
    CONSTRAINT routes_duration_check CHECK (estimated_minutes > 0)
);

-- Indexes for routes
CREATE INDEX idx_routes_origin ON routes(origin) WHERE deleted_at IS NULL;
CREATE INDEX idx_routes_destination ON routes(destination) WHERE deleted_at IS NULL;
CREATE INDEX idx_routes_origin_destination ON routes(origin, destination) WHERE is_active = true AND deleted_at IS NULL;
CREATE INDEX idx_routes_is_active ON routes(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_routes_deleted_at ON routes(deleted_at);

-- Comments
COMMENT ON TABLE routes IS 'Routes between cities/locations';
COMMENT ON COLUMN routes.distance_km IS 'Total distance in kilometers';
COMMENT ON COLUMN routes.estimated_minutes IS 'Estimated travel duration in minutes';
