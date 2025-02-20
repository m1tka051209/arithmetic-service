package task_manager

import (
	"testing"

	"github.com/mitka051209/arithmetic-service/orchestrator/models"
	"github.com/stretchr/testify/assert"
)

func TestTaskManagement(t *testing.T) {
	// Тестирование сохранения выражения
	tasks := []models.Task{
		{ID: "task1", Operation: "+", Arg1: 2, Arg2: 3},
		{ID: "task2", Operation: "*", Arg1: 4, Arg2: 5},
	}

	SaveExpression("expr1", tasks)

	// Проверка сохранения задач
	assert.Equal(t, 2, len(tasks), "Должно быть 2 задачи")

	// Проверка получения задачи
	task, exists := GetNextTask()
	assert.True(t, exists, "Задача должна существовать")
	assert.Equal(t, "task1", task.ID, "Неверный ID задачи")

	// Обновление статуса задачи
	SaveTaskResult("task1", 5)
	updatedTask, _ := tasks["task1"]
	assert.Equal(t, "completed", updatedTask.Status, "Статус должен быть 'completed'")
}
