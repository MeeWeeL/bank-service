CREATE TABLE credits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    account_id INTEGER REFERENCES accounts(id) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    rate DECIMAL(5,2) NOT NULL,
    period INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_credits_user_id ON credits(user_id);
CREATE INDEX idx_credits_account_id ON credits(account_id);
