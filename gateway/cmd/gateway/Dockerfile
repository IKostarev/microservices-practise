# указываем образ, на основе которого будем собирать приложение
FROM golang:1.21

# указываем рабочую папку внутри образа
WORKDIR /build

# установим переменные среды
ENV APP_HOST=0.0.0.0
ENV APP_PORT=3000
ENV LOG_INDEX=gateway
ENV LOG_IS_DEBUG=false
ENV LOG_TO_FILE=false

ENV USERS_HOST=users-service
ENV USERS_REST_PORT=3000
ENV USERS_GRPC_PORT=50000

ENV JAEGER_HOST=jaeger:6831
ENV JAEGER_LOG_SPANS=true

ENV REDIS_ADDRES=redis:6379
ENV REDIS_PASSWORD=redisPassword
ENV REDIS_JWTDB=0

# копируем файлы приложения в рабочую папку образа
COPY app /build/app
COPY cmd /build/cmd
COPY config /build/config
COPY internal /build/internal
COPY pkg /build/pkg

COPY go.mod /build/
COPY go.sum /build/


# скачиваем зависимости приложения
RUN go mod download

# запускаем тесты
RUN go test ./...

# запускаем в образе команды
# скачиваем зависимости приложения
# собирвем приложение
# делаем бинарный файл исполняемым
RUN go build -o cmd/gateway/gateway ./cmd/gateway/main.go && \
    chmod +x cmd/gateway/gateway

# указываем бинарный файл
ENTRYPOINT ["cmd/gateway/gateway"]

# открываем порты из образа
EXPOSE 3000/tcp
EXPOSE 50000/tcp
