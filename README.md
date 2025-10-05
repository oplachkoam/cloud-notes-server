# CloudNotes

Современный REST API сервер для управления заметками, построенный на Go с использованием
слоистой архитектуры и лучших практик разработки.

## Особенности

### Архитектура

- **Слоистая архитектура** - четкое разделение на слои: handlers, services, storage
- **Dependency Injection** - явная инверсия зависимостей для тестируемости
- **Interface Segregation** - использование интерфейсов для слабой связанности

### Безопасность

- **JWT аутентификация**
- **Хеширование паролей** с использованием bcrypt
- **Валидация входных данных** с помощью go-playground/validator
- **Middleware для безопасности** - проверка токенов и сессий
- **CORS защита** и другие security headers

### Технологический стек

- **Go 1.25** - основной язык разработки
- **PostgreSQL** - основная база данных
- **Redis** - кеширование и сессии
- **Chi Router** - быстрый HTTP роутер
- **Docker & Docker Compose** - контейнеризация
- **GolangCI-Lint** - статический анализ кода

### DevOps

- **Docker контейнеризация** - готовые к продакшену образы
- **Health checks** - мониторинг состояния сервисов
- **Миграции БД** - автоматическое управление схемой
- **Structured logging** - структурированное логирование
- **Environment configuration** - гибкая конфигурация через переменные окружения

## Структура проекта

```
cloud-notes-server/
├── cmd/                   # Точки входа приложения
│   ├── server/            # Основной сервер
│   └── migrator/          # Миграции БД
├── internal/              # Внутренняя логика приложения
│   ├── config/            # Конфигурация
│   ├── handlers/          # HTTP обработчики
│   ├── services/          # Бизнес-логика
│   ├── storage/           # Слой данных
│   ├── middleware/        # HTTP middleware
│   ├── security/          # Безопасность и JWT
│   └── logger/            # Логирование
├── migrations/            # SQL миграции
├── docker-compose.yml     # Docker окружение
└── Makefile               # Команды разработки
```

## Установка и запуск

### Предварительные требования

- Docker и Docker Compose
- Go 1.25+ (для разработки)

### Быстрый старт

1. **Клонируйте репозиторий**

```bash
git clone <repository-url>
cd cloud-notes-server
```

2. **Создайте файл окружения**

```bash
cp .env.example .env
# Отредактируйте .env файл с вашими настройками
```

3. **Запустите проект**

```bash
make start
```

### Доступные команды

```bash
make start       # Запуск всех сервисов
make start-d     # Запуск в фоновом режиме
make stop        # Остановка сервисов
make logs        # Просмотр логов сервера
make migrate     # Выполнение миграций
make lint        # Проверка кода линтером
make format      # Форматирование кода
make test        # Запуск тестов
make delete      # Удаление контейнеров
make delete-data # Удаление контейнеров и данных
```

## API Документация

### Аутентификация

#### Регистрация

```http
POST /api/auth/register
Content-Type: application/json

{
  "login": "user@example.com",
  "password": "securepassword",
  "first_name": "John",
  "timezone": "UTC"
}
```

#### Вход

```http
POST /api/auth/login
Content-Type: application/json

{
  "login": "user@example.com",
  "password": "securepassword"
}
```

#### Выход

```http
POST /api/auth/logout
Authorization: Bearer <access_token>
```

### Пользователи

#### Получение профиля

```http
GET /api/user/profile
Authorization: Bearer <access_token>
```

#### Обновление профиля

```http
PUT /api/user/profile
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "first_name": "John",
  "timezone": "Europe/Moscow"
}
```

### Заметки

#### Создание заметки

```http
POST /api/notes
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Моя заметка",
  "text": "Содержимое заметки",
  "pinned": false
}
```

#### Получение списка заметок

```http
GET /api/notes
Authorization: Bearer <access_token>
```

#### Обновление заметки

```http
PUT /api/notes/{note-id}
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "title": "Обновленная заметка",
  "text": "Новое содержимое",
  "pinned": true
}
```

#### Удаление заметки

```http
DELETE /api/notes/{note-id}
Authorization: Bearer <access_token>
```

## Планы развития

### TODO: Дальнейшие улучшения

- [ ] **Кеширование** - добавление Redis кеша для часто запрашиваемых данных
- [ ] **Система тегов** - возможность категоризации заметок с помощью тегов
- [ ] **Прикрепление файлов** - поддержка загрузки и прикрепления файлов к заметкам
- [ ] **Поиск по заметкам** - полнотекстовый поиск с использованием PostgreSQL
- [ ] **Экспорт/импорт** - возможность экспорта заметок в различные форматы
- [ ] **Совместное использование** - возможность делиться заметками с другими пользователями
- [ ] **API версионирование** - поддержка нескольких версий API
- [ ] **Rate limiting** - ограничение частоты запросов
- [ ] **Метрики и мониторинг** - интеграция с Prometheus/Grafana
- [ ] **Swagger документация** - автоматическая генерация API документации

## Автор

Телеграм: **@op1alex**

Создано как pet-проект для демонстрации современных практик разработки на Go.
