FROM golang:1.22.1

ENV TODO_PORT=7540
ENV TODO_DBFILE=data/scheduler.db
ENV TODO_PASSWORD=privet
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE ${TODO_PORT}

RUN go build -o /myapp

CMD ["./myapp"]