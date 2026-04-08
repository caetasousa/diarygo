# DiaryGo — Sistema de Contratacao de Diarista

Plataforma intermediadora entre clientes e diaristas. MVP sem pagamento online.
Documento completo de regras de negocio e esquema do banco: `README.md`
Plano de implementacao por etapas: `PLANO.md`
Detalhes adicionais em `docs/` (estrutura, regras, banco, comandos).

## Manutencao do README.md

**REGRA OBRIGATORIA:** toda mudanca no sistema deve ser refletida no `README.md`.
- Nova tabela ou coluna no banco → atualizar secao "Esquema do Banco de Dados"
- Nova regra de negocio → atualizar secao correspondente e o checklist
- Nova funcionalidade → atualizar secao "Funcionalidades por Fase"
- Nova entidade de dominio → atualizar tabela de entidades
- Mudanca de stack ou estrutura de pastas → atualizar secao "Stack e Estrutura do Projeto"

## Stack

Go | chi (router HTTP) | PostgreSQL | Flyway | pgAdmin | Git | Swagger (swaggo/swag)

## Git Flow

- Branches permanentes: `master` (producao) e `developer` (integracao)
- Cada etapa do PLANO.md = uma branch `feature/<nome>` criada a partir de `developer`
- Nunca commitar direto em `master` ou `developer`
- Merge de feature na `developer` com `--no-ff`
- Hotfix criado a partir de `master`; merged em `master` E `developer`
- `git pull --rebase` sempre ao atualizar branches (nunca merge de atualizacao)
- Convencao de commits: `feat:`, `fix:`, `chore:`, `test:`, `docs:`

## Estrategia de Repository

- Interfaces de repository definidas em `internal/domain/`
- Implementacoes in-memory em `internal/repository/memory/` (Etapas 0-11, sem banco)
- Implementacoes PostgreSQL em `internal/repository/postgres/` (Etapa 12+)
- `main.go` decide qual injetar — services e handlers nunca mudam ao trocar a implementacao
- Handlers HTTP usam chi; anotacoes swaggo em todos os handlers para geracao do Swagger

## Convencoes de Codigo

- Go idiomatico (effective Go, go vet, gofmt)
- Pacotes em minusculo, sem underscores
- Interfaces pequenas (1-3 metodos)
- Erros retornados explicitamente, nunca panic
- `context.Context` como primeiro parametro em funcoes de I/O
- Structs de dominio em `internal/domain/`, sem dependencia de frameworks

## Regras Criticas

1. **LC 150/2015:** max 2 visitas/semana mesma profissional no mesmo endereco — bloquear automaticamente
2. Agendamento minimo 24h de antecedencia
3. Cancelamento < 24h = penalizacao no score
4. Nota < 3.5 por 3 servicos consecutivos = suspensao
5. Atribuicao por avaliacao + proximidade geografica
6. Profissional tem 30 min para aceitar; apos, redireciona

## Testes

- Testes unitarios dos services: usam repositories in-memory, sem banco, sem Docker
- Testes de integracao dos repositories postgres: usam banco real (Etapa 13+)
- Nunca usar mocks para PostgreSQL — usar implementacao in-memory ou banco real
- Nomenclatura: `TestMetodo_Cenario_Resultado`

## Comandos Rapidos

```bash
go run cmd/api/main.go                    # rodar (sem banco nas Etapas 0-11)
go test ./...                             # testar
gofmt -w .                                # formatar
go vet ./...                              # verificar
swag init -g cmd/api/main.go -o docs/swagger  # gerar Swagger
docker-compose up -d                      # infra (Etapa 12+)
```

## Fase Atual: MVP

Implementacao em andamento conforme `PLANO.md` (13 etapas).
Etapas 0-11: logica completa com repositories in-memory, sem dependencia de banco.
Etapa 12: substituicao para PostgreSQL. Etapa 13: testes de integracao e Swagger.

NAO inclui no MVP: pagamento online, chat, indicacao, planos premium. Arquitetura preparada via interfaces.
