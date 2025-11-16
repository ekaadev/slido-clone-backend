CREATE TRIGGER after_vote_insert
    AFTER INSERT ON votes
    FOR EACH ROW
BEGIN
    UPDATE questions
    SET upvote_count = upvote_count + 1
    WHERE id = NEW.question_id;
END;
