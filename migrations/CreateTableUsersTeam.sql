CREATE TABLE IF NOT EXISTS "users_team" (
  "user_id" int NOT NULL,
  "team_id" int NOT NULL,
  "position" varchar,
  "role_id" int NOT NULL,
  "joined_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "left_at" timestamp with time zone,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
  FOREIGN KEY ("team_id") REFERENCES "team" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id")
);