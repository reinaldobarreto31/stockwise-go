# StockWise — Controle de Estoque

<div align="center">

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=flat&logo=jsonwebtokens&logoColor=white)
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
├── .env.example
├── go.mod
└── go.sum
```

### Fluxo de dados

```
Client → HTTP → Middleware (JWT/CORS) → Handler → Service → Repository → PostgreSQL
```

## Funcionalidades Implementadas

- [x] Registro e login de usuários com bcrypt + JWT
- [x] CRUD completo de produtos (GET, POST, PUT, DELETE)
- [x] Movimentações de estoque — entrada e saída com validação de saldo
- [x] Atualização atômica do estoque (transação SQL)
- [x] Middleware JWT: rotas protegidas vs. públicas
- [x] Middleware CORS
- [x] Migrations SQL prontas para uso

## Rodando localmente

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

A API estará disponível em `http://localhost:8080`.

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

## Endpoints

### Público

| Método | Rota    | Descrição     |
|--------|---------|---------------|
| GET    | `/health` | Health check |

### Auth (público)

| Método | Rota                  | Descrição              |
|--------|-----------------------|------------------------|
| POST   | `/api/auth/register`  | Registrar novo usuário |
| POST   | `/api/auth/login`     | Login — retorna JWT    |

### Produtos (requer `Authorization: Bearer <token>`)

| Método | Rota                    | Descrição           |
|--------|-------------------------|---------------------|
| GET    | `/api/products`         | Listar produtos     |
| GET    | `/api/products/{id}`    | Buscar por ID       |
| POST   | `/api/products`         | Criar produto       |
| PUT    | `/api/products/{id}`    | Atualizar produto   |
| DELETE | `/api/products/{id}`    | Remover produto     |

### Movimentações (requer `Authorization: Bearer <token>`)

| Método | Rota                              | Descrição                     |
|--------|-----------------------------------|-------------------------------|
| GET    | `/api/movements`                  | Listar todas as movimentações |
| GET    | `/api/movements?product_id={id}`  | Filtrar por produto           |
| POST   | `/api/movements`                  | Registrar movimentação        |

### Exemplo — Registro

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"João","email":"joao@exemplo.com","password":"senha123"}'
```

### Exemplo — Criar produto (autenticado)

```bash
TOKEN="<jwt-retornado-no-login>"

curl -X POST http://localhost:8080/api/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Teclado Mecânico","sku":"TEC-001","category":"Periféricos","price":349.90,"stock":20,"min_stock":5}'
```

### Exemplo — Registrar entrada de estoque

```bash
curl -X POST http://localhost:8080/api/movements \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"type":"in","quantity":50,"note":"Compra inicial"}'
```

## Autor

**Reinaldo Barreto** — [github.com/reinaldobarreto31](https://github.com/reinaldobarreto31)

---

> Próximos passos: Docker Compose, documentação Swagger, frontend React.
