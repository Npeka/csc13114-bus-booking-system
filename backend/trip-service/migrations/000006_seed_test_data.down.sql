-- Rollback seed data
-- Delete trips first (foreign key dependent)
DELETE FROM trips WHERE departure_time >= '2026-01-02' AND departure_time < '2026-01-09';

-- Note: We don't delete routes and buses as they may be used by other data
-- If you need to fully rollback, uncomment below:
-- DELETE FROM routes WHERE id IN (
--     '2d70fa4a-a146-4c19-92bc-e33a51a517c5',
--     '3449c443-2c0e-48f7-b9e9-48e0626e5fd2',
--     '1a3baca5-b105-43c1-91c5-5fad154a4dce',
--     '011269a3-48c9-4e25-85a7-a1f493a3e55f',
--     'aaa145b8-7280-4bb2-8172-2e8ac43b9443',
--     'bd656d0c-3bac-4ac4-8c09-2fff22e3e90d',
--     '78fdd0a1-0a2c-498f-9293-d1fca7599c15',
--     '5124fb9e-2fd4-436a-82fe-efc270a33b08',
--     '270cd14d-4bae-441d-9d69-d5ea14f8f86b',
--     'fa703cfd-716c-4892-879c-4dc4e24b1074'
-- );
-- DELETE FROM buses WHERE plate_number IN ('29A-12345', '30B-67890', '51C-11111', '29A-11111');
