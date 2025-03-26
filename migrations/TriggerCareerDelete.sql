CREATE FUNCTION delete_cascade_for_career() RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM users_career_tag AS uct WHERE uct.career_user_id = OLD.id AND uct.type = 1;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER career_delete_trigger
AFTER DELETE ON career
FOR EACH ROW EXECUTE FUNCTION delete_cascade_for_career();