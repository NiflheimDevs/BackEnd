CREATE TABLE "user_chat" (
  "chat_id" int,
  "user_id" int,
  FOREIGN KEY ("chat_id") REFERENCES "chat" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "user" ("id")
);