# Proyecto 1: Mini Blockchain en Go

## Contexto

Soy un Backend Developer con experiencia en Python (FastAPI, Django), Go, TypeScript (NestJS), Java (Spring Boot) y arquitectura hexagonal/microservicios. Estoy aprendiendo blockchain desde cero construyendo proyectos progresivos. Este es el **Proyecto 1 de 6** de mi ruta de aprendizaje.

**Objetivo de aprendizaje:** Entender cómo funciona una blockchain internamente: estructura de bloques, hashing criptográfico, encadenamiento, Proof of Work y validación de integridad. Sin depender de ninguna librería blockchain externa — Go puro.

---

## Descripción del proyecto

Construir una **blockchain funcional desde cero en Go** que incluya:

1. Estructura de bloque con hash SHA-256
2. Encadenamiento de bloques (cada bloque apunta al hash del anterior)
3. Algoritmo de Proof of Work (minería con dificultad configurable)
4. API REST con Gin para interactuar con la cadena
5. Validación de integridad de la cadena completa
6. Persistencia del estado en JSON
7. Tests unitarios con el paquete `testing` de Go

---

## Stack tecnológico

| Herramienta     | Versión    | Uso                                      |
|-----------------|------------|------------------------------------------|
| Go              | 1.22+      | Lenguaje principal                       |
| Gin             | v1.9       | Framework HTTP para la API REST          |
| crypto/sha256   | stdlib     | Hashing de bloques                       |
| encoding/json   | stdlib     | Serialización y persistencia             |
| testing         | stdlib     | Tests unitarios                          |
| Docker          | latest     | Containerización del servicio            |
| docker-compose  | v2         | Orquestación local                       |

---

## Arquitectura del proyecto

```
mini-blockchain/
├── cmd/
│   └── main.go                  # Entry point — inicializa blockchain y arranca Gin
├── internal/
│   ├── blockchain/
│   │   ├── block.go             # Struct Block, CalculateHash(), MineBlock()
│   │   ├── blockchain.go        # Struct Blockchain, AddBlock(), IsValid(), Genesis()
│   │   ├── blockchain_test.go   # Tests unitarios completos
│   │   └── persistence.go       # Guardar/cargar cadena en chain.json
│   └── api/
│       ├── handlers.go          # Handlers de Gin para cada endpoint
│       ├── router.go            # Setup de rutas
│       └── dto.go               # Request/Response structs
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

**Patrón de arquitectura:** Separación clara entre dominio (`internal/blockchain`) e infraestructura (`internal/api`). El paquete `blockchain` no debe importar nada de `api` — dependencia unidireccional.

---

## Estructura de datos

### Block

```go
type Block struct {
    Index     int    `json:"index"`
    Timestamp string `json:"timestamp"`
    Data      string `json:"data"`
    PrevHash  string `json:"prev_hash"`
    Hash      string `json:"hash"`
    Nonce     int    `json:"nonce"`
}
```

### Blockchain

```go
type Blockchain struct {
    Chain      []Block `json:"chain"`
    Difficulty int     `json:"difficulty"`
}
```

---

## Lógica de negocio requerida

### 1. Cálculo de hash (`block.go`)

El hash de un bloque debe ser el SHA-256 de la concatenación de todos sus campos:

```
hash = SHA256(index + timestamp + data + prev_hash + nonce)
```

Si **cualquier campo cambia**, el hash cambia completamente. Esto garantiza inmutabilidad.

### 2. Proof of Work (`block.go`)

El minero debe incrementar el `Nonce` hasta encontrar un hash que comience con `Difficulty` caracteres `"0"`:

```
Difficulty = 4  →  hash debe empezar con "0000"
Difficulty = 3  →  hash debe empezar con "000"
```

```go
func (b *Block) MineBlock(difficulty int) {
    target := strings.Repeat("0", difficulty)
    for !strings.HasPrefix(b.Hash, target) {
        b.Nonce++
        b.Hash = b.CalculateHash()
    }
}
```

### 3. Creación del bloque génesis (`blockchain.go`)

El primer bloque de la cadena (index=0) tiene `PrevHash = "0"` y `Data = "Genesis Block"`. Debe ser minado igual que los demás.

### 4. Agregar bloque (`blockchain.go`)

```go
func (bc *Blockchain) AddBlock(data string) Block {
    // 1. Obtener el último bloque de la cadena
    // 2. Crear nuevo bloque con PrevHash = último hash
    // 3. Minar el nuevo bloque
    // 4. Agregar a bc.Chain
    // 5. Persistir en JSON
    // 6. Retornar el bloque minado
}
```

### 5. Validación de cadena (`blockchain.go`)

```go
func (bc *Blockchain) IsValid() bool {
    // Para cada bloque i (desde i=1):
    //   a. Recalcular el hash del bloque i y comparar con block.Hash
    //   b. Verificar que block[i].PrevHash == block[i-1].Hash
    //   c. Verificar que el hash cumple la dificultad (Proof of Work)
    // Si cualquier verificación falla → return false
}
```

### 6. Persistencia (`persistence.go`)

- `SaveChain(bc Blockchain, path string) error` → serializar a JSON y escribir en archivo
- `LoadChain(path string) (Blockchain, error)` → leer JSON y deserializar

---

## API REST — Endpoints requeridos

| Método | Ruta       | Descripción                              | Request Body           | Response                   |
|--------|------------|------------------------------------------|------------------------|----------------------------|
| GET    | /chain     | Retorna la cadena completa               | —                      | `{ chain: [...], length: N }` |
| POST   | /mine      | Mina un nuevo bloque con los datos dados | `{ "data": "string" }` | El bloque recién minado    |
| GET    | /validate  | Valida la integridad de la cadena        | —                      | `{ "valid": bool }`        |
| GET    | /block/:id | Retorna un bloque por su índice          | —                      | El bloque solicitado       |
| GET    | /health    | Health check del servicio                | —                      | `{ "status": "ok" }`       |

### Ejemplo de respuesta GET /chain

```json
{
  "chain": [
    {
      "index": 0,
      "timestamp": "2026-03-13T10:00:00Z",
      "data": "Genesis Block",
      "prev_hash": "0",
      "hash": "0000a3f...",
      "nonce": 72914
    },
    {
      "index": 1,
      "timestamp": "2026-03-13T10:01:23Z",
      "data": "Primera transacción",
      "prev_hash": "0000a3f...",
      "hash": "0000b7c...",
      "nonce": 45231
    }
  ],
  "length": 2,
  "difficulty": 4
}
```

### Manejo de errores

Todos los errores deben retornar JSON con la estructura:

```json
{ "error": "descripción del error" }
```

Códigos HTTP: `400` para bad request, `404` para recurso no encontrado, `500` para errores internos.

---

## Tests requeridos (`blockchain_test.go`)

Deben cubrir los siguientes escenarios con el paquete `testing` estándar de Go:

```go
// 1. El bloque génesis se crea correctamente
func TestGenesisBlock(t *testing.T)

