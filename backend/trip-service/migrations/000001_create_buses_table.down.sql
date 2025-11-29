-- Drop buses table and its indexes
DROP INDEX IF EXISTS idx_buses_is_active;
DROP INDEX IF EXISTS idx_buses_deleted_at;
DROP INDEX IF EXISTS idx_buses_plate_number_active;
DROP TABLE IF EXISTS buses CASCADE;