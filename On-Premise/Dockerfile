FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD ./cmd ./cmd
ADD ./pkg ./pkg
RUN go build -o /app/on_premise ./cmd/on_premise/on_premise.go

FROM alpine as prod
WORKDIR /app
COPY --from=builder /app/on_premise /app/on_premise
ADD ./onPremiseFiles ./onPremiseFiles
CMD ["./on_premise"]