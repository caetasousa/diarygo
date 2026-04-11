---
name: owasp-security
description: Implement secure coding practices following OWASP Top 10:2025 for Go + chi + PostgreSQL (pgx). Use when preventing security vulnerabilities, implementing authentication, securing APIs, or conducting security reviews. Triggers on OWASP, security, XSS, SQL injection, CSRF, authentication security, secure coding, vulnerability, supply chain, exceptional conditions.
---

# OWASP Top 10:2025 — Go + chi + PostgreSQL

Guia de seguranca para aplicacoes Go usando chi (router HTTP) e PostgreSQL (pgx).
Baseado na versao **2025** do OWASP Top 10. Todos os exemplos usam as convencoes e entidades do projeto DiaryGo.

## OWASP Top 10:2025 — Resumo

| #        | Vulnerabilidade                          | Prevencao em Go                                              |
| -------- | ---------------------------------------- | ------------------------------------------------------------ |
| A01:2025 | Broken Access Control                    | Middleware chi + verificacao de ownership nos handlers        |
| A02:2025 | Security Misconfiguration                | Headers de seguranca via middleware, hardening de config      |
| A03:2025 | Software Supply Chain Failures           | `govulncheck`, SBOM, dependencias assinadas, go.sum          |
| A04:2025 | Cryptographic Failures                   | `golang.org/x/crypto/bcrypt`, `crypto/rand`, TLS 1.2+       |
| A05:2025 | Injection                                | Queries parametrizadas com pgx (`$1, $2...`), nunca concatenar |
| A06:2025 | Insecure Design                          | Threat modeling, rate limiting, regras de negocio no service  |
| A07:2025 | Authentication Failures                  | JWT curto, bcrypt cost >= 12, anti timing attack, MFA        |
| A08:2025 | Software or Data Integrity Failures      | Validacao de input, assinatura de tokens, go.sum integrity    |
| A09:2025 | Security Logging and Alerting Failures   | `log/slog` estruturado, alertas em eventos criticos          |
| A10:2025 | Mishandling of Exceptional Conditions    | Tratamento explicito de erros, fail closed, rollback          |

---

## A01:2025 — Broken Access Control

> 100% das aplicacoes testadas apresentaram vulnerabilidades. 32.654 CVEs mapeadas.
> SSRF (CWE-918) e CSRF (CWE-352) agora fazem parte desta categoria (em 2021 SSRF era A10).

### Middleware de Autorizacao por Tipo de Usuario

```go
package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UsuarioContextKey contextKey = "usuario"

type TipoUsuario string

const (
	TipoCliente      TipoUsuario = "CLIENTE"
	TipoProfissional TipoUsuario = "PROFISSIONAL"
	TipoAdmin        TipoUsuario = "ADMIN"
)

// UsuarioAuth armazena os dados extraidos do JWT validado.
type UsuarioAuth struct {
	ID   string
	Tipo TipoUsuario
}

// RequererTipo retorna um middleware chi que bloqueia acesso
// para usuarios cujo tipo nao esta na lista permitida.
// Implementa deny-by-default: se o tipo nao esta na lista, bloqueia.
func RequererTipo(tipos ...TipoUsuario) func(http.Handler) http.Handler {
	permitidos := make(map[TipoUsuario]bool, len(tipos))
	for _, t := range tipos {
		permitidos[t] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, ok := r.Context().Value(UsuarioContextKey).(*UsuarioAuth)
			if !ok || !permitidos[u.Tipo] {
				http.Error(w, `{"error":"acesso negado"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
```

### Uso com chi — Rotas protegidas por perfil

```go
r.Route("/api/v1", func(r chi.Router) {
	// Rotas publicas — nao exigem autenticacao
	r.Post("/auth/login", handler.Login)
	r.Post("/auth/registro/cliente", handler.RegistrarCliente)

	// Rotas autenticadas (qualquer tipo)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Autenticar) // valida JWT e injeta UsuarioAuth no contexto

		// Somente clientes
		r.With(middleware.RequererTipo(middleware.TipoCliente)).Route("/clientes/me", func(r chi.Router) {
			r.Get("/", handler.MeuPerfil)
			r.Get("/enderecos", handler.ListarEnderecos)
		})

		// Somente profissionais
		r.With(middleware.RequererTipo(middleware.TipoProfissional)).Route("/profissionais/me", func(r chi.Router) {
			r.Get("/", handler.MeuPerfilProfissional)
		})

		// Somente admin
		r.With(middleware.RequererTipo(middleware.TipoAdmin)).Route("/admin", func(r chi.Router) {
			r.Get("/profissionais", handler.ListarProfissionaisAdmin)
			r.Put("/profissionais/{id}/aprovar", handler.AprovarProfissional)
		})
	})
})
```

### IDOR — Verificacao de Ownership

```go
// BAD: retorna dados de qualquer usuario — IDOR
func GetSolicitacao(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s, _ := repo.BuscarPorID(r.Context(), id)
	json.NewEncoder(w).Encode(s) // qualquer um acessa qualquer solicitacao
}

// GOOD: verifica que o recurso pertence ao usuario autenticado
func GetSolicitacao(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value(middleware.UsuarioContextKey).(*middleware.UsuarioAuth)
	id := chi.URLParam(r, "id")

	s, err := repo.BuscarPorID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"nao encontrado"}`, http.StatusNotFound)
		return
	}

	// Ownership: cliente so ve suas proprias solicitacoes
	if u.Tipo == middleware.TipoCliente && s.ClienteID != u.ID {
		http.Error(w, `{"error":"acesso negado"}`, http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(s)
}
```

### SSRF — Agora parte de A01:2025

