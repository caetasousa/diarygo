# Banco de Dados

## Estrategia: In-Memory First

O banco de dados **nao e necessario para desenvolver ou testar** nas Etapas 0–11.
Toda a logica usa repositories in-memory (`internal/repository/memory/`).
O PostgreSQL e introduzido apenas na **Etapa 12**, substituindo as implementacoes in-memory.

## PostgreSQL (Etapa 12+)

- PostgreSQL 14+ como unico banco de dados
- Todas as migracoes via Flyway em `migrations/`
- Nomenclatura de migracoes: `V{numero}__{descricao}.sql` (ex: `V1__schema_inicial.sql`)
- Nunca alterar migracoes ja aplicadas — criar nova migracao para correcoes
- Usar pgAdmin para inspecao visual do banco em desenvolvimento
- Schema completo documentado na secao 17 do `README.md`

## Subir o banco (Etapa 12+)

```bash
docker-compose up -d
# PostgreSQL na porta 5432, pgAdmin em http://localhost:5050
# Flyway aplica automaticamente as migracoes de migrations/
```
