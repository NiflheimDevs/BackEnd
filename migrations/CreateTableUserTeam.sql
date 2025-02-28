CREATE TABLE "user_team" (
  "user_id" int,
  "team_id" int,
  FOREIGN KEY ("user_id") REFERENCES "user" ("id"),
  FOREIGN KEY ("team_id") REFERENCES "team" ("id")
);