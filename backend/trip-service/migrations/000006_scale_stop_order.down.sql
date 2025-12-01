-- Revert stop_order scaling by dividing by 100
UPDATE route_stops 
SET stop_order = stop_order / 100 
WHERE deleted_at IS NULL;
