FROM golang:1.25-alpine as builder

RUN apk add --no-cache \
      gcc \
      g++ \
      make \
      vips-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk add --no-cache \
    vips \
    ca-certificates \
    tzdata

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/uploads

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]