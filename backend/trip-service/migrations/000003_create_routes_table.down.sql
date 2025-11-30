-- Drop routes table and its indexes
DROP INDEX IF EXISTS idx_routes_is_active;
DROP INDEX IF EXISTS idx_routes_origin_destination;
DROP INDEX IF EXISTS idx_routes_deleted_at;
DROP TABLE IF EXISTS routes CASCADE;