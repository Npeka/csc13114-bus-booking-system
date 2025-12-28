-- Drop indexes first
DROP INDEX IF EXISTS idx_refunds_status_created;
DROP INDEX IF EXISTS idx_refunds_processed_by;
DROP INDEX IF EXISTS idx_refunds_created_at;
DROP INDEX IF EXISTS idx_refunds_pending;
DROP INDEX IF EXISTS idx_refunds_status;
DROP INDEX IF EXISTS idx_refunds_user_id;
DROP INDEX IF EXISTS idx_refunds_transaction_id;
DROP INDEX IF EXISTS idx_refunds_booking_id;
DROP INDEX IF EXISTS idx_refunds_unique_booking;

-- Drop table
DROP TABLE IF EXISTS refunds;
