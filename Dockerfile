# syntax=docker/dockerfile:1.7
FROM golang:1.26.10-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# distroless: no shell, no package manager, runs as uid 65532
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/api /app/api
COPY --from=build /src/seed /app/seed
# PORT is provided by the host (Railway); UPLOAD_DIR should be a mounted volume.
ENV UPLOAD_DIR=/data/uploads SEED_FILE=/app/seed/demo.sql
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/api"]
