# ---------- Stage 1: Build the Vue frontend ----------
FROM node:20-alpine AS webbuilder
WORKDIR /frontend
# Copy only package descriptors first for better caching
COPY frontend/package*.json ./
RUN npm ci --no-audit --no-fund
# Copy the rest of the frontend and build
COPY frontend/ ./
RUN npm run build

# ---------- Stage 2: Build the Go backend ----------
FROM golang:1.25-alpine AS gobuilder
WORKDIR /app

# System deps (optional: ca-certificates are required for outbound HTTPS like Gemini)
RUN apk add --no-cache ca-certificates tzdata

# Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the backend code
COPY . .

# We only need the dist from the frontend stage (copy to ./dist)
COPY --from=webbuilder /frontend/dist ./dist

# Build statically linked binary for minimal runtime base image
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./main.go

# ---------- Stage 3: Minimal runtime image ----------
FROM gcr.io/distroless/base-debian12 AS runtime
WORKDIR /app

# Copy binary and required assets/config
COPY --from=gobuilder /app/server ./server
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# If your app reads these at runtime, copy them in (optional)
# Remove if you embed config or read from env instead
COPY --from=gobuilder /app/config.yaml ./config.yaml
# If you need static files like JSONs under /data at runtime:
# COPY --from=gobuilder /app/data ./data

# The Vue build output must be present at /app/dist
COPY --from=gobuilder /app/dist ./dist

# Cloud Run listens on PORT env (default 8080). Expose for local runs.
EXPOSE 8080

# Non-root by default on distroless/base (UID 65532)
USER 65532:65532

# Start
ENTRYPOINT ["./server"]
