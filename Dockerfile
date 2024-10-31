FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files to the working directory and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Set environment variables
ENV DB_USER=admin \
    DB_PASSWORD=admin \
    DB_NAME=postgres \
    DB_HOST=db \
    DB_PORT=5432 \
    JWT_SECRET=secret

# Build the Go application from the cmd directory
RUN go build -o app ./cmd

# Expose port if necessary (adjust if your app uses a different port)
EXPOSE 8080

# Set the entry point to the compiled application
CMD ["./app"]
