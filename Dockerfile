# ── builder stage ──
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git nodejs npm

WORKDIR /app

COPY go.mod go.sum package.json package-lock.json ./
RUN go mod download && npm ci --quiet


COPY . .

RUN npx tailwindcss \
      -i api/static/input.css \
      -o api/static/styles.css \
      --minify

RUN go build -o app main.go


# ── runtime stage ─
FROM alpine:3.18

RUN apk add --no-cache \
      ca-certificates \
      docker-cli \
      clang \
      python3 \
      py3-pip \
      nodejs \
      npm && \
    pip3 install --no-cache-dir python-lsp-server && \
    pip3 install --no-cache-dir pylsp-mypy pylsp-rope && \
    npm install -g pyright sql-language-server

WORKDIR /app

COPY --from=builder /app/app ./
COPY --from=builder /app/api/static ./api/static
COPY --from=builder /app/api/templates ./api/templates

ENV PORT=":8080"

EXPOSE 8080

ENTRYPOINT ["./app"]

