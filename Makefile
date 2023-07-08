run:
	go run cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go

test:
	go test -v ./...

cover:
	go test -v ./... -cover

cover-html:
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out