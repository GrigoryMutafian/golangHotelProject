CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    number INT NOT NULL UNIQUE,
    room_count INT NOT NULL DEFAULT 1,
    is_occupied BOOLEAN NOT NULL DEFAULT FALSE,
    floor INT NOT NULL,
    sleeping_places INT NOT NULL DEFAULT 1,
    room_type VARCHAR(50) NOT NULL CHECK (room_type IN ('Standard', 'Deluxe', 'Suite')),
    need_cleaning BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO rooms (number, room_count, is_occupied, floor, sleeping_places, room_type, need_cleaning)
VALUES
    (1, 1, FALSE, 1, 2, 'Standard', FALSE),
    (2, 1, TRUE, 2, 4, 'Deluxe', TRUE),
    (3, 1, FALSE, 3, 3, 'Suite', FALSE);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    room_id INT NOT NULL,
    guest_id INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status BOOLEAN NOT NULL DEFAULT FALSE,
    CHECK (start_date < end_date),
    CONSTRAINT fk_bookings_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

INSERT INTO bookings (room_id, guest_id, start_date, end_date, status)
VALUES
    (1, 1,  '2025-10-17', '2025-11-17', TRUE),
    (2, 2,   '2030-10-31', '2030-11-20', FALSE),
    (3, 3, '2025-12-20', '2026-01-11', FALSE);
