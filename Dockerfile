# Use an official Go runtime as a parent image
FROM golang:1.22-alpine3.20

# Install libs
RUN apk add --no-cache \
    nano \
    supervisor \
    ca-certificates\
    && rm -rf /var/lib/apt/lists/*

# Install ZSH
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v1.2.0/zsh-in-docker.sh)" -- \
    -t robbyrussell

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/server/main.go

# Expose the port on which the app will listen
EXPOSE 8080

# Copy and prepare the entrypoint script
COPY ./docker/supervisord.prod.conf /etc/supervisor/conf.d/supervisord.conf
COPY ./docker/entrypoint.prod.sh ./docker/entrypoint.sh

# Set executable permission to entrypoint
RUN chmod +x ./docker/entrypoint.sh

# Start supervisord via the entrypoint script
ENTRYPOINT [ "./docker/entrypoint.sh" ]
