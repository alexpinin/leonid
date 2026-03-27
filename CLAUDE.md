# Leonid

Telegram bot powered by LLM (OpenAI/DeepSeek). Receives messages from Telegram users, processes them through an LLM, and responds back. Maintains conversation history, authentication, activation, and quota management.

## Tech Stack

- Go 1.25.5
- SQLite (modernc.org/sqlite)
- go-telegram/bot for Telegram API
- openai-go for LLM integration

## Project Structure

```
src/internal/
  cmd/main.go          — entry point
  bot/
    bot.go             — bot initialization, handler chain setup
    dto/config.go      — data transfer objects
    handler/           — chain of responsibility handlers
    repo/config.go     — database access (repository pattern)
    service/           — business logic (config, openai, quota)
  db/                  — database connection & query executor
  logger/              — structured logging (slog)
  testutil/            — test assertion helpers
db/
  init.sql             — schema definition
```

## Commands

- **Build:** `./build.sh` (compiles to `./leonid`)
- **Run:** fill env vars in `start.sh`, then `./start.sh`
- **Test:** `make test` (clears cache, runs with race detector)
- **Lint/Vet:** `make audit` (tidy, verify, vet, test)
- **Format:** `make tidy` (go mod tidy + go fmt)

## Environment Variables

- `DB_FILE` — path to SQLite database
- `BOT_TOKEN` — Telegram bot token
- `LLM_PROVIDER` — LLM provider name (`openai` or `deepseek`)
- `LLM_TOKEN` — LLM API key
- `LLM_MODEL` — LLM model name

## Architecture

- **Chain of Responsibility** handler pattern: InputGuard → ChatChecker → ChatActivator → AuthGuard → CallGuard → QuotaGuard → MessageSender
- **Repository pattern** for data access
- **Dependency injection** via interfaces
- Conversation history stored as JSON in SQLite, sliding window of last 10 messages

## Code Conventions

- All application code under `src/internal/` (unexported)
- Interfaces are small and focused (single-method where possible)
- Handler constructors are lowercase package-level functions (`newInputGuard()`)
- Error wrapping with `fmt.Errorf("context: %w", err)`
- Table-driven tests, co-located `*_test.go` files
- No panics; all errors propagated up the call stack
