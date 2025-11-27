-- Migration: Rollback search indexes

DROP INDEX IF EXISTS idx_operators_active;
DROP INDEX IF EXISTS idx_buses_active;
DROP INDEX IF EXISTS idx_buses_operator;
DROP INDEX IF EXISTS idx_routes_active;
DROP INDEX IF EXISTS idx_routes_operator;
DROP INDEX IF EXISTS idx_routes_origin_dest;
DROP INDEX IF EXISTS idx_trips_price;
DROP INDEX IF EXISTS idx_trips_departure_time;
DROP INDEX IF EXISTS idx_trips_status_active;
DROP INDEX IF EXISTS idx_trips_route_date;
