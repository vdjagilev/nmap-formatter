FROM golang:1.23-alpine AS builder

WORKDIR /src/
COPY . /src/

RUN apk add --update gcc musl-dev
RUN CGO_ENABLED=1 go build -o /bin/nmap-formatter

FROM golang:1.23-alpine

COPY --from=builder /bin/nmap-formatter /bin/nmap-formatter

ENTRYPOINT ["/bin/nmap-formatter"]
