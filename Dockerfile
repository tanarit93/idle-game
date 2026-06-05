# Use the official Golang image with Go 1.23
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy everything
COPY . .

# Run tidy and then the application at runtime to handle volume mounts
CMD ["sh", "-c", "go mod tidy && go run main.go"]
