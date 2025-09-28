package jobs

import (
    "log"
)

type Job func() error

var jobQueue = make(chan Job, 100) // fila com até 100 jobs

// StartWorker inicializa um worker que processa jobs em background
func StartWorker() {
    go func() {
        for job := range jobQueue {
            if err := job(); err != nil {
                log.Println("Erro no job:", err)
            } else {
                log.Println("Job executado com sucesso")
            }
        }
    }()
}

// AddJob adiciona um job na fila
func AddJob(job Job) {
    jobQueue <- job
}
