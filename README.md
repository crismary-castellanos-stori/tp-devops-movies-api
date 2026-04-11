# Movies API

API REST en Go que consume la API de TMDB (The Movie Database) para búsqueda y consulta de películas.

## Requisitos

- Go 1.26.1 o superior
- Cuenta en [TMDB](https://www.themoviedb.org/) para obtener una API Key

## Configuración

1. Clonar el repositorio:

```bash
git clone https://github.com/crismary-castellanos-stori/tp-devops-movies-api.git
cd tp-devops-movies-api
```

2. Crear el archivo de variables de entorno:

```bash
cp .env.example .env
```

3. Editar `.env` y configurar tu API Key de TMDB:

```
PORT=8080
TMDB_API_KEY=tu_api_key_aqui
TMDB_BASE_URL=https://api.themoviedb.org/3
```

Para obtener una API Key de TMDB:
- Crear una cuenta en https://www.themoviedb.org/
- Ir a Settings > API
- Solicitar una API Key (tipo "Developer")
- Usar el "API Read Access Token" (Bearer token)

## Ejecución

### Opción 1: Ejecutar directamente

```bash
go run api/main.go
```

### Opción 2: Compilar y ejecutar

```bash
go build -o movies-api api/main.go
./movies-api
```

El servidor iniciará en `http://localhost:8080` (o el puerto configurado en `PORT`).

### Opción 3: Ejecutar con Docker

Construir la imagen:

```bash
docker build -t movies-api:latest .
```

Ejecutar el contenedor usando el archivo `.env`:

```bash
docker run --rm -p 8080:8080 --env-file .env movies-api:latest
```

El servidor quedará disponible en `http://localhost:8080`.

### Opción 4: Usar Makefile

Comandos disponibles:

```bash
make help
make test
make run
make build
make docker-build
make docker-run
make compose-up
make compose-down
```

### Opción 5: Ejecutar con Docker Compose

Levantar la aplicación:

```bash
docker compose up --build
```

Detenerla:

```bash
docker compose down
```

## Endpoints

### Health Check

```bash
curl http://localhost:8080/health
```

Respuesta:
```json
{"status":"ok"}
```

### Buscar películas

```bash
curl "http://localhost:8080/movies/search?title=batman"
```

Respuesta:
```json
{
  "results": [
    {
      "id": 268,
      "title": "Batman",
      "overview": "Batman must face his most ruthless nemesis...",
      "release_date": "1989-06-23",
      "poster_path": "/tDexQyu6FWltcd0VhEDK7uib42f.jpg",
      "vote_average": 7.2
    }
  ],
  "total": 45
}
```

### Obtener película por ID

```bash
curl http://localhost:8080/movies/550
```

Respuesta:
```json
{
  "id": 550,
  "title": "Fight Club",
  "overview": "A ticking-Loss time bomb of a movie...",
  "release_date": "1999-10-15",
  "poster_path": "/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg",
  "vote_average": 8.4
}
```

## Estructura del Proyecto

```
tp-devops-movies-api/
├── api/
│   └── main.go              # Entry point de la aplicación
├── internal/
│   ├── config/
│   │   └── config.go        # Configuración desde variables de entorno
│   ├── handler/
│   │   ├── health.go        # Handler GET /health
│   │   ├── movies.go        # Handlers de películas
│   │   └── response.go      # Helpers para respuestas JSON
│   ├── service/
│   │   └── movies.go        # Lógica de negocio
│   ├── client/
│   │   └── tmdb.go          # Cliente HTTP para TMDB API
│   └── model/
│       └── movie.go         # Modelos de datos
├── .env.example             # Ejemplo de variables de entorno
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `PORT` | Puerto del servidor HTTP | `8080` |
| `TMDB_API_KEY` | API Key de TMDB (requerido) | - |
| `TMDB_BASE_URL` | URL base de la API de TMDB | `https://api.themoviedb.org/3` |

## Códigos de Estado

| Código | Descripción |
|--------|-------------|
| 200 | OK |
| 400 | Bad Request (parámetro faltante o inválido) |
| 404 | Not Found (película no encontrada) |
| 502 | Bad Gateway (error al comunicarse con TMDB) |
