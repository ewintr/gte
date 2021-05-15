test:
	go test ./...

deploy-service: test
	go build -o gte-daemon ./cmd/daemon/service.go
	scp gte-daemon zerocontent.org:/home/erik/gte
	ssh zerocontent.org sudo /bin/chown gte:gte /home/erik/gte
	ssh zerocontent.org sudo /usr/sbin/service gte-daemon stop
	ssh zerocontent.org sudo /bin/mv /home/erik/gte /usr/local/bin/gte
	ssh zerocontent.org sudo /usr/sbin/service gte-daemon start
