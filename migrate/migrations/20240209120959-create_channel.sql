
-- +migrate Up
CREATE TABLE channels (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at DATETIME DEFAULT (strftime('%F %R:%f')) NOT NULL,
  updated_at DATETIME DEFAULT (strftime('%F %R:%f')) NOT NULL
);

CREATE INDEX channels_name_idx ON channels(name);

CREATE TABLE channel_users (
  user_id TEXT NOT NULL,
  channel_id TEXT NOT NULL,
  created_at DATETIME DEFAULT (strftime('%F %R:%f')) NOT NULL,

  PRIMARY KEY (user_id, channel_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (channel_id) REFERENCES channels(id)
);

-- +migrate Down
DROP TABLE channel_users;
DROP TABLE channels;
