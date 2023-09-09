# Start from the official Go image
FROM golang:1.17

# Set the Current Working Directory inside the container
WORKDIR /app

# Update the package list and install bash
RUN apt update && apt install -y bash

# Copy everything from the current directory to the working directory
COPY . .

# build the project
RUN go build -o server .

EXPOSE 8080

# Command to run the application
CMD ["./server"]

