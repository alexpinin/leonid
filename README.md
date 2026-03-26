# Leonid

Telegram bot powered by LLM. Add the bot to a group chat, activate it with a password, and it responds when mentioned by nickname. Supports OpenAI and DeepSeek as LLM providers.

## Requirements

- Go 1.25.5+
- SQLite3

## Set up a Telegram bot

1. Create a bot with @BotFather
2. Get bot token
3. Change bot privacy with @BotFather and let it read all the messages

## Init DB

- Create a db file: `touch db/leonid.sqlite3`
- Create db schema: `sqlite3 db/leonid.sqlite3 < db/init.sql`
- Insert a config entry:
```sql
sqlite3 db/leonid.sqlite3 "INSERT INTO config (chat_id, pass, pass_valid_by, nicknames, system_prompt) VALUES (0, 'your-secret-pass', strftime('%s', '2027-01-01'), 'botname', 'You are a helpful assistant');"
```

## Environment variables

Copy `.env.example` to `.env` and set the values:

| Variable       | Description                                           |
|----------------|-------------------------------------------------------|
| `DB_FILE`      | Path to SQLite database file                          |
| `BOT_TOKEN`    | Telegram bot token from @BotFather                    |
| `LLM_PROVIDER` | `openai` or `deepseek`                                |
| `LLM_TOKEN`    | API key for the chosen provider                       |
| `LLM_MODEL`    | Model name (e.g. `gpt-5-mini-2025-08-07`, `deepseek-chat`) |

## Build

- Run `make build`

## Run

- Run `make start`

## Stop

- Run `make stop`
