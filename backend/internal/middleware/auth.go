package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/caetasousa/diarygo/internal/domain"
	"github.com/caetasousa/diarygo/internal/service"
)

type contextKey string

// UsuarioContextKey e a chave usada para armazenar o TokenPayload no contexto da request.
const UsuarioContextKey contextKey = "usuario_payload"

// Autenticar e um middleware chi que valida o JWT do header Authorization
// e injeta o TokenPayload no contexto da request.
// OWASP A10: fail closed — qualquer erro resulta em 401, nunca em acesso permitido.
func Autenticar(svc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"erro":"token nao fornecido"}`, http.StatusUnauthorized)
				return
			}

			// OWASP A07: formato esperado "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"erro":"formato de token invalido"}`, http.StatusUnauthorized)
				return
			}

			claims, err := svc.ValidarToken(parts[1])
			if err != nil {
				// OWASP A10: fail closed — nunca permitir acesso em caso de erro
				http.Error(w, `{"erro":"token invalido ou expirado"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UsuarioContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequererTipo retorna um middleware chi que verifica se o tipo do usuario autenticado
// esta na lista de tipos permitidos (deny-by-default — OWASP A01).
func RequererTipo(tipos ...domain.TipoUsuario) func(http.Handler) http.Handler {
	permitidos := make(map[domain.TipoUsuario]bool, len(tipos))
	for _, t := range tipos {
		permitidos[t] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload, ok := UsuarioDoContexto(r.Context())
			if !ok {
				// Fail closed: sem claims no contexto = nao autenticado
				http.Error(w, `{"erro":"nao autenticado"}`, http.StatusUnauthorized)
				return
			}

			if !permitidos[payload.Tipo] {
				http.Error(w, `{"erro":"acesso negado"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// UsuarioDoContexto extrai o TokenPayload do contexto da request.
// Retorna false se o payload nao estiver presente.
func UsuarioDoContexto(ctx context.Context) (*domain.TokenPayload, bool) {
	payload, ok := ctx.Value(UsuarioContextKey).(*domain.TokenPayload)
	return payload, ok && payload != nil
}
