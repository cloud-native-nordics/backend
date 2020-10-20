FROM golang:alpine as builder
ENV CGO_ENABLED=0
ENV GOOS=linux
RUN mkdir /app
COPY server.go go.mod go.sum /app/
WORKDIR /app
RUN go build -o server server.go

FROM scratch
WORKDIR /app
COPY --from=builder app/server ./server

CMD ["./server"]