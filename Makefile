test:
	go test -cover ./...

service-deploy: test
	flyctl deploy

cli-install: test
	go build -o gte ./cmd/cli/main.go
	mv gte ${HOME}/bin

app-build:
	cd cmd/android-app && fyne package -os android --icon ../../Icon.png --appID nl.ewintr.gte -name gte
	mv cmd/android-app/gte.apk .

app-run:
	go run -tags mobile ./cmd/android-app/main.go
