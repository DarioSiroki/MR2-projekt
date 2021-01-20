FROM golang:1.14

EXPOSE 6379

WORKDIR ./server
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["server"]