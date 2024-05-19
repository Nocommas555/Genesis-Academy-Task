FROM golang:1.22 AS stage1

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
RUN go mod verify

COPY server.go server.go

RUN CGO_ENABLED=0 GOOS=linux go build -a -o server


FROM gcr.io/distroless/base

# copy the binary, now go is unnecessary
COPY --from=stage1 /usr/src/app/server server

EXPOSE 8080

#CMD ["./server"]
