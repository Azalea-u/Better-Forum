-- User table --
CREATE TABLE IF NOT EXISTS User (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT UNIQUE NOT NULL,
    age INTEGER NOT NULL,
    gender TEXT CHECK(gender IN ('Male', 'Female', 'Other')) NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password BLOB NOT NULL, -- Hashed password
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP -- Track when the user profile is updated
);

-- Post table --
CREATE TABLE IF NOT EXISTS Post (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    category TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Track post updates
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Comment table --
CREATE TABLE IF NOT EXISTS Comment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Track comment updates
    parent_comment_id INTEGER, -- For nested comments
    FOREIGN KEY (post_id) REFERENCES Post(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES Comment(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Like table --
CREATE TABLE IF NOT EXISTS Like ( -- for both posts or comments
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    comment_id INTEGER,
    user_id INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK(value IN (1, -1)), -- 1 for like, -1 for dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES Post(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES Comment(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Message table --
CREATE TABLE IF NOT EXISTS Message (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    status TEXT DEFAULT 'UNREAD' CHECK(status IN ('UNREAD', 'READ')), -- Status for read/unread messages
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Track when message is read or updated
    FOREIGN KEY (sender_id) REFERENCES User(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES User(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Status table --
CREATE TABLE IF NOT EXISTS OnlineStatus (
    user_id INTEGER PRIMARY KEY,
    is_online BOOLEAN DEFAULT FALSE,
    last_active DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Indexes for performance optimization (useful for queries) --
CREATE INDEX IF NOT EXISTS idx_like_post_id ON Like(post_id);
CREATE INDEX IF NOT EXISTS idx_like_comment_id ON Like(comment_id);
CREATE INDEX IF NOT EXISTS idx_like_user_id ON Like(user_id);
CREATE INDEX IF NOT EXISTS idx_message_sender_id ON Message(sender_id);
CREATE INDEX IF NOT EXISTS idx_message_receiver_id ON Message(receiver_id);
CREATE INDEX IF NOT EXISTS idx_post_user_id ON Post(user_id);
CREATE INDEX IF NOT EXISTS idx_comment_post_id ON Comment(post_id);
CREATE INDEX IF NOT EXISTS idx_comment_user_id ON Comment(user_id);

-- Additional indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_message_status ON Message(status); -- Index for message status queries
CREATE INDEX IF NOT EXISTS idx_online_status ON OnlineStatus(is_online); -- Index for online users
