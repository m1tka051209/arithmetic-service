package task_manager

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
	
	"github.com/m1tka051209/arithmetic-service/orchestrator/models"
)

const (
	idLength = 8
	charset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type TaskManager struct {
	expressions   map[string]models.Expression
	tasks         map[string]models.Task
	mu            sync.RWMutex
	idMu          sync.Mutex
	rand          *rand.Rand
	operationTime map[string]time.Duration
}

func NewTaskManager() *TaskManager {
	src := rand.NewSource(time.Now().UnixNano())
	return &TaskManager{
		expressions: make(map[string]models.Expression),
		tasks:       make(map[string]models.Task),
		rand:        rand.New(src),
		operationTime: map[string]time.Duration{
			"+": 100 * time.Millisecond,
			"-": 100 * time.Millisecond,
			"*": 200 * time.Millisecond,
			"/": 200 * time.Millisecond,
		},
	}
}

func (tm *TaskManager) GenerateID() string {
	tm.idMu.Lock()
	defer tm.idMu.Unlock()

	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[tm.rand.Intn(len(charset))]
	}
	return string(b)
}

func (tm *TaskManager) ParseExpression(expr string) ([]models.Task, error) {
	tokens := strings.Fields(expr)
	if len(tokens)%2 == 0 || len(tokens) < 3 {
		return nil, fmt.Errorf("invalid expression format")
	}

	var tasks []models.Task
	for i := 1; i < len(tokens); i += 2 {
		task := models.Task{
			ID:        tm.GenerateID(),
			Arg1:      parseNumber(tokens[i-1]),
			Arg2:      parseNumber(tokens[i+1]),
			Operation: tokens[i],
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (tm *TaskManager) SaveExpression(id string, tasks []models.Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.expressions[id] = models.Expression{
		ID:     id,
		Status: "processing",
	}

	for _, t := range tasks {
		t.OperationTime = tm.operationTime[t.Operation]
		t.Status = "pending"
		tm.tasks[t.ID] = t
	}
}

func (tm *TaskManager) GetNextTask() (models.Task, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.tasks {
		if task.Status == "pending" {
			task.Status = "in_progress"
			tm.tasks[task.ID] = task
			return task, true
		}
	}
	return models.Task{}, false
}

func (tm *TaskManager) SaveTaskResult(taskID string, result float64) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if task, exists := tm.tasks[taskID]; exists {
		task.Result = result
		task.Status = "completed"
		tm.tasks[taskID] = task
	}
}

func (tm *TaskManager) GetAllExpressions() []models.Expression {
    tm.mu.RLock()
    defer tm.mu.RUnlock()

    result := make([]models.Expression, 0, len(tm.expressions))
    for _, expr := range tm.expressions {
        result = append(result, expr)
    }
    return result
}

func (tm *TaskManager) GetExpressionByID(id string) (models.Expression, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	expr, exists := tm.expressions[id]
	return expr, exists
}

func parseNumber(s string) float64 {
	var num float64
	_, err := fmt.Sscanf(s, "%f", &num)
	if err != nil {
		return 0
	}
	return num
}