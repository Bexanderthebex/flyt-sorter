VERSION ?= v0.1.0
PORT ?= 8080

build:
	go build .

test:
	go test -v ./...

bench-test:
	go test -bench=. -benchmem

docker-run:
	$(MAKE) build
	docker buildx build -t flight-sorter:$(VERSION) .
	docker run -d --name flight-sorter -p $(PORT):8080 flight-sorter:$(VERSION)