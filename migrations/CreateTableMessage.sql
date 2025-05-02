CREATE TABLE IF NOT EXISTS "message" (
  "id" serial PRIMARY KEY,
  "chat_id" int NOT NULL,
  "sender_id" int,
  "sent_time" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  "edit_time" TIMESTAMP WITH TIME ZONE,
  "content" text NOT NULL,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id"),
  FOREIGN KEY ("sender_id") REFERENCES "users" ("id")
);