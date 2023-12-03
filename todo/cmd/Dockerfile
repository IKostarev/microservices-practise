FROM golang:1.21

WORKDIR /build

ENV APP_PORT=":3000"

COPY app /build/app
COPY cmd /build/cmd
COPY config /build/config
COPY internal /build/internal

COPY go.mod /build/
COPY go.sum /build/

RUN go mod download && \
    go build -o cmd/todo ./cmd/main.go && \
    chmod +x cmd/todo

ENTRYPOINT ["cmd/todo"]

EXPOSE 3000/tcp
