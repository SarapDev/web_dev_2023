check: deps vet lint unit

unit:
	go test -v ./...

vet:
	go vet ./...

lint:
	golangci-lint run -v

deps:
	go mod verify

vendor:
	go mod vendor

hook:
	cp -f resources/hooks/pre-commit.sh .git/hooks/pre-commit

