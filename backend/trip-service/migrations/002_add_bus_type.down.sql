-- Remove bus_type column from buses table
DROP INDEX IF EXISTS idx_buses_bus_type;

ALTER TABLE buses
DROP CONSTRAINT IF EXISTS buses_bus_type_check;

ALTER TABLE buses 
DROP COLUMN IF EXISTS bus_type;
