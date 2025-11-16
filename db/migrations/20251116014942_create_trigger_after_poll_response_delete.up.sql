CREATE TRIGGER after_poll_response_delete
    AFTER DELETE ON poll_responses
    FOR EACH ROW
BEGIN
    UPDATE poll_options
    SET vote_count = vote_count - 1
    WHERE id = OLD.poll_option_id;
END;
