CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    sender_user_id INTEGER NOT NULL REFERENCES users (id),
    receiver_user_id INTEGER NOT NULL REFERENCES users (id),
    currency INTEGER NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);