FROM golang:alpine as builder

WORKDIR /app

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates

COPY server.go go.mod go.sum /app/
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o server server.go

FROM scratch

WORKDIR /app

COPY --from=builder app/server ./server

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./server"]