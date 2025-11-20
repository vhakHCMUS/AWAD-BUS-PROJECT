-- Seed data for development and testing

-- Insert admin user
INSERT INTO users (email, name, role, password_hash, is_active) VALUES
('admin@busbooking.com', 'System Admin', 'admin', '$2a$10$xQPQkjrXkOxGMXkzWZxs6eH7vZSNJ5d7lKqF8YqnQjGxQ4yxJxQ4G', true),
('passenger@example.com', 'John Doe', 'passenger', '$2a$10$xQPQkjrXkOxGMXkzWZxs6eH7vZSNJ5d7lKqF8YqnQjGxQ4yxJxQ4G', true);

-- Insert buses with seat layouts
INSERT INTO buses (license_plate, bus_type, manufacturer, model, year, seat_layout, amenities, status) VALUES
('29A-12345', 'sleeper', 'Thaco', 'Universe', 2023, 
 '{"rows": 12, "columns": 4, "total_seats": 40, "floors": 2, "layout": [["A1", "A2", "aisle", "A3", "A4"], ["B1", "B2", "aisle", "B3", "B4"]]}',
 ARRAY['wifi', 'ac', 'charging', 'blanket'], 'active'),
 
('30B-67890', 'seater', 'Hyundai', 'Universe Noble', 2022,
 '{"rows": 11, "columns": 4, "total_seats": 45, "floors": 1, "layout": [["A1", "A2", "aisle", "A3", "A4"]]}',
 ARRAY['wifi', 'ac', 'recliner'], 'active'),
 
('51C-11111', 'semi-sleeper', 'Thaco', 'TB120S', 2023,
 '{"rows": 10, "columns": 4, "total_seats": 36, "floors": 1, "layout": [["A1", "A2", "aisle", "A3", "A4"]]}',
 ARRAY['wifi', 'ac', 'usb_charging'], 'active');

-- Insert routes
INSERT INTO routes (name, from_city, to_city, distance, base_price, description, is_active) VALUES
('Hanoi - Ho Chi Minh', 'Hanoi', 'Ho Chi Minh City', 1710, 450000, 'Express route with minimal stops', true),
('Hanoi - Da Nang', 'Hanoi', 'Da Nang', 764, 300000, 'Scenic coastal route', true),
('Ho Chi Minh - Nha Trang', 'Ho Chi Minh City', 'Nha Trang', 448, 200000, 'Popular beach destination', true),
('Da Nang - Hoi An', 'Da Nang', 'Hoi An', 30, 50000, 'Short route to ancient town', true),
('Hanoi - Halong', 'Hanoi', 'Halong', 165, 120000, 'Popular tourist destination', true);

-- Insert trips (upcoming schedules)
INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id,
    r.id,
    NOW() + INTERVAL '1 day' + INTERVAL '6 hours',
    NOW() + INTERVAL '1 day' + INTERVAL '36 hours',
    1800, -- 30 hours
    500000,
    'scheduled',
    'Nguyen Van A',
    '0912345678'
FROM buses b, routes r
WHERE b.license_plate = '29A-12345' AND r.name = 'Hanoi - Ho Chi Minh'
LIMIT 1;

INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id,
    r.id,
    NOW() + INTERVAL '2 days' + INTERVAL '8 hours',
    NOW() + INTERVAL '2 days' + INTERVAL '22 hours',
    840, -- 14 hours
    320000,
    'scheduled',
    'Tran Van B',
    '0987654321'
FROM buses b, routes r
WHERE b.license_plate = '30B-67890' AND r.name = 'Hanoi - Da Nang'
LIMIT 1;

INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id,
    r.id,
    NOW() + INTERVAL '1 day' + INTERVAL '20 hours',
    NOW() + INTERVAL '2 days' + INTERVAL '5 hours',
    540, -- 9 hours
    220000,
    'scheduled',
    'Le Van C',
    '0901234567'
FROM buses b, routes r
WHERE b.license_plate = '51C-11111' AND r.name = 'Ho Chi Minh - Nha Trang'
LIMIT 1;

-- Initialize seats for all trips
INSERT INTO seats_status (trip_id, seat_number, status)
SELECT 
    t.id,
    'A' || generate_series(1, 40),
    'available'
FROM trips t
WHERE EXISTS (
    SELECT 1 FROM buses b 
    WHERE b.id = t.bus_id 
    AND b.license_plate = '29A-12345'
);

INSERT INTO seats_status (trip_id, seat_number, status)
SELECT 
    t.id,
    'S' || generate_series(1, 45),
    'available'
FROM trips t
WHERE EXISTS (
    SELECT 1 FROM buses b 
    WHERE b.id = t.bus_id 
    AND b.license_plate = '30B-67890'
);

INSERT INTO seats_status (trip_id, seat_number, status)
SELECT 
    t.id,
    'B' || generate_series(1, 36),
    'available'
FROM trips t
WHERE EXISTS (
    SELECT 1 FROM buses b 
    WHERE b.id = t.bus_id 
    AND b.license_plate = '51C-11111'
);

-- Sample booking (for testing)
INSERT INTO bookings (trip_id, user_id, contact_email, contact_phone, contact_name, seats, total_price, status, booking_code)
SELECT 
    t.id,
    u.id,
    'passenger@example.com',
    '0912345678',
    'John Doe',
    ARRAY['A1', 'A2'],
    1000000,
    'confirmed',
    'BK' || EXTRACT(EPOCH FROM NOW())::BIGINT
FROM trips t, users u
WHERE u.email = 'passenger@example.com'
LIMIT 1;

-- Mark sample booked seats
UPDATE seats_status
SET status = 'booked', booking_id = (SELECT id FROM bookings LIMIT 1)
WHERE trip_id = (SELECT id FROM trips LIMIT 1)
AND seat_number IN ('A1', 'A2');

-- Create sample tickets
INSERT INTO tickets (booking_id, passenger_name, passenger_email, passenger_phone, seat_number, ticket_code)
SELECT 
    b.id,
    'John Doe',
    'passenger@example.com',
    '0912345678',
    unnest(ARRAY['A1', 'A2']),
    'TK' || EXTRACT(EPOCH FROM NOW())::BIGINT || generate_series(1, 2)
FROM bookings b
LIMIT 1;
