package task_manager

import (
	"fmt"
	"math/rand"
	"regexp"
	// "strconv"
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
    // Удаляем все пробелы
    expr = strings.ReplaceAll(expr, " ", "")
    
    // Проверяем на валидные символы и структуру
    if !regexp.MustCompile(`^(\d+|(\d+\.\d+))([+\-*/](\d+|(\d+\.\d+)))+$`).MatchString(expr) {
        return nil, fmt.Errorf("неверный формат выражения")
    }

    // Разбиваем на токены с учетом чисел и операторов
    re := regexp.MustCompile(`(\d+\.?\d*)|([+\-*/])`)
    matches := re.FindAllString(expr, -1)
    
    if len(matches)%2 == 0 || len(matches) < 3 {
        return nil, fmt.Errorf("неверное количество токенов")
    }

    var tasks []models.Task
    for i := 1; i < len(matches); i += 2 {
        op := matches[i]
        if !isValidOperator(op) {
            return nil, fmt.Errorf("неподдерживаемая операция: %s", op)
        }

        arg1, err := parseNumber(matches[i-1])
        if err != nil {
            return nil, fmt.Errorf("ошибка парсинга левого операнда: %v", err)
        }

        arg2, err := parseNumber(matches[i+1])
        if err != nil {
            return nil, fmt.Errorf("ошибка парсинга правого операнда: %v", err)
        }

        tasks = append(tasks, models.Task{
            ID:            tm.GenerateID(),
            Arg1:          arg1,
            Arg2:          arg2,
            Operation:     op,
            OperationTime: tm.operationTime[op],
            Status:        "pending",
        })
    }
    return tasks, nil
}

func parseNumber(s string) (float64, error) {
    num, err := strconv.ParseFloat(s, 64)
    if err != nil {
        return 0, fmt.Errorf("неверный формат числа: %s", s)
    }
    return num, nil
}

func isValidOperator(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/"
}

func filterEmpty(tokens []string) []string {
	var result []string
	for _, t := range tokens {
		if t != "" {
			result = append(result, t)
		}
	}
	return result
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

func parseNumber(s string) (float64, error) {
	var num float64
	_, err := fmt.Sscanf(s, "%f", &num)
	if err != nil {
		return 0, fmt.Errorf("неверный формат числа: %s", s)
	}
	return num, nil
}

func ValidateTasks(tasks []models.Task) error {
	for _, task := range tasks {
		if task.Operation == "/" && task.Arg2 == 0 {
			return fmt.Errorf("деление на ноль в задаче %s", task.ID)
		}
	}
	return nil
}