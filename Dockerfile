FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o rinha-backend

FROM alpine:latest
COPY --from=builder /app/rinha-backend /app/
CMD ["/app/rinha-backend"]