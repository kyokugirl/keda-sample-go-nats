FROM golang:1.23 AS builder
ENV CGO_ENABLED=0

FROM builder AS builder-send
WORKDIR /send-src
COPY ./send/go.mod ./send/go.sum ./
RUN go mod download
COPY ./send/* ./
RUN go build -o /send

FROM builder AS builder-receive
WORKDIR /receive-src
COPY ./receive/go.mod ./receive/go.sum ./
RUN go mod download
COPY ./receive/* ./
RUN go build -o /receive

FROM scratch
COPY --from=builder-send /send /usr/local/bin/
COPY --from=builder-receive /receive /usr/local/bin/
CMD ["receive"]
