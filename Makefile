test:
	go test -cover ./...

deploy-service: test
	go build -o gte-daemon ./cmd/daemon/service.go
	scp gte-daemon ewintr.nl:/tmp/gte-daemon
	ssh -t erik@ewintr.nl /home/erik/bin/deploy-gte-daemon.sh

install-cli: test
	go build -o gte ./cmd/cli/main.go
	mv gte ${HOME}/bin
