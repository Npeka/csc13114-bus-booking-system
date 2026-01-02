-- Seed buses for testing
-- Only insert if not exists (using ON CONFLICT DO NOTHING)

INSERT INTO buses (id, plate_number, model, bus_type, seat_capacity, amenities, is_active, created_at, updated_at) VALUES
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '29A-12345', 'Mercedes-Benz Travego', 'vip', 45, ARRAY['wifi', 'ac', 'toilet', 'charging', 'blanket', 'water'], true, NOW(), NOW()),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '30B-67890', 'Hyundai Universe', 'standard', 40, ARRAY['wifi', 'ac', 'charging', 'water'], true, NOW(), NOW()),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', '51C-11111', 'Thaco TB120S', 'sleeper', 34, ARRAY['ac', 'charging', 'blanket'], true, NOW(), NOW()),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '29A-11111', 'Mercedes-Benz Travego', 'vip', 45, ARRAY['wifi', 'ac', 'toilet', 'charging', 'blanket', 'water'], true, NOW(), NOW())
ON CONFLICT (plate_number) DO NOTHING;

-- Seed seats for each bus
-- Bus 1: Mercedes-Benz Travego VIP (45 seats - 2 floor layout)
-- Floor 1: 3 columns (A, B, C), 8 rows = 24 seats
-- Floor 2: 3 columns (A, B, C), 7 rows = 21 seats
INSERT INTO seats (bus_id, seat_number, seat_type, floor, "row", "column", price_multiplier, is_available) VALUES
    -- Bus 1 (VIP) - Floor 1
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '1A', 'vip', 1, 1, 1, 1.2, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '1B', 'vip', 1, 1, 2, 1.2, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '1C', 'vip', 1, 1, 3, 1.2, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '2A', 'vip', 1, 2, 1, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '2B', 'vip', 1, 2, 2, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '2C', 'vip', 1, 2, 3, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '3A', 'vip', 1, 3, 1, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '3B', 'vip', 1, 3, 2, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '3C', 'vip', 1, 3, 3, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '4A', 'vip', 1, 4, 1, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '4B', 'vip', 1, 4, 2, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '4C', 'vip', 1, 4, 3, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '5A', 'vip', 1, 5, 1, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '5B', 'vip', 1, 5, 2, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '5C', 'vip', 1, 5, 3, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '6A', 'vip', 1, 6, 1, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '6B', 'vip', 1, 6, 2, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '6C', 'vip', 1, 6, 3, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '7A', 'vip', 1, 7, 1, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '7B', 'vip', 1, 7, 2, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '7C', 'vip', 1, 7, 3, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '8A', 'vip', 1, 8, 1, 0.9, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '8B', 'vip', 1, 8, 2, 0.9, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', '8C', 'vip', 1, 8, 3, 0.9, true),
    -- Bus 1 (VIP) - Floor 2
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U1A', 'vip', 2, 1, 1, 1.3, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U1B', 'vip', 2, 1, 2, 1.3, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U1C', 'vip', 2, 1, 3, 1.3, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U2A', 'vip', 2, 2, 1, 1.2, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U2B', 'vip', 2, 2, 2, 1.2, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U2C', 'vip', 2, 2, 3, 1.2, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U3A', 'vip', 2, 3, 1, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U3B', 'vip', 2, 3, 2, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U3C', 'vip', 2, 3, 3, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U4A', 'vip', 2, 4, 1, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U4B', 'vip', 2, 4, 2, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U4C', 'vip', 2, 4, 3, 1.1, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U5A', 'vip', 2, 5, 1, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U5B', 'vip', 2, 5, 2, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U5C', 'vip', 2, 5, 3, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U6A', 'vip', 2, 6, 1, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U6B', 'vip', 2, 6, 2, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U6C', 'vip', 2, 6, 3, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U7A', 'vip', 2, 7, 1, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U7B', 'vip', 2, 7, 2, 1.0, true),
    ('2799008e-13ce-4c22-90db-d1e0ead79ce7', 'U7C', 'vip', 2, 7, 3, 1.0, true),
    
    -- Bus 2: Hyundai Universe Standard (40 seats - single floor)
    -- Single floor: 4 columns (A, B, aisle, C, D), 10 rows = 40 seats
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '1A', 'standard', 1, 1, 1, 1.1, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '1B', 'standard', 1, 1, 2, 1.1, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '1C', 'standard', 1, 1, 3, 1.1, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '1D', 'standard', 1, 1, 4, 1.1, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '2A', 'standard', 1, 2, 1, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '2B', 'standard', 1, 2, 2, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '2C', 'standard', 1, 2, 3, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '2D', 'standard', 1, 2, 4, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '3A', 'standard', 1, 3, 1, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '3B', 'standard', 1, 3, 2, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '3C', 'standard', 1, 3, 3, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '3D', 'standard', 1, 3, 4, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '4A', 'standard', 1, 4, 1, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '4B', 'standard', 1, 4, 2, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '4C', 'standard', 1, 4, 3, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '4D', 'standard', 1, 4, 4, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '5A', 'standard', 1, 5, 1, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '5B', 'standard', 1, 5, 2, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '5C', 'standard', 1, 5, 3, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '5D', 'standard', 1, 5, 4, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '6A', 'standard', 1, 6, 1, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '6B', 'standard', 1, 6, 2, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '6C', 'standard', 1, 6, 3, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '6D', 'standard', 1, 6, 4, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '7A', 'standard', 1, 7, 1, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '7B', 'standard', 1, 7, 2, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '7C', 'standard', 1, 7, 3, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '7D', 'standard', 1, 7, 4, 1.0, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '8A', 'standard', 1, 8, 1, 0.9, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '8B', 'standard', 1, 8, 2, 0.9, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '8C', 'standard', 1, 8, 3, 0.9, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '8D', 'standard', 1, 8, 4, 0.9, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '9A', 'standard', 1, 9, 1, 0.9, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '9B', 'standard', 1, 9, 2, 0.9, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '9C', 'standard', 1, 9, 3, 0.9, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '9D', 'standard', 1, 9, 4, 0.9, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '10A', 'standard', 1, 10, 1, 0.85, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '10B', 'standard', 1, 10, 2, 0.85, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '10C', 'standard', 1, 10, 3, 0.85, true),
    ('d4955ce2-885e-4c0d-891c-09640bc5d835', '10D', 'standard', 1, 10, 4, 0.85, true),
    
    -- Bus 3: Thaco TB120S Sleeper (34 seats - 2 floor layout with sleeper beds)
    -- Floor 1: 2 columns (A, B), 9 rows = 18 beds
    -- Floor 2: 2 columns (A, B), 8 rows = 16 beds
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L1A', 'sleeper', 1, 1, 1, 1.2, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L1B', 'sleeper', 1, 1, 2, 1.2, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L2A', 'sleeper', 1, 2, 1, 1.1, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L2B', 'sleeper', 1, 2, 2, 1.1, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L3A', 'sleeper', 1, 3, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L3B', 'sleeper', 1, 3, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L4A', 'sleeper', 1, 4, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L4B', 'sleeper', 1, 4, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L5A', 'sleeper', 1, 5, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L5B', 'sleeper', 1, 5, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L6A', 'sleeper', 1, 6, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L6B', 'sleeper', 1, 6, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L7A', 'sleeper', 1, 7, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L7B', 'sleeper', 1, 7, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L8A', 'sleeper', 1, 8, 1, 0.9, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L8B', 'sleeper', 1, 8, 2, 0.9, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L9A', 'sleeper', 1, 9, 1, 0.9, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'L9B', 'sleeper', 1, 9, 2, 0.9, true),
    -- Bus 3 - Floor 2
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U1A', 'sleeper', 2, 1, 1, 1.3, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U1B', 'sleeper', 2, 1, 2, 1.3, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U2A', 'sleeper', 2, 2, 1, 1.2, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U2B', 'sleeper', 2, 2, 2, 1.2, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U3A', 'sleeper', 2, 3, 1, 1.1, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U3B', 'sleeper', 2, 3, 2, 1.1, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U4A', 'sleeper', 2, 4, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U4B', 'sleeper', 2, 4, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U5A', 'sleeper', 2, 5, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U5B', 'sleeper', 2, 5, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U6A', 'sleeper', 2, 6, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U6B', 'sleeper', 2, 6, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U7A', 'sleeper', 2, 7, 1, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U7B', 'sleeper', 2, 7, 2, 1.0, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U8A', 'sleeper', 2, 8, 1, 0.9, true),
    ('e6fa78ac-e16a-4b26-a479-0d885d30ca73', 'U8B', 'sleeper', 2, 8, 2, 0.9, true),
    
    -- Bus 4: Mercedes-Benz Travego VIP (45 seats - same layout as Bus 1)
    -- Floor 1: 3 columns (A, B, C), 8 rows = 24 seats
    ('39acceb1-e047-4142-bd99-85ef780a850f', '1A', 'vip', 1, 1, 1, 1.2, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '1B', 'vip', 1, 1, 2, 1.2, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '1C', 'vip', 1, 1, 3, 1.2, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '2A', 'vip', 1, 2, 1, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '2B', 'vip', 1, 2, 2, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '2C', 'vip', 1, 2, 3, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '3A', 'vip', 1, 3, 1, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '3B', 'vip', 1, 3, 2, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '3C', 'vip', 1, 3, 3, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '4A', 'vip', 1, 4, 1, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '4B', 'vip', 1, 4, 2, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '4C', 'vip', 1, 4, 3, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '5A', 'vip', 1, 5, 1, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '5B', 'vip', 1, 5, 2, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '5C', 'vip', 1, 5, 3, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '6A', 'vip', 1, 6, 1, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '6B', 'vip', 1, 6, 2, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '6C', 'vip', 1, 6, 3, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '7A', 'vip', 1, 7, 1, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '7B', 'vip', 1, 7, 2, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '7C', 'vip', 1, 7, 3, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '8A', 'vip', 1, 8, 1, 0.9, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '8B', 'vip', 1, 8, 2, 0.9, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', '8C', 'vip', 1, 8, 3, 0.9, true),
    -- Bus 4 - Floor 2
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U1A', 'vip', 2, 1, 1, 1.3, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U1B', 'vip', 2, 1, 2, 1.3, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U1C', 'vip', 2, 1, 3, 1.3, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U2A', 'vip', 2, 2, 1, 1.2, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U2B', 'vip', 2, 2, 2, 1.2, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U2C', 'vip', 2, 2, 3, 1.2, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U3A', 'vip', 2, 3, 1, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U3B', 'vip', 2, 3, 2, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U3C', 'vip', 2, 3, 3, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U4A', 'vip', 2, 4, 1, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U4B', 'vip', 2, 4, 2, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U4C', 'vip', 2, 4, 3, 1.1, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U5A', 'vip', 2, 5, 1, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U5B', 'vip', 2, 5, 2, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U5C', 'vip', 2, 5, 3, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U6A', 'vip', 2, 6, 1, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U6B', 'vip', 2, 6, 2, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U6C', 'vip', 2, 6, 3, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U7A', 'vip', 2, 7, 1, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U7B', 'vip', 2, 7, 2, 1.0, true),
    ('39acceb1-e047-4142-bd99-85ef780a850f', 'U7C', 'vip', 2, 7, 3, 1.0, true)
