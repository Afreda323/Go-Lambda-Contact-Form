build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/sendMail sendMail/main.go sendMail/userData.go sendMail/helpers.go sendMail/email.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
