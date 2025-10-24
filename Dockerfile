# ---- Build stage ----
FROM golang:1.24.7-bookworm AS builder

WORKDIR /app

# Build-time argument: "api" or "collector"
ARG PORT=8080
ENV SERVICE=exporter
ENV PORT=${PORT}

COPY pkg pkg
# Copy only what's needed first
COPY go.mod  ./
RUN go mod tidy
# Copy the rest
COPY . .

# Build only the selected binary
RUN echo "Building ${SERVICE}..." && \
    CGO_ENABLED=1 GOOS=linux go build -o bin/${SERVICE} ./cmd/${SERVICE}

# ---- Runtime stage ----
FROM debian:bookworm-slim AS runner
WORKDIR /app

# Copy the single binary
ENV SERVICE=exporter
COPY --from=builder /app/bin/${SERVICE} ./${SERVICE}

# Copy config if needed
COPY config ./config
COPY .ver .ver

# Optional: expose ports (only if the API needs it)
EXPOSE $PORT

# Run the chosen service
ENTRYPOINT ["sh", "-c", "./${SERVICE}"]
