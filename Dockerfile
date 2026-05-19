FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
ARG VERSION=dev
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X nezdemos-api/internal/buildinfo.Version=${VERSION}" -o /out/nezdemos-api ./cmd/nezdemos-api

FROM alpine:3.22

RUN addgroup -S -g 1000 app && adduser -S -u 1000 -G app app
COPY --from=builder /out/nezdemos-api /usr/local/bin/nezdemos-api
USER app

ENTRYPOINT ["nezdemos-api"]
