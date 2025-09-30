run:
	go run ./cmd/api/ --config=./configs/config.yaml

build:
	go build ./cmd/api/

test:
	go test ./...

test-cover:
	go test ./... -coverprofile cover.test.tmp -coverpkg ./...
	type cover.test.tmp | findstr -v "mocks" > cover.test 
	del cover.test.tmp 
	go tool cover -func cover.test 