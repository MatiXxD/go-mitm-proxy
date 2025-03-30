# Go mitmproxy

## Как запустить

Сначала надо сгенерировать сертификаты:

```bash
make gen
```

Перед тем как запускать эту команду надо убедиться, что в директории проекта есть папка `certs`, если нету создать: `mkdir certs`.

После того, как сертификаты сгенерированы можно запустить проект через `docker compose`:

```bash
make docker-compose-build && make docker-compose-up
```

Чтобы остановить контейнеры:

```bash
make docker-compose-stop
```

## Примеры запросов

**HTTP запрос:**

```bash
curl -x http://127.0.0.1:8080 http://mail.ru
```

**HTTPS запрос:**

```bash
curl -k -x http://127.0.0.1:8080 https://mail.ru
```