```go
package service

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var allowedHosts = map[string]bool{
	"viacep.com.br":      true, // consulta de CEP
	"api.diarygo.com.br": true,
}

// ValidarURLExterna verifica se uma URL e segura para requisicao server-side.
// Previne SSRF (CWE-918) bloqueando IPs privados e hosts fora da allowlist.
func ValidarURLExterna(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("URL invalida")
	}

	// Permitir apenas HTTPS
	if u.Scheme != "https" {
		return fmt.Errorf("apenas HTTPS e permitido")
	}

	// Verificar allowlist
	host := strings.ToLower(u.Hostname())
	if !allowedHosts[host] {
		return fmt.Errorf("host nao permitido: %s", host)
	}

	// Bloquear IPs privados (anti DNS rebinding)
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("falha ao resolver DNS: %w", err)
	}

	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("IP privado bloqueado")
		}
	}

	return nil
}
```

### CSRF — Agora parte de A01:2025

```go
// Para APIs REST com JWT no header Authorization, CSRF e mitigado naturalmente:
// - Tokens Bearer nao sao enviados automaticamente pelo browser (diferente de cookies)
// - Se usar cookies para sessao, aplicar SameSite=Strict e verificar header Origin

// CORS restritivo para a API:
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"https://diarygo.com.br":     true,
			"https://app.diarygo.com.br": true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "86400")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

### Checklist A01

- [ ] Deny-by-default: rotas sem middleware explicito sao bloqueadas
- [ ] Todo handler verifica `UsuarioAuth.Tipo` antes de operar
- [ ] Recursos acessados por ID sempre validam ownership (`ClienteID == u.ID`)
- [ ] UUIDs como identificadores — nunca IDs sequenciais expostos
- [ ] Admin nao criado via endpoint publico
- [ ] JWT tokens invalidados server-side apos logout
- [ ] CORS restritivo: somente origens permitidas
- [ ] SSRF: allowlist de hosts + bloqueio de IPs privados
- [ ] Testes de autorizacao em testes unitarios e de integracao

---

## A02:2025 — Security Misconfiguration

> Subiu de #5 (2021) para #2 (2025). 100% das aplicacoes testadas apresentaram alguma misconfiguracao.
> CWE-16 (Configuration) e CWE-611 (XXE) sao as mais notaveis.

### Headers de Seguranca via Middleware chi

```go
// SecurityHeaders aplica headers OWASP recomendados.
// Deve ser o primeiro middleware aplicado apos Logger e Recoverer.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "0") // desativado; usar CSP
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		w.Header().Set("Permissions-Policy", "geolocation=(self)")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Cache-Control", "no-store") // dados sensiveis nunca cacheados
		w.Header().Del("Server")                    // nao expor stack
		next.ServeHTTP(w, r)
	})
}
```

### Erro sem Vazamento de Stack Trace

```go
// RespostaErro retorna erro JSON sem expor detalhes internos em producao.
// CWE-209: Error messages exposing sensitive information.
func RespostaErro(w http.ResponseWriter, status int, mensagemUsuario string, errInterno error) {
	if errInterno != nil {
		slog.Error("erro interno",
			"status", status,
			"mensagem", mensagemUsuario,
			"erro", errInterno.Error(),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Em producao, nunca expor detalhes do erro
	if os.Getenv("APP_ENV") == "production" {
		fmt.Fprintf(w, `{"error":"%s"}`, mensagemUsuario)
		return
	}

	// Em dev, incluir detalhes para debug
	fmt.Fprintf(w, `{"error":"%s","detail":"%v"}`, mensagemUsuario, errInterno)
}
```

### Desabilitar Funcionalidades Desnecessarias

```go
func setupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(SecurityHeaders)

	// Swagger SOMENTE em dev — nunca em producao
	if os.Getenv("APP_ENV") != "production" {
		r.Get("/swagger/*", httpSwagger.Handler())
	}

	// Nao expor directory listing, debug endpoints, pprof em producao
	// Nao usar credenciais default em nenhum ambiente

	return r
}
```

### Configuracao Segura de Ambiente

```bash
# config/app.env — NUNCA commitar com segredos reais
# Este arquivo deve estar no .gitignore
PORT=8080
APP_ENV=development
JWT_SECRET=mude-para-um-valor-seguro-com-no-minimo-32-caracteres

# Etapa 12+
DATABASE_URL=postgres://diarygo:senha@localhost:5432/diarygo?sslmode=require
```

### Hardening do PostgreSQL

```go
// Pool de conexoes com configuracao segura
config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
if err != nil {
	log.Fatal(err)
}

// Limitar conexoes para evitar resource exhaustion
config.MaxConns = 25
config.MinConns = 5
config.MaxConnLifetime = 1 * time.Hour
config.MaxConnIdleTime = 30 * time.Minute

// sslmode=require obrigatorio em producao
// Nunca usar sslmode=disable em producao
```

### Checklist A02

- [ ] Headers de seguranca aplicados em TODAS as respostas via middleware
- [ ] Nunca expor stack trace em producao (CWE-209)
- [ ] Swagger UI desabilitado em producao
- [ ] `config/app.env` no `.gitignore`
- [ ] PostgreSQL com `sslmode=require` em producao
- [ ] Sem credenciais default em nenhum ambiente
- [ ] Sem portas, servicos ou features desnecessarios habilitados
- [ ] Cache desabilitado para respostas com dados sensiveis
- [ ] Pool de conexoes com limites configurados

---

## A03:2025 — Software Supply Chain Failures

> **NOVA categoria em 2025.** Maior taxa de incidencia media (5.19%) entre todas as categorias.
> Incorpora partes do antigo A06:2021 (Vulnerable Components) e adiciona ataques a cadeia de suprimentos.
> Exemplos reais: SolarWinds (2019), Log4Shell (2021), Bybit $1.5B theft (2025).

### Auditoria de Dependencias Go

```bash
# Verificar vulnerabilidades conhecidas nas dependencias
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Manter dependencias atualizadas
go get -u ./...
go mod tidy

# Verificar modulos desatualizados
go list -m -u all

# Verificar integridade do go.sum (contra tampering)
go mod verify
```

### go.sum como Garantia de Integridade

```go
// go.sum contem hashes criptograficos de cada dependencia.
// O Go verifica automaticamente a integridade ao baixar modulos.
// NUNCA deletar go.sum — ele e sua protecao contra supply chain attacks.
// SEMPRE commitar go.sum no repositorio.

// Usar o Go module proxy e checksum database (padrao):
// GONOSUMCHECK=  (vazio = verificar tudo)
// GOFLAGS=-mod=readonly  (impede alteracoes nao intencionais)
```

### SBOM — Software Bill of Materials

```bash
# Gerar SBOM em formato CycloneDX
go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest
cyclonedx-gomod mod -output sbom.json -json

# Ou usar syft (mais completo)
# syft . -o cyclonedx-json > sbom.json
```

### CI Pipeline — Verificacao Automatica

```yaml
# .github/workflows/security.yml
name: Security Checks
on: [push, pull_request]
jobs:
  vuln-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: govulncheck
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...
      - name: Verify module integrity
        run: go mod verify
      - name: Check for unused dependencies
        run: go mod tidy && git diff --exit-code go.mod go.sum
```

### Checklist A03

- [ ] `govulncheck ./...` executado na CI antes de merge
- [ ] `go mod verify` na CI para verificar integridade
- [ ] `go.sum` commitado no repositorio (nunca no .gitignore)
- [ ] `go mod tidy` em todo commit que altera dependencias
- [ ] Dependabot ou Renovate configurado no repositorio GitHub
- [ ] Nunca usar versoes com CVEs conhecidas
- [ ] Preferir dependencias de fontes confiáveis e mantidas
- [ ] SBOM gerado para cada release
- [ ] CI/CD com acesso restrito e logs de auditoria

---

## A04:2025 — Cryptographic Failures

> Caiu de #2 (2021) para #4 (2025). 32 CWEs, 2.185 CVEs.
> Destaque para preparacao para criptografia pos-quantica (meta 2030).
> OWASP recomenda: "encrypt all data in transit with protocols >= TLS 1.2 only, with forward secrecy."

### Hashing de Senha com bcrypt

```go
package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// OWASP 2025 recomenda: Argon2, yescrypt, scrypt, PBKDF2-HMAC-SHA-512, ou bcrypt.
// bcrypt e a escolha mais simples e segura para Go. Cost >= 12.
const bcryptCost = 12

// HashSenha gera o hash bcrypt da senha em texto plano.
func HashSenha(senha string) (string, error) {
	if len(senha) < 8 {
		return "", errors.New("senha deve ter no minimo 8 caracteres")
	}
	if len(senha) > 72 { // limite do bcrypt
		return "", errors.New("senha deve ter no maximo 72 caracteres")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerificarSenha compara senha em texto plano com o hash armazenado.
func VerificarSenha(senha, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
}
```

### Geracao de Tokens Seguros (CSPRNG)

```go
package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GerarTokenSeguro usa crypto/rand (CSPRNG) — nunca math/rand.
// CWE-338: Use of Cryptographically Weak PRNG.
func GerarTokenSeguro(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("falha ao gerar token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
```

### JWT — Secret via Variavel de Ambiente

```go
package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET nao configurado — defina em config/app.env")
	}
	if len(secret) < 32 {
		panic("JWT_SECRET deve ter no minimo 32 caracteres")
	}
	jwtSecret = []byte(secret)
}

type TokenClaims struct {
	UsuarioID string `json:"sub"`
	Tipo      string `json:"tipo"`
	jwt.RegisteredClaims
}

// GerarJWT cria um access token com expiracao curta (15 min).
func GerarJWT(usuarioID, tipo string) (string, error) {
	claims := TokenClaims{
		UsuarioID: usuarioID,
		Tipo:      tipo,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "diarygo",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidarJWT valida o token e retorna as claims.
// Impede ataque "alg: none" aceitando SOMENTE HMAC.
func ValidarJWT(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritmo nao permitido: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalido")
	}
	return claims, nil
}
```

### HSTS — Forcar HTTPS

```go
// Ja aplicado no SecurityHeaders middleware (A02):
// w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

// Em producao, o reverse proxy (nginx/caddy) deve redirecionar HTTP -> HTTPS.
// A aplicacao Go nunca deve aceitar conexoes HTTP diretas em producao.
```

### Checklist A04

- [ ] Senhas com bcrypt cost >= 12; nunca MD5/SHA para senhas (CWE-327)
- [ ] JWT_SECRET via env, nunca hardcoded, minimo 32 chars (CWE-798)
- [ ] Tokens aleatorios com `crypto/rand`, nunca `math/rand` (CWE-338)
- [ ] JWT com expiracao curta (15 min access, 7 dias refresh)
- [ ] Algoritmo de assinatura validado no parse (bloquear `alg: none`)
- [ ] TLS 1.2+ obrigatorio em producao com forward secrecy
- [ ] HSTS habilitado (max-age >= 2 anos)
- [ ] Cache-Control: no-store para respostas com dados sensiveis
- [ ] Dados sensiveis descartados assim que nao forem mais necessarios

---

## A05:2025 — Injection

> Caiu de #3 (2021) para #5 (2025). 37 CWEs, 62.445 CVEs.
> SQL Injection (CWE-89): 14.000+ CVEs. XSS (CWE-79): 30.000+ CVEs.
> "The preferred approach uses a safe API which avoids using the interpreter entirely."

### SQL Injection — Queries Parametrizadas com pgx

```go
// BAD: concatenacao de string — SQL INJECTION (CWE-89)
func (r *repo) BuscarPorEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	query := "SELECT id, email, senha_hash FROM usuarios WHERE email = '" + email + "'"
	// ATAQUE: email = "'; DROP TABLE usuarios;--"
	row := r.pool.QueryRow(ctx, query)
	// ...
}

// GOOD: query parametrizada — pgx usa $1, $2... nativamente
func (r *repo) BuscarPorEmail(ctx context.Context, email string) (*domain.Usuario, error) {
	query := "SELECT id, email, senha_hash, tipo FROM usuarios WHERE email = $1"
	row := r.pool.QueryRow(ctx, query, email)

	var u domain.Usuario
	err := row.Scan(&u.ID, &u.Email, &u.SenhaHash, &u.Tipo)
	if err != nil {
		return nil, fmt.Errorf("usuario nao encontrado: %w", err)
	}
	return &u, nil
}
```

### Filtros Dinamicos com Seguranca

```go
// Construir WHERE dinamico sem concatenar input do usuario.
// Nomes de tabelas e colunas NAO podem ser parametrizados — usar allowlist.
func (r *repo) ListarServicos(ctx context.Context, filtros FiltroServico) ([]domain.Servico, error) {
	query := "SELECT id, status, data_agendada FROM servicos WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if filtros.Status != "" {
		// Validar contra allowlist — nomes de colunas/valores de enum nao parametrizaveis
		statusValidos := map[string]bool{
			"AGENDADO": true, "EM_ANDAMENTO": true, "CONCLUIDO": true, "CANCELADO": true,
		}
		if !statusValidos[filtros.Status] {
			return nil, fmt.Errorf("status invalido: %s", filtros.Status)
		}
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, filtros.Status)
		argIdx++
	}

	if !filtros.DataInicio.IsZero() {
		query += fmt.Sprintf(" AND data_agendada >= $%d", argIdx)
		args = append(args, filtros.DataInicio)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY data_agendada DESC LIMIT $%d", argIdx)
	args = append(args, filtros.Limite)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servicos []domain.Servico
	for rows.Next() {
		var s domain.Servico
		if err := rows.Scan(&s.ID, &s.Status, &s.DataAgendada); err != nil {
			return nil, err
		}
		servicos = append(servicos, s)
	}
	return servicos, rows.Err()
}
```

### XSS — Agora parte de A05:2025 (CWE-79: 30.000+ CVEs)

```go
// Para APIs REST que retornam JSON, XSS e mitigado por:
// 1. Content-Type: application/json (browser nao interpreta como HTML)
// 2. CSP: default-src 'none' (bloqueia execucao de scripts)
// 3. X-Content-Type-Options: nosniff (impede MIME sniffing)

// Se a API retornar HTML em algum caso (ex: emails, templates):
import "html/template"

// GOOD: html/template escapa automaticamente
tmpl := template.Must(template.New("email").Parse(`
	<p>Ola, {{.Nome}}</p>
`))

// BAD: text/template NAO escapa — nunca usar com input de usuario
import "text/template" // PERIGO para conteudo HTML
```

### Command Injection (CWE-78)

```go
// BAD: entrada do usuario no shell
func Resize(userInput string) error {
	cmd := exec.Command("sh", "-c", "convert "+userInput+" output.png")
	return cmd.Run()
}

// GOOD: argumentos separados, sem shell
func Resize(inputPath string) error {
	cmd := exec.Command("convert", inputPath, "output.png")
	return cmd.Run()
}
```

### Checklist A05

- [ ] pgx com `$1, $2...` em TODAS as queries — nunca concatenar input (CWE-89)
- [ ] Nomes de tabelas/colunas validados contra allowlist
- [ ] `html/template` para conteudo HTML, nunca `text/template` (CWE-79)
- [ ] Content-Type, CSP e nosniff aplicados em todas as respostas
- [ ] `exec.Command` com argumentos separados, nunca `sh -c` (CWE-78)
- [ ] Queries complexas: preferir views/functions no Postgres

---

## A06:2025 — Insecure Design

> Caiu de #4 (2021) para #6 (2025). 39 CWEs, 7.647 CVEs.
> "An insecure design cannot be fixed by a perfect implementation."
> Foco: threat modeling, rate limiting, regras de negocio como barreiras.

### Rate Limiting com httprate (chi-compatible)

```go
package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Rate limit global: 100 req/min por IP
	r.Use(httprate.LimitByIP(100, time.Minute))

	r.Route("/auth", func(r chi.Router) {
		// Rate limit restrito para autenticacao: 10 req/min por IP
		// Previne credential stuffing e password spray (OWASP A07)
		r.Use(httprate.LimitByIP(10, time.Minute))
		r.Post("/login", handler.Login)
		r.Post("/registro/cliente", handler.RegistrarCliente)
		r.Post("/registro/profissional", handler.RegistrarProfissional)
	})

	return r
}
```

### Validacao de Input Estruturada

```go
package handler

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

func ValidarEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("email invalido: %s", email)
	}
	return nil
}

