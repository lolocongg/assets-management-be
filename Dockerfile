# ---------- build stage ----------
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o app ./cmd/api && \
    go build -ldflags="-s -w" -o migrate ./cmd/migrate

# ---------- runtime ----------
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache tzdata

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder --chown=app:app /app/app .
COPY --from=builder --chown=app:app /app/migrate .
COPY --from=builder --chown=app:app /app/internal/database/migrations ./internal/database/migrations

USER app

EXPOSE 8080

CMD ["./app"]