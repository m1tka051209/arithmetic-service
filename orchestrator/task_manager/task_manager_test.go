package task_manager

import (
	"testing"
	// "time"

	"github.com/m1tka051209/arithmetic-service/orchestrator/models"
	"github.com/stretchr/testify/assert"
)

func TestShuntingYard(t *testing.T) {
	// tm := NewTaskManager()
	expr := "2 + 3 * (4 - 1)"
	rpn, err := shuntingYard(expr)
	assert.NoError(t, err)
	assert.Equal(t, []string{"2", "3", "4", "1", "-", "*", "+"}, rpn)
}

func TestExpressionCompletion(t *testing.T) {
	tm := NewTaskManager()
	tasks, err := tm.ParseExpression("(5 + 3) * 2")
	assert.NoError(t, err)

	for _, task := range tasks {
		tm.SaveTaskResult(task.ID, calculateTask(task))
	}

	expr, exists := tm.GetExpressionByID(tasks[0].ExpressionID)
	assert.True(t, exists)
	assert.Equal(t, "completed", expr.Status)
	assert.Equal(t, 16.0, expr.Result)
}

func TestDivisionByZero(t *testing.T) {
	tm := NewTaskManager()
	_, err := tm.ParseExpression("5 / 0")
	assert.ErrorContains(t, err, "division by zero")
}

// ... (остальной код теста без изменений)

// Добавить мок для методов TaskManager если нужно
func TestGetExpressionByID(t *testing.T) {
    tm := NewTaskManager()
    exprID := "test123"
    tm.expressions[exprID] = models.Expression{
        ID:     exprID,
        Status: "processing",
    }

    expr, exists := tm.GetExpressionByID(exprID)
    assert.True(t, exists)
    assert.Equal(t, "processing", expr.Status)
}

func TestTaskManagement(t *testing.T) {
	tm := NewTaskManager()
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

	exprID := tm.GenerateID()
	tm.SaveExpression(exprID, tasksToSave)

	task, exists := tm.GetNextTask()
	assert.True(t, exists)
	assert.Equal(t, "+", task.Operation)

	tm.SaveTaskResult(task.ID, 5.0)
	updatedTask, exists := tm.tasks[task.ID]
	assert.True(t, exists)
	assert.Equal(t, "completed", updatedTask.Status)

	task2, exists := tm.GetNextTask()
	assert.True(t, exists)
	assert.Equal(t, "*", task2.Operation)
}