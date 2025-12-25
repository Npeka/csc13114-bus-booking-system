-- Create trips table
CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    bus_id UUID NOT NULL REFERENCES buses(id) ON DELETE RESTRICT,
    departure_time TIMESTAMPTZ NOT NULL,
    arrival_time TIMESTAMPTZ NOT NULL,
    base_price DECIMAL(10,2) NOT NULL CHECK (base_price >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Create indexes
CREATE INDEX idx_trips_route_id ON trips(route_id);
CREATE INDEX idx_trips_bus_id ON trips(bus_id);
CREATE INDEX idx_trips_departure_time ON trips(departure_time);
CREATE INDEX idx_trips_status ON trips(status);
CREATE INDEX idx_trips_is_active ON trips(is_active);
CREATE INDEX idx_trips_deleted_at ON trips(deleted_at);

-- Composite index for common queries
CREATE INDEX idx_trips_route_departure ON trips(route_id, departure_time) WHERE deleted_at IS NULL;
CREATE INDEX idx_trips_bus_date_range ON trips(bus_id, departure_time, arrival_time) WHERE deleted_at IS NULL;