func ValidarSenha(senha string) error {
	n := utf8.RuneCountInString(senha)
	if n < 8 {
		return fmt.Errorf("senha deve ter no minimo 8 caracteres")
	}
	if n > 72 { // limite do bcrypt
		return fmt.Errorf("senha deve ter no maximo 72 caracteres")
	}
	// OWASP 2025/NIST 800-63b: NAO forcar rotacao, NAO exigir complexidade artificial.
	// Verificar contra lista de senhas comprometidas (haveibeenpwned.com) e recomendado.
	return nil
}

var cpfRegex = regexp.MustCompile(`^\d{11}$`)

func ValidarCPF(cpf string) error {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	if !cpfRegex.MatchString(cpf) {
		return fmt.Errorf("CPF deve conter 11 digitos")
	}
	if strings.Count(cpf, string(cpf[0])) == 11 {
		return fmt.Errorf("CPF invalido")
	}
	if !validarDigitosCPF(cpf) {
		return fmt.Errorf("CPF invalido")
	}
	return nil
}

type RegistroRequest struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"`
	CPF   string `json:"cpf"`
}

func (rr RegistroRequest) Validar() error {
	if strings.TrimSpace(rr.Nome) == "" {
		return fmt.Errorf("nome e obrigatorio")
	}
	if err := ValidarEmail(rr.Email); err != nil {
		return err
	}
	if err := ValidarSenha(rr.Senha); err != nil {
		return err
	}
	if err := ValidarCPF(rr.CPF); err != nil {
		return err
	}
	return nil
}
```

### Regras de Negocio como Barreira — Threat Modeling

```go
// LC 150/2015: max 2 visitas/semana mesma profissional no mesmo endereco.
// Esta regra e LEGAL — nao pode ser contornada por nenhum ator, incluindo admin.
// Threat model: um atacante poderia tentar criar multiplas solicitacoes
// para o mesmo endereco com a mesma profissional. O service BLOQUEIA.
func (s *MatchingService) VerificarLimiteLegal(
	ctx context.Context,
	profissionalID, enderecoID string,
	semana time.Time,
) error {
	count, err := s.servicoRepo.ContarNaSemana(ctx, profissionalID, enderecoID, semana)
	if err != nil {
		return fmt.Errorf("erro ao verificar limite legal: %w", err)
	}
	if count >= 2 {
		return ErrLimiteLegalExcedido
	}
	return nil
}

