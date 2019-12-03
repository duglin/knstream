all: client stream .stream

APP_IMAGE ?= duglin/stream

client: client.go
	go build -o client client.go

stream: stream.go
	go build -o stream stream.go

.stream: stream.go Dockerfile stream
	go build -o /dev/null stream.go  # Fail fast for compilation errors
	docker build -t $(APP_IMAGE) .   # Do the real build and create image
	docker push $(APP_IMAGE)
	touch .stream                    # Just to keep `make` happy

test: .stream client
	-kn service delete stream > /dev/null 2>&1 && sleep 5
	cat service.yaml | sed 's^IMAGE^$(APP_IMAGE)^' | kubectl apply -f -
	sleep 5
	./client $(shell kn service describe stream -o jsonpath='{.status.url}'):80

test-local: stream client
	./stream &
	sleep 2
	./client http://localhost:8080

clean:
	-rm -f stream client out .stream
	-kn service delete stream > /dev/null 2>&1
