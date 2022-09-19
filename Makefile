test:
	go test -cover ./...

deploy-service: test
	flyctl deploy

install-cli: test
	go build -o gte ./cmd/cli/main.go
	mv gte ${HOME}/bin
