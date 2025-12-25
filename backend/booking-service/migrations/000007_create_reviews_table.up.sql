-- Create reviews table in booking-service
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    -- Trip reference
    trip_id UUID NOT NULL,
    
    -- User & booking verification
    user_id UUID NOT NULL,
    booking_id UUID NOT NULL,
    
    -- Review content
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    
    -- Moderation
    is_verified BOOLEAN NOT NULL DEFAULT true,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    admin_notes TEXT,
    
    -- Check constraints
    CONSTRAINT chk_reviews_status CHECK (status IN ('active', 'hidden', 'flagged', 'removed'))
);

-- Prevent duplicate reviews per booking (using partial unique index)
CREATE UNIQUE INDEX uq_reviews_booking ON reviews(booking_id) WHERE deleted_at IS NULL;

-- Create indexes
CREATE INDEX idx_reviews_deleted_at ON reviews(deleted_at);
CREATE INDEX idx_reviews_trip_id ON reviews(trip_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reviews_user_id ON reviews(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reviews_booking_id ON reviews(booking_id);
CREATE INDEX idx_reviews_rating ON reviews(rating) WHERE deleted_at IS NULL;
CREATE INDEX idx_reviews_status ON reviews(status) WHERE deleted_at IS NULL AND status = 'active';

-- Composite index for trip review queries
CREATE INDEX idx_reviews_trip_rating ON reviews(trip_id, rating DESC) WHERE deleted_at IS NULL AND status = 'active';
