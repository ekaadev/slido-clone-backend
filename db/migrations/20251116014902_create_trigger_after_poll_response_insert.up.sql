CREATE OR REPLACE FUNCTION fn_after_poll_response_insert()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE poll_options
    SET vote_count = vote_count + 1
    WHERE id = NEW.poll_option_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_poll_response_insert
    AFTER INSERT ON poll_responses
    FOR EACH ROW
    EXECUTE FUNCTION fn_after_poll_response_insert();
