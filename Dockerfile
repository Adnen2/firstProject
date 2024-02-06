# Use the official Golang image as the base image
FROM golang:1.21.6-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the application files into the container
COPY . .

# Build the application
RUN go build -mod=vendor -o main .

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
