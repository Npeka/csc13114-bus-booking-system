-- Drop performance indexes

DROP INDEX IF EXISTS idx_trips_list_covering;
DROP INDEX IF EXISTS idx_route_stops_route_order;
DROP INDEX IF EXISTS idx_seats_count_by_bus;
DROP INDEX IF EXISTS idx_seats_bus_type;
DROP INDEX IF EXISTS idx_buses_active_composite;
DROP INDEX IF EXISTS idx_buses_amenities;
DROP INDEX IF EXISTS idx_trips_list_order;
DROP INDEX IF EXISTS idx_trips_departure_status;
DROP INDEX IF EXISTS idx_trips_search_composite;
DROP INDEX IF EXISTS idx_routes_active_origin_dest;
DROP INDEX IF EXISTS idx_routes_destination_trgm;
DROP INDEX IF EXISTS idx_routes_origin_trgm;

DROP EXTENSION IF EXISTS pg_trgm;
