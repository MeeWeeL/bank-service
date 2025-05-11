CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    account_id INTEGER REFERENCES accounts(id) NOT NULL,
    encrypted_data TEXT NOT NULL,
    hmac TEXT NOT NULL,
    cvv_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cards_user_id ON cards(user_id);
CREATE INDEX idx_cards_hmac ON cards(hmac);