CREATE TABLE IF NOT EXISTS users (
    id BLOB PRIMARY KEY,
    oauth_id TEXT NOT NULL UNIQUE,
    oauth_provider TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    mashboard_email TEXT NOT NULL,
    email_rss_link TEXT DEFAULT "",
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feeds  (
    feed_id BLOB PRIMARY KEY,
    title TEXT DEFAULT "",
    link TEXT NOT NULL,
    image TEXT DEFAULT "",
    alt_text TEXT DEFAULT "", 
    media_type TEXT DEFAULT "",
    categories TEXT DEFAULT "",
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
    alt_text TEXT DEFAULT "", 
    media_type TEXT DEFAULT "",
    categories TEXT DEFAULT "",
    pub_date TEXT DEFAULT CURRENT_TIMESTAMP,
    guid TEXT UNIQUE,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(feed_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS feed_follows (
    user_id BLOB NOT NULL,
    feed_id BLOB NOT NULL,
    user_feed_name TEXT DEFAULT "",
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_id) REFERENCES feeds(feed_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, feed_id)
);

CREATE TABLE IF NOT EXISTS saved_feeds (
    user_id BLOB NOT NULL,
    feed_item_id INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feed_item_id) REFERENCES feed_items(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, feed_item_id)
);

CREATE TABLE IF NOT EXISTS collections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_collections (
    user_id BLOB NOT NULL,
    collection_id INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, collection_id)
);

