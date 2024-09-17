FROM golang:1.22

COPY server.go go.mod .

RUN go build

EXPOSE 8080

CMD ["./server"]
