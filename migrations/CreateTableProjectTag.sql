CREATE TABLE IF NOT EXISTS "project_tag" (
  "project_id" int NOT NULL,
  "tag_id" int NOT NULL,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("project_id") REFERENCES "project" ("id") ON DELETE CASCADE,
  UNIQUE ("tag_id", "project_id")
);