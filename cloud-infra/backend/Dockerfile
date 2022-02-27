
FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD ./cmd ./cmd
ADD ./pkg ./pkg
RUN go build -o /app/backend ./cmd/backend/backend.go

FROM alpine as prod
WORKDIR /app
COPY --from=builder /app/backend /app/backend
EXPOSE 12345
CMD ["./backend"]
