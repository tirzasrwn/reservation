ALTER TABLE room_restrictions
ADD CONSTRAINT room_restrictions_reservations_id_fk
FOREIGN KEY (reservation_id)
REFERENCES reservations(id)
ON DELETE CASCADE
ON UPDATE CASCADE;

CREATE INDEX idx_reservations_email
ON reservations (email);

CREATE INDEX idx_reservations_last_name
ON reservations (last_name);
