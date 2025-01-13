CREATE SCHEMA IF NOT EXISTS msg;

CREATE TABLE IF NOT EXISTS msg.user (
  id SERIAL PRIMARY KEY,
  nickname VARCHAR(64) UNIQUE NOT NULL,
  password_hash VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS msg.friendship (
  user_id1 INT NOT NULL REFERENCES msg.user(id) ON DELETE CASCADE,
  user_id2 INT NOT NULL REFERENCES msg.user(id) ON DELETE CASCADE,
  user1 INT GENERATED ALWAYS AS (LEAST(user_id1, user_id2)) STORED,
  user2 INT GENERATED ALWAYS AS (GREATEST(user_id1, user_id2)) STORED,
  PRIMARY KEY (user1, user2),
  CHECK (user_id1 != user_id2)
);

CREATE TABLE IF NOT EXISTS msg.friendship_invitation (
  from_user_id INT NOT NULL REFERENCES msg.user(id) ON DELETE CASCADE,
  to_user_id INT NOT NULL REFERENCES msg.user(id) ON DELETE CASCADE,
  user1 INT GENERATED ALWAYS AS (LEAST(from_user_id, to_user_id)) STORED,
  user2 INT GENERATED ALWAYS AS (GREATEST(from_user_id, to_user_id)) STORED,
  PRIMARY KEY (user1, user2),
  CHECK (from_user_id != to_user_id)
);

CREATE TABLE IF NOT EXISTS msg.chat (
  id SERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS msg.chat_membership (
  chat_id INT NOT NULL REFERENCES msg.chat(id) ON DELETE CASCADE,
  user_id INT NOT NULL REFERENCES msg.user(id) ON DELETE CASCADE,
  PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS msg.chat_invitation (
  chat_id INT NOT NULL REFERENCES msg.chat(id) ON DELETE CASCADE,
  user_id INT NOT NULL REFERENCES msg.user(id) ON DELETE CASCADE,
  PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS msg.message (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL REFERENCES msg.user(id) ON DELETE CASCADE,
  chat_id INT NOT NULL REFERENCES msg.chat(id) ON DELETE CASCADE,
  msg_text TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON msg.message(chat_id);
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON msg.message(user_id);
