INSERT INTO config (pass, nicknames, system_prompt)
VALUES (
'Hello Leonid 123',
'Leo,Leonid',
'You are a helpful assistant' )
ON CONFLICT (pass) DO UPDATE
    SET pass = excluded.pass,
        nicknames = excluded.nicknames,
        system_prompt = excluded.system_prompt;