// 2. El hash calculado es determinístico (mismos inputs → mismo output)
func TestHashDeterminism(t *testing.T)

// 3. Proof of Work: el hash minado cumple la dificultad
func TestProofOfWork(t *testing.T)

// 4. Cambiar cualquier campo invalida el hash (efecto avalancha)
func TestHashChangesOnDataMutation(t *testing.T)

// 5. Cadena válida retorna IsValid() = true
func TestValidChain(t *testing.T)

// 6. Modificar un bloque pasado invalida la cadena
func TestTamperedChainIsInvalid(t *testing.T)

// 7. La cadena persiste y se recarga correctamente
func TestPersistenceRoundtrip(t *testing.T)
```

Comando de ejecución: `go test ./... -v`

---

## Docker

### Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o blockchain ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/blockchain .
EXPOSE 8080
CMD ["./blockchain"]
```

### docker-compose.yml

```yaml
services:
  blockchain:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data   # Para persistir chain.json fuera del container
    environment:
      - DIFFICULTY=4
      - PORT=8080
```

---

## Variables de entorno (`.env.example`)

```env
DIFFICULTY=4
PORT=8080
CHAIN_FILE=./data/chain.json
```

---

## Criterios de completitud (Definition of Done)

- [ ] `go test ./... -v` → todos los tests pasan, sin errores
- [ ] `go vet ./...` → sin warnings
- [ ] API responde correctamente a los 5 endpoints documentados
- [ ] Modificar un bloque pasado invalida `GET /validate` → `{ "valid": false }`
- [ ] `docker compose up` levanta el servicio en el puerto 8080
- [ ] Persistencia: reiniciar el contenedor mantiene la cadena anterior
- [ ] README con: descripción, cómo correrlo, cómo correr tests, y diagrama ASCII de la arquitectura
- [ ] Código en GitHub con commits atómicos y descriptivos

---

## Flujo de desarrollo sugerido

1. **Primero el dominio:** implementar y testear `block.go` y `blockchain.go` sin nada de HTTP
2. **Luego la API:** una vez que los tests del dominio pasan, construir los handlers de Gin
3. **Después la persistencia:** agregar `persistence.go` y el test de roundtrip
4. **Por último Docker:** containerizar una vez que todo funciona localmente

---

## Conceptos de blockchain que vas a entender al terminar

| Concepto              | Cómo lo implementas en este proyecto                         |
|-----------------------|--------------------------------------------------------------|
| Bloque                | `struct Block` con índice, datos, hash y nonce              |
| Hash criptográfico    | `crypto/sha256` — cualquier cambio → hash completamente distinto |
| Encadenamiento        | `PrevHash` en cada bloque apunta al hash del anterior       |
| Inmutabilidad         | Cambiar un bloque rompe todos los hashes siguientes         |
| Proof of Work         | Loop de nonce hasta que el hash empiece con N ceros         |
| Dificultad            | `strings.Repeat("0", difficulty)` — más ceros = más trabajo |
| Génesis               | Primer bloque con `PrevHash = "0"` — punto de partida       |
| Validación            | Recalcular hashes y verificar encadenamiento en toda la cadena |

---

## Recursos de referencia

- [Jeiwan - Build Blockchain in Go](https://github.com/Jeiwan/blockchain_go) — la referencia principal
- [Go crypto/sha256 docs](https://pkg.go.dev/crypto/sha256)
- [Gin quickstart](https://gin-gonic.com/docs/quickstart/)
- Bitcoin Whitepaper sección 2-4 (la lógica de bloques y PoW)
