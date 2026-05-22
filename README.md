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

## CI/CD

### CI

El repositorio incluye un workflow de GitHub Actions que:

- corre los tests
- compila la aplicación
- valida el build de la imagen Docker

### Publicación a Docker Hub

También se incluye un workflow para publicar la imagen Docker en Docker Hub cuando hay cambios en `main`.

Para que funcione, en GitHub tenés que crear estos secrets del repositorio:

- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

La imagen se publica como:

```text
<DOCKERHUB_USERNAME>/movies-api
```

El workflow publica siempre estas tags:

- `latest`
- una versión semántica, por ejemplo `v0.2.0`
- una tag única por commit, por ejemplo `sha-1a2b3c4`

La versión semántica se calcula automáticamente en cada merge a `main`:

- ramas o PRs `feat/...` incrementan `minor`
- ramas o PRs `fix/...` incrementan `patch`

El archivo `VERSION` solo se usa como versión base inicial si todavía no existe ninguna tag `v*` en el repositorio.

### Despliegue en Render

El despliegue se realiza con un Web Service de Render configurado con `Existing Image`.

Pasos manuales:

1. Crear una cuenta en Render.
2. Crear un `Web Service`.
3. Elegir `Existing Image` en lugar de conectar el repositorio.
4. Usar como imagen inicial `docker.io/<DOCKERHUB_USERNAME>/movies-api:latest`.
5. Seleccionar el plan deseado y crear el servicio.
6. Configurar las variables de entorno:
   - `TMDB_API_KEY`
   - `TMDB_BASE_URL`
   - `METRICS_USERNAME`
   - `METRICS_PASSWORD`

Render provee automáticamente la variable `PORT`, y la aplicación ya está preparada para escuchar ese puerto.

Después de crear el servicio:

1. Ir a `Settings` del servicio.
2. Copiar la URL del `Deploy Hook`.
3. Guardar esa URL en GitHub como secret `RENDER_DEPLOY_HOOK_URL`.

Una vez configurado, cada publicación de imagen sobre `main` dispara automáticamente el deploy en Render usando la misma etiqueta `sha-<commit>` que se publicó en Docker Hub.

### Monitoreo con Grafana Cloud

La API expone un endpoint Prometheus en:

```text
/metrics
```

Ese endpoint incluye métricas de runtime de Go y métricas HTTP de la aplicación:

- total de requests por método, ruta y código de estado
- duración de requests
- cantidad de requests en curso

El endpoint `/metrics` está protegido con Basic Auth usando estas variables de entorno:

- `METRICS_USERNAME`
- `METRICS_PASSWORD`

Pasos manuales para conectarlo con Grafana Cloud:

1. Crear una cuenta en Grafana Cloud.
2. Ir a `Connections` > `Add new connection`.
3. Buscar `Metrics Endpoint`.
4. Crear un scrape job con la URL pública:

```text
https://movies-api-latest.onrender.com/metrics
```

5. Completar autenticación `Basic` con el mismo usuario y contraseña configurados en Render.
6. Probar la conexión y guardar el scrape job.

Grafana Cloud hará el scrape automáticamente cada 60 segundos.

Consultas útiles para el dashboard:

```promql
sum by (route, status) (rate(movies_api_http_requests_total[5m]))
```

```promql
histogram_quantile(0.95, sum by (le, route) (rate(movies_api_http_request_duration_seconds_bucket[5m])))
```

```promql
movies_api_http_requests_in_flight
```

Para generar tráfico de prueba antes de la presentación:

```bash
make traffic
```

O directamente contra Render:

```bash
make traffic-render
```

Si querés sobrescribir la URL manualmente:

```bash
BASE_URL=https://movies-api-latest.onrender.com make traffic
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
