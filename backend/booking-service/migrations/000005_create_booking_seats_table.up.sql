-- Create booking_seats table
CREATE TABLE IF NOT EXISTS booking_seats (
    -- Standard fields
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Business fields
    booking_id UUID NOT NULL,
    seat_id UUID NOT NULL,
    seat_number VARCHAR(10) NOT NULL,
    seat_type VARCHAR(50) NOT NULL,
    floor INT NOT NULL DEFAULT 1,
    price DECIMAL(10,2) NOT NULL,
    price_multiplier DECIMAL(3,2) NOT NULL DEFAULT 1.0,
    
    -- Optional passenger info
    passenger_name VARCHAR(255),
    passenger_id VARCHAR(50),
    passenger_phone VARCHAR(20),
    
    -- Foreign keys
    CONSTRAINT fk_booking_seats_booking FOREIGN KEY (booking_id) 
        REFERENCES bookings(id) ON UPDATE CASCADE ON DELETE CASCADE,
    
    -- Constraints
    CONSTRAINT uq_booking_seats_booking_seat UNIQUE (booking_id, seat_id)
);

-- Create indexes
CREATE INDEX idx_booking_seats_booking_id ON booking_seats(booking_id);
CREATE INDEX idx_booking_seats_seat_id ON booking_seats(seat_id);
CREATE INDEX idx_booking_seats_deleted_at ON booking_seats(deleted_at);

-- Add updated_at trigger
CREATE TRIGGER update_booking_seats_updated_at BEFORE UPDATE ON booking_seats
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
