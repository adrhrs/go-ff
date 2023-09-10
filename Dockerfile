# Use an official Go runtime as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application files to the container's working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port your Go application listens on
EXPOSE 8080

# Define the command to run your application when the container starts
CMD ["./main"]
