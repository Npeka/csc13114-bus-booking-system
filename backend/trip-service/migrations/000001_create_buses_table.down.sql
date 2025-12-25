-- Drop buses table
DROP INDEX IF EXISTS idx_buses_deleted_at;
DROP INDEX IF EXISTS idx_buses_is_active;
DROP INDEX IF EXISTS idx_buses_bus_type;
DROP INDEX IF EXISTS idx_buses_plate_number;
DROP TABLE IF EXISTS buses CASCADE;
