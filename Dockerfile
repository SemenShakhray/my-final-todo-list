FROM golang:1.22.1

ENV TODO_PORT=7540
ENV TODO_DBFILE=data/scheduler.db
ENV TODO_PASSWORD=privet

WORKDIR /app

COPY . .

RUN go mod tidy

EXPOSE 7540

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /myapp

CMD ["./myapp"]