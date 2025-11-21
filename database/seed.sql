-- Seed data for VietBusBooking - Vietnamese intercity bus booking system
-- This includes real Vietnamese cities, popular bus operators, and realistic routes

-- Insert admin user and test passengers
INSERT INTO users (email, name, phone, role, password_hash, is_active) VALUES
('admin@vietbusbooking.com', 'Quản trị viên hệ thống', '0901234567', 'admin', '$2a$10$xQPQkjrXkOxGMXkzWZxs6eH7vZSNJ5d7lKqF8YqnQjGxQ4yxJxQ4G', true),
('nguyenvana@gmail.com', 'Nguyễn Văn A', '0912345678', 'passenger', '$2a$10$xQPQkjrXkOxGMXkzWZxs6eH7vZSNJ5d7lKqF8YqnQjGxQ4yxJxQ4G', true),
('tranthib@gmail.com', 'Trần Thị B', '0987654321', 'passenger', '$2a$10$xQPQkjrXkOxGMXkzWZxs6eH7vZSNJ5d7lKqF8YqnQjGxQ4yxJxQ4G', true);

-- Insert buses from popular Vietnamese operators with realistic seat layouts
INSERT INTO buses (license_plate, bus_type, manufacturer, model, year, operator_name, seat_layout, amenities, status) VALUES
-- Phương Trang (Futa Bus Lines) - Premium Limousine 29 seats
('51F-123.45', 'limousine', 'Mercedes-Benz', 'Sprinter 519', 2024, 'Phương Trang FUTA', 
 '{"rows": 7, "columns": 4, "total_seats": 29, "floors": 2, "type": "29_limousine"}',
 ARRAY['wifi', 'ac', 'massage_seat', 'blanket', 'entertainment', 'water', 'usb_charging'], 'active'),

-- Mai Linh Express - Sleeper 45 beds
('29M-678.90', 'sleeper', 'Thaco', 'Mobihome MB', 2023, 'Mai Linh Express',
 '{"rows": 12, "columns": 4, "total_seats": 45, "floors": 2, "type": "45_sleeper"}',
 ARRAY['wifi', 'ac', 'blanket', 'pillow', 'water', 'usb_charging'], 'active'),

-- Thành Bưởi - Seater 40 seats
('50T-111.22', 'seater', 'Hyundai', 'Universe Noble', 2023, 'Thành Bưởi',
 '{"rows": 10, "columns": 4, "total_seats": 40, "floors": 1, "type": "40_seater"}',
 ARRAY['wifi', 'ac', 'recliner', 'usb_charging'], 'active'),

-- Hoàng Long - Limousine 16 seats
('51H-333.44', 'limousine', 'Ford', 'Transit Limousine', 2024, 'Hoàng Long',
 '{"rows": 4, "columns": 4, "total_seats": 16, "floors": 1, "type": "16_limousine"}',
 ARRAY['wifi', 'ac', 'massage_seat', 'entertainment', 'water', 'snack', 'usb_charging'], 'active'),

-- Hưng Thành - Semi-sleeper 36 seats
('29H-555.66', 'semi-sleeper', 'Thaco', 'TB120S', 2022, 'Hưng Thành',
 '{"rows": 9, "columns": 4, "total_seats": 36, "floors": 1, "type": "36_semi_sleeper"}',
 ARRAY['wifi', 'ac', 'blanket', 'water'], 'active');

-- Insert major Vietnamese routes with realistic distances and prices
INSERT INTO routes (name, from_city, to_city, distance, base_price, description, is_active) VALUES
-- North to South routes
('Hà Nội - TP. Hồ Chí Minh', 'Hà Nội', 'TP. Hồ Chí Minh', 1710, 450000, 'Tuyến xuyên Việt - Giường nằm cao cấp, ít điểm dừng', true),
('Hà Nội - Đà Nẵng', 'Hà Nội', 'Đà Nẵng', 764, 280000, 'Tuyến ven biển - Xe limousine đời mới', true),
('Hà Nội - Vinh', 'Hà Nội', 'Vinh', 319, 180000, 'Tuyến Bắc Trung - Ghế ngồi thương gia', true),
('Hà Nội - Huế', 'Hà Nội', 'Huế', 688, 250000, 'Di sản văn hóa - Xe giường nằm', true),
('Hà Nội - Nha Trang', 'Hà Nội', 'Nha Trang', 1278, 380000, 'Đến biển Nha Trang - Giường nằm cao cấp', true),

-- Southern routes
('TP. Hồ Chí Minh - Đà Lạt', 'TP. Hồ Chí Minh', 'Đà Lạt', 308, 150000, 'Lên thành phố ngàn hoa - Limousine', true),
('TP. Hồ Chí Minh - Nha Trang', 'TP. Hồ Chí Minh', 'Nha Trang', 448, 200000, 'Biển Nha Trang - Giường nằm', true),
('TP. Hồ Chí Minh - Vũng Tàu', 'TP. Hồ Chí Minh', 'Vũng Tàu', 125, 100000, 'Biển gần Sài Gòn - Limousine', true),
('TP. Hồ Chí Minh - Cần Thơ', 'TP. Hồ Chí Minh', 'Cần Thơ', 169, 120000, 'Miền Tây sông nước - Xe ghế ngồi', true),
('TP. Hồ Chí Minh - Phú Quốc', 'TP. Hồ Chí Minh', 'Phú Quốc', 520, 280000, 'Đảo ngọc - Xe + phà', true),

