-- Drop routes table
DROP INDEX IF EXISTS idx_routes_deleted_at;
DROP INDEX IF EXISTS idx_routes_is_active;
DROP INDEX IF EXISTS idx_routes_origin_destination;
DROP INDEX IF EXISTS idx_routes_destination;
DROP INDEX IF EXISTS idx_routes_origin;
DROP TABLE IF EXISTS routes CASCADE;
