## test: runs all tests

test:
	@go test -v ./...

## cover: opens coverage in browser

cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

coverage:
	@go test -cover ./...
