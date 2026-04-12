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
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/caetasousa/diarygo/internal/config"
	"github.com/caetasousa/diarygo/internal/handler"
	mw "github.com/caetasousa/diarygo/internal/middleware"
	"github.com/caetasousa/diarygo/internal/repository/memory"
	"github.com/caetasousa/diarygo/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func main() {
	// Configurar logger estruturado — OWASP A09
	cfg, err := config.Carregar()
	if err != nil {
		slog.Error("falha ao carregar configuracao", "erro", err)
		os.Exit(1)
	}

	var logHandler slog.Handler
	if cfg.Env == "production" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	slog.SetDefault(slog.New(logHandler))

	// Injecao de dependencias
	usuarioRepo := memory.NewUsuarioRepository()
	authService := service.NewAuthService(usuarioRepo, cfg.JWTSecret, cfg.JWTExpiry, cfg.Env)
	authHandler := handler.NewAuthHandler(authService)

	r := chi.NewRouter()

	// Middlewares globais
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(securityHeaders)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "diarygo"}) //nolint:errcheck
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Rotas de autenticacao — com rate limiting restrito (OWASP A06)
		r.Route("/auth", func(r chi.Router) {
			r.Use(httprate.LimitByIP(10, time.Minute))
			r.Mount("/", authHandler.Routes())
		})

		// Rotas protegidas — requerem JWT valido
		r.Group(func(r chi.Router) {
			r.Use(mw.Autenticar(authService))

			// GET /api/v1/me — retorna o payload do token (util para testes)
			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				payload, ok := mw.UsuarioDoContexto(r.Context())
				if !ok {
					handler.RespostaErro(w, http.StatusUnauthorized, "nao autenticado")
					return
				}
				handler.RespostaJSON(w, http.StatusOK, payload)
			})
		})
	})

	slog.Info("DiaryGo API iniciada", "porta", cfg.Port, "env", cfg.Env)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("servidor encerrado com erro", "erro", err)
		os.Exit(1)
	}
}

// securityHeaders aplica os headers de seguranca OWASP A02 em todas as respostas.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "0") // usar CSP em vez disso
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		w.Header().Set("Permissions-Policy", "geolocation=(self)")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Del("Server")
		next.ServeHTTP(w, r)
	})
}
