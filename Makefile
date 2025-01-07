.PHONY: build
build:
	go build -o ./build/app ./cmd/main.go

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test
test:
	go test -coverprofile ./test/cover.out ./...
	go tool cover -html ./test/cover.out -o ./test/cover.html
