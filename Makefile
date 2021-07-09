pull:
	git pull

test:
	go test -cover ./...

deploy-service: pull test
	go build -o gte-daemon ./cmd/daemon/service.go
	scp gte-daemon ewintr.nl:/home/erik/bin/gte-daemon
	ssh ewintr.nl /home/erik/bin/deploy-gte-daemon.sh

install-cli: pull test
	go build -o gte ./cmd/cli/main.go
	mv gte ${HOME}/bin
