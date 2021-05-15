test:
	go test ./...

deploy-service: test
	go build -o gte-daemon ./cmd/daemon/service.go
	scp gte-daemon ewintr.nl:/home/erik/gte
	ssh ewintr.nl sudo /bin/chown gte:gte /home/erik/gte
	ssh ewintr.nl sudo /usr/sbin/service gte-daemon stop
	ssh ewintr.nl sudo /bin/mv /home/erik/gte /usr/local/bin/gte
	ssh ewintr.nl sudo /usr/sbin/service gte-daemon start
