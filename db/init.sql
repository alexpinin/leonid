CREATE TABLE IF NOT EXISTS pass
(
    pass     text    NOT NULL UNIQUE,
    valid_by integer NOT NULL
);

CREATE TABLE IF NOT EXISTS chat
(
    chat_id integer NOT NULL UNIQUE
);
