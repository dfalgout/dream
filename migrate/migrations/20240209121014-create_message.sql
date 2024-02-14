
-- +migrate Up
CREATE TABLE messages (
  id TEXT PRIMARY KEY,
  message TEXT NOT NULL,
  user_id TEXT NOT NULL,
  channel_id TEXT NOT NULL,
  created_at DATETIME DEFAULT (strftime('%F %R:%f')) NOT NULL,

  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (channel_id) REFERENCES channels(id)
);

-- +migrate Down
DROP TABLE messages;
