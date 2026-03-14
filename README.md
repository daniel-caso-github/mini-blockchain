# Mini Blockchain

A minimal blockchain implementation in Go. Built from scratch to understand blockchain internals: blocks, SHA-256 hashing, chain linking, Proof of Work, and validation.

## Features

- Block creation with SHA-256 hashing
- Proof of Work mining with configurable difficulty
- Chain validation (hash integrity + linking)
- SQLite persistence (pure Go, no CGO)
- REST API (Gin)
- Thread-safe with `sync.RWMutex`

## Architecture

```
mini-blockchain/
├── cmd/
│   └── main.go              ← entry point
├── internal/
│   ├── blockchain/           ← domain (no external deps)
│   │   ├── block.go          ← Block struct, CalculateHash(), MineBlock()
│   │   ├── blockchain.go     ← Blockchain struct, AddBlock(), IsValid()
│   │   ├── persistence.go    ← Store (SQLite persistence)
│   │   └── blockchain_test.go
│   └── api/                  ← HTTP infrastructure (depends on blockchain)
│       ├── router.go         ← Gin route setup
│       ├── handlers.go       ← endpoint handlers
│       └── dto.go            ← request/response structs
├── config/
│   └── config.go             ← env var loading
├── Dockerfile
├── docker-compose.yml
└── Taskfile.yml
```

### Dependency Flow

```
cmd/main.go → config.Load()
            → blockchain.New()
            → api.SetupRouter(bc)
            → router.Run()

api/ ──depends──▶ internal/blockchain/
internal/blockchain/ ──depends──▶ stdlib + modernc.org/sqlite
```

## Prerequisites

- [Go](https://go.dev/dl/) 1.21+
- [Task](https://taskfile.dev/installation/) — task runner (`brew install go-task` en macOS)
- [Node.js](https://nodejs.org/) 18+ (para el frontend)
- [Docker](https://docs.docker.com/get-docker/) (opcional, para contenedores)

## Quick Start

```bash
task run
```

Or with Docker:

```bash
task docker:up
```

## Tasks

| Task | Descripción |
|------|-------------|
| `task build` | Compila el binario |
| `task run` | Ejecuta sin compilar |
| `task test` | Corre todos los tests |
| `task vet` | Análisis estático |
| `task lint` | Alias de vet |
| `task swagger` | Genera/actualiza documentación Swagger |
| `task clean` | Limpia binario y datos |
| `task dev` | Inicia backend + frontend concurrentemente |
| `task frontend:install` | Instala dependencias del frontend |
| `task frontend:dev` | Inicia servidor de desarrollo frontend (puerto 5173) |
| `task frontend:build` | Build del frontend para producción |
| `task docker:build` | Construye imagen Docker |
| `task docker:up` | Levanta servicio en background |
| `task docker:down` | Baja el servicio |
| `task docker:logs` | Sigue los logs |

## API

| Method | Route | Description |
|--------|-------|-------------|
| GET | /health | Health check |
| GET | /chain | Get full chain |
| POST | /mine | Mine a block (`{"data": "..."}`) |
| GET | /validate | Validate chain integrity |
| GET | /block/:id | Get block by index |

## Swagger / API Docs

Con el servidor corriendo, abre [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html) para la documentación interactiva.

Para regenerar los docs después de modificar anotaciones:

```bash
task swagger
```

## Configuration

| Env Var | Default | Description |
|---------|---------|-------------|
| PORT | 8080 | Server port |
| DIFFICULTY | 2 | PoW difficulty (leading zeros) |
| DB_PATH | data/blockchain.db | SQLite database path |

## Tests

```bash
task test
```
