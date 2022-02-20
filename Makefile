all: client stream .stream

export GO111MODULE=off

client: client.go
	go build -o client client.go

stream: stream.go
	go build -o stream stream.go

.stream: Dockerfile stream
	go build -o /dev/null stream.go  # Fail fast for compilation errors
	docker build -t duglin/stream .  # Do the real build and create image
	docker push duglin/stream
	touch .stream                    # Just to keep `make` happy

runserver: stream
	./stream

runclient: client
	./client

deploy: .stream
	-kn service delete stream > /dev/null 2>&1 && sleep 3
	kn service create stream --image duglin/stream

test: deploy client
	sleep 5
	./client $(shell kn service describe stream -o jsonpath='{.status.url}'):80

clean:
	-rm -f stream client .stream
	-kn service delete stream > /dev/null 2>&1
