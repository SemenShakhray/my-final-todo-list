# Итоговое задания курса.
# Проект Go веб-сервера, который реализует функциональность простейшего планировщика задач

1. В браузере указывается адрес http://localhost:7540/. Используются переменные окружения TODO_PORT, TODO_DBFILE, TODO_PASSWORD

2. В tests/settings.go используются следующие параметры:
var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `` //Токен возвращается из api/signin и хранинтся в cookie "token"

3. Docker-образ создается командой 
``` go
docker build --tag my_app:v1 .
```

4. Запуск контейнера осуществляется командой
```go
docker run -d --name myapp -p 7540:7540 -e TODO_DBFILE=/data/scheduler.db -e TODO_PASSWORD=privet -v data/scheduler.db  my_app:v1
```

