run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go --output docs

dev-run: swagger run

COVERAGE_EXCLUDE=mocks|main.go|test

test:
	go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html