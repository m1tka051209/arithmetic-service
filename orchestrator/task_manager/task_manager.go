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
    "log"
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
            "+": getDurationFromEnv("TIME_ADDITION_MS", 100),
            "-": getDurationFromEnv("TIME_SUBTRACTION_MS", 100),
            "*": getDurationFromEnv("TIME_MULTIPLICATION_MS", 200),
            "/": getDurationFromEnv("TIME_DIVISION_MS", 200),
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
    expr = strings.ReplaceAll(expr, " ", "")
    
    // Улучшенное регулярное выражение
    re := regexp.MustCompile(`(-?\d+\.?\d*)|([+*\/-])`)
    matches := re.FindAllString(expr, -1)

    // Логирование токенов
    log.Printf("Original tokens: %v", matches)

    // Обработка унарных операторов
    var cleaned []string
    for i := 0; i < len(matches); i++ {
        if (matches[i] == "+" || matches[i] == "-") && (i == 0 || isOperator(matches[i-1])) {
            if i+1 < len(matches) && isNumber(matches[i+1]) {
                cleaned = append(cleaned, matches[i]+matches[i+1])
                i++
                continue
            }
        }
        cleaned = append(cleaned, matches[i])
    }

    // Логирование после обработки
    log.Printf("Cleaned tokens: %v", cleaned)

    // Проверка количества токенов
    if len(cleaned) < 3 || len(cleaned)%2 == 0 {
        return nil, fmt.Errorf("неверное количество токенов: %d", len(cleaned))
    }

    var tasks []models.Task
    exprID := tm.GenerateID()
    
    for i := 1; i < len(cleaned); i += 2 {
        op := cleaned[i]
        if !isValidOperator(op) {
            return nil, fmt.Errorf("неподдерживаемая операция: %s", op)
        }

        arg1, err := parseNumber(cleaned[i-1])
        if err != nil {
            return nil, fmt.Errorf("ошибка парсинга аргумента 1: %v", err)
        }

        arg2, err := parseNumber(cleaned[i+1])
        if err != nil {
            return nil, fmt.Errorf("ошибка парсинга аргумента 2: %v", err)
        }

        tasks = append(tasks, models.Task{
            ID:            tm.GenerateID(),
            Arg1:          arg1,
            Arg2:          arg2,
            Operation:     op,
            OperationTime: tm.operationTime[op],
            Status:        "pending",
            ExpressionID:  exprID,
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

func isOperator(s string) bool {
    return s == "+" || s == "-" || s == "*" || s == "/"
}

func isNumber(s string) bool {
    _, err := strconv.ParseFloat(s, 64)
    return err == nil
}

// Остальные методы (SaveExpression, GetNextTask и т.д.) остаются без изменений