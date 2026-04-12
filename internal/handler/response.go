package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

// ErroResponse e a estrutura padrao de resposta de erro.
type ErroResponse struct {
	Erro string `json:"erro"`
}

// RespostaJSON serializa data como JSON e escreve na resposta com o status informado.
func RespostaJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("erro ao serializar resposta JSON", "erro", err)
	}
}

// RespostaErro escreve uma resposta de erro JSON padronizada.
// Em producao, nunca expoe detalhes internos (OWASP A02/A10).
func RespostaErro(w http.ResponseWriter, status int, mensagem string) {
	if os.Getenv("ENV") != "production" || status >= 500 {
		slog.Error("resposta de erro", "status", status, "mensagem", mensagem)
	}
	RespostaJSON(w, status, ErroResponse{Erro: mensagem})
}
