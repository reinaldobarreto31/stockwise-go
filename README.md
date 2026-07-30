# StockWise — Controle de Estoque

<div align="center">

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=flat&logo=jsonwebtokens&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker&logoColor=white)
![Status](https://img.shields.io/badge/status-funcional-brightgreen?style=flat)

**API RESTful de controle de estoque — Go + PostgreSQL + JWT**

</div>

---

## Sobre o Projeto

StockWise é um sistema de controle de estoque com backend em Go puro. Expõe uma API RESTful com autenticação JWT, CRUD completo de produtos e registro transacional de movimentações de estoque (entrada/saída).

## Stack

| Camada     | Tecnologia                        |
|------------|-----------------------------------|
| Backend    | Go 1.22 + net/http + chi          |
| Banco      | PostgreSQL 16 (lib/pq)            |
| Auth       | JWT — golang-jwt/jwt v5           |
| Senhas     | bcrypt — golang.org/x/crypto      |
| Config     | godotenv                          |
| Container  | Docker + Docker Compose           |

## Arquitetura

```
stockwise-go/
├── cmd/
│   └── api/
│       └── main.go              # Entrypoint — wiring de dependências + HTTP server
├── internal/
│   ├── db/
│   │   └── db.go                # Conexão PostgreSQL
│   ├── handler/                 # Camada HTTP — decode, validate, encode
│   │   ├── auth.go
│   │   ├── product.go
│   │   ├── movement.go
│   │   └── helpers.go           # writeJSON / writeError
│   ├── middleware/
│   │   ├── auth.go              # Validação JWT + injeção de contexto
│   │   └── cors.go              # CORS headers
│   ├── model/                   # Structs de domínio
│   │   ├── user.go
│   │   ├── product.go
│   │   └── movement.go
│   ├── repository/              # Camada de dados — SQL puro (lib/pq)
│   │   ├── user.go
│   │   ├── product.go
│   │   └── movement.go
│   └── service/                 # Regras de negócio
│       ├── auth.go
│       ├── product.go
│       └── movement.go
├── db/
│   └── migrations/              # SQL migrations (aplicar em ordem)
│       ├── 001_create_users.sql
│       ├── 002_create_products.sql
│       └── 003_create_movements.sql
├── Dockerfile                   # Multi-stage build (builder + alpine runner)
├── docker-compose.yml           # Go API + PostgreSQL — sobe tudo com um comando
├── entrypoint.sh                # Aguarda o banco, aplica migrations, inicia a API
├── .env.example
├── go.mod
└── go.sum
```

### Fluxo de dados

```
Client → HTTP → Middleware (JWT/CORS) → Handler → Service → Repository → PostgreSQL
```

## 🐳 Getting Started — Docker (recomendado)

> Requer apenas **Docker** e **Docker Compose**. Sem Go, sem PostgreSQL, sem configuração manual.

```bash
# 1. Clonar o repositório
git clone https://github.com/reinaldobarreto31/stockwise-go.git
cd stockwise-go

# 2. Configurar variáveis de ambiente
cp .env.docker.example .env.docker
# Edite .env.docker se precisar alterar a senha padrão

# 3. Subir tudo com um comando
docker compose up --build
```

O Docker Compose irá:
1. Subir o PostgreSQL 16 e aguardar ele ficar saudável
2. Aplicar automaticamente as 3 migrations SQL
3. Iniciar a API Go na porta **8080**

A API estará disponível em `http://localhost:8080`.

**Parar e remover os containers:**
```bash
docker compose down
```

**Parar e remover os containers + dados do banco:**
```bash
docker compose down -v
```

## Rodando localmente (sem Docker)

### Pré-requisitos

- Go 1.22+
- PostgreSQL 16+

### Passos

```bash
# Clonar o repositório
git clone https://github.com/reinaldobarreto31/stockwise-go.git
cd stockwise-go

# Configurar variáveis de ambiente
cp .env.example .env
# Edite .env com suas credenciais PostgreSQL

# Aplicar migrations (na ordem)
psql $DATABASE_URL -f db/migrations/001_create_users.sql
psql $DATABASE_URL -f db/migrations/002_create_products.sql
psql $DATABASE_URL -f db/migrations/003_create_movements.sql

# Baixar dependências e rodar
go mod download
go run cmd/api/main.go
```

## Variáveis de Ambiente

```env
# Banco de dados
DATABASE_URL=postgres://user:password@localhost:5432/stockwise?sslmode=disable

# JWT
JWT_SECRET=sua-chave-secreta-de-pelo-menos-32-caracteres
JWT_EXPIRATION_HOURS=24

# Servidor
PORT=8080
ENV=development
```

> No Docker Compose essas variáveis já estão configuradas automaticamente.

## Endpoints

### Público

| Método | Rota      | Descrição     |
|--------|-----------|---------------|
| GET    | `/health` | Health check  |
| GET    | `/api`    | Info da API   |

### Auth (público)

| Método | Rota               | Descrição           |
|--------|--------------------|---------------------|
| POST   | `/api/auth/register` | Cadastro de usuário |
| POST   | `/api/auth/login`    | Login — retorna JWT |

### Produtos (requer JWT)

| Método | Rota                  | Descrição              |
|--------|-----------------------|------------------------|
| GET    | `/api/products`       | Listar produtos        |
| POST   | `/api/products`       | Criar produto          |
| GET    | `/api/products/{id}`  | Buscar produto por ID  |
| PUT    | `/api/products/{id}`  | Atualizar produto      |
| DELETE | `/api/products/{id}`  | Remover produto        |

### Movimentações (requer JWT)

| Método | Rota                    | Descrição                    |
|--------|-------------------------|------------------------------|
| GET    | `/api/movements`        | Listar movimentações         |
| POST   | `/api/movements`        | Registrar entrada ou saída   |
| GET    | `/api/movements/{id}`   | Buscar movimentação por ID   |

### Exemplo rápido

```bash
# Health check
curl http://localhost:8080/health

# Registrar usuário
curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"secret123"}'

# Login — guarde o token retornado
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123"}' | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Criar produto (requer JWT)
curl -s -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Notebook","sku":"NB-001","category":"Electronics","price":2999.99,"stock":10}'
```

## Funcionalidades Implementadas

- [x] Registro e login de usuários com bcrypt + JWT
- [x] CRUD completo de produtos (GET, POST, PUT, DELETE)
- [x] Movimentações de estoque — entrada e saída com validação de saldo
- [x] Atualização atômica do estoque (transação SQL)
- [x] Middleware JWT: rotas protegidas vs. públicas
- [x] Middleware CORS
- [x] Migrations SQL prontas para uso
- [x] Docker + Docker Compose (sobe tudo com um comando)
- [x] Multi-stage Dockerfile (imagem final ~20 MB)

## Licença

MIT
