CREATE TABLE IF NOT EXISTS "user_chat" (
  "chat_id" int NOT NULL,
  "user_id" int NOT NULL,
  "role_id" int NOT NULL,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);