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

CREATE INDEX IF NOT EXISTS idx_user_nickname ON User(nickname);
CREATE INDEX IF NOT EXISTS idx_user_email ON User(email);

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

CREATE INDEX IF NOT EXISTS idx_post_user_id ON Post(user_id);
CREATE INDEX IF NOT EXISTS idx_post_category ON Post(category);

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

CREATE INDEX IF NOT EXISTS idx_comment_post_id ON Comment(post_id);
CREATE INDEX IF NOT EXISTS idx_comment_user_id ON Comment(user_id);
CREATE INDEX IF NOT EXISTS idx_comment_parent_id ON Comment(parent_comment_id);

-- Like table --
CREATE TABLE IF NOT EXISTS Like (
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

CREATE INDEX IF NOT EXISTS idx_like_post_id ON Like(post_id);
CREATE INDEX IF NOT EXISTS idx_like_comment_id ON Like(comment_id);
CREATE INDEX IF NOT EXISTS idx_like_user_id ON Like(user_id);

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

CREATE INDEX IF NOT EXISTS idx_message_sender_id ON Message(sender_id);
CREATE INDEX IF NOT EXISTS idx_message_receiver_id ON Message(receiver_id);
CREATE INDEX IF NOT EXISTS idx_message_status ON Message(status);

-- OnlineStatus table with session management --
CREATE TABLE IF NOT EXISTS OnlineStatus (
    id INTEGER PRIMARY KEY AUTOINCREMENT,   -- Unique session identifier
    user_id INTEGER NOT NULL,              -- Foreign key to User table
    session_id TEXT NOT NULL UNIQUE,       -- Unique session identifier for each login
    ip_address BLOB,                       -- IP address of the user
    user_agent BLOB,                       -- User agent string (browser/device)
    is_online BOOLEAN DEFAULT FALSE,       -- Online status
    last_active DATETIME DEFAULT CURRENT_TIMESTAMP, -- Last activity timestamp
    login_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- Login timestamp
    logout_time DATETIME,                 -- Logout timestamp
    FOREIGN KEY (user_id) REFERENCES User(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT unique_session_per_user UNIQUE (user_id, session_id)
);

CREATE INDEX IF NOT EXISTS idx_online_status_user_id ON OnlineStatus(user_id);
CREATE INDEX IF NOT EXISTS idx_online_status_is_online ON OnlineStatus(is_online);
CREATE INDEX IF NOT EXISTS idx_online_status_last_active ON OnlineStatus(last_active);

-- Additional Indexes for Optimization --
CREATE INDEX IF NOT EXISTS idx_user_created_at ON User(created_at);
CREATE INDEX IF NOT EXISTS idx_post_created_at ON Post(created_at);
CREATE INDEX IF NOT EXISTS idx_comment_created_at ON Comment(created_at);
CREATE INDEX IF NOT EXISTS idx_like_created_at ON Like(created_at);
CREATE INDEX IF NOT EXISTS idx_message_created_at ON Message(created_at);

