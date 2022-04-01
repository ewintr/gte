test:
	go test -cover ./...

deploy-service: test
	go build -o gte-daemon ./cmd/daemon/service.go
	scp gte-daemon blok:/tmp/gte-daemon
	ssh -t blok /home/erik/bin/deploy-gte-daemon.sh

install-cli: test
	go build -o gte ./cmd/cli/main.go
	mv gte ${HOME}/bin