-- Central routes
('Đà Nẵng - Hội An', 'Đà Nẵng', 'Hội An', 30, 50000, 'Phố cổ Hội An - Xe limousine', true),
('Đà Nẵng - Huế', 'Đà Nẵng', 'Huế', 95, 80000, 'Cố đô Huế - Ghế ngồi', true),
('Đà Nẵng - Quy Nhơn', 'Đà Nẵng', 'Quy Nhơn', 303, 150000, 'Biển Quy Nhơn - Giường nằm', true),

-- Northern routes
('Hà Nội - Hạ Long', 'Hà Nội', 'Hạ Long', 165, 120000, 'Vịnh Hạ Long - Limousine', true),
('Hà Nội - Sapa', 'Hà Nội', 'Sapa', 376, 200000, 'Thị trấn mây - Giường nằm', true),
('Hà Nội - Ninh Bình', 'Hà Nội', 'Ninh Bình', 95, 80000, 'Tràng An - Ghế ngồi', true),
('Hà Nội - Hải Phòng', 'Hà Nội', 'Hải Phòng', 120, 90000, 'Cảng Hải Phòng - Limousine', true),

-- Cross-country routes  
('Đà Lạt - Nha Trang', 'Đà Lạt', 'Nha Trang', 140, 100000, 'Núi xuống biển - Limousine', true),
('Huế - Hội An', 'Huế', 'Hội An', 125, 90000, 'Di sản miền Trung - Ghế ngồi', true),
('Nha Trang - Quy Nhơn', 'Nha Trang', 'Quy Nhơn', 238, 130000, 'Duyên hải miền Trung - Giường nằm', true);

-- Insert realistic trip schedules (next 7 days)
-- Hà Nội - TP. Hồ Chí Minh (Premium Limousine - Phương Trang)
INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id, r.id,
    CURRENT_DATE + INTERVAL '1 day' + INTERVAL '20 hours',
    CURRENT_DATE + INTERVAL '2 days' + INTERVAL '26 hours',
    1800, -- 30 hours
    550000,
    'scheduled',
    'Phạm Văn Đức',
    '0901234567'
FROM buses b, routes r
WHERE b.license_plate = '51F-123.45' AND r.name = 'Hà Nội - TP. Hồ Chí Minh';

-- Hà Nội - Đà Nẵng (Sleeper - Mai Linh)
INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id, r.id,
    CURRENT_DATE + INTERVAL '1 day' + INTERVAL '21 hours',
    CURRENT_DATE + INTERVAL '2 days' + INTERVAL '12 hours',
    900, -- 15 hours
    320000,
    'scheduled',
    'Lê Minh Tuấn',
    '0912345678'
FROM buses b, routes r
WHERE b.license_plate = '29M-678.90' AND r.name = 'Hà Nội - Đà Nẵng';

-- TP. HCM - Đà Lạt (Limousine - Hoàng Long)
INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id, r.id,
    CURRENT_DATE + INTERVAL '1 day' + INTERVAL '7 hours',
    CURRENT_DATE + INTERVAL '1 day' + INTERVAL '14 hours',
    420, -- 7 hours
    180000,
    'scheduled',
    'Nguyễn Hoàng Long',
    '0923456789'
FROM buses b, routes r
WHERE b.license_plate = '51H-333.44' AND r.name = 'TP. Hồ Chí Minh - Đà Lạt';

-- TP. HCM - Vũng Tàu (Seater - Thành Bưởi)
INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id, r.id,
    CURRENT_DATE + INTERVAL '1 day' + INTERVAL '6 hours',
    CURRENT_DATE + INTERVAL '1 day' + INTERVAL '8 hours' + INTERVAL '30 minutes',
    150, -- 2.5 hours
    120000,
    'scheduled',
    'Trần Văn Bình',
    '0934567890'
FROM buses b, routes r
WHERE b.license_plate = '50T-111.22' AND r.name = 'TP. Hồ Chí Minh - Vũng Tàu';

-- Đà Nẵng - Hội An (Limousine - Hoàng Long)
INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id, r.id,
    CURRENT_DATE + INTERVAL '1 day' + INTERVAL '8 hours',
    CURRENT_DATE + INTERVAL '1 day' + INTERVAL '9 hours',
    60, -- 1 hour
    60000,
    'scheduled',
    'Võ Văn Cường',
    '0945678901'
FROM buses b, routes r
WHERE b.license_plate = '51H-333.44' AND r.name = 'Đà Nẵng - Hội An';

-- Hà Nội - Hạ Long (Semi-sleeper - Hưng Thành)
INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id, r.id,
    CURRENT_DATE + INTERVAL '2 days' + INTERVAL '7 hours',
    CURRENT_DATE + INTERVAL '2 days' + INTERVAL '10 hours' + INTERVAL '30 minutes',
    210, -- 3.5 hours
    140000,
    'scheduled',
    'Đỗ Văn Em',
    '0956789012'
FROM buses b, routes r
WHERE b.license_plate = '29H-555.66' AND r.name = 'Hà Nội - Hạ Long';

-- TP. HCM - Nha Trang (Sleeper - Mai Linh)
INSERT INTO trips (bus_id, route_id, departure_time, arrival_time, duration, price, status, driver_name, driver_phone)
SELECT 
    b.id, r.id,
    CURRENT_DATE + INTERVAL '2 days' + INTERVAL '22 hours',
    CURRENT_DATE + INTERVAL '3 days' + INTERVAL '7 hours',
    540, -- 9 hours
    230000,
    'scheduled',
    'Phan Văn Phúc',
    '0967890123'
FROM buses b, routes r
WHERE b.license_plate = '29M-678.90' AND r.name = 'TP. Hồ Chí Minh - Nha Trang';

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
