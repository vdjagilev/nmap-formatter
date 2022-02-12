FROM golang:1.17.7-alpine AS build

WORKDIR /src/
COPY . /src/
RUN CGO_ENABLED=0 go build -o /bin/nmap-formatter

FROM scratch
COPY --from=build /bin/nmap-formatter /bin/nmap-formatter

ENTRYPOINT ["/bin/nmap-formatter"]
