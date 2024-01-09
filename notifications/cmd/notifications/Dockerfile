# указываем образ, на основе которого будем собирать приложение
FROM golang:1.21

# указываем рабочую папку внутри образа
WORKDIR /build

# установим переменные среды
ENV LOG_INDEX=notifications
ENV LOG_IS_DEBUG=false
ENV LOG_TO_FILE=false

ENV SMTP_HOST=smtp.yandex.ru
ENV SMTP_PORT=25
ENV SMTP_USER=course.notifications@yandex.ru
ENV SMTP_PASSWORD=ulyakmzpjokpcpiz

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
    go build -o cmd/notifications/notifications ./cmd/notifications/main.go && \
    chmod +x cmd/notifications/notifications

# указываем бинарный файл
ENTRYPOINT ["cmd/notifications/notifications"]

# открываем порты из образа
EXPOSE 5672/tcp
EXPOSE 465/tcp


