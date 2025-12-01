-- Scale stop_order by 100 to allow flexible ordering
-- This provides space between stops for future insertions without reordering
UPDATE route_stops 
SET stop_order = stop_order * 100 
WHERE deleted_at IS NULL;
