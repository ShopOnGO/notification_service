#!/bin/sh

echo "⏳ Ожидание MongoDB на $MONGO_URI..."

# Парсим хост и порт даже если есть авторизация
host=$(echo "$MONGO_URI" | sed -E 's|.*@([^:/]+).*|\1|')
port=$(echo "$MONGO_URI" | sed -E 's|.*:([0-9]+)/.*|\1|')

echo "Ждём подключения к $host:$port..."

while ! nc -z "$host" "$port"; do
  echo "MongoDB ещё не готов — ждём..."
  sleep 2
done

echo "✅ MongoDB доступен, запускаем сервис..."

exec "$@"
