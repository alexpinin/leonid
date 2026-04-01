# Leonid

Telegram bot powered by LLM. Add the bot to a group chat, activate it with a password, and it responds when mentioned by nickname. Supports OpenAI and DeepSeek as LLM providers.

## Requirements

- Go 1.26.1
- SQLite3

## Set up a Telegram bot

1. Create a bot with @BotFather
2. Get bot token
3. Change bot privacy with @BotFather and let it read all the messages

## Init DB

- Create a db file and schema: `make db/init`
- Insert a config entry: `make db/insert [file=insert.example.sql]`. File is optional, default is `insert.example.sql`

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

## TODO

- **Unbounded `sync.Map` growth:** `OpenAIService.chatLocks` stores a mutex per chat ID and never evicts entries. (`service/openai.go`)
- **No user feedback on quota exceeded:** When quota is exceeded the error is logged but the Telegram user receives no message.
- **No message length guard yet:** probably makes sense to limit
