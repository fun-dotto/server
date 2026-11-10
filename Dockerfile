# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.27.3-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

COPY . .

# TODO: BFF廃止後は、admin-apiとapp-apiと各種Jobのみがビルド対象となるように変更する
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  set -eux; \
  mkdir -p /out/bin; \
  for cmd in academic-api admin-api announcement-api app-api app-bff-api build-class-change-notifications-job dispatch-notifications-job user-api; do \
  CGO_ENABLED=0 GOOS=linux \
  go build -tags timetzdata -trimpath -ldflags='-s -w' \
  -o /out/bin/${cmd} ./cmd/${cmd}; \
  done

# ---- Runtime stage ----
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/bin/ /bin/

USER nonroot:nonroot

# 各 Cloud Run Service / Job は command で /bin/<name> を指定して起動する。
ENTRYPOINT []
