-- Add performance indexes for common queries

-- 1. Routes: Add text search indexes for ILIKE queries (origin, destination)
-- Use trigram index for faster ILIKE '%search%' queries
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_routes_origin_trgm ON routes USING gin (origin gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_routes_destination_trgm ON routes USING gin (destination gin_trgm_ops);

-- Composite index for active routes search
CREATE INDEX IF NOT EXISTS idx_routes_active_origin_dest ON routes(is_active, origin, destination) 
WHERE deleted_at IS NULL AND is_active = true;

-- 2. Trips: Add composite indexes for SearchTrips queries
-- This query joins trips with routes and buses, filters by date, status, and active
CREATE INDEX IF NOT EXISTS idx_trips_search_composite ON trips(is_active, status, departure_time) 
WHERE deleted_at IS NULL AND is_active = true;

-- Composite index for date + status queries (avoid function-based index)
CREATE INDEX IF NOT EXISTS idx_trips_departure_status ON trips(departure_time, status, is_active) 
WHERE deleted_at IS NULL AND is_active = true;

-- Composite index for ListTrips with ordering
CREATE INDEX IF NOT EXISTS idx_trips_list_order ON trips(departure_time) 
WHERE deleted_at IS NULL;

-- 3. Buses: Add index for amenities array queries (ANY operator)
CREATE INDEX IF NOT EXISTS idx_buses_amenities ON buses USING gin (amenities);

-- Composite index for active buses
CREATE INDEX IF NOT EXISTS idx_buses_active_composite ON buses(is_active, created_at) 
WHERE deleted_at IS NULL AND is_active = true;

-- 4. Seats: Add composite index for seat type filtering
CREATE INDEX IF NOT EXISTS idx_seats_bus_type ON seats(bus_id, seat_type) 
WHERE deleted_at IS NULL;

-- Index for counting seats by bus
CREATE INDEX IF NOT EXISTS idx_seats_count_by_bus ON seats(bus_id) 
WHERE deleted_at IS NULL;

-- 5. Route Stops: Add index for ordering queries
CREATE INDEX IF NOT EXISTS idx_route_stops_route_order ON route_stops(route_id, stop_order) 
WHERE deleted_at IS NULL AND is_active = true;

-- 6. Add covering index for trips list query (includes preload fields)
CREATE INDEX IF NOT EXISTS idx_trips_list_covering ON trips(id, route_id, bus_id, departure_time, arrival_time, base_price, status, is_active) 
WHERE deleted_at IS NULL;
