-- Drop seats table
DROP INDEX IF EXISTS idx_seats_deleted_at;
DROP INDEX IF EXISTS idx_seats_bus_available;
DROP INDEX IF EXISTS idx_seats_is_available;
DROP INDEX IF EXISTS idx_seats_seat_number;
DROP INDEX IF EXISTS idx_seats_bus_id;
DROP TABLE IF EXISTS seats CASCADE;
