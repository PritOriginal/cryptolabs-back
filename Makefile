run:
	go run ./cmd/api/ --config=./configs/config.yaml

build:
	go build ./cmd/api/

test:
	go test ./...

test-cover:
	go test ./... -coverprofile cover.test.tmp -coverpkg ./...
	cat cover.test.tmp | grep -v "mocks" > cover.test 
	rm cover.test.tmp 
	go tool cover -func cover.test 