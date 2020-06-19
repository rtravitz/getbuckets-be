run:
	go run cmd/server/main.go --db-disable-tls=1

watch:
	CompileDaemon --build='go build -o server cmd/server/main.go' --command='./server'

# Admin
migrate:
	go run cmd/admin/main.go --db-disable-tls=1 migrate
