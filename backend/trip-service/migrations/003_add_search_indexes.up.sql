-- Migration: Add performance indexes for trip search
-- Description: Optimize search queries with proper indexes

-- Trips table indexes
CREATE INDEX IF NOT EXISTS idx_trips_route_date ON trips(route_id, departure_time);
CREATE INDEX IF NOT EXISTS idx_trips_status_active ON trips(status, is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_trips_departure_time ON trips(departure_time);
CREATE INDEX IF NOT EXISTS idx_trips_price ON trips(base_price);

-- Routes table indexes
CREATE INDEX IF NOT EXISTS idx_routes_origin_dest ON routes(origin, destination);
CREATE INDEX IF NOT EXISTS idx_routes_operator ON routes(operator_id) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_routes_active ON routes(is_active);

-- Buses table indexes
CREATE INDEX IF NOT EXISTS idx_buses_operator ON buses(operator_id) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_buses_active ON buses(is_active);

-- Operators table indexes
CREATE INDEX IF NOT EXISTS idx_operators_active ON operators(is_active);

-- Comments
COMMENT ON INDEX idx_trips_route_date IS 'Optimize trip search by route and date';
COMMENT ON INDEX idx_trips_status_active IS 'Filter active trips by status';
COMMENT ON INDEX idx_routes_origin_dest IS 'Optimize route search by origin and destination';
