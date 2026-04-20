default: fmt lint build generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

vet:
	go vet ./...

tidy:
	go mod tidy
	cd tools && go mod tidy

clean:
	rm -f terraform-provider-jumpcloud
	rm -rf dist/

.PHONY: fmt lint test testacc build install generate vet tidy clean
