FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /server /server
EXPOSE 8080
CMD ["/server"]
