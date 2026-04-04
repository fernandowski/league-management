FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp ./apps/api


# Development
FROM golang:1.25 AS dev
WORKDIR /workspace
RUN adduser --disabled-password --gecos '' appuser
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY go.mod go.sum ./
RUN go mod download
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
