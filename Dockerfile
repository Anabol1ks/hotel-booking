# Используем базовый образ Go для сборки приложения
FROM golang:1.23 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы приложения в контейнер
COPY . .

# Загружаем зависимости и собираем приложение
RUN go mod download
RUN go build -o app .

# Используем минимальный образ для выполнения приложения
FROM debian:bullseye-slim

# Устанавливаем необходимые инструменты
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем приложение из предыдущего этапа
COPY --from=builder /app/app .

# Копируем скрипт ожидания готовности базы данных
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Устанавливаем переменные среды
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=postgres
ENV POSTGRES_HOST=db
ENV POSTGRES_PORT=5432

# Открываем порт для приложения
EXPOSE 8080

# Команда запуска: ждем готовности базы данных и запускаем приложение
CMD ["/wait-for-it.sh", "db:5432", "--", "./app"]
