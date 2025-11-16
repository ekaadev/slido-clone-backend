CREATE TRIGGER after_vote_delete
    AFTER DELETE ON votes
    FOR EACH ROW
BEGIN
    UPDATE questions
    SET upvote_count = upvote_count - 1
    WHERE id = OLD.question_id;
END;
