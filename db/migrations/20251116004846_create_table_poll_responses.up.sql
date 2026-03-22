CREATE TABLE poll_responses (
    id BIGSERIAL PRIMARY KEY,
    poll_id BIGINT NOT NULL,
    participant_id BIGINT NOT NULL,
    poll_option_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_poll_responses_poll FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE,
    CONSTRAINT fk_poll_responses_participant FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE,
    CONSTRAINT fk_poll_responses_option FOREIGN KEY (poll_option_id) REFERENCES poll_options(id) ON DELETE CASCADE,
    CONSTRAINT unique_poll_response UNIQUE (poll_id, participant_id)
);

CREATE INDEX idx_poll_responses_poll ON poll_responses (poll_id);
CREATE INDEX idx_poll_responses_participant ON poll_responses (participant_id);
CREATE INDEX idx_poll_responses_option ON poll_responses (poll_option_id);
CREATE INDEX idx_poll_responses_created_at ON poll_responses (created_at DESC);
