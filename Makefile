start-db:
	docker compose -f docker-compose.dev.yaml up -d

stop-db:
	docker compose -f docker-compose.dev.yaml down

start:
	SERVER_ADDRESS=":8082" DATABASE_DSN="host=localhost port=5433 user=postgres password=postgres dbname=gcpfaqs-svc-accesss_test sslmode=disable" go run cmd/rest/main.go