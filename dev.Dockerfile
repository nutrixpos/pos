FROM golang:alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Expose the necessary port (if any)
EXPOSE 8000

# Command to run your application
CMD ["go", "run", "."]