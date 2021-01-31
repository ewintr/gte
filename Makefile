test:
	go test ./...

deploy:
	go build -o gte-process-inbox ./cmd/process-inbox/main.go
	go build -o gte-generate-recurring ./cmd/generate-recurring/main.go
	scp gte-* zerocontent.org:bin/
