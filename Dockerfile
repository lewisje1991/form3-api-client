FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build ./...

ENTRYPOINT ["go", "test", "-v", "./...", "-coverprofile", "cover.out"]