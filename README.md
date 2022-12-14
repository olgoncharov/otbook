# Otbook

Прототип социальной сети для курса [Highload Architect](https://otus.ru/lessons/highloadarchitect/) от OTUS. Реализован в виде REST API и имеет следующие возможности:

* Регистрация профиля пользователя
* JWT авторизация
* Редактирование собственного профиля
* Просмотр профилей
* Поиск профилей по имени и фамилии
* Добавление в друзья других пользователей
* Просмотр друзей пользователя
* Создание постов
* Просмотр ленты постов

## Стек технологий
* Go 1.18
* MySQL 8.0
* Redis 7.0

## Инструкция по запуску

Для запуска приложения необходим `docker compose`.

1. Склонируйте репозиторий

```bash
git clone git@github.com:olgoncharov/otbook.git
```

2. Перейдите в склонированную папку с проектом

```bash
cd otbook
```

3. Выполните запуск docker контейнеров

```bash
docker compose up -d
```

4. Включите репликацию с помощью скрипта `scripts/replication/start_repl`

```bash
./scripts/replication/start_repl
```


Документация к API доступна по адресу http://localhost:7000/.

В папке `/examples/postman/` находится коллекция Postman с примерами запросов. В целях тестирования база уже наполнена данными о 1 000 000 профилях пользователей. Для выполнения запросов, доступных только авторизованным пользователям, необходимо получить JWT токен (endpoint `/api/v1/auth/login/`) и передавать его в заголовке `Authorization` в виде `Bearer <access token>`

По адресу http://localhost:8080/ доступна админка базы данных.