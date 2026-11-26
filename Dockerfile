FROM golang:1.23-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /pulsewatch ./cmd/api

FROM alpine:3.20

RUN adduser -D -H -u 10001 appuser
USER appuser

COPY --from=builder /pulsewatch /pulsewatch

EXPOSE 8080

ENTRYPOINT ["/pulsewatch"]
