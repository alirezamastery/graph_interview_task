# Todo list API

---

Go 1.21+ required

## Configure environment variables

For convenience, I commited the .env files in git.

If you want to run the webserver using docker, no change is needed. but if you want run it directly from command line
copy the `env/.go.env` file into `src` and change these variables:

```
DOCKER=false
DB_HOST=127.0.0.1
DB_PORT=5450
```

---

## Setup

to run the webserver using docker just use this command: `sudo ocker compose up --build`

if you want to run the webserver directly from command line:

```shell
sudo docker compose up db -d
cd src
go run .
```

---

## Swagger

Swagger UI:

- `http://127.0.0.1:8000/swagger/index.html`

The raw OpenAPI JSON (useful for debugging):

- `http://127.0.0.1:8000/swagger/doc.json`

---

## Prometheus

Metrics endpoint:

- `http://127.0.0.1:8000/metrics`

---

## Unit Tests

From src run this command:

```bash
go test ./tests -cover -v
```

---

## Testing with curl

### 1) Create a todo (POST)

```bash
curl -i -X POST "http://127.0.0.1:8000/api/task/todos" \
  -H "Content-Type: application/json" \
  -d '{"title":"test 1","description":"test desc","is_done":false}'
```

### 2) List todos (GET)

```bash
curl -i "http://127.0.0.1:8000/api/task/todos"
```

If your list endpoint supports pagination and filters, examples:

```bash
curl -i "http://127.0.0.1:8000/api/task/todos?page=1&page_size=20"
curl -i "http://127.0.0.1:8000/api/task/todos?is_done=true"
```

### 3) Get todo by ID (GET)

```bash
curl -i "http://127.0.0.1:8000/api/task/todos/1"
```

### 4) Update todo (PATCH)

Example: change is_done

```bash
curl -i -X PATCH "http://127.0.0.1:8000/api/task/todos/1" \
  -H "Content-Type: application/json" \
  -d '{"is_done":true}'
```

Example: update title/description

```bash
curl -i -X PATCH "http://127.0.0.1:8000/api/task/todos/1" \
  -H "Content-Type: application/json" \
  -d '{"title":"test 1 updated","description":"desc udpated"}'
```

### 5) Delete todo (DELETE)

```bash
curl -i -X DELETE "http://127.0.0.1:8000/api/task/todos/1"
```

---
