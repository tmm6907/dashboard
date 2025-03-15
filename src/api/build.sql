CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    oauth_id TEXT NOT NULL,
    oauth_provider TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    mashboard_email TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feeds  (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id BLOB NOT NULL UNIQUE,
    title TEXT NOT NULL,
    link TEXT NOT NULL,
    image TEXT DEFAULT "",
    media_type TEXT DEFAULT "",
    description TEXT DEFAULT "",
    language TEXT DEFAULT "",
    last_build_date TEXT DEFAULT CURRENT_TIMESTAMP,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feed_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feed_id BLOB NOT NULL,
    title TEXT NOT NULL,
    link TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT "",
    image TEXT DEFAULT "", 
    media_type TEXT DEFAULT "",
    pub_date TEXT DEFAULT CURRENT_TIMESTAMP,
    guid TEXT UNIQUE,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(feed_id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_feed_item_feed_id ON feed_items(feed_id);