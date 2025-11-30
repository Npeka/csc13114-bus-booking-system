-- Create trips table
CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    route_id UUID NOT NULL,
    bus_id UUID NOT NULL,
    departure_time TIMESTAMP WITH TIME ZONE NOT NULL,
    arrival_time TIMESTAMP WITH TIME ZONE NOT NULL,
    base_price DECIMAL(10,2) NOT NULL CHECK (base_price >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    -- Foreign keys WITHOUT CASCADE DELETE (trip deletion doesn't affect route/bus)
    CONSTRAINT fk_trips_route 
        FOREIGN KEY (route_id) 
        REFERENCES routes(id) 
        ON UPDATE CASCADE,
    
    CONSTRAINT fk_trips_bus 
        FOREIGN KEY (bus_id) 
        REFERENCES buses(id) 
        ON UPDATE CASCADE,
    
    -- Check constraints
    CONSTRAINT chk_trips_times CHECK (arrival_time > departure_time),
    CONSTRAINT chk_trips_status CHECK (status IN ('scheduled', 'in_progress', 'completed', 'cancelled'))
);

-- Create indexes
CREATE INDEX idx_trips_deleted_at ON trips(deleted_at);
CREATE INDEX idx_trips_route_id ON trips(route_id);
CREATE INDEX idx_trips_bus_id ON trips(bus_id);
CREATE INDEX idx_trips_departure_time ON trips(departure_time);
CREATE INDEX idx_trips_status ON trips(status);
CREATE INDEX idx_trips_is_active ON trips(is_active) WHERE is_active = true;

-- Composite index for searching trips by route and time
CREATE INDEX idx_trips_route_departure ON trips(route_id, departure_time) WHERE deleted_at IS NULL;

-- Composite index for searching trips by bus and time
CREATE INDEX idx_trips_bus_departure ON trips(bus_id, departure_time) WHERE deleted_at IS NULL;