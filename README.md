# arithmetic-service

# Арифметический сервис

Этот проект представляет собой распределённый вычислитель арифметических выражений. Пользователь может отправлять выражения на вычисление, а система обрабатывает их в фоновом режиме, возвращая результат по запросу. Вычисления выполняются агентами, которые могут масштабироваться для увеличения производительности.

---

## 🚀 Запуск проекта

### С помощью Go (локально)
1. Убедитесь, что установлен Go (версия 1.21+).
2. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/m1tka051209/arithmetic-service.git
   cd arithmetic-service
3. Запустите оркестратор:
   go run ./orchestrator/main.go
4. Запустите агент (в другом терминале):
   COMPUTING_POWER=3 go run ./agent/main.go
## Архитектура системы

### Оркестратор:

Принимает выражения от пользователей.

Разбивает выражения на задачи (например, 2 + 3, 3 * 4).

Управляет выполнением задач агентами.

### Агент:

Получает задачи от оркестратора.

Выполняет арифметические операции.

Возвращает результаты.


## Примеры использования

1. Добавление выражения для вычисления

200 - 
curl -X POST -H "Content-Type: application/json" -d '{
  "expression": "2 + 3 * 4"
}' http://localhost:8080/api/v1/calculate


{
  "id": "abc123"
}


422 - 
curl -X POST -H "Content-Type: application/json" -d '{
  "expression": "2 + * 4"
}' http://localhost:8080/api/v1/calculate


{
  "error": "invalid expression"
}


2. Получение списка всех выражений

200 -
curl --location 'http://localhost:8080/api/v1/expressions'


{
  "expressions": [
    {
      "id": "abc123",
      "status": "processing",
      "result": 0
    },
    {
      "id": "def456",
      "status": "completed",
      "result": 14
    }
  ]
}


3. Получение выражения по ID

200 -

curl --location 'http://localhost:8080/api/v1/expressions/abc123'

{
  "expression": {
    "id": "abc123",
    "status": "processing",
    "result": 0
  }
}


404 -

curl --location 'http://localhost:8080/api/v1/expressions/invalid_id'

{
  "error": "expression not found"
}

4. Получение задачи для выполнения (агент)

200 -

curl --location 'http://localhost:8080/internal/task'


{
  "task": {
    "id": "task123",
    "arg1": 2,
    "arg2": 3,
    "operation": "+",
    "operation_time": 1000
  }
}

404 -

curl --location 'http://localhost:8080/internal/task'


{
  "error": "no tasks available"
}

5. Отправка результата выполнения задачи (агент)

200 -

curl -X POST -H "Content-Type: application/json" -d '{
  "id": "task123",
  "result": 5
}' http://localhost:8080/internal/task


Пустое тело ответа.

422 -

curl -X POST -H "Content-Type: application/json" -d '{
  "id": "task123",
  "result": "invalid"
}' http://localhost:8080/internal/task


{
    "error":"invalid request body"
}


6. Ошибка сервера (500):


curl --location 'http://localhost:8080/api/v1/expressions'


{
  "error": "internal server error"
}
