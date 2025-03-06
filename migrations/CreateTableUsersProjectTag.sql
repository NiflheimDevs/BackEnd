CREATE TABLE IF NOT EXISTS "users_project_tag" (
  "project_user_id" int NOT NULL,
  "tag_id" int NOT NULL,
  "type" int NOT NULL,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
  FOREIGN KEY ("project_user_id") REFERENCES "project" ("id"),
  FOREIGN KEY ("project_user_id") REFERENCES "users" ("id")
);