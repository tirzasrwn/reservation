CREATE TABLE room_restrictions (
  id SERIAL PRIMARY KEY,
  start_date DATE,
  end_date DATE,
  room_id INTEGER,
  reservation_id INTEGER,
  restriction_id INTEGER,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
