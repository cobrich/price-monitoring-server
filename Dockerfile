FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./out/server ./cmd/app

EXPOSE 8080

CMD ["./out/server"]
