CREATE TABLE "message" (
  "id" serial PRIMARY KEY,
  "chat_id" int,
  "sender_id" int,
  "content" text,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id"),
  FOREIGN KEY ("sender_id") REFERENCES "user" ("id")
);