deps:
	go get ./...

build:
	./build.sh

run:
	go run cmd/main.go -port 8081

test: 
	go test -test.v -race -cover ./internal/...

update_abi:
	cd contract && solc --abi --bin OrderContract.sol -o build && ../tools/abigen --bin=./build/OrderContract.bin --abi=./build/OrderContract.abi --pkg=contract --out=contract.go