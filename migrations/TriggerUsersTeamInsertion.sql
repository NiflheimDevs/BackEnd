CREATE FUNCTION unique_check_users_team_insertion() RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM users_team
        WHERE user_id = NEW.user_id
          AND team_id = NEW.team_id
          AND left_at IS NULL
    ) THEN
        RAISE EXCEPTION 'User % is already in team %', NEW.user_id, NEW.team_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER unique_check_users_team_insertion_trigger
BEFORE INSERT ON users_team
FOR EACH ROW EXECUTE FUNCTION unique_check_users_team_insertion();