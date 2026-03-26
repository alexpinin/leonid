# leonid

## Set up a Telegram bot
1. Create a bot with @BotFather
2. Get bot token
3. Change bot privacy with @BotFather and let it read all the messages

## Init DB
- Create a db file: `touch db/leonid.sqlite3`
- Create db schema: `sqlite3 db/leonid.sqlite3 < db/init.sql`

## Build
- Run `make build`

## Run
- Copy `.env.example` to `.env` and set the environment variables
- Run `make start`

## Stop
- Run `make stop`
