FROM golang:1.25.6

COPY . /app
WORKDIR /app

RUN apt-get update && apt-get install -y postgresql-client make

RUN go mod download
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]