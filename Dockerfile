FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o firestore-go .

FROM gcr.io/distroless/base:latest

WORKDIR /app
COPY --from=builder /app/firestore-go .

EXPOSE 8080

ENTRYPOINT ["/app/firestore-go"]
