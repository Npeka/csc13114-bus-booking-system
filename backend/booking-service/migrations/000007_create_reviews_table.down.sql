-- Drop reviews table
DROP INDEX IF EXISTS idx_reviews_trip_rating;
DROP INDEX IF EXISTS idx_reviews_status;
DROP INDEX IF EXISTS idx_reviews_rating;
DROP INDEX IF EXISTS idx_reviews_booking_id;
DROP INDEX IF EXISTS idx_reviews_user_id;
DROP INDEX IF EXISTS idx_reviews_trip_id;
DROP INDEX IF EXISTS idx_reviews_deleted_at;
DROP TABLE IF EXISTS reviews;
