FROM golang:latest as builder

RUN go get github.com/dinumathai/loadtest/server

FROM alpine:latest
WORKDIR /

COPY --from=builder /go/bin/server server

RUN chmod a+x /server
CMD ["./server"]
