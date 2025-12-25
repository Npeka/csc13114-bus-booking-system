-- Drop trips table and related indexes
DROP INDEX IF EXISTS idx_trips_bus_date_range;
DROP INDEX IF EXISTS idx_trips_route_departure;
DROP INDEX IF EXISTS idx_trips_deleted_at;
DROP INDEX IF EXISTS idx_trips_is_active;
DROP INDEX IF EXISTS idx_trips_status;
DROP INDEX IF EXISTS idx_trips_departure_time;
DROP INDEX IF EXISTS idx_trips_bus_id;
DROP INDEX IF EXISTS idx_trips_route_id;
DROP TABLE IF EXISTS trips CASCADE;
