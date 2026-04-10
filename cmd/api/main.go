// @title           DiaryGo API
// @version         1.0
// @description     Plataforma intermediadora entre clientes e diaristas.
// @termsOfService  http://swagger.io/terms/

// @contact.name   DiaryGo Support
// @contact.email  suporte@diarygo.com.br

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// A05: Security Misconfiguration — cabeçalhos de segurança HTTP
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "0") // desativado — usar CSP em vez disso
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'none'")
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok","service":"diarygo"}`)
	})

	r.Route("/api/v1", func(r chi.Router) {
		// rotas registradas nas próximas etapas
	})

	log.Printf("DiaryGo API iniciada na porta %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
