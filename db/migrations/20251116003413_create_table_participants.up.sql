CREATE TABLE participants (
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL,
    user_id BIGINT NULL,
    display_name VARCHAR(100) NOT NULL,
    xp_score INT NOT NULL DEFAULT 0,
    is_anonymous BOOLEAN NOT NULL DEFAULT TRUE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_participants_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT fk_participants_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT unique_user_room UNIQUE (room_id, user_id)
);

CREATE INDEX idx_participants_room ON participants (room_id);
CREATE INDEX idx_participants_xp_score ON participants (xp_score DESC);
CREATE INDEX idx_participants_joined_at ON participants (joined_at DESC);
