-- Add bus_type column to buses table
ALTER TABLE buses 
ADD COLUMN bus_type VARCHAR(20) NOT NULL DEFAULT 'standard';

-- Add check constraint for valid bus types
ALTER TABLE buses
ADD CONSTRAINT buses_bus_type_check 
CHECK (bus_type IN ('standard', 'vip', 'sleeper', 'double_decker'));

-- Add index for bus_type filtering
CREATE INDEX idx_buses_bus_type ON buses(bus_type);
