FROM golang:1.24.0 AS builder

WORKDIR /app

COPY ./src/go.mod ./src/go.sum ./

RUN go mod download

COPY ./src .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM scratch

WORKDIR /app

COPY --from=builder /app/main .

CMD ["/app/main"]
