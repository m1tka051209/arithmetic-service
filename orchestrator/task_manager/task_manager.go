package task_manager

import (
    "sync"
    "time"
    "github.com/mitka051209/arithmetic-service/config"
	"github.com/mitka051209/arithmetic-service/models"
)

var (
    conf = config.Load() // Инициализация конфига

    operationTime = map[string]time.Duration{
        "+": time.Duration(conf.AdditionTime) * time.Millisecond,
        "-": time.Duration(conf.SubtractionTime) * time.Millisecond,
        "*": time.Duration(conf.MultiplicationTime) * time.Millisecond,
        "/": time.Duration(conf.DivisionTime) * time.Millisecond,
    }

    expressions = make(map[string]models.Expression)
    tasks       = make(map[string]models.Task)
    mu          sync.RWMutex
)