-- Drop indexes
DROP INDEX IF EXISTS idx_route_stops_deleted_at;
DROP INDEX IF EXISTS idx_route_stops_location;
DROP INDEX IF EXISTS idx_route_stops_route_order;
DROP INDEX IF EXISTS idx_route_stops_route_id;

-- Drop table
DROP TABLE IF EXISTS route_stops;
