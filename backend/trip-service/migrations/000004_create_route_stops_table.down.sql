-- Drop route_stops table and its indexes
DROP INDEX IF EXISTS idx_route_stops_is_active;
DROP INDEX IF EXISTS idx_route_stops_deleted_at;
DROP INDEX IF EXISTS idx_route_stops_location;
DROP INDEX IF EXISTS idx_route_stops_route;
DROP INDEX IF EXISTS idx_route_stops_route_order_active;
DROP TABLE IF EXISTS route_stops CASCADE;