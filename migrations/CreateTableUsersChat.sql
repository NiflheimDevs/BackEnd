CREATE TABLE IF NOT EXISTS "users_chat" (
  "chat_id" int NOT NULL,
  "user_id" int NOT NULL,
  "role_id" int NOT NULL,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
  UNIQUE ("chat_id", "user_id")
);