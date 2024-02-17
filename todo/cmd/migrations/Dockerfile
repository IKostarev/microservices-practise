FROM golang:1.21

WORKDIR /build

ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DBNAME=postgres

COPY cmd/migrations /build/cmd/migrations
COPY migrations /build/migrations
COPY config /build/config
COPY pkg /build/pkg

COPY go.mod /build/
COPY go.sum /build/

RUN go mod download && \
        go build -o cmd/migrations/migrations ./cmd/migrations/main.go && \
        chmod +x cmd/migrations/migrations

CMD ["cmd/migrations/migrations"]