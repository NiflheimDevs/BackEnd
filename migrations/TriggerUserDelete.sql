CREATE FUNCTION delete_cascade_for_user() RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM users_career_tag AS ust WHERE uct.career_user_id = OLD.id AND uct.type = 0;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_delete_trigger
BEFORE DELETE ON users
FOR EACH ROW EXECUTE FUNCTION delete_cascade_for_user();