CREATE TABLE "user_project_tag" (
  "project_user_id" int,
  "tag_id" int,
  "type" int,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
  FOREIGN KEY ("project_user_id") REFERENCES "project" ("id"),
  FOREIGN KEY ("project_user_id") REFERENCES "user" ("id")
);