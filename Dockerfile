FROM golang:1.14

WORKDIR /app
COPY go.mod .
COPY go.sum .

WORKDIR /app/server
COPY server/ .

RUN go install -v ./...
RUN go build

CMD ["server"]