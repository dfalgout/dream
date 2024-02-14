
-- +migrate Up
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  full_name TEXT,
  is_admin BOOLEAN DEFAULT FALSE NOT NULL,
  verify_code TEXT,
  verify_expires_at DATETIME,
  created_at DATETIME DEFAULT (strftime('%F %R:%f')) NOT NULL,
  updated_at DATETIME DEFAULT (strftime('%F %R:%f')) NOT NULL
);

CREATE INDEX users_email_idx ON users(email);

-- +migrate Down
DROP TABLE users;
