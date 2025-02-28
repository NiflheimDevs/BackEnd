CREATE TABLE "bid" (
  "id" serial PRIMARY KEY,
  "user_team_id" int,
  "projet_id" int,
  "value" numeric,
  "expected_time" datetime,
  "type" int,
  FOREIGN KEY ("projet_id") REFERENCES "project" ("id"),
  FOREIGN KEY ("user_team_id") REFERENCES "user" ("id"),
  FOREIGN KEY ("user_team_id") REFERENCES "team" ("id"),
  FOREIGN KEY ("id") REFERENCES "project" ("selected_bid_id")
);