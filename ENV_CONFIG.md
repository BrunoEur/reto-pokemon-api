# Variables de Entorno

Este documento describe las variables de entorno disponibles para configurar la aplicación.

## Variables Disponibles

### PORT
- **Descripción**: Puerto en el que se ejecutará el servidor
- **Valor por defecto**: `8080`
- **Ejemplo**: `PORT=3000`

### POKEAPI_BASE_URL
- **Descripción**: URL base de la PokeAPI
- **Valor por defecto**: `https://pokeapi.co/api/v2`
- **Ejemplo**: `POKEAPI_BASE_URL=https://pokeapi.co/api/v2`
- **Uso**: Útil para apuntar a una instancia diferente de PokeAPI o para testing con un mock server

### ENV
- **Descripción**: Entorno de ejecución
- **Valor por defecto**: `development`
- **Valores posibles**: `development`, `staging`, `production`
- **Ejemplo**: `ENV=production`

## Configuración Local

### 1. Crear archivo .env

Copia el archivo de ejemplo:

```bash
cp .env.example .env
```

### 2. Editar las variables

Edita el archivo `.env` con tus valores:

```bash
PORT=8080
POKEAPI_BASE_URL=https://pokeapi.co/api/v2
ENV=development
```

### 3. Ejecutar la aplicación

Las variables se cargarán automáticamente:

```bash
# Con make
make run-local

# O directamente
go run cmd/server/main.go
```

## Configuración con Docker

### Pasar variables al contenedor

```bash
# Usando -e para variables individuales
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e POKEAPI_BASE_URL=https://pokeapi.co/api/v2 \
  pokemon-api:latest

# Usando --env-file
docker run -p 8080:8080 --env-file .env pokemon-api:latest
```

### Docker Compose

Crea un archivo `docker-compose.yml`:

```yaml
version: '3.8'
services:
  pokemon-api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - POKEAPI_BASE_URL=https://pokeapi.co/api/v2
      - ENV=production
    # O usar env_file
    # env_file:
    #   - .env
```

## Configuración en AWS ECS

### Task Definition

En tu `task-definition.json`, añade las variables de entorno:

```json
{
  "containerDefinitions": [
    {
      "name": "pokemon-api",
      "environment": [
        {
          "name": "PORT",
          "value": "8080"
        },
        {
          "name": "POKEAPI_BASE_URL",
          "value": "https://pokeapi.co/api/v2"
        },
        {
          "name": "ENV",
          "value": "production"
        }
      ]
    }
  ]
}
```

### Usando AWS Systems Manager Parameter Store

Para valores sensibles, usa Parameter Store:

```json
{
  "containerDefinitions": [
    {
      "name": "pokemon-api",
      "secrets": [
        {
          "name": "POKEAPI_BASE_URL",
          "valueFrom": "arn:aws:ssm:region:account-id:parameter/pokemon-api/pokeapi-url"
        }
      ]
    }
  ]
}
```

## Testing

Para tests, puedes usar diferentes valores:

```bash
# Test con URL mock
POKEAPI_BASE_URL=http://localhost:3001 go test ./...

# O en el código de tests, puedes usar:
os.Setenv("POKEAPI_BASE_URL", "http://mock-pokeapi.test")
```

## Validación

Para verificar que las variables se están cargando correctamente, puedes añadir logs en el inicio de la aplicación o verificar los endpoints.

## Seguridad

⚠️ **Importante**: 
- Nunca commitas el archivo `.env` al repositorio
- El archivo `.env` ya está incluido en `.gitignore`
- Para producción, usa servicios como AWS Secrets Manager o Parameter Store
- No incluyas valores sensibles en el Dockerfile

## Troubleshooting

### La aplicación no lee las variables

1. Verifica que el archivo `.env` existe en el directorio raíz
2. Verifica que las variables están correctamente definidas (sin espacios alrededor del `=`)
3. En Docker, asegúrate de pasar las variables con `-e` o `--env-file`

### La aplicación usa valores por defecto

Si no se definen las variables de entorno, la aplicación usará los valores por defecto:
- PORT: `8080`
- POKEAPI_BASE_URL: `https://pokeapi.co/api/v2`
- ENV: `development`
