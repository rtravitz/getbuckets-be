FROM golang:1.14 as dev

WORKDIR /app
COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

CMD CompileDaemon --build='go build cmd/server/main.go' --command='./main'
