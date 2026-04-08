FROM golang:1.25 AS builder

WORKDIR /app

COPY apps/api/go.mod apps/api/go.sum ./apps/api/

WORKDIR /app/apps/api
RUN go mod download

WORKDIR /app
COPY . .

WORKDIR /app/apps/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp .


# Development
FROM golang:1.25 AS dev
WORKDIR /workspace
RUN adduser --disabled-password --gecos '' appuser
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
WORKDIR /workspace/apps/api
RUN go mod download
WORKDIR /workspace
COPY . .
WORKDIR /workspace/apps/api
USER appuser
EXPOSE 8080

CMD ["air"]

# Production Stage
FROM alpine:latest AS prod

# Set the working directory
WORKDIR /app

# Copy the built Go binary from the builder stage
COPY --from=builder /app/myapp .

# Expose the port
EXPOSE 8080

# Command to run the production app
CMD ["./myapp"]
