run:
	go run ./cmd/api/ --config=./configs/config.yaml

build:
	go build ./cmd/api/

test:
	go test ./...

test-cover:
	go test ./... -coverprofile cover.out.tmp -coverpkg ./...
	cat cover.out.tmp | grep -v "mocks" > cover.out 
	rm cover.out.tmp 
	go tool cover -func cover.out 