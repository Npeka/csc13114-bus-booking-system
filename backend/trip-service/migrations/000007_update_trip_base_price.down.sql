-- Revert base_price update by multiplying by 10 for all trips
UPDATE trips
SET base_price = base_price * 10
WHERE base_price IS NOT NULL;
