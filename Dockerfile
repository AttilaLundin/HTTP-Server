FROM ubuntu:latest
LABEL authors="attila"

# Use the official Golang image to create a build artifact.
# This is a multi-stage build. In the first stage we build the binary.
FROM golang:latest as builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Clone the master branch of the repository.
RUN git clone https://github.com/AttilaLundin/HTTP-Server.git .

# Build the Go app.
RUN go build -o http-server .

# Use a Docker multi-stage build to create a lean production image.
# Start from scratch to keep the final image clean and small.
FROM scratch

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /app/http-server .

# This container exposes port 8080 to the outside world.
EXPOSE 8080

# Run the binary program produced by `go install`.
ENTRYPOINT ["./http-server"]

#docker run -p 8080:8080 -e PORT=8080 http-server