# Civitas API

OpenAPI spec:

https://github.com/baditaflorin/civitas/blob/main/api/openapi.yaml

Default local backend:

http://localhost:8080

## Examples

```sh
curl -s http://localhost:8080/healthz
curl -s http://localhost:8080/readyz
curl -s http://localhost:8080/api/v1/version
curl -s http://localhost:8080/api/v1/processors
```

Create a case:

```sh
curl -s http://localhost:8080/api/v1/cases \
  -H 'content-type: application/json' \
  -d '{"title":"Harbor contracts","description":"procurement dump"}'
```

Upload evidence:

```sh
curl -s http://localhost:8080/api/v1/cases/case_id/documents \
  -F 'file=@sample.txt'
```

Search a case:

```sh
curl -s 'http://localhost:8080/api/v1/cases/case_id/search?q=contract'
```

Generate a safe publishing export:

```sh
curl -s http://localhost:8080/api/v1/cases/case_id/exports \
  -H 'content-type: application/json' \
  -d '{"format":"markdown"}'
```
