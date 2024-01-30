# GoURL
URL shortener service written in Go

Following [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

## Docker

```sh
docker build -t go-url -f deployments/Dockerfile .
docker run -p 8080:8080 go-url
```

**Compose**
```sh
docker compose -f deployments/docker-compose.yaml up --build
```

## TODO

Use alpine JS for front-end

- [ ] implement readiness probe (/readyz endpoint)
- [ ] UI
  - Alpine.js / htmx for js in html
  - Templ