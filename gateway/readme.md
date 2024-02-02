### GATEWAY

Для генерации документации выполните в консоли в директории которая содержит этот файл команду:
```shell
  cd internal/api/rest && swag init --parseDependency --parseInternal --parseDepth 5 -g server.go -o ./../docs
```

Документация доступна по адресу: http://127.0.0.1:3000/docs/swagger/index.html#/

[Main README](../README.md)