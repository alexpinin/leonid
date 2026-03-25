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
- `OPENAI_LLM_TOKEN` — OpenAI API key
- `DEEPSEEK_LLM_TOKEN` — DeepSeek API key

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

## Known Issues

### TODO

- Add the message length guard

### Bugs

1. **README references `run.sh` but file is `start.sh`** — README.md:16
2. **start.sh env vars are out of sync with main.go** — start.sh exports OPENAI_LLM_TOKEN/DEEPSEEK_LLM_TOKEN but main.go expects LLM_PROVIDER, LLM_TOKEN, LLM_MODEL. Bot will fail to start with current start.sh.
3. **callGuard nickname case mismatch** — message is lowercased but nicknames from DB are not; reply-to check lowers replyToNickname but not the nickname. Matches fail when nicknames have uppercase. `call_guard.go:30-46`

### Code Quality

12. **Duplicate configRepo interface** — exported in service/config.go:29, unexported in service/openai.go:33. Can drift.
13. **MockQueryExecutor passes nil tx** — code under test panics if it actually uses the tx. `executor.go:51-53`
14. **Logger strips slog structured logging** — only exposes Info(string)/Error(string), no fields or context. `logger.go`
15. **testutil.Equal missing t.Helper()** — failure traces point to equal.go, not the test. `equal.go:8`
16. **Unnecessary fmt.Sprintf** — `logger.Info(fmt.Sprintf("Starting bot"))` should be `logger.Info("Starting bot")`. `main.go:28`
18. **Unwrapped error from json.Unmarshal** — in conversationHistory, unmarshal error returned without context. `openai.go:110-112`
