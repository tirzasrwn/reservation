ALTER TABLE room_restrictions
DROP CONSTRAINT room_restrictions_reservations_id_fk;

DROP INDEX IF EXISTS reservations.reservations_email_idx;
DROP INDEX IF EXISTS reservations.reservations_last_name_idx;
