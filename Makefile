test:
	go test -race -cover ./...

build: test
	env GOOS=linux CGO_ENABLED=0 go build ./cmd/server

update-deps:
	dep ensure -update

docker-build:
	docker build -t heppu/contact-storage .

docker-run: docker-build
	docker run --rm -it	-p 8000:8000 heppu/contact-storage

.PHONY: test update-deps docker-build docker-run
