# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.26.3-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=0.3.0
ARG COMMIT_SHA=dev

WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata
RUN mkdir -p /out /empty
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
  go build -trimpath -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT_SHA}" \
  -o /out/civitas ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

ARG VERSION=0.3.0
ARG COMMIT_SHA=dev
ARG CREATED=unknown

LABEL org.opencontainers.image.title="Civitas" \
  org.opencontainers.image.description="Civic investigation OS backend API" \
  org.opencontainers.image.source="https://github.com/baditaflorin/civitas" \
  org.opencontainers.image.revision="${COMMIT_SHA}" \
  org.opencontainers.image.version="${VERSION}" \
  org.opencontainers.image.created="${CREATED}" \
  org.opencontainers.image.licenses="MIT"

WORKDIR /
COPY --from=builder --chown=nonroot:nonroot /empty /var/lib/civitas
COPY --from=builder /out/civitas /civitas

ENV CIVITAS_ADDR=:8080
EXPOSE 8080
USER nonroot:nonroot
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD ["/civitas", "healthcheck"]
ENTRYPOINT ["/civitas"]
