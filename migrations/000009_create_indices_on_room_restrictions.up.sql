CREATE INDEX room_restrictions_start_end_idx ON room_restrictions (start_date, end_date);
CREATE INDEX room_restrictions_room_id_idx ON room_restrictions (room_id);
CREATE INDEX room_restrictions_reservation_id_idx ON room_restrictions (reservation_id);
