FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /projek ./cmd/api

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /projek .
COPY --from=builder /app/migrations ./migrations

ENV PORT=8080
ENV NODE_ENV=production

EXPOSE 8080

CMD ["./projek"]
