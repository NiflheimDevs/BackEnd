CREATE TABLE "user_team" (
  "user_id" int NOT NULL,
  "team_id" int NOT NULL,
  "position" varchar,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id"),
  FOREIGN KEY ("team_id") REFERENCES "team" ("id")
);