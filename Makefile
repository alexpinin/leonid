.PHONY: test
test:
	@echo 'Running tests...'
	@go clean -testcache && go test -vet=off -race ./...
	@echo 'Done'

.PHONY: tidy
tidy:
	@echo 'Tidying dependencies...'
	@go mod tidy
	@echo 'Formatting .go files...'
	@go fmt ./...
	@echo 'Done'

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
