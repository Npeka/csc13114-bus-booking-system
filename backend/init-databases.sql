-- Initialize databases for each microservice
-- This script runs automatically when PostgreSQL container starts for the first time

-- Create user service database
CREATE DATABASE bus_booking_user;

-- Create trip service database
CREATE DATABASE bus_booking_trip;

-- Create booking service database
CREATE DATABASE bus_booking_booking;

-- Create payment service database
CREATE DATABASE bus_booking_payment;

-- Grant all privileges to postgres user (already has them, but being explicit)
GRANT ALL PRIVILEGES ON DATABASE bus_booking_user TO postgres;
GRANT ALL PRIVILEGES ON DATABASE bus_booking_trip TO postgres;
GRANT ALL PRIVILEGES ON DATABASE bus_booking_booking TO postgres;
GRANT ALL PRIVILEGES ON DATABASE bus_booking_payment TO postgres;
