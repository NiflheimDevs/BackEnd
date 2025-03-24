CREATE TABLE IF NOT EXISTS "users_role" (
  "user_id" int NOT NULL,
  "role_id" int NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id") ON DELETE CASCADE
);