CREATE TABLE IF NOT EXISTS config
(
    id                   INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
    chat_id              INTEGER NULL UNIQUE, -- update on runtime by the pass
    pass                 TEXT    NOT NULL UNIQUE,
    pass_valid_by        INTEGER NOT NULL DEFAULT (strftime('%s', 'now', '+1 month')),
    chat_activated_at    INTEGER NOT NULL DEFAULT 0,
    nicknames            TEXT    NOT NULL,
    system_prompt        TEXT    NOT NULL,
    conversation_context TEXT    NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS quota
(
    id              INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
    chat_id         INTEGER NOT NULL UNIQUE,
    last_reset_date INTEGER NOT NULL DEFAULT 0,
    remaining       INTEGER NOT NULL DEFAULT 0
);
