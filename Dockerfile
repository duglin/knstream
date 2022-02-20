FROM golang
COPY stream.go /src/
WORKDIR /src/
ENV GO111MODULE=off
RUN go get -d .
RUN find / -name gorilla -print
RUN go build -o /stream stream.go
ENTRYPOINT [ "/stream" ]
