-- Create routes table
CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    distance_km INTEGER NOT NULL CHECK (distance_km >= 1),
    estimated_minutes INTEGER NOT NULL CHECK (estimated_minutes >= 1),
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Create indexes
CREATE INDEX idx_routes_deleted_at ON routes(deleted_at);
CREATE INDEX idx_routes_origin_destination ON routes(origin, destination) WHERE deleted_at IS NULL;
CREATE INDEX idx_routes_is_active ON routes(is_active) WHERE is_active = true;