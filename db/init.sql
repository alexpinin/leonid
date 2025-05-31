CREATE TABLE IF NOT EXISTS pass_phrases (
    phrase text NOT NULL UNIQUE,
    valid_by integer NOT NULL
);

CREATE TABLE IF NOT EXISTS chats (
    chat_id integer NOT NULL UNIQUE
);
