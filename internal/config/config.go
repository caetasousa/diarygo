package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"
)

// Config armazena as configuracoes carregadas de variaveis de ambiente.
type Config struct {
	Port                 string
	JWTSecret            string
	JWTExpirationMinutes int
	JWTExpiry            time.Duration
	Env                  string
}

// Carregar le as variaveis de ambiente e retorna a configuracao validada.
func Carregar() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		if env == "production" {
			return nil, errors.New("JWT_SECRET deve ter no minimo 32 caracteres em producao")
		}
		slog.Warn("JWT_SECRET fraco — use no minimo 32 caracteres em producao",
			"tamanho_atual", len(jwtSecret))
		// Em development, preencher com padding para atingir minimo
		if jwtSecret == "" {
			jwtSecret = "chave-desenvolvimento-insegura-troque-em-producao"
		}
	}

	expirationMinutes := 15
	if v := os.Getenv("JWT_EXPIRATION_MINUTES"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			slog.Warn("JWT_EXPIRATION_MINUTES invalido, usando 15 minutos", "valor", v)
		} else {
			expirationMinutes = n
		}
	}

	return &Config{
		Port:                 port,
		JWTSecret:            jwtSecret,
		JWTExpirationMinutes: expirationMinutes,
		JWTExpiry:            time.Duration(expirationMinutes) * time.Minute,
		Env:                  env,
	}, nil
}