// Agendamento minimo 24h de antecedencia — impede abuso de solicitacoes de ultima hora.
// Cancelamento < 24h = penalizacao no score — desincentiva cancelamentos abusivos.
```

### Limite de Payload no Body

```go
// Limitar tamanho do body para evitar ataques de payload grande (CWE-770).
func LimitarBody(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// Uso: r.Use(LimitarBody(1 << 20)) // 1 MB max
```

### Checklist A06

- [ ] Rate limiting global (100/min) e restrito em /auth/* (10/min)
- [ ] Toda entrada validada no handler antes de chegar ao service
- [ ] `http.MaxBytesReader` em handlers que recebem body
- [ ] Regras legais (LC 150/2015) no service, nunca bypassaveis
- [ ] Threat modeling documentado para fluxos criticos (login, matching, pagamento futuro)
- [ ] Plausibility checks: data no passado? UUID invalido? Status impossivel?

---

## A07:2025 — Authentication Failures

> Manteve #7. 36 CWEs, 7.147 CVEs.
> OWASP 2025 enfatiza: MFA obrigatorio, NIST 800-63b, protecao contra credential stuffing.
> "Hybrid password attacks" (Password1!, Password2!) sao uma ameaca crescente.

### Middleware de Autenticacao JWT

```go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/caetasousa/diarygo/internal/service"
)

