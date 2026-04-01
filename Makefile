.PHONY: help
help: # default command position
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


# =================================================================================================
# Environment
# =================================================================================================

## db/init: init database
.PHONY: db/init
db/init:
	@touch db/leonid.sqlite3
	@sqlite3 db/leonid.sqlite3 < db/init.sql

## db/insert: insert data into the database from a file [file=db/insert.example.sql]
.PHONY: db/insert
db/insert: file ?= 'db/insert.example.sql'
db/insert:
	@sqlite3 db/leonid.sqlite3 < ${file}

## docker/up: docker compose up (if needed)
.PHONY: docker/up
docker/up:
	@docker compose -f compose.dev.yaml up --detach

## docker/down: docker compose up (if needed)
.PHONY: docker/down
docker/down:
	@docker compose -f compose.dev.yaml down

# =================================================================================================
# Application
# =================================================================================================

## build: build project binary
.PHONY: build
build:
	@go build -o leonid leonid/src/cmd
	@chmod +x leonid

## start: run project binary detached from terminal
.PHONY: start
start: file ?= .env
start:
	@set -a && source ${file} && set +a && nohup ./leonid > leonid.log 2>&1 &
	@echo "Started leonid detached"

## status: check running status
.PHONY: status
status:
	@ps ax | grep leonid | grep -v grep

## stop: stop running leonid process
.PHONY: stop
stop:
	@pkill -f leonid || echo "No process found"


# =================================================================================================
# Development
# =================================================================================================

## test: run tests
.PHONY: test
test:
	@echo 'Running tests...'
	@go clean -testcache && go test -vet=off -race ./...
	@echo 'Done'

## tidy: tidy dependencies
.PHONY: tidy
tidy:
	@echo 'Tidying dependencies...'
	@go mod tidy
	@echo 'Formatting .go files...'
	@go fmt ./...
	@echo 'Done'

## audit: audit project
.PHONY: audit
audit:
	@echo 'Checking dependencies...'
	@go mod tidy -diff
	@go mod verify
	@echo 'Vetting code...'
	@go vet ./...
	@echo 'Running tests...'
	@make test
	@echo 'Done'
