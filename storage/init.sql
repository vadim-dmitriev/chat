CREATE TABLE IF NOT EXISTS USERS (
    'user_id' INTEGER,
    'username' TEXT NOT NULL UNIQUE,
    'password' TEXT NOT NULL,

    PRIMARY KEY ('user_id')
);

CREATE TABLE IF NOT EXISTS CONVERSATIONS (
    'conversation_id' INTEGER,
    'name' TEXT,
    'is_dialog' INTEGER NOT NULL,

    PRIMARY KEY ('conversation_id')
);

CREATE TABLE  IF NOT EXISTS MESSAGES (
    'message_id' INTEGER,
    'value' TEXT NOT NULL,
    'user_id' INTEGER NOT NULL,
    'conversation_id' INTEGER NOT NULL,
    'time' DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY ('message_id'),
    FOREIGN KEY ('user_id') REFERENCES users ('user_id'),
    FOREIGN KEY ('conversation_id') REFERENCES conversations ('conversation_id')
);

CREATE TABLE  IF NOT EXISTS MEMBERS (
    'member_id' INTEGER,
    'user_id' INTEGER NOT NULL,
    'conversation_id' INTEGER NOT NULL,

    
    PRIMARY KEY ('member_id'),
    FOREIGN KEY ('user_id') REFERENCES users ('user_id'),
    FOREIGN KEY ('conversation_id') REFERENCES conversations ('conversation_id')
);

CREATE TABLE IF NOT EXISTS COOKIES (
    'cookie_id' INTEGER,
    'user_id' INTEGER NOT NULL,
    'value' TEXT,

    PRIMARY KEY ('cookie_id'),
    FOREIGN KEY ('user_id') REFERENCES users ('user_id')
);
