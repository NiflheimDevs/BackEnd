--tyep: 0 for user, 1 for team and 2 for group.
CREATE TABLE IF NOT EXISTS "member_role" (
  "user_id" int NOT NULL,
  "origin_id" int,
  "type" int NOT NULL,
  "role_id" int NOT NULL,
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("role_id") REFERENCES "role" ("id"),
  UNIQUE ("user_id" , "origin_id" , "type")
);

CREATE OR REPLACE FUNCTION enforce_fk_constraint_on_member_role() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.type = 1 THEN
        IF NOT EXISTS (SELECT 1 FROM "team" WHERE "id" = NEW.origin_id) THEN
            RAISE EXCEPTION 'Invalid reference: % does not exist in team', NEW.origin_id;
        END IF;
    ELSIF NEW.type = 2 THEN
        IF NOT EXISTS (SELECT 1 FROM chat WHERE id = NEW.origin_id) THEN
            RAISE EXCEPTION 'Invalid reference: % does not exist in chat', NEW.origin_id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_fk_member_role
BEFORE INSERT OR UPDATE member_role
FOR EACH ROW EXECUTE FUNCTION enforce_fk_constraint_on_member_role();