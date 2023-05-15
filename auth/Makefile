DB_URL=postgresql://postgres:password@localhost:5431/go_grpc?sslmode=disable

migrate_create:
	migrate create -ext sql -dir db/migrations -seq ${MIGRATE_NAME}
migrate_up:
	migrate -path db/migrations -database "${DB_URL}" -verbose up
migrate_down:
	migrate -path db/migrations -database "${DB_URL}" -verbose down
sqlc:
	sqlc generate
tests:
	go test -v -cover ./...
start:
	go run main.
mocks:
	mockgen -package mockdb --destination db/mocks/store.go github.com/ShadrackAdwera/go-gRPC/db/sqlc TxStore
	mockgen -package mockworkers --destination workers/mocks/distributor.go github.com/ShadrackAdwera/go-gRPC/workers Distributor
proto:
	rm -rf pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto
evans:
	evans --host localhost --port 5000 -r repl
pg_exec_it:
	docker container exec -it postgres psql -U postgres

.PHONY: migrate_create migrate_up migrate_down sqlc tests start proto evans
