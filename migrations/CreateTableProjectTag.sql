CREATE TABLE IF NOT EXISTS "project_tag" (
  "project_id" int NOT NULL,
  "tag_id" int NOT NULL,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
  FOREIGN KEY ("project_id") REFERENCES "project" ("id"),
  UNIQUE ("tag_id", "project_id")
);