CREATE TABLE rooms (
    id BIGSERIAL PRIMARY KEY,
    room_code VARCHAR(20) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    presenter_id BIGINT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'closed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_rooms_presenter FOREIGN KEY (presenter_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_rooms_room_code ON rooms (room_code);
CREATE INDEX idx_rooms_status ON rooms (status);
CREATE INDEX idx_rooms_presenter ON rooms (presenter_id);
CREATE INDEX idx_rooms_created_at ON rooms (created_at DESC);
