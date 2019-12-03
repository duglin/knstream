FROM golang:alpine
COPY stream.go /
RUN go build -o /stream /stream.go
ENTRYPOINT [ "/stream" ]
