package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"challenge_enube/internal/usecase"
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
	out, _ := os.Create(tmpPath)
	defer out.Close()
	io.Copy(out, file)

	go func() {
		start := time.Now()
		log.Println("Início da importação:", start.Format("2006-01-02 15:04:05"))
	
		if err := h.importUseCase.ImportFromExcel(tmpPath); err != nil {
			log.Printf("Erro no import: %v", err)
		} else {
			end := time.Now()
			log.Println("Importação concluída com sucesso!")
			log.Println("Fim da importação:", end.Format("2006-01-02 15:04:05"))
			log.Println("Tempo total de processamento:", end.Sub(start))
		}
	}()
	

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Arquivo recebido! Processamento em background.")
}
