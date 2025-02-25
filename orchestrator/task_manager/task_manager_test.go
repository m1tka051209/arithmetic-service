package task_manager

import (
    "testing"
    // "time"
    
    "github.com/m1tka051209/arithmetic-service/orchestrator/models"
    "github.com/stretchr/testify/assert"
)

func TestTaskManagement(t *testing.T) {
    // Создаем менеджер задач
    tm := NewTaskManager()

    // Тест сохранения выражения
    tasksToSave := []models.Task{
        {
            ID:        tm.GenerateID(),
            Arg1:      2,
            Arg2:      3,
            Operation: "+",
        },
        {
            ID:        tm.GenerateID(),
            Arg1:      4,
            Arg2:      5,
            Operation: "*",
        },
    }

    // Сохраняем выражение
    exprID := tm.GenerateID()
    tm.SaveExpression(exprID, tasksToSave)

    // Проверяем получение задачи
    task, exists := tm.GetNextTask()
    assert.True(t, exists, "Задача должна существовать")
    assert.Equal(t, "+", task.Operation, "Неверная операция")

    // Сохраняем результат
    tm.SaveTaskResult(task.ID, 5.0)

    // Проверяем обновление статуса
    updatedTask, exists := tm.tasks[task.ID]
    assert.True(t, exists, "Задача должна существовать")
    assert.Equal(t, "completed", updatedTask.Status, "Статус должен быть 'completed'")

    // Проверяем получение следующей задачи
    task2, exists := tm.GetNextTask()
    assert.True(t, exists, "Вторая задача должна существовать")
    assert.Equal(t, "*", task2.Operation, "Неверная операция второй задачи")
}

func TestParseExpression(t *testing.T) {
    tm := NewTaskManager()
    tasks, err := tm.ParseExpression("2 + 3 * 4")
    assert.NoError(t, err, "Ошибка парсинга выражения")
    assert.Len(t, tasks, 2, "Должно быть 2 задачи")
    assert.Equal(t, "+", tasks[0].Operation, "Первая операция должна быть '+'")
    assert.Equal(t, "*", tasks[1].Operation, "Вторая операция должна быть '*'")
}