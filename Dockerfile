# ── builder stage ──
FROM golang:1.23-alpine AS builder

# нам нужны git, nodejs и npm только для сборки статики
RUN apk add --no-cache git nodejs npm

WORKDIR /app

# скачиваем go-модули и npm-зависимости
COPY go.mod go.sum package.json package-lock.json ./
RUN go mod download && npm ci --quiet


# копируем весь код
COPY . .

# собираем CSS через tailwindcss
RUN npx tailwindcss \
      -i api/static/input.css \
      -o api/static/styles.css \
      --minify

# компилируем бинарник
RUN go build -o app main.go


# ── runtime stage ─
FROM alpine:3.18

# runtime-зависимости и LSP-сервера
RUN apk add --no-cache \
      ca-certificates \
      docker-cli \
      clang \
      python3 \
      py3-pip \
      nodejs \
      npm && \
    # Python LSP
    pip3 install --no-cache-dir python-lsp-server && \
    pip3 install --no-cache-dir pylsp-mypy pylsp-rope && \
    # JS LSP-серверы
    npm install -g pyright sql-language-server

WORKDIR /app

# копируем Go-приложение и статику
COPY --from=builder /app/app ./
COPY --from=builder /app/api/static ./api/static
COPY --from=builder /app/api/templates ./api/templates

# порт и прочие переменные задаются в docker-compose
ENV PORT=":8080"

EXPOSE 8080

ENTRYPOINT ["./app"]

