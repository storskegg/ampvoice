# use official Golang image
FROM golang:1.26-alpine

# set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download and install the dependencies
RUN go mod vendor

# Build the Go app
RUN go build -o api-service ./cmd/api-service

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./api-service"]
