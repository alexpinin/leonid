CREATE TABLE IF NOT EXISTS config
(
    chat_id              integer NOT NULL UNIQUE PRIMARY KEY,
    pass                 text    NOT NULL UNIQUE,
    pass_valid_by        integer NOT NULL DEFAULT 0,
    chat_activated_at    integer NOT NULL DEFAULT 0,
    nicknames            text    NOT NULL DEFAULT '',
    system_prompt        text    NOT NULL DEFAULT '',
    conversation_context text    NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS quota
(
    chat_id         integer NOT NULL UNIQUE PRIMARY KEY,
    last_reset_date integer NOT NULL DEFAULT 0,
    remaining       integer NOT NULL DEFAULT 0
);
