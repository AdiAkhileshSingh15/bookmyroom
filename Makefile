run:
	go run cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go cmd/web/send-mail.go -dbname=bookmyroom -dbuser=postgres -dbpass=adi123

test:
	go test -v ./...

cover:
	go test -v ./... -cover

cover-html:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

build:
	go build -o bookmyroom cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go cmd/web/send-mail.go