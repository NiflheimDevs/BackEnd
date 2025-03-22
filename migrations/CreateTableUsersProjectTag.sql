CREATE TABLE IF NOT EXISTS "users_career_tag" (
  "career_user_id" int NOT NULL,
  "tag_id" int NOT NULL,
  "type" int NOT NULL,
  "level" int,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
  FOREIGN KEY ("career_user_id") REFERENCES "career" ("id"),
  FOREIGN KEY ("career_user_id") REFERENCES "users" ("id")
);