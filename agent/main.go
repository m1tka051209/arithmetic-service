package main

import (
	"log"
	"os"

	"github.com/m1tka051209/arithmetic-service/agent/worker"
)

func main() {
	log.Println("Агент запущен")
	worker.StartWorkers(os.Getenv("COMPUTING_POWER"))
}
