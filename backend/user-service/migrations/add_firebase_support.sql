-- Migration to add Firebase UID support to users table
-- This migration adds firebase_uid column if it doesn't exist

DO $$
BEGIN
    -- Add firebase_uid column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='users' AND column_name='firebase_uid') THEN
        ALTER TABLE users ADD COLUMN firebase_uid VARCHAR(255);
        CREATE INDEX IF NOT EXISTS idx_users_firebase_uid ON users(firebase_uid);
    END IF;
    
    -- Add email_verified column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='users' AND column_name='email_verified') THEN
        ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;
    END IF;
END $$;

-- Update comments
COMMENT ON COLUMN users.firebase_uid IS 'Firebase Authentication UID for OAuth2 users';
COMMENT ON COLUMN users.email_verified IS 'Whether the user email has been verified';