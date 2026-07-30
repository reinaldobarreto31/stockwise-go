# StockWise — Controle de Estoque

<div align="center">

![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react&logoColor=black)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-Auth-000000?style=flat&logo=jsonwebtokens&logoColor=white)
![Status](https://img.shields.io/badge/status-em%20desenvolvimento-yellow?style=flat)

**Sistema de controle de estoque com backend em Go e frontend React.js**

</div>

---

## Sobre o Projeto

StockWise é um sistema de controle de estoque desenvolvido com foco em performance e simplicidade. O backend em Go expõe uma API RESTful com autenticação JWT, enquanto o frontend em React.js oferece uma interface moderna e responsiva para gestão de produtos e movimentações.

## Stack

| Camada     | Tecnologia                        |
|------------|-----------------------------------|
| Backend    | Go 1.22 + net/http + chi          |
| Frontend   | React 18 + Vite + Tailwind CSS    |
| Banco      | PostgreSQL 16                     |
| Auth       | JWT (golang-jwt/jwt)              |
| Migrations | golang-migrate                    |
| Container  | Docker + Docker Compose           |

## Arquitetura

```
stockwise-go/
├── cmd/
│   └── api/
│       └── main.go          # Entrypoint — HTTP server
├── internal/
│   ├── handler/             # HTTP handlers (controllers)
│   │   ├── auth.go
│   │   ├── product.go
│   │   └── movement.go
│   ├── model/               # Domain models / structs
│   │   ├── user.go
│   │   ├── product.go
│   │   └── movement.go
│   ├── repository/          # Data access layer (PostgreSQL)
│   │   ├── user.go
│   │   ├── product.go
│   │   └── movement.go
│   └── service/             # Business logic
│       ├── auth.go
│       ├── product.go
│       └── movement.go
├── db/
│   └── migrations/          # SQL migration files
├── go.mod
├── go.sum
└── README.md
```

### Fluxo de dados

```
Client (React) → HTTP Request → Handler → Service → Repository → PostgreSQL
                                       ↑
                              JWT Middleware
```

## Funcionalidades Planejadas

### MVP (v1.0)
- [x] Estrutura base do projeto
- [ ] Autenticação JWT (login, registro, refresh token)
- [ ] CRUD completo de produtos (nome, SKU, categoria, preço, estoque)
- [ ] Registro de movimentações (entrada / saída)
- [ ] Alertas de estoque mínimo
- [ ] Dashboard com resumo de estoque

### v2.0
- [ ] Relatórios em PDF/Excel
- [ ] Multi-tenant (múltiplas empresas)
- [ ] Histórico de preços
- [ ] API de fornecedores
- [ ] Frontend React completo

## Rodando localmente

### Pré-requisitos

- Go 1.22+
- PostgreSQL 16+
- Docker (opcional)

### Com Docker Compose

```bash
# Clonar o repositório
git clone https://github.com/reinaldobarreto31/stockwise-go.git
cd stockwise-go

# Subir banco e aplicação
docker compose up -d
```

### Sem Docker

```bash
# Instalar dependências
go mod download

# Configurar variáveis de ambiente
cp .env.example .env
# Editar .env com suas credenciais PostgreSQL

# Executar migrations
# (em breve — usando golang-migrate)

# Rodar a API
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`.

## Variáveis de Ambiente

```env
# Banco de dados
DATABASE_URL=postgres://user:password@localhost:5432/stockwise?sslmode=disable

# JWT
JWT_SECRET=sua-chave-secreta-aqui
JWT_EXPIRATION_HOURS=24

# Servidor
PORT=8080
ENV=development
```

## Endpoints (planejados)

### Auth
| Método | Rota              | Descrição            |
|--------|-------------------|----------------------|
| POST   | `/api/auth/login` | Login do usuário     |
| POST   | `/api/auth/register` | Registro         |
| POST   | `/api/auth/refresh` | Refresh token     |

### Produtos
| Método | Rota                  | Descrição             |
|--------|-----------------------|-----------------------|
| GET    | `/api/products`       | Listar produtos       |
| GET    | `/api/products/:id`   | Buscar produto        |
| POST   | `/api/products`       | Criar produto         |
| PUT    | `/api/products/:id`   | Atualizar produto     |
| DELETE | `/api/products/:id`   | Remover produto       |

### Movimentações
| Método | Rota                    | Descrição               |
|--------|-------------------------|-------------------------|
| GET    | `/api/movements`        | Listar movimentações    |
| POST   | `/api/movements`        | Registrar movimentação  |

## Autor

**Reinaldo Barreto** — [github.com/reinaldobarreto31](https://github.com/reinaldobarreto31)

---

> Projeto em desenvolvimento ativo. Em breve com frontend React e documentação Swagger.
