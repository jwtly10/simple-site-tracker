FROM golang:1.21

WORKDIR /app

COPY . .

COPY .env_empty .env

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
