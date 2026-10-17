CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL DEFAULT 'dam',
    title TEXT NOT NULL,
    amount INTEGER NOT NULL,
    type TEXT NOT NULL,
    due_date TEXT,
    group_id TEXT,
    status TEXT DEFAULT 'PENDING'
);