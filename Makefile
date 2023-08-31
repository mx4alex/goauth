.PHONY: run build
run: 
	go run ./cmd/app/main.go
build:
	go build -o ./build/goauth ./cmd/app/main.go
swag:
	swag init -g cmd/app/main.go