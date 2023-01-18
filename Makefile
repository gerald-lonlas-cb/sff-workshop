deps:
	go get ./...
	npm install

build:
	./build.sh

run:
	go run cmd/main.go -port 8081
