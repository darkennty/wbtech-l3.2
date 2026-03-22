# URL Shortener

Сервис для сокращения URL с аналитикой переходов. Проект написан на Go с использованием Gin, PostgreSQL и Redis.

## Функциональность

- **Создание коротких ссылок** - генерация случайного короткого кода или использование пользовательского
- **Редирект** - перенаправление с короткой ссылки на оригинальный URL
- **Аналитика** - сбор и отображение статистики переходов:
  - Общее количество кликов
  - Время переходов
  - User-Agent браузера
  - Агрегация по дням, месяцам или User-Agent

## Технологический стек

- **Язык**: Go 1.25
- **Веб-фреймворк**: [Gin](https://gin-gonic.com/) + [wbf](https://github.com/wb-go/wbf)
- **База данных**: PostgreSQL
- **Кэширование**: Redis
- **Контейнеризация**: Docker, Docker Compose
- **Миграции**: [golang-migrate](https://github.com/golang-migrate/migrate)

## Структура проекта

```
WBTech_L3.2/
├── cmd/app/              # Точка входа в приложение
│   └── main.go
├── internal/             # Внутренние пакеты
│   ├── app/              # Инициализация и запуск приложения
│   ├── api/              # HTTP API
│   │   ├── handler/      # Обработчики запросов
│   │   └── server/       # Настройка HTTP сервера
│   ├── cache/            # Кэширование (Redis)
│   ├── config/           # Конфигурация приложения
│   ├── model/            # Модели данных
│   ├── repository/       # Слой доступа к данным
│   └── service/          # Бизнес-логика
├── migrations/           # Миграции базы данных
│   ├── 20260319194604_init.down.sql
│   ├── 20260319194604_init.up.sql
│   └── create_db.sql
├── web/                  # Фронтенд
│   ├── index.html
│   ├── app.js
│   └── styles.css
├── docker-compose.yml    # Docker Compose конфигурация
├── go.mod               # Зависимости Go
├── go.sum               # Суммы зависимостей
└── .env                 # Переменные окружения (создать из .env.example)
```

## Быстрый старт

### Предварительные требования

- [Docker](https://www.docker.com/) и Docker Compose
- [Go](https://golang.org/dl/) 1.25+ (для локальной разработки)

### Запуск

1. Создайте файл `.env` на основе примера из `.env.example`

2. Запустите сервисы:

```bash
docker-compose up -d
```

3. Запустите приложение:

```bash
go mod download
go run cmd/app/main.go
```

## API Endpoints

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/shorten` | Создание короткой ссылки |
| GET | `/s/:short_url` | Редирект на оригинальный URL |
| GET | `/analytics/:short_url` | Получение статистики |
| GET | `/` | Веб-интерфейс |

### Создание короткой ссылки

**Запрос:**

```http
POST /shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url",
  "desired_short_url": "custom123"  // опционально
}
```

**Ответ:**

```json
{
  "short_url": "custom123"
}
```

### Получение аналитики

**Запрос:**

```http
GET /analytics/custom123?aggregateBy=day
```

Параметр `aggregateBy` (опционально):
- `day` - группировка по дням
- `month` - группировка по месяцам
- `user_agent` - группировка по User-Agent

## Конфигурация

Конфигурация загружается из переменных окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `APP_HTTP_ADDR` | Адрес HTTP сервера | `8080` |
| `APP_POSTGRES_USER` | Пользователь PostgreSQL | `postgres` |
| `APP_POSTGRES_PASS` | Пароль PostgreSQL | `postgres` |
| `APP_POSTGRES_HOST` | Хост PostgreSQL | `localhost` |
| `APP_POSTGRES_PORT` | Порт PostgreSQL | `5432` |
| `APP_POSTGRES_DB` | Имя базы данных | `delayed_notifier` |
| `APP_POSTGRES_SSL_MODE` | Режим SSL PostgreSQL | `disable` |
| `APP_REDIS_ADDR` | Адрес Redis | `localhost:6379` |
| `APP_REDIS_PASSWORD` | Пароль Redis | `` |
| `APP_REDIS_DB` | Номер базы Redis | `0` |
| `APP_RETRY_BASE_DELAY` | Базовая задержка повторов | `1m` |
| `APP_RETRY_MAX` | Максимальное количество повторов | `5` |
