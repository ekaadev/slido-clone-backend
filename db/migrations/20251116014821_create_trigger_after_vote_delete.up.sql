CREATE OR REPLACE FUNCTION fn_after_vote_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE questions
    SET upvote_count = upvote_count - 1
    WHERE id = OLD.question_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_vote_delete
    AFTER DELETE ON votes
    FOR EACH ROW
    EXECUTE FUNCTION fn_after_vote_delete();
