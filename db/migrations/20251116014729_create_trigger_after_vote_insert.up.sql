CREATE OR REPLACE FUNCTION fn_after_vote_insert()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE questions
    SET upvote_count = upvote_count + 1
    WHERE id = NEW.question_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_vote_insert
    AFTER INSERT ON votes
    FOR EACH ROW
    EXECUTE FUNCTION fn_after_vote_insert();
