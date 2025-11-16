CREATE TRIGGER after_poll_response_insert
    AFTER INSERT ON poll_responses
    FOR EACH ROW
BEGIN
    UPDATE poll_options
    SET vote_count = vote_count + 1
    WHERE id = NEW.poll_option_id;
END;
