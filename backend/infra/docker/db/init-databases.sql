-- Initialize databases for each microservice
-- This script runs automatically when PostgreSQL container starts for the first time

-- Create user service database
CREATE DATABASE user_service_db;

-- Create trip service database
CREATE DATABASE trip_service_db;

-- Create booking service database
CREATE DATABASE booking_service_db;

-- Create payment service database
CREATE DATABASE payment_service_db;

-- Grant all privileges to postgres user (already has them, but being explicit)
GRANT ALL PRIVILEGES ON DATABASE user_service_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE trip_service_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE booking_service_db TO postgres;
GRANT ALL PRIVILEGES ON DATABASE payment_service_db TO postgres;
