package handlers

import (
    "fmt"
    "io"       // para io.Copy
    "net/http"
    "os"       // para os.Create, os.Remove, etc.

    "challenge_enube/internal/usecase"

    "challenge_enube/internal/jobs"
)

type ImportHandler struct {
    importUseCase *usecase.ImportUseCase
}

func NewImportHandler(uc *usecase.ImportUseCase) *ImportHandler {
    return &ImportHandler{importUseCase: uc}
}
func (h *ImportHandler) HandlerExcel(w http.ResponseWriter, r *http.Request) {
    file, _, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Erro ao ler arquivo", http.StatusBadRequest)
        return
    }
    defer file.Close()

    tmpPath := "/tmp/upload.xlsx"
    out, err := os.Create(tmpPath)
    if err != nil {
        http.Error(w, "Erro ao criar arquivo temporário", http.StatusInternalServerError)
        return
    }
    defer out.Close()

    if _, err = io.Copy(out, file); err != nil {
        http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
        return
    }

    // Adiciona job na fila
    jobs.AddJob(func() error {
        return h.importUseCase.ImportFromExcel(tmpPath)
    })

    // Responde imediatamente
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "Arquivo recebido! Processamento em background.")
}
