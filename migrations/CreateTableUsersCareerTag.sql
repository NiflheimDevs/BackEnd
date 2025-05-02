CREATE TABLE IF NOT EXISTS "users_career_tag" (
  "career_user_id" int NOT NULL,
  "tag_id" int NOT NULL,
  "type" int NOT NULL,
  "level" int DEFAULT -1,
  FOREIGN KEY ("tag_id") REFERENCES "tag" ("id"),
  FOREIGN KEY ("career_user_id") REFERENCES "career" ("id"),
  FOREIGN KEY ("career_user_id") REFERENCES "users" ("id"),
  UNIQUE ("career_user_id" , "tag_id" , "type")
);

CREATE OR REPLACE FUNCTION enforce_fk_constraint_on_users_career_tag() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.type = 1 THEN
        IF NOT EXISTS (SELECT 1 FROM career WHERE id = NEW.career_user_id) THEN
            RAISE EXCEPTION 'Invalid reference: % does not exist in career', NEW.career_user_id;
        END IF;
    ELSIF NEW.type = 0 THEN
        IF NOT EXISTS (SELECT 1 FROM users WHERE id = NEW.career_user_id) THEN
            RAISE EXCEPTION 'Invalid reference: % does not exist in users', NEW.career_user_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_fk_users_career_tag
BEFORE INSERT OR UPDATE ON users_career_tag
FOR EACH ROW EXECUTE FUNCTION enforce_fk_constraint_on_users_career_tag();