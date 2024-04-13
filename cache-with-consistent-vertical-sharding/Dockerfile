# Use the official Go image as a parent image
FROM golang:1.21-alpine

# Set the Current Working Directory inside the container
WORKDIR /distributed-cache

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Command to run the executable
# CMD 'cd cache; go test -v ./... -bench=DataDistribution -run="^#"'
CMD ["tail", "-f", "/dev/null"]

