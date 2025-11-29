-- Drop seats table and its indexes
DROP INDEX IF EXISTS idx_seats_is_available;
DROP INDEX IF EXISTS idx_seats_deleted_at;
DROP INDEX IF EXISTS idx_seats_bus_id;
DROP INDEX IF EXISTS idx_seats_bus_position_active;
DROP INDEX IF EXISTS idx_seats_bus_seat_number_active;
DROP TABLE IF EXISTS seats CASCADE;