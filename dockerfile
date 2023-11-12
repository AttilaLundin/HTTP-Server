
#golang image from latest the tag as base of the image
FROM golang:latest

#Set the working directory
WORKDIR /app

# Copy everything from the current directory to app
COPY . /app/

#Build the main
RUN go build main.go

#Listen to port 5431
EXPOSE 5431

# Command to run when the container start
CMD ["./main"]


