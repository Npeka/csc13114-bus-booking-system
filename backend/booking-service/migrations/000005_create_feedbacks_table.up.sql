-- Create feedbacks table
CREATE TABLE IF NOT EXISTS feedbacks (
    -- Standard fields
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    
    -- Business fields
    booking_id UUID NOT NULL,
    user_id UUID NOT NULL,
    trip_id UUID NOT NULL,
    rating INTEGER NOT NULL,
    comment TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign keys
    CONSTRAINT fk_feedbacks_booking FOREIGN KEY (booking_id) 
        REFERENCES bookings(id) ON UPDATE CASCADE ON DELETE CASCADE,
    
    -- Constraints
    CONSTRAINT chk_feedback_rating CHECK (rating >= 1 AND rating <= 5)
);

-- Create indexes
CREATE INDEX idx_feedbacks_booking_id ON feedbacks(booking_id);
CREATE INDEX idx_feedbacks_user_id ON feedbacks(user_id);
CREATE INDEX idx_feedbacks_trip_id ON feedbacks(trip_id);
CREATE INDEX idx_feedbacks_rating ON feedbacks(rating);
CREATE INDEX idx_feedbacks_deleted_at ON feedbacks(deleted_at);

-- Add updated_at trigger
CREATE TRIGGER update_feedbacks_updated_at BEFORE UPDATE ON feedbacks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
