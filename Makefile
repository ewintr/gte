test:
	go test -cover ./...

deploy-service: test
	GOARCH=386 go build -o gte-daemon ./cmd/daemon/service.go
	scp gte-daemon wittekastje:/home/erik/bin/gte-daemon
	ssh wittekastje /home/erik/bin/deploy-gte-daemon.sh

install-cli: test
	go build -o gte ./cmd/cli/main.go
	mv gte ${HOME}/bin
