FROM golang:alpine

WORKDIR /go/bin

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./generator ./cmd/server/main.go

ENTRYPOINT ["./generator"]