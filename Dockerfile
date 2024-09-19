FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp .


# Development
FROM golang:1.23.1 AS dev
WORKDIR /app
RUN adduser --disabled-password --gecos '' appuser
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
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
