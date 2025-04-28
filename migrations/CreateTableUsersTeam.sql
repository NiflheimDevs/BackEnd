CREATE TABLE IF NOT EXISTS "users_team" (
  "user_id" int NOT NULL,
  "team_id" int NOT NULL,
  "position" varchar,
  "role_id" int,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("team_id") REFERENCES "team" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id")
);