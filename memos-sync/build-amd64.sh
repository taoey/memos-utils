docker run --rm \
  --platform=linux/amd64 \
  -v "$PWD":/app \
  -w /app \
  golang:1.24-alpine \
  sh -c "
    apk add --no-cache gcc musl-dev pkgconfig sqlite-dev && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags '-extldflags \"-static\"' -o app-linux-amd64 .
  "