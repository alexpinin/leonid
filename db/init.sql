CREATE TABLE IF NOT EXISTS config
(
    id                   text    NOT NULL UNIQUE,
    pass                 text    NOT NULL UNIQUE,
    pass_valid_by        integer NOT NULL DEFAULT 0,
    chat_id              integer NOT NULL DEFAULT 0,
    chat_activated_at    integer NOT NULL DEFAULT 0,
    nicknames            text    NOT NULL DEFAULT '',
    system_prompt        text    NOT NULL DEFAULT '',
    conversation_context text    NOT NULL DEFAULT '{}'
);
