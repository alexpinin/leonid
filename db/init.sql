CREATE TABLE IF NOT EXISTS config
(
    id                text    NOT NULL UNIQUE,
    pass              text    NOT NULL UNIQUE,
    pass_valid_by     integer NOT NULL,
    chat_id           integer NOT NULL,
    chat_activated_at integer NOT NULL,
    nicknames         text    NOT NULL,
    system_prompt     text    NOT NULL
);
