CREATE TABLE IF NOT EXISTS config
(
    id                   integer NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
    chat_id              integer NOT NULL UNIQUE,
    pass                 text    NOT NULL UNIQUE,
    pass_valid_by        integer NOT NULL DEFAULT 0,
    chat_activated_at    integer NOT NULL DEFAULT 0,
    nicknames            text    NOT NULL DEFAULT '',
    system_prompt        text    NOT NULL DEFAULT '',
    conversation_context text    NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS quota
(
    id              integer NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
    chat_id         integer NOT NULL UNIQUE,
    last_reset_date integer NOT NULL DEFAULT 0,
    remaining       integer NOT NULL DEFAULT 0
);
