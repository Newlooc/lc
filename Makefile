fun:
	go mod vendor
	go mod tidy
	go build -o _bin/fun cmd/main.go
