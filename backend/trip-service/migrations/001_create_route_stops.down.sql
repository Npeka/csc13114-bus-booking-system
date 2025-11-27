-- Migration: Rollback route_stops table

DROP INDEX IF EXISTS idx_route_stops_deleted;
DROP INDEX IF EXISTS idx_route_stops_location;
DROP INDEX IF EXISTS idx_route_stops_route;

DROP TABLE IF EXISTS route_stops;
