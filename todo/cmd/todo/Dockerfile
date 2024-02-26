FROM golang:1.21

WORKDIR /build

ENV APP_PORT=3000
ENV LOG_INDEX=todos
ENV LOG_IS_DEBUG=false
ENV LOG_TO_FILE=false
ENV LOG_TO_FILE=false

ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DBNAME=postgres

ENV RABBITMQ_USER=notifications
ENV RABBITMQ_PASSWORD=notifications
ENV RABBITMQ_HOST=rabbitmq
ENV RABBITMQ_PORT=5672

ENV JAEGER_HOST=jaeger:6831
ENV JAEGER_LOG_SPANS=true

ENV REDIS_ADDRES=redis:6379
ENV REDIS_PASSWORD=password
ENV REDIS_DB=2

COPY app /build/app
COPY cmd /build/cmd
COPY config /build/config
COPY internal /build/internal
COPY pkg /build/pkg

COPY go.mod /build/
COPY go.sum /build/

RUN go mod download && \
    go build -o cmd/todo/todo ./cmd/todo/main.go && \
    chmod +x cmd/todo/todo

ENTRYPOINT ["cmd/todo/todo"]

EXPOSE 3000/tcp
