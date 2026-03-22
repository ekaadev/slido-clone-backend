CREATE TABLE polls (
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL,
    question TEXT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'closed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    activated_at TIMESTAMPTZ NULL,
    closed_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_polls_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

CREATE INDEX idx_polls_room ON polls (room_id);
CREATE INDEX idx_polls_status ON polls (status);
CREATE INDEX idx_polls_room_status ON polls (room_id, status);
CREATE INDEX idx_polls_activated_at ON polls (activated_at DESC);
CREATE INDEX idx_polls_created_at ON polls (created_at DESC);
