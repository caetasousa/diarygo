# Comandos Uteis

```bash
# Rodar a aplicacao (sem banco nas Etapas 0-11)
go run cmd/api/main.go

# Rodar testes
go test ./...
go test -v ./...
go test -cover ./...

# Formatar e verificar
gofmt -w .
go vet ./...

# Gerar documentacao Swagger (requer swag instalado)
swag init -g cmd/api/main.go -o docs/swagger

# Instalar swag (uma vez)
go install github.com/swaggo/swag/cmd/swag@latest

# Subir infraestrutura — PostgreSQL + pgAdmin + Flyway (Etapa 12+)
docker-compose up -d

# Aplicar migracoes manualmente
flyway -configFiles=config/flyway.conf migrate
```
