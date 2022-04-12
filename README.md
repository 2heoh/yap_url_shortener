# go-musthave-shortener-tpl
Шаблон репозитория для практического трека "Go в веб-разработке"

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Postgres

docker run -d -e POSTGRES_PASSWORD=postgres -e PGDATA=/var/lib/postgresql/data/pgdata -v ${PWD}/db:/var/lib/postgresql/data -p5432:5432 postgres


DATABASE_DSN=postgres://postgres:postgres@localhost:5432/links go run cmd/shortener/main.go