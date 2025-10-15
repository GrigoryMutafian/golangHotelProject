CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    number INT  NOT NULL UNIQUE,
    room_count INT,
    is_occupied BOOLEAN DEFAULT FALSE,
    floor INT,
    sleeping_places INT,
    room_type VARCHAR(50) CHECK (room_type IN ('Standard', 'Deluxe', 'Suite')),
    need_cleaning BOOLEAN DEFAULT FALSE
);  

INSERT INTO rooms (number, room_count, is_occupied, floor, sleeping_places, room_type, need_cleaning)
VALUES
    (1, 1, FALSE, 1, 2, 'Standard', FALSE),
    (2, 1, TRUE, 2, 4, 'Deluxe', TRUE),
    (3, 1, FALSE, 3, 3, 'Suite', FALSE);