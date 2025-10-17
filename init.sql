CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    number INT  NOT NULL UNIQUE,
    room_count INT DEFAULT 1,
    is_occupied BOOLEAN DEFAULT FALSE,
    floor INT NOT NULL,
    sleeping_places INT NOT NULL DEFAULT 1,
    room_type  NOT NULL VARCHAR(50) CHECK (room_type IN ('Standard', 'Deluxe', 'Suite')),
    need_cleaning BOOLEAN DEFAULT FALSE
);  

INSERT INTO rooms (number, room_count, is_occupied, floor, sleeping_places, room_type, need_cleaning)
VALUES
    (1, 1, FALSE, 1, 2, 'Standard', FALSE),
    (2, 1, TRUE, 2, 4, 'Deluxe', TRUE),
    (3, 1, FALSE, 3, 3, 'Suite', FALSE);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    room_id INT NOT NULL UNIQUE,
    guest_name VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    arrival_status BOOLEAN DEFAULT FALSE,
    CHECK (start_date < end_date)
);

INSERT INTO bookings (guest_name, start_date, end_date, arrival_status)
VALUES
    ('John', 1, '2025-10-17', '2025-11-17', true),
    ('Ann', 2, '2030-10-31', '2030-11-20', false),
    ('Maria', 3, '2025-12-20', '2026-01-11', false);