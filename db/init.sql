CREATE TABLE IF NOT EXISTS pass
(
    pass          text    NOT NULL UNIQUE,
    valid_by      integer NOT NULL,
    nicknames     text    NOT NULL,
    system_prompt text    NOT NULL
);

CREATE TABLE IF NOT EXISTS chat
(
    id            integer NOT NULL UNIQUE,
    nicknames     text    NOT NULL,
    system_prompt text    NOT NULL
);
