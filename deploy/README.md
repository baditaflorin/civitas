# Civitas Deployment

Frontend:

https://baditaflorin.github.io/civitas/

Backend image:

ghcr.io/baditaflorin/civitas:latest

## Prerequisites

- Docker Engine and Docker Compose plugin
- DNS pointing to the server
- Let's Encrypt certificate mounted at `/etc/letsencrypt/live/civitas`
- GHCR access if the image is private; public images can be pulled anonymously

## First Deploy

```sh
cd deploy
cp .env.example .env
docker compose pull
docker compose up -d
```

The nginx service publishes TLS on host port `25342` and proxies to the internal Go API at `app:8080`.

## TLS

One common first-time certificate flow:

```sh
sudo certbot certonly --standalone -d your-domain.example
sudo mkdir -p /etc/letsencrypt/live/civitas
sudo cp /etc/letsencrypt/live/your-domain.example/fullchain.pem /etc/letsencrypt/live/civitas/fullchain.pem
sudo cp /etc/letsencrypt/live/your-domain.example/privkey.pem /etc/letsencrypt/live/civitas/privkey.pem
```

Adjust `deploy/nginx/nginx.conf` `server_name` before production use.

## Update

```sh
cd deploy
docker compose pull
docker compose up -d
```

## Rollback

Pin the image tag in `deploy/docker-compose.yml`, then:

```sh
docker compose pull
docker compose up -d
```

## Logs

```sh
docker compose logs -f app
docker compose logs -f nginx
```

## Metrics

```sh
docker compose --profile observability up -d prometheus
```

Prometheus is exposed only on `127.0.0.1:9090`. nginx blocks public `/metrics`.
