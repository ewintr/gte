test:
	go test ./...

deploy-service: test
	go build -o gte-daemon ./cmd/daemon/service.go
	scp gte-daemon ewintr.nl:/home/erik/bin/gte-daemon
	ssh ewintr.nl /home/erik/bin/deploy-gte-daemon.sh