ON CONFLICT (bus_id, seat_number) DO NOTHING;

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

-- Seed route stops for each route
INSERT INTO route_stops (route_id, stop_order, stop_type, location, address, latitude, longitude, offset_minutes, is_active) VALUES
    -- Route: Hà Nội → TP. Hồ Chí Minh (2d70fa4a-a146-4c19-92bc-e33a51a517c5)
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 1, 'pickup', 'Bến xe Mỹ Đình', '20 Phạm Hùng, Mỹ Đình, Nam Từ Liêm, Hà Nội', 21.0285, 105.7823, 0, true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 2, 'pickup', 'Bến xe Giáp Bát', '6 Giải Phóng, Hoàng Mai, Hà Nội', 20.9789, 105.8417, 30, true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 3, 'both', 'Bến xe Vinh', 'TP. Vinh, Nghệ An', 18.6796, 105.6814, 300, true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 4, 'both', 'Bến xe Đà Nẵng', '33 Điện Biên Phủ, Đà Nẵng', 16.0678, 108.2120, 840, true),
    ('2d70fa4a-a146-4c19-92bc-e33a51a517c5', 5, 'dropoff', 'Bến xe Miền Đông', '292 Đinh Bộ Lĩnh, Bình Thạnh, TP.HCM', 10.8161, 106.7138, 1800, true),
    
    -- Route: TP. Hồ Chí Minh → Đà Lạt (3449c443-2c0e-48f7-b9e9-48e0626e5fd2)
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 1, 'pickup', 'Bến xe Miền Đông', '292 Đinh Bộ Lĩnh, Bình Thạnh, TP.HCM', 10.8161, 106.7138, 0, true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 2, 'pickup', 'Ngã tư Thủ Đức', 'Quốc lộ 1A, TP. Thủ Đức, TP.HCM', 10.8489, 106.7539, 20, true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 3, 'both', 'Bảo Lộc', 'TP. Bảo Lộc, Lâm Đồng', 11.5486, 107.8079, 180, true),
    ('3449c443-2c0e-48f7-b9e9-48e0626e5fd2', 4, 'dropoff', 'Bến xe Đà Lạt', '1 Tô Hiến Thành, Phường 3, Đà Lạt', 11.9456, 108.4514, 360, true),
    
    -- Route: Hà Nội → Đà Nẵng (1a3baca5-b105-43c1-91c5-5fad154a4dce)
    ('1a3baca5-b105-43c1-91c5-5fad154a4dce', 1, 'pickup', 'Bến xe Nước Ngầm', 'Phan Trọng Tuệ, Thanh Trì, Hà Nội', 20.9648, 105.8532, 0, true),
    ('1a3baca5-b105-43c1-91c5-5fad154a4dce', 2, 'both', 'Bến xe Vinh', 'TP. Vinh, Nghệ An', 18.6796, 105.6814, 300, true),
    ('1a3baca5-b105-43c1-91c5-5fad154a4dce', 3, 'both', 'Bến xe Huế', 'Bùi Thị Xuân, TP. Huế', 16.4637, 107.5909, 600, true),
    ('1a3baca5-b105-43c1-91c5-5fad154a4dce', 4, 'dropoff', 'Bến xe Đà Nẵng', '33 Điện Biên Phủ, Đà Nẵng', 16.0678, 108.2120, 840, true),
    
    -- Route: TP. Hồ Chí Minh → Nha Trang (011269a3-48c9-4e25-85a7-a1f493a3e55f)
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', 1, 'pickup', 'Bến xe Miền Đông', '292 Đinh Bộ Lĩnh, Bình Thạnh, TP.HCM', 10.8161, 106.7138, 0, true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', 2, 'both', 'Bến xe Phan Thiết', '51 Trường Chinh, Phan Thiết, Bình Thuận', 10.9304, 108.1020, 180, true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', 3, 'both', 'Cam Ranh', 'TP. Cam Ranh, Khánh Hòa', 11.9215, 109.1590, 360, true),
    ('011269a3-48c9-4e25-85a7-a1f493a3e55f', 4, 'dropoff', 'Bến xe phía Nam Nha Trang', '23 Tháng 10, Vĩnh Hải, Nha Trang', 12.2584, 109.1822, 480, true),
    
    -- Route: Hà Nội → Hải Phòng (aaa145b8-7280-4bb2-8172-2e8ac43b9443)
    ('aaa145b8-7280-4bb2-8172-2e8ac43b9443', 1, 'pickup', 'Bến xe Gia Lâm', 'Quốc lộ 5, Gia Lâm, Hà Nội', 21.0381, 105.9321, 0, true),
    ('aaa145b8-7280-4bb2-8172-2e8ac43b9443', 2, 'pickup', 'Cầu Chui', 'Quốc lộ 5, Long Biên, Hà Nội', 21.0456, 105.9123, 15, true),
    ('aaa145b8-7280-4bb2-8172-2e8ac43b9443', 3, 'dropoff', 'Bến xe Niệm Nghĩa', '1 Cầu Rào II, Lê Chân, Hải Phòng', 20.8449, 106.6881, 120, true),
    
    -- Route: TP. Hồ Chí Minh → Cần Thơ (bd656d0c-3bac-4ac4-8c09-2fff22e3e90d)
    ('bd656d0c-3bac-4ac4-8c09-2fff22e3e90d', 1, 'pickup', 'Bến xe Miền Tây', 'Kinh Dương Vương, Bình Tân, TP.HCM', 10.7424, 106.6186, 0, true),
    ('bd656d0c-3bac-4ac4-8c09-2fff22e3e90d', 2, 'both', 'Bến xe Long An', 'Quốc lộ 1A, Tân An, Long An', 10.5359, 106.4037, 60, true),
    ('bd656d0c-3bac-4ac4-8c09-2fff22e3e90d', 3, 'both', 'Bến xe Vĩnh Long', 'TP. Vĩnh Long, Vĩnh Long', 10.2539, 105.9571, 120, true),
    ('bd656d0c-3bac-4ac4-8c09-2fff22e3e90d', 4, 'dropoff', 'Bến xe 91B Cần Thơ', '91B Nguyễn Văn Linh, Ninh Kiều, Cần Thơ', 10.0452, 105.7469, 180, true),
    
    -- Route: Hà Nội → Hạ Long (78fdd0a1-0a2c-498f-9293-d1fca7599c15)
    ('78fdd0a1-0a2c-498f-9293-d1fca7599c15', 1, 'pickup', 'Bến xe Mỹ Đình', '20 Phạm Hùng, Mỹ Đình, Nam Từ Liêm, Hà Nội', 21.0285, 105.7823, 0, true),
    ('78fdd0a1-0a2c-498f-9293-d1fca7599c15', 2, 'both', 'TP. Hải Dương', 'Bến xe Hải Dương, Hải Dương', 20.9406, 106.3234, 60, true),
    ('78fdd0a1-0a2c-498f-9293-d1fca7599c15', 3, 'dropoff', 'Bến xe Bãi Cháy', 'Hạ Long, Quảng Ninh', 20.9698, 107.0474, 180, true),
    
    -- Route: Đà Nẵng → TP. Hồ Chí Minh (5124fb9e-2fd4-436a-82fe-efc270a33b08)
    ('5124fb9e-2fd4-436a-82fe-efc270a33b08', 1, 'pickup', 'Bến xe Đà Nẵng', '33 Điện Biên Phủ, Đà Nẵng', 16.0678, 108.2120, 0, true),
    ('5124fb9e-2fd4-436a-82fe-efc270a33b08', 2, 'both', 'Bến xe Quy Nhơn', 'Bình Định', 13.7765, 109.2236, 360, true),
    ('5124fb9e-2fd4-436a-82fe-efc270a33b08', 3, 'both', 'Bến xe Nha Trang', 'Khánh Hòa', 12.2584, 109.1822, 600, true),
    ('5124fb9e-2fd4-436a-82fe-efc270a33b08', 4, 'dropoff', 'Bến xe Miền Đông', '292 Đinh Bộ Lĩnh, Bình Thạnh, TP.HCM', 10.8161, 106.7138, 960, true),
    
    -- Route: Hà Nội → Sapa (270cd14d-4bae-441d-9d69-d5ea14f8f86b)
    ('270cd14d-4bae-441d-9d69-d5ea14f8f86b', 1, 'pickup', 'Bến xe Mỹ Đình', '20 Phạm Hùng, Mỹ Đình, Nam Từ Liêm, Hà Nội', 21.0285, 105.7823, 0, true),
    ('270cd14d-4bae-441d-9d69-d5ea14f8f86b', 2, 'both', 'Bến xe Lào Cai', 'TP. Lào Cai, Lào Cai', 22.4856, 103.9707, 300, true),
    ('270cd14d-4bae-441d-9d69-d5ea14f8f86b', 3, 'dropoff', 'Bến xe Sapa', 'Thị trấn Sapa, Lào Cai', 22.3364, 103.8438, 360, true),
    
    -- Route: Hà Nội → Cao Bằng (fa703cfd-716c-4892-879c-4dc4e24b1074)
    ('fa703cfd-716c-4892-879c-4dc4e24b1074', 1, 'pickup', 'Bến xe Gia Lâm', 'Quốc lộ 5, Gia Lâm, Hà Nội', 21.0381, 105.9321, 0, true),
    ('fa703cfd-716c-4892-879c-4dc4e24b1074', 2, 'both', 'Bến xe Thái Nguyên', 'TP. Thái Nguyên, Thái Nguyên', 21.5942, 105.8481, 120, true),
    ('fa703cfd-716c-4892-879c-4dc4e24b1074', 3, 'both', 'Bến xe Bắc Kạn', 'TP. Bắc Kạn, Bắc Kạn', 22.1469, 105.8348, 240, true),
    ('fa703cfd-716c-4892-879c-4dc4e24b1074', 4, 'dropoff', 'Bến xe Cao Bằng', 'TP. Cao Bằng, Cao Bằng', 22.6657, 106.2622, 360, true)
ON CONFLICT DO NOTHING;

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
