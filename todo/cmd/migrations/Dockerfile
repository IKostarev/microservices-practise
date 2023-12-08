FROM golang:1.21

WORKDIR /build

ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5430
ENV POSTGRES_USER=ibs_test
ENV POSTGRES_PASSWORD=ibs_test
ENV POSTGRES_DBNAME=ibs_test

COPY app /build/app
COPY cmd /build/cmd
COPY config /build/config
COPY internal /build/internal
COPY pkg /build/pkg

COPY go.mod /build/
COPY go.sum /build/

RUN go mod download && \
        go build -o cmd/migrations ./cmd/migrations/main.go && \
        chmod +x cmd/migrations

CMD ["cmd/migrations"]