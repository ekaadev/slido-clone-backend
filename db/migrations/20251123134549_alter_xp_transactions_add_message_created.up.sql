ALTER TABLE xp_transactions
    MODIFY COLUMN source_type
    ENUM('poll','question_created','upvote_received','presenter_validated','message_created')
    NOT NULL;
