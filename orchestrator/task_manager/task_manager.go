package task_manager

import (
	"errors"
	"sync"

	"github.com/mitka051209/arithmetic-service/orchestrator/models/models.go"
)

type TaskManager struct {
	expressions map[string]*models.Expression
	tasks       map[string]*models.Task
	taskQueue   []string // Queue of task IDs
	mu          sync.Mutex
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		expressions: make(map[string]*models.Expression),
		tasks:       make(map[string]*models.Task),
		taskQueue:   make([]string, 0),
	}
}

func (tm *TaskManager) AddExpression(expression *models.Expression) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.expressions[expression.ID] = expression
}

func (tm *TaskManager) GetExpression(id string) (*models.Expression, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	expression, ok := tm.expressions[id]
	return expression, ok
}

func (tm *TaskManager) GetAllExpressions() []models.Expression {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	expressions := make([]models.Expression, 0, len(tm.expressions))
	for _, expression := range tm.expressions {
		expressions = append(expressions, *expression)
	}
	return expressions
}

func (tm *TaskManager) AddTask(task *models.Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks[task.ID] = task
	tm.taskQueue = append(tm.taskQueue, task.ID)
}

func (tm *TaskManager) GetNextTask() (*models.Task, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if len(tm.taskQueue) == 0 {
		return nil, false
	}

	taskID := tm.taskQueue[0]
	tm.taskQueue = tm.taskQueue[1:] // Remove the task from the queue

	task, ok := tm.tasks[taskID]
	if !ok {
		return nil, false
	}

	return task, true
}

func (tm *TaskManager) CompleteTask(taskID string, result float64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, ok := tm.tasks[taskID]
	if !ok {
		return errors.New("task not found")
	}

	expression, ok := tm.expressions[task.ExpressionID]
	if !ok {
		return errors.New("expression not found")
	}

	expression.Result = &result // Save result
	expression.Status = "completed" // Update Status

	delete(tm.tasks, taskID) // Remove task

	return nil
}