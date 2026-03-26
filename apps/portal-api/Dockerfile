FROM golang:1.23-alpine AS builder

WORKDIR /src
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/portal-api ./apps/portal-api

FROM alpine:3.21

WORKDIR /app
RUN adduser -D -H -u 10001 portal

COPY --from=builder /out/portal-api /app/portal-api
COPY docs /app/docs

USER portal
EXPOSE 8080
ENTRYPOINT ["/app/portal-api"]
