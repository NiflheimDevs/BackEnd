CREATE TABLE IF NOT EXISTS "user_role" (
  "user_id" int NOT NULL,
  "role_id" int NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id"),
  FOREIGN KEY ("role_id") REFERENCES "role" ("id")
);