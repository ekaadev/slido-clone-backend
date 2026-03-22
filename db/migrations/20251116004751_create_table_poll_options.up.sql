CREATE TABLE poll_options (
    id BIGSERIAL PRIMARY KEY,
    poll_id BIGINT NOT NULL,
    option_text VARCHAR(255) NOT NULL,
    vote_count INT NOT NULL DEFAULT 0,
    "order" SMALLINT NOT NULL,

    CONSTRAINT fk_poll_options_poll FOREIGN KEY (poll_id) REFERENCES polls(id) ON DELETE CASCADE
);

CREATE INDEX idx_poll_options_poll ON poll_options (poll_id);
CREATE INDEX idx_poll_options_order ON poll_options (poll_id, "order");
