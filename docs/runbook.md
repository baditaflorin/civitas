# Civitas Runbook

## Health

```sh
curl -fsS http://localhost:8080/healthz
curl -fsS http://localhost:8080/readyz
```

## Logs

Docker logs:

```sh
docker compose -f deploy/docker-compose.yml logs -f app
docker compose -f deploy/docker-compose.yml logs -f nginx
```

Backend logs are JSON on stdout. They include method, path, status, duration, and request id. They must not include evidence body text.

## Metrics

Prometheus scrape target inside Compose:

http://app:8080/metrics

Local Prometheus with profile:

```sh
docker compose -f deploy/docker-compose.yml --profile observability up -d
```

Public nginx blocks `/metrics`.

## Storage

Default local storage:

`./storage`

Production storage volume:

`civitas-storage`

Back up:

```sh
docker run --rm -v civitas_civitas-storage:/data -v "$PWD":/backup alpine \
  tar czf /backup/civitas-storage.tgz -C /data .
```

Restore:

```sh
docker run --rm -v civitas_civitas-storage:/data -v "$PWD":/backup alpine \
  sh -c 'cd /data && tar xzf /backup/civitas-storage.tgz'
```

## Resource Sizing

The Go API skeleton runs comfortably with 1 CPU and 512 MB RAM. Native OCR, transcription, NLP, and local LLM adapters need more capacity. A practical v1 server target is 2 to 8 CPU cores, 8 to 32 GB RAM, and enough encrypted disk for evidence dumps and derived artifacts.

## Escalation

1. Check `/readyz`.
2. Check `docker compose ps`.
3. Check app logs for storage or processor errors.
4. Check disk pressure on the evidence volume.
5. Temporarily disable optional processor adapters by removing them from the backend image.
