FROM golang:latest

COPY . .

RUN go build -o /go/bin/server github.com/dinumathai/loadtest/server
RUN go build -o /go/bin/client github.com/dinumathai/loadtest/client

CMD ["/go/bin/server"]
