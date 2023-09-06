FROM golang:1.20 as build

COPY . /src

WORKDIR /src
RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o kv-store

FROM scratch

COPY --from=build kv-store .
EXPOSE 8080

CMD ["./kv-store"]
