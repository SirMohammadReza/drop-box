# ---------- Stage 1: Builder ----------
FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app

# Install git (needed if using private modules)
RUN apk add --no-cache git

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api/main.go


# ---------- Stage 2: Runtime ----------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .

EXPOSE 50051

CMD ["./app"]