# Pokemon API 

Backend en Go con Clean Architecture que consume la API de PokeAPI.

## Características

- ✨ Clean Architecture (Domain, Application, Infrastructure, Delivery)
- 🔄 Hot reload con Air para desarrollo
- 📊 Api rest para Pokemon
- 🔄 Sincronización con PokeAPI
- 🧪 Tests unitarios
- 📝 API REST bien documentada en readme

## Arquitectura

```
├── cmd/
│   └── server/          # Entry point para desarrollo local
├── internal/
│   ├── domain/          # Entidades, interfaces, reglas de negocio
│   ├── application/     # Casos de uso
│   ├── infrastructure/  # Implementaciones externas (PokeAPI)
│   └── delivery/        # Handlers HTTP 
├── bin/                 # Binarios compilados
└── tmp/                 # Archivos temporales de Air
```

## Endpoints API

### Health Check
- `GET /health` - Status del servicio

### Pokemon
- `GET /api/v1/pokemon` - Listar todos los Pokemon (con filtros)
- `GET /api/v1/pokemon/{id}` - Obtener Pokemon por ID
- `GET /api/v1/pokemon/name/{name}` - Obtener Pokemon por nombre


## Instalación y Configuración

### Prerrequisitos

```bash
# Instalar GO Lasted (Necesario)
go version

# Instalar Air (Necesario)
go install github.com/air-verse/air@latest

# Herramientas de desarrollo(Opcional)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Instalación

```bash
# Clonar el proyecto
git clone <repository-url>
cd reto-pokemon-api

# Instalar dependencias
make deps

```

## Desarrollo

### Servidor Local con Hot Reload

```bash
# Desarrollo con Air (hot reload)
make dev

# O directamente con Go
make run-local
```

El servidor estará disponible en `http://localhost:8080`

## Testing

```bash
# Ejecutar tests
make test

# Tests con coverage
make test-coverage

# Linting
make lint

# Security check
make security
```

## Ejemplos de Uso


### Obtener todos los Pokemon

```bash
curl http://localhost:8080/api/v1/pokemon
```

### Filtrar Pokemon

```bash
# Por nombre
curl "http://localhost:8080/api/v1/pokemon/name/pikachu"

# Con paginación
curl "http://localhost:8080/api/v1/pokemon?limit=10&offset=0"
```



