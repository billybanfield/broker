clean:
	rm -rf ./bin

test:
	GO111MODULE=on go test ./... -test.v

build-client:
	GO111MODULE=on go build -o ./bin/client ./cmd/client/

build-server:
	GO111MODULE=on go build -o ./bin/server ./cmd/server/
