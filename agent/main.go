package main

import (
    "os"
    "strconv"
    
    "github.com/m1tka051209/arithmetic-service/agent/worker"
)

func main() {
    // Получаем значение переменной окружения и конвертируем в int
    powerStr := os.Getenv("COMPUTING_POWER")
    power, err := strconv.Atoi(powerStr)
    if err != nil || power < 1 {
        power = 1 // Значение по умолчанию
    }
    
    worker.StartWorkers(power)
}