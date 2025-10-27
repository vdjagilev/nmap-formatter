FROM golang:1.25-alpine AS builder

WORKDIR /src/
COPY . /src/

RUN apk add --update gcc musl-dev
RUN CGO_ENABLED=1 go build -o /bin/nmap-formatter

FROM alpine:3.22

COPY --from=builder /bin/nmap-formatter /bin/nmap-formatter

ENTRYPOINT ["/bin/nmap-formatter"]
