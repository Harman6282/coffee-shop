up:
	docker-compose up -d
down:
	docker-compose down

gateway:
	cd services/api-gateway && go run ./cmd

user:
	cd services/user-service && go run .

order:
	cd services/order-service && go run .

migrate-up:
	migrate -database ${DB_URL} -path ${MIGRATIONS_PATH} up

migrate-down:
	migrate -database ${DB_URL} -path ${MIGRATIONS_PATH} down