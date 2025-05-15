CREATE FUNCTION delete_cascade_for_career() RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM users_career_tag AS uct WHERE uct.career_user_id = OLD.id AND uct.type = 1;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER career_delete_trigger
BEFORE DELETE ON career
FOR EACH ROW EXECUTE FUNCTION delete_cascade_for_career();

CREATE FUNCTION delete_cascade_for_user() RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM users_career_tag AS uct WHERE uct.career_user_id = OLD.id AND uct.type = 0;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_delete_trigger
BEFORE DELETE ON users
FOR EACH ROW EXECUTE FUNCTION delete_cascade_for_user();

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

CREATE FUNCTION unique_check_users_team_insertion() RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM "users_team"
        WHERE "user_id" = NEW.user_id
          AND "team_id" = NEW.team_id
          AND "left_at" IS NULL
    ) THEN
        RAISE EXCEPTION 'User % is already in team %', NEW.user_id, NEW.team_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER unique_check_users_team_insertion_trigger
BEFORE INSERT ON users_team
FOR EACH ROW EXECUTE FUNCTION unique_check_users_team_insertion();