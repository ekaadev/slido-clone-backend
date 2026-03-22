CREATE OR REPLACE FUNCTION fn_after_poll_response_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE poll_options
    SET vote_count = vote_count - 1
    WHERE id = OLD.poll_option_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_poll_response_delete
    AFTER DELETE ON poll_responses
    FOR EACH ROW
    EXECUTE FUNCTION fn_after_poll_response_delete();
