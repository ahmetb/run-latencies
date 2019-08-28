FROM golang:1-alpine
WORKDIR /src
COPY ../probe_server .
RUN env CGO_ENABLED=0 go build -tags netgo -o probe_server.out .

FROM alpine:latest
COPY --from=0 /src/probe_server.out /bin/probe_server
ENV PORT 8080
ENTRYPOINT ["/bin/probe_server"]
