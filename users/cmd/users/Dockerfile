# указываем образ, на основе которого будем собирать приложение
FROM golang:1.21

# указываем рабочую папку внутри образа
WORKDIR /build

# установим переменные среды
ENV APP_PORT=3000
ENV LOG_INDEX=users
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

# копируем файлы приложения в рабочую папку образа
COPY app /build/app
COPY cmd /build/cmd
COPY config /build/config
COPY internal /build/internal
COPY pkg /build/pkg

COPY go.mod /build/
COPY go.sum /build/

# запускаем в образе команды
# скачиваем зависимости приложения
# собирвем приложение
# делаем бинарный файл исполняемым
RUN go mod download && \
    go build -o cmd/users/users ./cmd/users/main.go && \
    chmod +x cmd/users/users

# указываем бинарный файл
ENTRYPOINT ["cmd/users/users"]

# открываем порты из образа
EXPOSE 3000/tcp