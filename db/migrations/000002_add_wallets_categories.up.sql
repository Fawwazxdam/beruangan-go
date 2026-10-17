CREATE TABLE IF NOT EXISTS wallets (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    balance INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL
);

ALTER TABLE transactions ADD COLUMN category_id TEXT;
ALTER TABLE transactions ADD COLUMN wallet_id TEXT;
ALTER TABLE transactions ADD COLUMN transaction_date TEXT;