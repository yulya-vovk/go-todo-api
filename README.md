# To-Do API на Go

Простое RESTful API для управления задачами, написанное на чистом Go без фреймворков.

## Возможности
- GET /tasks — получить список задач
- POST /tasks — добавить новую задачу
- PUT /tasks/:id — обновить задачу
- DELETE /tasks/:id — удалить задачу
- Данные сохраняются в `data/tasks.json`

## Как запустить
1. Клонируй репозиторий
2. Создай папку `data`
3. Запусти: `go run main.go`
4. Сервер будет доступен на `http://localhost:8080`

## Пример запроса
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Новая задача","done":false}'