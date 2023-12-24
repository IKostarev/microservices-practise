FROM golang:1.21 as builder

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./

RUN go build -o ./bin/app cmd/migrations/main.go


FROM alpine:3.15

COPY --from=builder /usr/local/src/bin/app /

CMD ["/app"]