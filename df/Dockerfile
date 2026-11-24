FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/cook_book_api ./cmd/main.go

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/cook_book_api /app/cook_book_api

RUN mkdir -p /app/logs && chown -R app:app /app

USER app

EXPOSE 8080

CMD ["/app/cook_book_api"]
