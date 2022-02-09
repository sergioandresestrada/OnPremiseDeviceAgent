
FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /app/backend

FROM alpine as prod
WORKDIR /app
COPY --from=builder /app/backend /app/backend
EXPOSE 12345
CMD ["./backend"]
