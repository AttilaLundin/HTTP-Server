FROM golang:latest

WORKDIR /app
COPY . /app/


RUN go build main.go
EXPOSE 5431

CMD ["./main"]


