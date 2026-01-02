-- Rollback seed data
-- Delete trips first (foreign key dependent)
DELETE FROM trips WHERE departure_time >= '2026-01-02' AND departure_time < '2026-01-09';

-- Delete route_stops for seeded routes
DELETE FROM route_stops WHERE route_id IN (
    '2d70fa4a-a146-4c19-92bc-e33a51a517c5',
    '3449c443-2c0e-48f7-b9e9-48e0626e5fd2',
    '1a3baca5-b105-43c1-91c5-5fad154a4dce',
    '011269a3-48c9-4e25-85a7-a1f493a3e55f',
    'aaa145b8-7280-4bb2-8172-2e8ac43b9443',
    'bd656d0c-3bac-4ac4-8c09-2fff22e3e90d',
    '78fdd0a1-0a2c-498f-9293-d1fca7599c15',
    '5124fb9e-2fd4-436a-82fe-efc270a33b08',
    '270cd14d-4bae-441d-9d69-d5ea14f8f86b',
    'fa703cfd-716c-4892-879c-4dc4e24b1074'
);

-- Delete seats for seeded buses
DELETE FROM seats WHERE bus_id IN (
    '2799008e-13ce-4c22-90db-d1e0ead79ce7',
    'd4955ce2-885e-4c0d-891c-09640bc5d835',
    'e6fa78ac-e16a-4b26-a479-0d885d30ca73',
    '39acceb1-e047-4142-bd99-85ef780a850f'
);

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
