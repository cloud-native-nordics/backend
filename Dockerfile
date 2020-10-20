FROM golang:alpine as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN mkdir /app

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates

COPY server.go go.mod go.sum /app/
WORKDIR /app
RUN go build -o server server.go

FROM scratch
WORKDIR /app
COPY --from=builder app/server ./server
# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./server"]