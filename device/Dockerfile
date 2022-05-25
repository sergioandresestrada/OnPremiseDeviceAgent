FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD ./cmd ./cmd
ADD ./pkg ./pkg
ADD ./files ./files
ADD ./receivedFiles ./receivedFiles

RUN go build -o /app/device ./cmd/device/device.go

FROM alpine as prod
WORKDIR /app
COPY --from=builder /app/device /app/device
EXPOSE 12345
CMD ["./device"]