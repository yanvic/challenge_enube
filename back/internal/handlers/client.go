package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "📊 Bem-vindo à tela principal!")
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "📥 Aqui é o importador de dados!")
}
