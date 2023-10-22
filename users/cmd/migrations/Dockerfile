# указываем образ, на основе которого будем собирать приложение
FROM golang:1.21

# указываем рабочую папку внутри образа
WORKDIR /build

ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DBNAME=postgres

# копируем файлы приложения в рабочую папку образа
COPY cmd/migrations /build/cmd/migrations
COPY migrations /build/migrations
COPY config /build/config
COPY pkg /build/pkg

COPY go.mod /build/
COPY go.sum /build/

# запускаем в образе команды
# скачиваем зависимости приложения
# собирвем приложение
# делаем бинарный файл исполняемым
RUN go mod download && \
        go build -o cmd/migrations/migrations ./cmd/migrations/main.go && \
        chmod +x cmd/migrations/migrations

# указываем бинарный файл
CMD ["cmd/migrations/migrations"]
