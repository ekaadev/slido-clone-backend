ALTER TABLE xp_transactions
    DROP CONSTRAINT IF EXISTS xp_transactions_source_type_check;

ALTER TABLE xp_transactions
    ADD CONSTRAINT xp_transactions_source_type_check
    CHECK (source_type IN ('poll', 'question_created', 'upvote_received', 'presenter_validated'));