// Autenticar valida o JWT no header Authorization e injeta
// UsuarioAuth no contexto da requisicao.
func Autenticar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"token nao fornecido"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, `{"error":"formato de token invalido"}`, http.StatusUnauthorized)
			return
		}

		claims, err := service.ValidarJWT(parts[1])
		if err != nil {
			http.Error(w, `{"error":"token invalido ou expirado"}`, http.StatusUnauthorized)
			return
		}

		u := &UsuarioAuth{
			ID:   claims.UsuarioID,
			Tipo: TipoUsuario(claims.Tipo),
		}

		ctx := context.WithValue(r.Context(), UsuarioContextKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

### Login com Protecao contra Timing Attack

```go
func (s *AuthService) Login(ctx context.Context, email, senha string) (string, error) {
	u, err := s.repo.BuscarPorEmail(ctx, email)
	if err != nil {
		// Executar bcrypt mesmo com usuario inexistente
		// para evitar timing attack que revela emails validos.
		// O tempo de resposta e identico nos dois caminhos.
		bcrypt.CompareHashAndPassword(
			[]byte("$2a$12$000000000000000000000000000000000000000000000000000000"),
			[]byte(senha),
		)
		return "", ErrCredenciaisInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.SenhaHash), []byte(senha)); err != nil {
		return "", ErrCredenciaisInvalidas
	}

	// Mensagem de erro GENERICA — nunca "email nao encontrado" vs "senha incorreta"
	// CWE-204: Observable Response Discrepancy
	return service.GerarJWT(u.ID.String(), string(u.Tipo))
}
```

### Verificacao contra Senhas Comprometidas (NIST 800-63b)

```go
// OWASP 2025 recomenda: validar senhas contra listas de credenciais comprometidas.
// Implementacao simplificada usando k-Anonymity da API HaveIBeenPwned:
import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SenhaComprometida(senha string) (bool, error) {
	hash := fmt.Sprintf("%X", sha1.Sum([]byte(senha)))
	prefix := hash[:5]
	suffix := hash[5:]

	resp, err := http.Get("https://api.pwnedpasswords.com/range/" + prefix)
	if err != nil {
		// Falhar aberto: se a API estiver fora, permitir (nao bloquear login)
		return false, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return strings.Contains(string(body), suffix), nil
}
```

### Checklist A07

- [ ] JWT com expiracao de 15 min (access), 7 dias (refresh)
- [ ] Mensagem de erro generica no login (CWE-204)
- [ ] bcrypt nos dois caminhos para evitar timing attack
- [ ] Rate limiting em /auth/*: 10 req/min por IP
- [ ] Bearer token no header Authorization, nunca na URL
- [ ] Sessoes invalidadas server-side apos logout
- [ ] Verificar senhas contra lista de comprometidas (NIST 800-63b)
- [ ] NAO forcar rotacao periodica de senha (NIST 800-63b)
- [ ] MFA recomendado para admin (implementar em fase futura)

---

## A08:2025 — Software or Data Integrity Failures

> Manteve #8. 14 CWEs, 3.331 CVEs.
> CWE-502 (Deserialization of Untrusted Data) e CWE-829 (Untrusted Control Sphere).
> Em Go: desserializacao insegura e menos comum, mas validacao de input continua critica.

### Validacao de Dados de Entrada

```go
// Validar TUDO que vem do cliente antes de processar.
// Nunca confiar em IDs, tipos ou status vindos do frontend.
func CriarSolicitacao(w http.ResponseWriter, r *http.Request) {
	var req SolicitacaoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespostaErro(w, http.StatusBadRequest, "JSON invalido", err)
		return
	}

	// Validar campos obrigatorios
	if req.CategoriaID == "" {
		RespostaErro(w, http.StatusBadRequest, "categoria e obrigatoria", nil)
		return
	}

	// Validar que UUID e valido (evitar input malformado)
	if _, err := uuid.Parse(req.CategoriaID); err != nil {
		RespostaErro(w, http.StatusBadRequest, "ID de categoria invalido", nil)
		return
	}

	// Validar regra de negocio: agendamento minimo 24h
	if time.Until(req.DataAgendada) < 24*time.Hour {
		RespostaErro(w, http.StatusBadRequest, "agendamento minimo 24h de antecedencia", nil)
		return
	}

	// ... prosseguir
}
```

### go.sum e go mod verify — Integridade de Dependencias

```bash
# go.sum garante que os modulos baixados sao identicos aos que voce testou.
# go mod verify confere os hashes de todos os modulos no cache local.

# Na CI:
go mod verify
# Saida esperada: "all modules verified"
# Se falhar: um modulo foi adulterado no cache ou no proxy.
```

### Protecao contra Desserializacao Insegura

```go
// Em Go, json.Unmarshal e relativamente seguro (nao executa codigo).
// Mas validar os dados APOS desserializar e obrigatorio.

// BAD: confiar no JSON sem validar
var req struct {
	Tipo   string `json:"tipo"`
	Status string `json:"status"`
}
json.NewDecoder(r.Body).Decode(&req)
// Atacante envia: {"tipo": "ADMIN", "status": "APROVADO"}

// GOOD: ignorar campos que o usuario nao pode definir
// O tipo vem do JWT, o status e definido pela logica de negocio
u := r.Context().Value(middleware.UsuarioContextKey).(*middleware.UsuarioAuth)
// u.Tipo e u.ID vem do token validado, nao do body
```

### Checklist A08

- [ ] Validar todo input no handler antes de passar ao service
- [ ] UUIDs parseados com `uuid.Parse`
- [ ] JWT: validar algoritmo explicitamente (bloquear `alg: none`)
- [ ] Nunca confiar em campos de tipo/status/role vindos do frontend
- [ ] `go.sum` commitado e `go mod verify` na CI
- [ ] Usar somente pacotes de fontes confiáveis com verificacao de assinatura

---

## A09:2025 — Security Logging and Alerting Failures

> Manteve #9. Nome atualizado para incluir "Alerting" — alerta e tao importante quanto log.
> 5 CWEs. Destaque: CWE-532 (Sensitive Info in Log), CWE-117 (Log Injection).
> Caso real: provedor de saude infantil teve breach nao detectado por 7 ANOS.

### Logging Estruturado com slog

```go
package main

import (
	"log/slog"
	"os"
)

func setupLogger() {
	var handler slog.Handler

	if os.Getenv("APP_ENV") == "production" {
		// JSON em producao para processamento por ferramentas (ELK, Loki, Grafana)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Texto legivel em desenvolvimento
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	slog.SetDefault(slog.New(handler))
}
```

### Eventos de Seguranca — O que Logar

```go
package service

import (
	"log/slog"
	"net/http"
)

// LogarEventoSeguranca registra eventos criticos com contexto suficiente
// para analise forense, mas sem dados sensiveis (CWE-532).
func LogarEventoSeguranca(evento string, r *http.Request, attrs ...slog.Attr) {
	baseAttrs := []slog.Attr{
		slog.String("evento", evento),
		slog.String("ip", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("metodo", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	}
	baseAttrs = append(baseAttrs, attrs...)

	args := make([]any, len(baseAttrs))
	for i, a := range baseAttrs {
		args[i] = a
	}

	slog.Warn("SEGURANCA", args...)
}

// === Eventos que DEVEM ser logados ===

// Login com falha
LogarEventoSeguranca("login_falhou", r,
	slog.String("email", email),
)

// Login com sucesso
LogarEventoSeguranca("login_sucesso", r,
	slog.String("usuario_id", u.ID),
)

// Acesso negado (tentativa de IDOR)
LogarEventoSeguranca("acesso_negado", r,
	slog.String("usuario_id", u.ID),
	slog.String("recurso_id", recursoID),
)

// Profissional suspensa por nota baixa
LogarEventoSeguranca("profissional_suspensa", r,
	slog.String("profissional_id", profID),
	slog.Float64("nota_media", notaMedia),
)

// Tentativa de violar limite LC 150/2015
LogarEventoSeguranca("limite_legal_bloqueado", r,
	slog.String("profissional_id", profID),
	slog.String("endereco_id", endID),
	slog.Int("servicos_na_semana", count),
)

// Status de servico alterado (auditoria)
LogarEventoSeguranca("status_alterado", r,
	slog.String("servico_id", servicoID),
	slog.String("de", statusAnterior),
	slog.String("para", statusNovo),
)

// Erro interno (500)
LogarEventoSeguranca("erro_interno", r,
	slog.String("handler", "CriarSolicitacao"),
	slog.String("erro", err.Error()),
)
```

### Protecao contra Log Injection (CWE-117)

```go
// BAD: input do usuario direto no log — atacante injeta newlines
slog.Info("Login: " + userInput) // userInput = "admin\n[WARN] Fake log entry"

// GOOD: slog com campos estruturados — input e encodado automaticamente
slog.Info("login tentativa", "email", userInput)
// slog JSON output: {"msg":"login tentativa","email":"admin\\n[WARN] Fake log entry"}
// O campo e tratado como dado, nao como formatacao.
```

### O que NUNCA Logar (CWE-532)

```go
// NUNCA LOGAR:
// - Senhas (nem texto plano, nem hash)
// - Tokens JWT completos
// - CPF completo — mascarar
// - Dados bancarios
// - Numeros de cartao

func MascararCPF(cpf string) string {
	if len(cpf) < 11 {
		return "***"
	}
	return "***.***.***-" + cpf[9:]
}

// Logar com CPF mascarado:
slog.Info("cliente registrado",
	"cpf", MascararCPF(req.CPF),
	"email", req.Email,
)
```

### Alerting — Detectar Ataques em Tempo Real

```go
// OWASP 2025 enfatiza: logging sem alerting nao previne breaches.
// Implementar alertas para padroes suspeitos:

// 1. Multiplas falhas de login do mesmo IP (credential stuffing)
//    -> httprate ja bloqueia, mas logar o evento para dashboard de monitoramento

// 2. Tentativas de IDOR repetidas (usuario tentando acessar recursos de outros)
//    -> contar acessos negados por usuario; se > threshold, alertar

// 3. Erros 500 em rajada (possivel ataque ou falha sistemica)
//    -> monitorar taxa de erros 500 por minuto

// Integracao: logs JSON -> Loki/ELK -> Grafana alerting -> Slack/email
// Honeytokens: criar endpoints falsos (/admin/debug, /backup) que geram alerta imediato
```

### Checklist A09

- [ ] `log/slog` como biblioteca de logging (stdlib Go)
- [ ] JSON em producao, texto em dev
- [ ] Logar: login (sucesso E falha), acesso negado, mudanca de status, erros 500
- [ ] Logar: violacoes de regras legais (LC 150/2015)
- [ ] Nunca logar: senhas, tokens, CPF completo, dados bancarios (CWE-532)
- [ ] slog estruturado para prevenir log injection (CWE-117)
- [ ] request_id em todos os logs para correlacao
- [ ] Alertas configurados para padroes de ataque (credential stuffing, IDOR em massa)
- [ ] Logs append-only com integridade (nao permitir alteracao)
- [ ] Plano de resposta a incidentes documentado (NIST 800-61r2)

---

## A10:2025 — Mishandling of Exceptional Conditions

> **NOVA categoria em 2025** (substituiu SSRF, que agora faz parte de A01).
> 24 CWEs, 3.416 CVEs. Max incidence rate: 20.67%.
> CWE-476 (NULL pointer dereference), CWE-209 (Error info disclosure), CWE-636 (Failing open).
> "Catch every possible system error directly at the place where they occur."

### Tratamento Explicito de Erros em Go

```go
// Go ja favorece tratamento explicito de erros (sem exceptions).
// A regra: SEMPRE verificar err. Nunca ignorar com _.

// BAD: ignorar erro — pode levar a estado inconsistente
user, _ := repo.BuscarPorID(ctx, id)
// user pode ser nil -> panic (CWE-476: NULL pointer dereference)

// GOOD: tratar cada erro no ponto onde ocorre
user, err := repo.BuscarPorID(ctx, id)
if err != nil {
	RespostaErro(w, http.StatusNotFound, "usuario nao encontrado", err)
	return
}
```

### Fail Closed — Nunca Fail Open (CWE-636)

```go
// REGRA CRITICA: em caso de erro, NEGAR acesso. Nunca permitir.

// BAD: fail open — se a verificacao falhar, permite acesso
func Autenticar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := service.ValidarJWT(extractToken(r))
		if err != nil {
			// "Vou deixar passar porque o servico de auth pode estar fora"
			next.ServeHTTP(w, r) // PERIGO: acesso sem autenticacao
			return
		}
		// ...
	})
}

// GOOD: fail closed — erro = acesso negado
func Autenticar(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := service.ValidarJWT(extractToken(r))
		if err != nil {
			http.Error(w, `{"error":"nao autorizado"}`, http.StatusUnauthorized)
			return // BLOQUEIA — fail closed
		}
		// ...
	})
}
```

### Rollback de Transacoes em Caso de Erro

```go
// OWASP 2025: "Roll back every part of the transaction" upon failure.
// Previne estados inconsistentes e fraudes financeiras (race conditions).

func (s *ServicoService) ConcluirServico(ctx context.Context, servicoID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transacao: %w", err)
	}
	// DEFER ROLLBACK: se qualquer etapa falhar, TUDO e revertido
	defer tx.Rollback(ctx)

	// 1. Atualizar status do servico
	_, err = tx.Exec(ctx,
		"UPDATE servicos SET status = 'CONCLUIDO', data_conclusao = NOW() WHERE id = $1",
		servicoID,
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar servico: %w", err)
	}

	// 2. Registrar transacao financeira
	_, err = tx.Exec(ctx,
		"INSERT INTO transacoes (servico_id, status, metodo) VALUES ($1, 'REGISTRADA', 'DIRETO_EXTERNO')",
		servicoID,
	)
	if err != nil {
		return fmt.Errorf("erro ao registrar transacao: %w", err)
	}

	// 3. Somente se TUDO deu certo, commitar
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao commitar: %w", err)
	}

	return nil
}
```

### Recoverer — Safety Net Global

```go
// chi middleware.Recoverer captura panics e retorna 500 em vez de derrubar o servidor.
// Mas NUNCA depender dele como tratamento primario — tratar erros explicitamente.

r.Use(middleware.Recoverer) // safety net — ultimo recurso

// Se um panic acontecer, logar com detalhes para investigar:
r.Use(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("PANIC recuperado",
					"panic", fmt.Sprintf("%v", rec),
					"path", r.URL.Path,
					"method", r.Method,
					"ip", r.RemoteAddr,
				)
				http.Error(w, `{"error":"erro interno"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
})
```

### Resource Exhaustion — Liberacao de Recursos

```go
// CWE-770: Allocation of Resources Without Limits.
// Sempre fechar recursos (rows, files, connections) em caso de erro.

// GOOD: defer garante fechamento mesmo em caso de erro
rows, err := pool.Query(ctx, query, args...)
if err != nil {
	return nil, err
}
defer rows.Close() // SEMPRE defer Close apos verificar err

// GOOD: timeout no contexto para evitar queries infinitas
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

row := pool.QueryRow(ctx, query, args...)
```

### Checklist A10

- [ ] SEMPRE verificar `err` — nunca ignorar com `_` em operacoes de I/O
- [ ] Fail closed: erro = negar acesso, nunca permitir (CWE-636)
- [ ] Transacoes com `defer tx.Rollback(ctx)` — commit somente no sucesso
- [ ] `defer rows.Close()` / `defer resp.Body.Close()` em todo recurso
- [ ] `context.WithTimeout` em queries e chamadas externas
- [ ] `middleware.Recoverer` como safety net, nao como tratamento primario
- [ ] Erros logados com contexto suficiente, mas sem dados sensiveis
- [ ] Respostas de erro genéricas em producao (nunca stack trace)
- [ ] Panic logado com detalhes (path, method, IP) para investigacao

---

## Checklist Geral Pre-Deploy — OWASP Top 10:2025

```markdown
### A01 — Access Control
- [ ] Deny-by-default em todas as rotas
- [ ] Middleware RequererTipo em todas as rotas protegidas
- [ ] Verificacao de ownership em acesso a recurso por ID
- [ ] UUIDs como identificadores (nunca sequenciais)
- [ ] JWT invalidado apos logout
- [ ] CORS restritivo (allowlist de origens)
- [ ] SSRF: allowlist de hosts + bloqueio de IPs privados

### A02 — Security Misconfiguration
- [ ] Headers de seguranca em todas as respostas
- [ ] Swagger desabilitado em producao
- [ ] Sem credenciais default
- [ ] PostgreSQL com sslmode=require
- [ ] config/app.env no .gitignore

### A03 — Supply Chain
- [ ] govulncheck na CI
- [ ] go mod verify na CI
- [ ] go.sum commitado
- [ ] Dependabot/Renovate configurado
- [ ] SBOM gerado para releases

### A04 — Cryptographic Failures
- [ ] bcrypt cost >= 12
- [ ] JWT_SECRET via env, >= 32 chars, nunca hardcoded
- [ ] crypto/rand para tokens (nunca math/rand)
- [ ] TLS 1.2+ com forward secrecy
- [ ] HSTS habilitado

### A05 — Injection
- [ ] Queries parametrizadas com pgx ($1, $2...)
- [ ] html/template para conteudo HTML
- [ ] exec.Command com argumentos separados
- [ ] CSP + nosniff em todas as respostas

### A06 — Insecure Design
- [ ] Rate limiting global + restrito em /auth/*
- [ ] http.MaxBytesReader em handlers com body
- [ ] Regras legais (LC 150/2015) no service
- [ ] Validacao completa de input no handler

### A07 — Authentication
- [ ] JWT 15 min access / 7 dias refresh
- [ ] Mensagem generica no login (anti enumeration)
- [ ] Protecao contra timing attack
- [ ] Senhas verificadas contra lista de comprometidas

### A08 — Data Integrity
- [ ] Input validado apos desserializacao
- [ ] Campos de tipo/status/role nunca do frontend
- [ ] go.sum + go mod verify

### A09 — Logging & Alerting
- [ ] slog estruturado (JSON em producao)
- [ ] Login, acesso negado, mudanca de status logados
- [ ] Senhas/tokens/CPF nunca logados
- [ ] Alertas para padroes de ataque
- [ ] request_id em todos os logs

### A10 — Exceptional Conditions
- [ ] Todo err verificado, nunca ignorado
- [ ] Fail closed (erro = negar acesso)
- [ ] Transacoes com defer Rollback
- [ ] Recursos fechados com defer
- [ ] context.WithTimeout em I/O
- [ ] Recoverer como safety net

### Regras DiaryGo
- [ ] LC 150/2015: max 2 visitas/semana
- [ ] Agendamento minimo 24h
- [ ] Cancelamento < 24h penaliza score
- [ ] Nota < 3.5 por 3 consecutivos = suspensao
- [ ] Profissional tem 30 min para aceitar
```

## Recursos

- **OWASP Top 10:2025**: https://owasp.org/Top10/2025/
- **OWASP Cheat Sheets**: https://cheatsheetseries.owasp.org/
- **OWASP Go Secure Coding**: https://owasp.org/www-project-go-secure-coding-practices-guide/
- **NIST 800-63b (Authentication)**: https://pages.nist.gov/800-63-4/sp800-63b.html
- **NIST 800-61r2 (Incident Response)**: https://csrc.nist.gov/pubs/sp/800/61/r2/final
- **govulncheck**: https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
- **Go crypto/rand**: https://pkg.go.dev/crypto/rand
- **pgx (queries parametrizadas)**: https://pkg.go.dev/github.com/jackc/pgx/v5
- **HaveIBeenPwned API**: https://haveibeenpwned.com/API/v3
