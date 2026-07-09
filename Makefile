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
