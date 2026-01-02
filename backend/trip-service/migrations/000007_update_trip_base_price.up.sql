-- Update base_price by dividing by 10 for all trips
UPDATE trips
SET base_price = base_price / 10
WHERE base_price IS NOT NULL;
