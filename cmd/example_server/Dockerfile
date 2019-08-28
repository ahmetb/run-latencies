FROM golang:1-alpine
WORKDIR /src
COPY . .
RUN env CGO_ENABLED=0 go build -tags netgo -o example_server.out .

FROM alpine:latest
COPY --from=0 /src/example_server.out /bin/example_server
ENTRYPOINT ["/bin/example_server"]
