# указываем образ, на основе которого будем собирать приложение
FROM golang:1.21

# указываем рабочую папку внутри образа
WORKDIR /build

# установим переменные среды
ENV APP_PORT=3000

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
    go build -o cmd/users ./cmd/main.go && \
    chmod +x cmd/users

# указываем бинарный файл
ENTRYPOINT ["cmd/users"]

# открываем порты из образа
EXPOSE 3000/tcp