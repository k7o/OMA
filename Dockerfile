# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/ui
COPY /ui .
RUN npm ci && npm run build

# Stage 2: Build backend
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app
COPY . .
COPY --from=frontend-builder /app/ui/dist ./ui/dist
RUN go mod download && go build -o main .

# Stage 3: Final stage
FROM alpine:3.20.3 AS final

WORKDIR /app
COPY --from=backend-builder /app/main .
COPY --from=backend-builder /app/transport ./transport

EXPOSE 6060

CMD ["./main"]