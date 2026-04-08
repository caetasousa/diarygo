# Estrutura do Projeto

```
diarygo/
├── cmd/
│   └── api/                    # main.go — entrypoint do servidor HTTP (chi)
├── internal/
│   ├── domain/                 # Entidades, interfaces de repository, regras de negocio
│   ├── handler/                # Handlers HTTP com anotacoes swaggo
│   ├── service/                # Use cases / logica de aplicacao
│   ├── repository/
│   │   ├── memory/             # Implementacoes in-memory (Etapas 0-11)
│   │   └── postgres/           # Implementacoes PostgreSQL (Etapa 12+)
│   └── middleware/             # Auth JWT, logging, etc.
├── pkg/                        # Codigo reutilizavel e exportavel
├── migrations/                 # Scripts SQL do Flyway (V1__*, V2__*, etc.)
├── config/                     # Configuracoes e env (app.env, flyway.conf)
├── docs/                       # Documentacao auxiliar (este diretorio)
│   └── swagger/                # Gerado por `swag init` — nao editar manualmente
├── docker-compose.yml          # PostgreSQL + pgAdmin + Flyway (usado a partir da Etapa 12)
├── README.md                   # Regras de negocio e esquema do banco
├── PLANO.md                    # Plano de implementacao por etapas
├── CLAUDE.md                   # Instrucoes para o Claude
├── go.mod
└── go.sum
```
