CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL,
    participant_id BIGINT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_messages_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT fk_messages_participant FOREIGN KEY (participant_id) REFERENCES participants(id) ON DELETE CASCADE
);

CREATE INDEX idx_messages_room ON messages (room_id);
CREATE INDEX idx_messages_participant ON messages (participant_id);
CREATE INDEX idx_messages_room_created ON messages (room_id, created_at DESC);
