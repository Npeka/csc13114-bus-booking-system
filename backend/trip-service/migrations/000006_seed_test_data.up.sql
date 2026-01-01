-- Seed buses for testing
-- Only insert if not exists (using ON CONFLICT DO NOTHING)

INSERT INTO buses (id, plate_number, model, bus_type, seat_capacity, amenities, is_active, created_at, updated_at) VALUES
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '29A-12345', 'Mercedes-Benz Travego', 'vip', 45, ARRAY['wifi', 'ac', 'toilet', 'charging', 'blanket', 'water'], true, NOW(), NOW()),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '30B-67890', 'Hyundai Universe', 'standard', 40, ARRAY['wifi', 'ac', 'charging', 'water'], true, NOW(), NOW()),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', '51C-11111', 'Thaco TB120S', 'sleeper', 34, ARRAY['ac', 'charging', 'blanket'], true, NOW(), NOW()),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '29A-11111', 'Mercedes-Benz Travego', 'vip', 45, ARRAY['wifi', 'ac', 'toilet', 'charging', 'blanket', 'water'], true, NOW(), NOW())
ON CONFLICT (plate_number) DO NOTHING;

-- Seed routes for testing
INSERT INTO routes (id, origin, destination, distance_km, estimated_minutes, is_active, created_at, updated_at) VALUES
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 'Hà Nội', 'TP. Hồ Chí Minh', 1710, 1800, true, NOW(), NOW()),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'TP. Hồ Chí Minh', 'Đà Lạt', 308, 360, true, NOW(), NOW()),
    ('1a3baca5-b105-43c1-91c5-5fad154a4dce', 'Hà Nội', 'Đà Nẵng', 763, 840, true, NOW(), NOW()),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', 'TP. Hồ Chí Minh', 'Nha Trang', 448, 480, true, NOW(), NOW()),
    ('aaa145b8-7280-4bb2-8172-2e8ac43b9443', 'Hà Nội', 'Hải Phòng', 103, 120, true, NOW(), NOW()),
    ('bd656d0c-3bac-4ac4-8c09-2fff22e3e90d', 'TP. Hồ Chí Minh', 'Cần Thơ', 169, 180, true, NOW(), NOW()),
    ('78fdd0a1-0a2c-498f-9293-d1fca7599c15', 'Hà Nội', 'Hạ Long', 156, 180, true, NOW(), NOW()),
    ('5124fb9e-2fd4-436a-82fe-efc270a33b08', 'Đà Nẵng', 'TP. Hồ Chí Minh', 947, 960, true, NOW(), NOW()),
    ('270cd14d-4bae-441d-9d69-d5ea14f8f86b', 'Hà Nội', 'Sapa', 315, 360, true, NOW(), NOW()),
    ('fa703cfd-716c-4892-879c-4dc4e24b1074', 'Hà Nội', 'Cao Bằng', 272, 360, true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Seed trips for January 2-8, 2026 (for chatbot testing)

-- January 2, 2026 trips
INSERT INTO trips (route_id, bus_id, departure_time, arrival_time, base_price, status, is_active) VALUES
    -- Hà Nội → TP. Hồ Chí Minh
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 06:00:00+07', '2026-01-03 12:00:00+07', 450000, 'scheduled', true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-02 18:00:00+07', '2026-01-04 00:00:00+07', 480000, 'scheduled', true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-02 22:00:00+07', '2026-01-04 04:00:00+07', 500000, 'scheduled', true),
    -- TP. Hồ Chí Minh → Đà Lạt
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-02 06:00:00+07', '2026-01-02 12:00:00+07', 250000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-02 08:00:00+07', '2026-01-02 14:00:00+07', 260000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 14:00:00+07', '2026-01-02 20:00:00+07', 240000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', '39acceb1-e047-4142-bd99-85ef780a850f', '2026-01-02 20:00:00+07', '2026-01-03 02:00:00+07', 280000, 'scheduled', true),
    -- Hà Nội → Đà Nẵng
    ('1a3baca5-b105-43c1-91c5-5fad154a4dce', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-02 07:00:00+07', '2026-01-02 21:00:00+07', 320000, 'scheduled', true),
    ('1a3baca5-b105-43c1-91c5-5fad154a4dce', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 20:00:00+07', '2026-01-03 10:00:00+07', 350000, 'scheduled', true),
    -- TP. Hồ Chí Minh → Nha Trang
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-02 08:00:00+07', '2026-01-02 16:00:00+07', 200000, 'scheduled', true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-02 14:00:00+07', '2026-01-02 22:00:00+07', 220000, 'scheduled', true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 22:00:00+07', '2026-01-03 06:00:00+07', 230000, 'scheduled', true),
    -- Hà Nội → Hải Phòng
    ('aaa145b8-7280-4bb2-8172-2e8ac43b9443', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 06:30:00+07', '2026-01-02 08:30:00+07', 120000, 'scheduled', true),
    ('aaa145b8-7280-4bb2-8172-2e8ac43b9443', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-02 10:00:00+07', '2026-01-02 12:00:00+07', 110000, 'scheduled', true),
    -- TP. Hồ Chí Minh → Cần Thơ
    ('bd656d0c-3bac-4ac4-8c09-2fff22e3e90d', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-02 07:30:00+07', '2026-01-02 10:30:00+07', 150000, 'scheduled', true),
    ('bd656d0c-3bac-4ac4-8c09-2fff22e3e90d', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 13:00:00+07', '2026-01-02 16:00:00+07', 160000, 'scheduled', true),
    -- Hà Nội → Hạ Long
    ('78fdd0a1-0a2c-498f-9293-d1fca7599c15', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-02 08:00:00+07', '2026-01-02 11:00:00+07', 180000, 'scheduled', true),
    ('78fdd0a1-0a2c-498f-9293-d1fca7599c15', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-02 15:00:00+07', '2026-01-02 18:00:00+07', 170000, 'scheduled', true),
    -- Đà Nẵng → TP. Hồ Chí Minh
    ('5124fb9e-2fd4-436a-82fe-efc270a33b08', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 19:00:00+07', '2026-01-03 11:00:00+07', 380000, 'scheduled', true),
    ('5124fb9e-2fd4-436a-82fe-efc270a33b08', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-02 21:00:00+07', '2026-01-03 13:00:00+07', 400000, 'scheduled', true),
    -- Hà Nội → Sapa
    ('270cd14d-4bae-441d-9d69-d5ea14f8f86b', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-02 20:00:00+07', '2026-01-03 02:00:00+07', 280000, 'scheduled', true),
    ('270cd14d-4bae-441d-9d69-d5ea14f8f86b', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 06:00:00+07', '2026-01-02 12:00:00+07', 300000, 'scheduled', true),
    -- Hà Nội → Cao Bằng
    ('fa703cfd-716c-4892-879c-4dc4e24b1074', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-02 05:00:00+07', '2026-01-02 11:00:00+07', 260000, 'scheduled', true),
    ('fa703cfd-716c-4892-879c-4dc4e24b1074', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-02 16:00:00+07', '2026-01-02 22:00:00+07', 270000, 'scheduled', true);

-- January 3, 2026 trips
INSERT INTO trips (route_id, bus_id, departure_time, arrival_time, base_price, status, is_active) VALUES
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-03 06:00:00+07', '2026-01-04 12:00:00+07', 450000, 'scheduled', true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-03 18:00:00+07', '2026-01-05 00:00:00+07', 480000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-03 06:00:00+07', '2026-01-03 12:00:00+07', 250000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-03 14:00:00+07', '2026-01-03 20:00:00+07', 240000, 'scheduled', true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-03 08:00:00+07', '2026-01-03 16:00:00+07', 210000, 'scheduled', true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-03 20:00:00+07', '2026-01-04 04:00:00+07', 220000, 'scheduled', true),
    ('bd656d0c-3bac-4ac4-8c09-2fff22e3e90d', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-03 07:30:00+07', '2026-01-03 10:30:00+07', 155000, 'scheduled', true),
    ('aaa145b8-7280-4bb2-8172-2e8ac43b9443', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-03 06:30:00+07', '2026-01-03 08:30:00+07', 115000, 'scheduled', true);

-- January 4, 2026 trips
INSERT INTO trips (route_id, bus_id, departure_time, arrival_time, base_price, status, is_active) VALUES
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-04 06:00:00+07', '2026-01-05 12:00:00+07', 450000, 'scheduled', true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-04 18:00:00+07', '2026-01-06 00:00:00+07', 480000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-04 06:00:00+07', '2026-01-04 12:00:00+07', 250000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-04 14:00:00+07', '2026-01-04 20:00:00+07', 240000, 'scheduled', true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-04 08:00:00+07', '2026-01-04 16:00:00+07', 210000, 'scheduled', true),
    ('1a3baca5-b105-43c1-91c5-5fad154a4dce', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-04 07:00:00+07', '2026-01-04 21:00:00+07', 330000, 'scheduled', true);

-- January 5, 2026 trips
INSERT INTO trips (route_id, bus_id, departure_time, arrival_time, base_price, status, is_active) VALUES
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-05 06:00:00+07', '2026-01-06 12:00:00+07', 450000, 'scheduled', true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-05 18:00:00+07', '2026-01-07 00:00:00+07', 480000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-05 06:00:00+07', '2026-01-05 12:00:00+07', 250000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-05 14:00:00+07', '2026-01-05 20:00:00+07', 240000, 'scheduled', true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-05 08:00:00+07', '2026-01-05 16:00:00+07', 210000, 'scheduled', true);

-- January 6, 2026 trips
INSERT INTO trips (route_id, bus_id, departure_time, arrival_time, base_price, status, is_active) VALUES
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-06 06:00:00+07', '2026-01-07 12:00:00+07', 450000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-06 06:00:00+07', '2026-01-06 12:00:00+07', 250000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-06 14:00:00+07', '2026-01-06 20:00:00+07', 240000, 'scheduled', true);

-- January 7, 2026 trips
INSERT INTO trips (route_id, bus_id, departure_time, arrival_time, base_price, status, is_active) VALUES
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-07 06:00:00+07', '2026-01-08 12:00:00+07', 450000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-07 06:00:00+07', '2026-01-07 12:00:00+07', 250000, 'scheduled', true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-07 08:00:00+07', '2026-01-07 16:00:00+07', 210000, 'scheduled', true);

-- January 8, 2026 trips
INSERT INTO trips (route_id, bus_id, departure_time, arrival_time, base_price, status, is_active) VALUES
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', '2799008e-13ce-4c22-90db-d1e0ead79ce7', '2026-01-08 06:00:00+07', '2026-01-09 12:00:00+07', 450000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'd4955ce2-885e-4c0d-891c-09640bc5d835', '2026-01-08 06:00:00+07', '2026-01-08 12:00:00+07', 250000, 'scheduled', true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 'e6fa78ac-e16a-4b26-a479-0d885d30ca73', '2026-01-08 14:00:00+07', '2026-01-08 20:00:00+07', 240000, 'scheduled', true);
