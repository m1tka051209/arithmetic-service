package task_manager

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
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
			"+": getDurationFromEnv("TIME_ADDITION_MS", 1000),
			"-": getDurationFromEnv("TIME_SUBTRACTION_MS", 1000),
			"*": getDurationFromEnv("TIME_MULTIPLICATION_MS", 2000),
			"/": getDurationFromEnv("TIME_DIVISION_MS", 2000),
		},
	}
}

func getDurationFromEnv(envVar string, defaultVal int) time.Duration {
	valStr := os.Getenv(envVar)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return time.Duration(defaultVal) * time.Millisecond
	}
	return time.Duration(val) * time.Millisecond
}

func shuntingYard(expr string) ([]string, error) {
	var output []string
	var stack []string
	tokens := regexp.MustCompile(`(-?\d+\.?\d*)|([+\-*/()])`).FindAllString(expr, -1)
	precedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	for _, token := range tokens {
		switch {
		case isNumber(token):
			output = append(output, token)
		case token == "(":
			stack = append(stack, token)
		case token == ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1]
		default:
			for len(stack) > 0 && precedence[token] <= precedence[stack[len(stack)-1]] && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return output, nil
}

func (tm *TaskManager) ParseExpression(expr string) ([]models.Task, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	rpn, err := shuntingYard(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid expression: %w", err)
	}

	var stack []float64
	var tasks []models.Task
	exprID := tm.GenerateID()

	for _, token := range rpn {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return nil, fmt.Errorf("invalid expression")
			}
			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			task := models.Task{
				ID:            tm.GenerateID(),
				Arg1:          arg1,
				Arg2:          arg2,
				Operation:     token,
				OperationTime: tm.operationTime[token],
				Status:        "pending",
				ExpressionID:  exprID,
			}
			tasks = append(tasks, task)
			stack = append(stack, calculateTask(task))
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid expression")
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.expressions[exprID] = models.Expression{
		ID:     exprID,
		Status: "processing",
		Result: stack[0],
	}

	for _, task := range tasks {
		tm.tasks[task.ID] = task
	}

	return tasks, nil
}

// GetAllExpressions возвращает список всех выражений
func (tm *TaskManager) GetAllExpressions() []models.Expression {
    tm.mu.RLock()
    defer tm.mu.RUnlock()

    expressions := make([]models.Expression, 0, len(tm.expressions))
    for _, expr := range tm.expressions {
        // Проверяем статус выражения
        allTasksCompleted := true
        for _, task := range tm.tasks {
            if task.ExpressionID == expr.ID && task.Status != "completed" {
                allTasksCompleted = false
                break
            }
        }
        if allTasksCompleted {
            expr.Status = "completed"
        }
        expressions = append(expressions, expr)
    }
    return expressions
}


// GetExpressionByID возвращает выражение по ID
func (tm *TaskManager) GetExpressionByID(id string) (models.Expression, bool) {
    tm.mu.RLock()
    defer tm.mu.RUnlock()

    expr, exists := tm.expressions[id]
    return expr, exists
}

func calculateTask(task models.Task) float64 {
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return 0
		}
		return task.Arg1 / task.Arg2
	default:
		return task.Arg1
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

func (tm *TaskManager) SaveExpression(id string, tasks []models.Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.expressions[id] = models.Expression{
		ID:     id,
		Status: "processing",
	}

	for _, t := range tasks {
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

		// Проверка завершения всех задач выражения
		exprID := task.ExpressionID
		allCompleted := true
		for _, t := range tm.tasks {
			if t.ExpressionID == exprID && t.Status != "completed" {
				allCompleted = false
				break
			}
		}

		if allCompleted {
			expr := tm.expressions[exprID]
			expr.Status = "completed"
			tm.expressions[exprID] = expr
		}
	}
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}